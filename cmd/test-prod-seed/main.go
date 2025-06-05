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

	fmt.Println("=== Testing Production Seeds ===")

	// Production auth seed (minimal)
	authSQL := `INSERT INTO blog_auth (id, username, password) VALUES (1, 'admin', 'admin123')`

	fmt.Println("Executing production auth seed...")
	_, err = db.Exec(authSQL)
	if err != nil {
		log.Printf("Failed to execute auth seed: %v", err)
	} else {
		fmt.Println("✅ Production auth seed executed successfully")
	}

	// Production tags seed (minimal)
	tagsSQL := `INSERT INTO blog_tag (name, created_on, created_by, state) VALUES ('General', UNIX_TIMESTAMP(), 'system', 1)`

	fmt.Println("Executing production tags seed...")
	_, err = db.Exec(tagsSQL)
	if err != nil {
		log.Printf("Failed to execute tags seed: %v", err)
	} else {
		fmt.Println("✅ Production tags seed executed successfully")
	}

	fmt.Println("\n=== Production Environment Summary ===")
	fmt.Println("• 1 admin user (essential for management)")
	fmt.Println("• 1 general tag (minimal category structure)")
	fmt.Println("• No test data (production-safe)")
}