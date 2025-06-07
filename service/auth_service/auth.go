package auth_service

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/jinzhu/gorm"
)

type Auth struct {
	Username string
	Email    string
	Password string
}

type SocialAuth struct {
	Email      string
	FirstName  string
	LastName   string
	AvatarURL  string
	Provider   string
	ProviderID string
}

// Check validates traditional username/password authentication
func (a *Auth) Check() (bool, error) {
	if a.Username != "" {
		return models.CheckAuth(a.Username, a.Password)
	}
	if a.Email != "" {
		return models.CheckAuthByEmail(a.Email, a.Password)
	}
	return false, nil
}

// GetUserByEmail retrieves user by email
func GetUserByEmail(email string) (*models.Auth, error) {
	return models.GetAuthByEmail(email)
}

// GetUserByProvider retrieves user by social provider
func GetUserByProvider(provider, providerID string) (*models.Auth, error) {
	return models.GetAuthByProvider(provider, providerID)
}

// CreateSocialUser creates a new user from social login
func (sa *SocialAuth) CreateUser() (*models.Auth, error) {
	return models.CreateSocialAuth(sa.Email, sa.FirstName, sa.LastName, sa.AvatarURL, sa.Provider, sa.ProviderID)
}

// RegisterLocalUser creates a new local user
func RegisterLocalUser(username, email, password string) (*models.Auth, error) {
	return models.CreateLocalAuth(username, email, password)
}

// AuthenticateOrCreateSocialUser handles social login authentication
func AuthenticateOrCreateSocialUser(socialAuth *SocialAuth) (*models.Auth, error) {
	// First try to find user by provider
	user, err := GetUserByProvider(socialAuth.Provider, socialAuth.ProviderID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// If user found, update login time and return
	if user != nil {
		err = user.UpdateLastLogin()
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	// Try to find user by email (in case they registered locally first)
	user, err = GetUserByEmail(socialAuth.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// If user exists with same email but different provider, we could link accounts
	// For now, we'll create a new user
	if user == nil {
		// Create new social user
		user, err = socialAuth.CreateUser()
		if err != nil {
			return nil, err
		}
	}

	// Update login time
	err = user.UpdateLastLogin()
	if err != nil {
		return nil, err
	}

	return user, nil
}
