package main

import (
	"fmt"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

func init() {
	setting.Setup()
	models.Setup()
}

func main() {
	fmt.Println("Testing authentication...")
	
	// Test with correct credentials
	fmt.Println("\n=== Testing with correct credentials ===")
	isValid, err := models.CheckAuth("admin", "admin123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("admin/admin123 -> Valid: %v\n", isValid)
	}
	
	isValid, err = models.CheckAuth("testuser", "test123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("testuser/test123 -> Valid: %v\n", isValid)
	}
	
	// Test with incorrect credentials
	fmt.Println("\n=== Testing with incorrect credentials ===")
	isValid, err = models.CheckAuth("admin", "wrongpassword")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("admin/wrongpassword -> Valid: %v\n", isValid)
	}
	
	isValid, err = models.CheckAuth("nonexistent", "password")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("nonexistent/password -> Valid: %v\n", isValid)
	}
	
	// Test with any random credentials (this should fail now)
	fmt.Println("\n=== Testing with random credentials ===")
	isValid, err = models.CheckAuth("random", "random")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("random/random -> Valid: %v\n", isValid)
	}
}