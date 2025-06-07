package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func main() {
	// Wait a moment for server to start
	time.Sleep(2 * time.Second)
	
	fmt.Println("Testing /auth endpoint...")
	
	// Test with correct credentials
	fmt.Println("\n=== Testing with correct credentials ===")
	testLogin("admin", "admin123")
	testLogin("testuser", "test123")
	
	// Test with incorrect credentials
	fmt.Println("\n=== Testing with incorrect credentials ===")
	testLogin("admin", "wrongpassword")
	testLogin("nonexistent", "password")
	testLogin("random", "random")
}

func testLogin(username, password string) {
	// Prepare form data
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	
	// Make POST request
	resp, err := http.PostForm("http://localhost:8000/auth", data)
	if err != nil {
		fmt.Printf("Error making request for %s/%s: %v\n", username, password, err)
		return
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response for %s/%s: %v\n", username, password, err)
		return
	}
	
	// Parse JSON response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("Error parsing JSON for %s/%s: %v\n", username, password, err)
		return
	}
	
	// Pretty print result
	prettyJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s/%s -> Status: %d\n%s\n", username, password, resp.StatusCode, string(prettyJSON))
}