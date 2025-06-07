package config

import (
	"os"
)

// OAuthConfig holds OAuth provider configurations
type OAuthConfig struct {
	Google   ProviderConfig
	GitHub   ProviderConfig
	Facebook ProviderConfig
}

// ProviderConfig holds configuration for a single OAuth provider
type ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Enabled      bool
}

// GetOAuthConfig returns OAuth configuration from environment variables
func GetOAuthConfig() *OAuthConfig {
	baseURL := getEnvOrDefault("BASE_URL", "http://localhost:8000")
	
	return &OAuthConfig{
		Google: ProviderConfig{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  baseURL + "/auth/google/callback",
			Enabled:      os.Getenv("GOOGLE_CLIENT_ID") != "",
		},
		GitHub: ProviderConfig{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  baseURL + "/auth/github/callback",
			Enabled:      os.Getenv("GITHUB_CLIENT_ID") != "",
		},
		Facebook: ProviderConfig{
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
			RedirectURL:  baseURL + "/auth/facebook/callback",
			Enabled:      os.Getenv("FACEBOOK_CLIENT_ID") != "",
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}