package main

import (
	"database/sql"
	"fmt"
	"log"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	_ "modernc.org/sqlite"
)

func main() {
	// 连接数据库
	db, err := sql.Open("sqlite", "miko_email.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// 设置数据库编码
	_, err = db.Exec("PRAGMA encoding = 'UTF-8'")
	if err != nil {
		log.Printf("Failed to set UTF-8 encoding: %v", err)
	}

	// 查询所有邮件
	rows, err := db.Query(`
		SELECT id, subject, body 
		FROM emails 
		ORDER BY id
	`)
	if err != nil {
		log.Fatal("Failed to query emails:", err)
	}
	defer rows.Close()

	fmt.Println("=== 开始修复所有邮件编码 ===")

	// GBK到UTF-8的转换器
	decoder := simplifiedchinese.GBK.NewDecoder()

	fixedCount := 0

	for rows.Next() {
		var id int
		var subject, body string

		err = rows.Scan(&id, &subject, &body)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		needsFix := false
		var fixedSubject, fixedBody string

		// 检查主题是否需要修复
		if !utf8.ValidString(subject) || containsReplacementChar(subject) {
			// 尝试从GBK解码
			subjectBytes := []byte(subject)
			utf8Subject, _, err := transform.Bytes(decoder, subjectBytes)
			if err != nil {
				fmt.Printf("邮件ID %d 主题解码失败: %v\n", id, err)
				fixedSubject = subject // 保持原样
			} else {
				fixedSubject = string(utf8Subject)
				if fixedSubject != subject {
					needsFix = true
					fmt.Printf("邮件ID %d 主题修复: %q -> %q\n", id, subject, fixedSubject)
				} else {
					fixedSubject = subject
				}
			}
		} else {
			fixedSubject = subject
		}

		// 检查内容是否需要修复
		if !utf8.ValidString(body) || containsReplacementChar(body) {
			// 尝试从GBK解码
			bodyBytes := []byte(body)
			utf8Body, _, err := transform.Bytes(decoder, bodyBytes)
			if err != nil {
				fmt.Printf("邮件ID %d 内容解码失败: %v\n", id, err)
				fixedBody = body // 保持原样
			} else {
				fixedBody = string(utf8Body)
				if fixedBody != body {
					needsFix = true
					bodyPreview := body
					if len(body) > 50 {
						bodyPreview = body[:50] + "..."
					}
					fixedBodyPreview := fixedBody
					if len(fixedBody) > 50 {
						fixedBodyPreview = fixedBody[:50] + "..."
					}
					fmt.Printf("邮件ID %d 内容修复: %q -> %q\n", id, bodyPreview, fixedBodyPreview)
				} else {
					fixedBody = body
				}
			}
		} else {
			fixedBody = body
		}

		// 如果需要修复，更新数据库
		if needsFix {
			_, err = db.Exec(`
				UPDATE emails 
				SET subject = ?, body = ?, updated_at = datetime('now')
				WHERE id = ?
			`, fixedSubject, fixedBody, id)

			if err != nil {
				log.Printf("Failed to update email %d: %v", id, err)
			} else {
				fmt.Printf("邮件ID %d 修复完成\n", id)
				fixedCount++
			}
		}
	}

	fmt.Printf("\n=== 修复完成，共修复 %d 封邮件 ===\n", fixedCount)

	// 验证修复结果
	fmt.Println("\n=== 验证修复结果 ===")
	rows2, err := db.Query(`
		SELECT id, subject, body 
		FROM emails 
		ORDER BY created_at DESC 
		LIMIT 10
	`)
	if err != nil {
		log.Fatal("Failed to query emails for verification:", err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var id int
		var subject, body string

		err = rows2.Scan(&id, &subject, &body)
		if err != nil {
			log.Printf("Failed to scan verification row: %v", err)
			continue
		}

		fmt.Printf("\n验证邮件ID: %d\n", id)
		fmt.Printf("主题: %s (UTF8有效: %v, 无乱码: %v)\n", subject, utf8.ValidString(subject), !containsReplacementChar(subject))
		bodyPreview := body
		if len(body) > 50 {
			bodyPreview = body[:50] + "..."
		}
		fmt.Printf("内容: %s (UTF8有效: %v, 无乱码: %v)\n", bodyPreview, utf8.ValidString(body), !containsReplacementChar(body))
	}
}

// containsReplacementChar 检查字符串是否包含Unicode替换字符（乱码标志）
func containsReplacementChar(s string) bool {
	for _, r := range s {
		if r == '\uFFFD' { // Unicode替换字符
			return true
		}
	}
	return false
}
