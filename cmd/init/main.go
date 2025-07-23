package main

import (
	"database/sql"
	"log"
	"os"

	"miko-email/internal/database"
	_ "modernc.org/sqlite"
)

func main() {
	// 初始化数据库
	db, err := database.Init("./miko_email.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// 执行初始化域名脚本
	if err := initDomains(db); err != nil {
		log.Fatal("Failed to initialize domains:", err)
	}

	log.Println("Database initialized successfully!")
}

func initDomains(db *sql.DB) error {
	// 读取初始化脚本
	sqlContent, err := os.ReadFile("scripts/init_domains.sql")
	if err != nil {
		return err
	}

	// 执行SQL脚本
	_, err = db.Exec(string(sqlContent))
	return err
}
