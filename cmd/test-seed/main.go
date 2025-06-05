package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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

	fmt.Println("=== Testing Direct SQL Execution ===")

	// Test auth seed
	authSQL := `INSERT INTO blog_auth (id, username, password) VALUES 
(1, 'admin', 'admin123'),
(2, 'testuser', 'test123'),
(3, 'developer', 'dev123')`

	fmt.Println("Executing auth seed...")
	_, err = db.Exec(authSQL)
	if err != nil {
		log.Printf("Failed to execute auth seed: %v", err)
	} else {
		fmt.Println("✅ Auth seed executed successfully")
	}

	// Test tags seed
	tagsSQL := `INSERT INTO blog_tag (name, created_on, created_by, state) VALUES 
('Technology', UNIX_TIMESTAMP(), 'system', 1),
('Programming', UNIX_TIMESTAMP(), 'system', 1),
('Go', UNIX_TIMESTAMP(), 'system', 1),
('Web Development', UNIX_TIMESTAMP(), 'system', 1),
('API', UNIX_TIMESTAMP(), 'system', 1),
('Testing', UNIX_TIMESTAMP(), 'system', 1)`

	fmt.Println("Executing tags seed...")
	_, err = db.Exec(tagsSQL)
	if err != nil {
		log.Printf("Failed to execute tags seed: %v", err)
	} else {
		fmt.Println("✅ Tags seed executed successfully")
	}

	// Now test file reading and parsing
	fmt.Println("\n=== Testing File Reading ===")
	
	content, err := os.ReadFile("seeds/development/001_seed_auth.up.sql")
	if err != nil {
		log.Printf("Failed to read seed file: %v", err)
		return
	}

	fmt.Printf("File content:\n%s\n", string(content))
	
	// Split content by semicolons to handle multiple statements
	statements := strings.Split(string(content), ";")
	
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			fmt.Printf("Skipping statement %d: empty or comment\n", i)
			continue
		}
		
		fmt.Printf("Executing statement %d: %s\n", i, stmt[:min(50, len(stmt))])
		_, err = db.Exec(stmt)
		if err != nil {
			log.Printf("Failed to execute statement %d: %v\nStatement: %s", i, err, stmt)
		} else {
			fmt.Printf("✅ Statement %d executed successfully\n", i)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}