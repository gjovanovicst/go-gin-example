package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

func main() {
	// Load configuration
	setting.Setup()

	// Create database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("=== Database Verification ===")
	
	// Check auth table
	fmt.Println("\n--- Auth Table Data ---")
	rows, err := db.Query("SELECT id, username, password FROM blog_auth")
	if err != nil {
		log.Printf("Failed to query auth table: %v", err)
	} else {
		defer rows.Close()
		count := 0
		for rows.Next() {
			var id int
			var username, password string
			if err := rows.Scan(&id, &username, &password); err != nil {
				log.Printf("Failed to scan row: %v", err)
				continue
			}
			fmt.Printf("ID: %d, Username: %s, Password: %s\n", id, username, password)
			count++
		}
		fmt.Printf("Total auth records: %d\n", count)
	}

	// Check tags table
	fmt.Println("\n--- Tags Table Data ---")
	rows, err = db.Query("SELECT id, name, created_by, state FROM blog_tag")
	if err != nil {
		log.Printf("Failed to query tags table: %v", err)
	} else {
		defer rows.Close()
		count := 0
		for rows.Next() {
			var id int
			var name, createdBy string
			var state int
			if err := rows.Scan(&id, &name, &createdBy, &state); err != nil {
				log.Printf("Failed to scan row: %v", err)
				continue
			}
			fmt.Printf("ID: %d, Name: %s, Created By: %s, State: %d\n", id, name, createdBy, state)
			count++
		}
		fmt.Printf("Total tag records: %d\n", count)
	}

	// Check tables exist
	fmt.Println("\n--- Available Tables ---")
	rows, err = db.Query("SHOW TABLES")
	if err != nil {
		log.Printf("Failed to show tables: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err != nil {
				log.Printf("Failed to scan table name: %v", err)
				continue
			}
			fmt.Printf("Table: %s\n", tableName)
		}
	}
}