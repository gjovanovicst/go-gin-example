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

	fmt.Println("=== Clearing existing seed data ===")
	
	// Clear auth table
	_, err = db.Exec("DELETE FROM blog_auth")
	if err != nil {
		log.Printf("Failed to clear auth table: %v", err)
	} else {
		fmt.Println("✅ Cleared blog_auth table")
	}

	// Clear tags table
	_, err = db.Exec("DELETE FROM blog_tag")
	if err != nil {
		log.Printf("Failed to clear tags table: %v", err)
	} else {
		fmt.Println("✅ Cleared blog_tag table")
	}

	// Clear articles table
	_, err = db.Exec("DELETE FROM blog_article")
	if err != nil {
		log.Printf("Failed to clear articles table: %v", err)
	} else {
		fmt.Println("✅ Cleared blog_article table")
	}

	fmt.Println("✅ All seed data cleared successfully!")
}