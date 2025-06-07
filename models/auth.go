package models

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"
)

type Auth struct {
	ID                int       `gorm:"primary_key" json:"id"`
	Username          *string   `gorm:"size:50" json:"username"`
	Password          *string   `gorm:"size:60" json:"password"`
	Email             string    `gorm:"size:255;unique_index" json:"email"`
	FirstName         string    `gorm:"size:100" json:"first_name"`
	LastName          string    `gorm:"size:100" json:"last_name"`
	AvatarURL         string    `gorm:"size:500" json:"avatar_url"`
	Provider          string    `gorm:"size:50;default:'local'" json:"provider"`
	ProviderID        string    `gorm:"size:255" json:"provider_id"`
	IsEmailVerified   bool      `gorm:"default:false" json:"is_email_verified"`
	LastLogin         *time.Time `json:"last_login"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName sets the table name for the Auth model
func (Auth) TableName() string {
	return "blog_auth"
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

	// Check if password is nil (social login user)
	if auth.Password == nil {
		return false, nil
	}

	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(*auth.Password), []byte(password))
	if err != nil {
		return false, nil
	}

	return true, nil
}

// CheckAuthByEmail checks if authentication information exists by email and password is correct
func CheckAuthByEmail(email, password string) (bool, error) {
	var auth Auth
	err := db.Where("email = ?", email).First(&auth).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	// Check if password is nil (social login user)
	if auth.Password == nil {
		return false, nil
	}

	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(*auth.Password), []byte(password))
	if err != nil {
		return false, nil
	}

	return true, nil
}

// GetAuthByEmail retrieves auth record by email
func GetAuthByEmail(email string) (*Auth, error) {
	var auth Auth
	err := db.Where("email = ?", email).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// GetAuthByProvider retrieves auth record by provider and provider ID
func GetAuthByProvider(provider, providerID string) (*Auth, error) {
	var auth Auth
	err := db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// CreateSocialAuth creates a new auth record for social login
func CreateSocialAuth(email, firstName, lastName, avatarURL, provider, providerID string) (*Auth, error) {
	auth := Auth{
		Email:           email,
		FirstName:       firstName,
		LastName:        lastName,
		AvatarURL:       avatarURL,
		Provider:        provider,
		ProviderID:      providerID,
		IsEmailVerified: true, // Social logins are typically email verified
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err := db.Create(&auth).Error
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

// CreateLocalAuth creates a new auth record for local registration
func CreateLocalAuth(username, email, password string) (*Auth, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	auth := Auth{
		Username:        &username,
		Email:           email,
		Password:        &hashedPassword,
		Provider:        "local",
		IsEmailVerified: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = db.Create(&auth).Error
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

// UpdateLastLogin updates the last login timestamp
func (auth *Auth) UpdateLastLogin() error {
	now := time.Now()
	auth.LastLogin = &now
	auth.UpdatedAt = now
	return db.Save(auth).Error
}

// GetUserIdentifier returns the best identifier for the user (email or username)
func (auth *Auth) GetUserIdentifier() string {
	if auth.Email != "" {
		return auth.Email
	}
	if auth.Username != nil {
		return *auth.Username
	}
	return ""
}

// GetDisplayName returns the display name for the user
func (auth *Auth) GetDisplayName() string {
	if auth.FirstName != "" && auth.LastName != "" {
		return auth.FirstName + " " + auth.LastName
	}
	if auth.FirstName != "" {
		return auth.FirstName
	}
	if auth.Username != nil {
		return *auth.Username
	}
	return auth.Email
}
