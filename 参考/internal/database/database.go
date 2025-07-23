package database

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Init 初始化数据库
func Init(dbPath, adminEmail, adminPassword string) (*sql.DB, error) {
	// SQLite默认使用UTF-8编码
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("创建表失败: %v", err)
	}

	// 运行数据库迁移
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("运行数据库迁移失败: %v", err)
	}

	if err := migrateDatabase(db); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %v", err)
	}

	if err := createDefaultAdmin(db, adminEmail, adminPassword); err != nil {
		return nil, fmt.Errorf("创建默认管理员失败: %v", err)
	}

	return db, nil
}

// createTables 创建数据库表
func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			plain_password TEXT,
			name TEXT NOT NULL,
			is_admin BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS emails (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			message_id TEXT UNIQUE NOT NULL,
			from_addr TEXT NOT NULL,
			to_addr TEXT NOT NULL,
			subject TEXT NOT NULL,
			body TEXT NOT NULL,
			html_body TEXT,
			is_read BOOLEAN DEFAULT FALSE,
			is_deleted BOOLEAN DEFAULT FALSE,
			is_sent BOOLEAN DEFAULT FALSE,
			user_id INTEGER,
			size INTEGER DEFAULT 0,
			attachments TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		`CREATE TABLE IF NOT EXISTS domains (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			dns_verified BOOLEAN DEFAULT FALSE,
			mx_record TEXT DEFAULT '',
			last_verified DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS mailboxes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			email TEXT UNIQUE NOT NULL,
			domain_id INTEGER NOT NULL,
			password TEXT NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			is_current BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
		)`,

		`CREATE INDEX IF NOT EXISTS idx_emails_user_id ON emails(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_emails_to_addr ON emails(to_addr)`,
		`CREATE INDEX IF NOT EXISTS idx_emails_from_addr ON emails(from_addr)`,
		`CREATE INDEX IF NOT EXISTS idx_emails_created_at ON emails(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_mailboxes_user_id ON mailboxes(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_mailboxes_email ON mailboxes(email)`,
		`CREATE INDEX IF NOT EXISTS idx_mailboxes_domain_id ON mailboxes(domain_id)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("执行SQL失败 [%s]: %v", query, err)
		}
	}

	return nil
}

// migrateDatabase 数据库迁移
func migrateDatabase(db *sql.DB) error {
	// 检查是否需要添加plain_password字段
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		return err
	}
	defer rows.Close()

	hasPlainPassword := false
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var defaultValue interface{}

		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		if err != nil {
			return err
		}

		if name == "plain_password" {
			hasPlainPassword = true
			break
		}
	}

	// 如果没有plain_password字段，添加它
	if !hasPlainPassword {
		_, err = db.Exec("ALTER TABLE users ADD COLUMN plain_password TEXT")
		if err != nil {
			return fmt.Errorf("添加plain_password字段失败: %v", err)
		}
	}

	return nil
}

// createDefaultAdmin 创建默认管理员账户
func createDefaultAdmin(db *sql.DB, adminEmail, adminPassword string) error {
	// 检查是否已存在管理员
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE is_admin = TRUE").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // 已存在管理员
	}

	// 使用传入的管理员信息创建默认管理员
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO users (email, password, name, is_admin, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, adminEmail, string(hashedPassword), "管理员", true, time.Now(), time.Now())

	return err
}

// runMigrations 运行数据库迁移
func runMigrations(db *sql.DB) error {
	// 检查是否需要为 mailboxes 表添加 password 字段
	var columnExists bool
	err := db.QueryRow(`
		SELECT COUNT(*) > 0
		FROM pragma_table_info('mailboxes')
		WHERE name = 'password'
	`).Scan(&columnExists)

	if err != nil {
		return fmt.Errorf("检查 password 字段失败: %v", err)
	}

	if !columnExists {
		// 添加 password 字段
		_, err = db.Exec(`ALTER TABLE mailboxes ADD COLUMN password TEXT DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("添加 password 字段失败: %v", err)
		}

		// 为现有的邮箱生成随机密码
		rows, err := db.Query(`SELECT id FROM mailboxes WHERE password = '' OR password IS NULL`)
		if err != nil {
			return fmt.Errorf("查询现有邮箱失败: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				continue
			}

			// 生成随机密码
			password := generateRandomPasswordForMigration(12)
			_, err = db.Exec(`UPDATE mailboxes SET password = ? WHERE id = ?`, password, id)
			if err != nil {
				return fmt.Errorf("更新邮箱密码失败: %v", err)
			}
		}
	}

	// 检查是否需要为 domains 表添加 user_id 字段
	var userIdColumnExists bool
	err = db.QueryRow(`
		SELECT COUNT(*) > 0
		FROM pragma_table_info('domains')
		WHERE name = 'user_id'
	`).Scan(&userIdColumnExists)

	if err != nil {
		return fmt.Errorf("检查 user_id 字段失败: %v", err)
	}

	if !userIdColumnExists {
		// 添加 user_id 字段
		_, err = db.Exec(`ALTER TABLE domains ADD COLUMN user_id INTEGER`)
		if err != nil {
			return fmt.Errorf("添加 user_id 字段失败: %v", err)
		}

		// 添加 description 字段
		_, err = db.Exec(`ALTER TABLE domains ADD COLUMN description TEXT DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("添加 description 字段失败: %v", err)
		}
	}

	return nil
}

// generateRandomPasswordForMigration 为数据库迁移生成随机密码
func generateRandomPasswordForMigration(length int) string {
	const (
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		numbers   = "0123456789"
	)

	allChars := uppercase + lowercase + numbers
	password := make([]byte, length)

	// 确保至少包含每种类型的字符
	charSets := []string{uppercase, lowercase, numbers}
	for i, charset := range charSets {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[n.Int64()]
	}

	// 填充剩余长度
	for i := 3; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		password[i] = allChars[n.Int64()]
	}

	// 打乱密码字符顺序
	for i := length - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password)
}
