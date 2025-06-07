package models

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"
)

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `gorm:"size:50" json:"username"`
	Password string `gorm:"size:60" json:"password"`
}

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckAuth checks if authentication information exists and password is correct
func CheckAuth(username, password string) (bool, error) {
	var auth Auth
	err := db.Where("username = ?", username).First(&auth).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(password))
	if err != nil {
		return false, nil
	}

	return true, nil
}
