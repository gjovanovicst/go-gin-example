package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	passwords := []string{"staging_admin_2024", "staging_test_123"}
	
	for _, password := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Error hashing %s: %v\n", password, err)
			continue
		}
		fmt.Printf("%s -> %s\n", password, string(hash))
	}
}