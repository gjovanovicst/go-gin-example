package oauth_service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// OAuthProvider represents the OAuth provider configuration
type OAuthProvider struct {
	Config *oauth2.Config
	Name   string
}

// GoogleUserInfo represents Google user information
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// GitHubUserInfo represents GitHub user information
type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// FacebookUserInfo represents Facebook user information
type FacebookUserInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   struct {
		Data struct {
			Height       int    `json:"height"`
			IsSilhouette bool   `json:"is_silhouette"`
			URL          string `json:"url"`
			Width        int    `json:"width"`
		} `json:"data"`
	} `json:"picture"`
}

// SocialUserInfo standardized user information from social providers
type SocialUserInfo struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Name      string
	AvatarURL string
	Provider  string
}

// OAuthService manages OAuth providers
type OAuthService struct {
	providers map[string]*OAuthProvider
}

// NewOAuthService creates a new OAuth service
func NewOAuthService() *OAuthService {
	return &OAuthService{
		providers: make(map[string]*OAuthProvider),
	}
}

// AddGoogleProvider adds Google OAuth provider
func (s *OAuthService) AddGoogleProvider(clientID, clientSecret, redirectURL string) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	s.providers["google"] = &OAuthProvider{Config: config, Name: "google"}
}

// AddGitHubProvider adds GitHub OAuth provider
func (s *OAuthService) AddGitHubProvider(clientID, clientSecret, redirectURL string) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	s.providers["github"] = &OAuthProvider{Config: config, Name: "github"}
}

// AddFacebookProvider adds Facebook OAuth provider
func (s *OAuthService) AddFacebookProvider(clientID, clientSecret, redirectURL string) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}
	s.providers["facebook"] = &OAuthProvider{Config: config, Name: "facebook"}
}

// GetAuthURL returns the OAuth authorization URL for a provider
func (s *OAuthService) GetAuthURL(provider, state string) (string, error) {
	p, exists := s.providers[provider]
	if !exists {
		return "", fmt.Errorf("provider %s not configured", provider)
	}
	return p.Config.AuthCodeURL(state), nil
}

// ExchangeCodeForToken exchanges authorization code for token
func (s *OAuthService) ExchangeCodeForToken(provider, code string) (*oauth2.Token, error) {
	p, exists := s.providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider %s not configured", provider)
	}
	return p.Config.Exchange(context.Background(), code)
}

// GetUserInfo retrieves user information from the OAuth provider
func (s *OAuthService) GetUserInfo(provider string, token *oauth2.Token) (*SocialUserInfo, error) {
	switch provider {
	case "google":
		return s.getGoogleUserInfo(token)
	case "github":
		return s.getGitHubUserInfo(token)
	case "facebook":
		return s.getFacebookUserInfo(token)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// getGoogleUserInfo retrieves user info from Google
func (s *OAuthService) getGoogleUserInfo(token *oauth2.Token) (*SocialUserInfo, error) {
	client := s.providers["google"].Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	return &SocialUserInfo{
		ID:        userInfo.ID,
		Email:     userInfo.Email,
		FirstName: userInfo.GivenName,
		LastName:  userInfo.FamilyName,
		Name:      userInfo.Name,
		AvatarURL: userInfo.Picture,
		Provider:  "google",
	}, nil
}

// getGitHubUserInfo retrieves user info from GitHub
func (s *OAuthService) getGitHubUserInfo(token *oauth2.Token) (*SocialUserInfo, error) {
	client := s.providers["github"].Config.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GitHubUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	// GitHub might not return email in the user endpoint, so we try the emails endpoint
	if userInfo.Email == "" {
		email, _ := s.getGitHubUserEmail(client)
		userInfo.Email = email
	}

	// Parse name into first and last name
	firstName, lastName := parseFullName(userInfo.Name)

	return &SocialUserInfo{
		ID:        fmt.Sprintf("%d", userInfo.ID),
		Email:     userInfo.Email,
		FirstName: firstName,
		LastName:  lastName,
		Name:      userInfo.Name,
		AvatarURL: userInfo.AvatarURL,
		Provider:  "github",
	}, nil
}

// getGitHubUserEmail gets primary email from GitHub
func (s *OAuthService) getGitHubUserEmail(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}

	if err := json.Unmarshal(data, &emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	if len(emails) > 0 {
		return emails[0].Email, nil
	}

	return "", fmt.Errorf("no email found")
}

// getFacebookUserInfo retrieves user info from Facebook
func (s *OAuthService) getFacebookUserInfo(token *oauth2.Token) (*SocialUserInfo, error) {
	client := s.providers["facebook"].Config.Client(context.Background(), token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,first_name,last_name,picture")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo FacebookUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	return &SocialUserInfo{
		ID:        userInfo.ID,
		Email:     userInfo.Email,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		Name:      userInfo.Name,
		AvatarURL: userInfo.Picture.Data.URL,
		Provider:  "facebook",
	}, nil
}

// parseFullName splits a full name into first and last name
func parseFullName(fullName string) (string, string) {
	if fullName == "" {
		return "", ""
	}

	names := strings.Fields(fullName)
	if len(names) == 0 {
		return "", ""
	}
	if len(names) == 1 {
		return names[0], ""
	}

	return names[0], strings.Join(names[1:], " ")
}