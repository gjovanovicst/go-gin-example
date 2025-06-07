package api

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/config"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/EDDYCJY/go-gin-example/service/oauth_service"
)

type auth struct {
	Username string `valid:"MaxSize(50)"`
	Email    string `valid:"Email; MaxSize(255)"`
	Password string `valid:"Required; MaxSize(50)"`
}

type registerRequest struct {
	Username string `json:"username" valid:"Required; MaxSize(50)"`
	Email    string `json:"email" valid:"Required; Email; MaxSize(255)"`
	Password string `json:"password" valid:"Required; MinSize(6); MaxSize(50)"`
}

var oauthService *oauth_service.OAuthService

// InitOAuthService initializes OAuth providers
func InitOAuthService() {
	oauthService = oauth_service.NewOAuthService()
	
	// Load OAuth configuration
	cfg := config.GetOAuthConfig()
	
	// Add enabled providers
	if cfg.Google.Enabled {
		oauthService.AddGoogleProvider(cfg.Google.ClientID, cfg.Google.ClientSecret, cfg.Google.RedirectURL)
	}
	if cfg.GitHub.Enabled {
		oauthService.AddGitHubProvider(cfg.GitHub.ClientID, cfg.GitHub.ClientSecret, cfg.GitHub.RedirectURL)
	}
	if cfg.Facebook.Enabled {
		oauthService.AddFacebookProvider(cfg.Facebook.ClientID, cfg.Facebook.ClientSecret, cfg.Facebook.RedirectURL)
	}
}

// @Summary Login
// @Accept application/json
// @Produce  json
// @Param credentials body auth true "Login credentials (username/email and password)"
// @Success 200 {object} map[string]interface{} "{"access_token": "jwt_token", "token_type": "Bearer", "user": {...}}"
// @Failure 400 {object} app.Response
// @Failure 401 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [post]
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	var credentials auth
	if err := c.ShouldBindJSON(&credentials); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	// Validate that either username or email is provided
	if credentials.Username == "" && credentials.Email == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{
			"error": "Username or email is required",
		})
		return
	}

	if credentials.Password == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{
			"error": "Password is required",
		})
		return
	}

	ok, _ := valid.Valid(&credentials)
	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{
		Username: credentials.Username,
		Email:    credentials.Email,
		Password: credentials.Password,
	}
	
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	// Get user info for response
	identifier := credentials.Username
	if identifier == "" {
		identifier = credentials.Email
	}

	var user *models.Auth
	if credentials.Email != "" {
		user, _ = auth_service.GetUserByEmail(credentials.Email)
	} else {
		// Note: You might want to add GetUserByUsername function
		user, _ = auth_service.GetUserByEmail(credentials.Username) // Fallback
	}

	token, err := util.GenerateToken(identifier, credentials.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	// Update last login
	if user != nil {
		user.UpdateLastLogin()
	}

	response := map[string]interface{}{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   3600,
	}

	if user != nil {
		response["user"] = map[string]interface{}{
			"id":           user.ID,
			"email":        user.Email,
			"username":     user.Username,
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"avatar_url":   user.AvatarURL,
			"provider":     user.Provider,
			"display_name": user.GetDisplayName(),
		}
	}

	appG.Response(http.StatusOK, e.SUCCESS, response)
}

// @Summary Logout
// @Produce  json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} app.Response
// @Failure 401 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth/logout [post]
func Logout(c *gin.Context) {
	appG := app.Gin{C: c}
	
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		appG.Response(http.StatusUnauthorized, e.INVALID_PARAMS, nil)
		return
	}

	// Extract token from "Bearer <token>"
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		appG.Response(http.StatusUnauthorized, e.INVALID_PARAMS, nil)
		return
	}

	// Parse token to get username
	claims, err := util.ParseToken(token)
	if err != nil {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	// Invalidate token in Redis
	err = util.InvalidateToken(claims.Username)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"message": "Successfully logged out",
	})
}

// @Summary Register new user
// @Accept application/json
// @Produce  json
// @Param user body registerRequest true "User registration data"
// @Success 201 {object} map[string]interface{} "{"access_token": "jwt_token", "token_type": "Bearer", "user": {...}}"
// @Failure 400 {object} app.Response
// @Failure 409 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth/register [post]
func Register(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	ok, _ := valid.Valid(&req)
	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	// Check if user already exists
	existingUser, _ := auth_service.GetUserByEmail(req.Email)
	if existingUser != nil {
		appG.Response(http.StatusConflict, e.ERROR_EXIST_TAG, map[string]string{
			"error": "User with this email already exists",
		})
		return
	}

	// Create new user
	user, err := auth_service.RegisterLocalUser(req.Username, req.Email, req.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	// Generate token
	token, err := util.GenerateToken(req.Email, req.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusCreated, e.SUCCESS, map[string]interface{}{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   3600,
		"user": map[string]interface{}{
			"id":           user.ID,
			"email":        user.Email,
			"username":     user.Username,
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"avatar_url":   user.AvatarURL,
			"provider":     user.Provider,
			"display_name": user.GetDisplayName(),
		},
	})
}

// @Summary Get OAuth authorization URL
// @Produce  json
// @Param provider path string true "OAuth provider (google, github, facebook)"
// @Success 200 {object} map[string]string "{"auth_url": "https://..."}"
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth/{provider} [get]
func GetOAuthURL(c *gin.Context) {
	appG := app.Gin{C: c}
	provider := c.Param("provider")

	if oauthService == nil {
		InitOAuthService()
	}

	// Generate random state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	// Store state in session or cache (simplified for demo)
	// In production, you should store this securely
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	authURL, err := oauthService.GetAuthURL(provider, state)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{
			"error": err.Error(),
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"auth_url": authURL,
	})
}

// @Summary Handle OAuth callback
// @Produce  json
// @Param provider path string true "OAuth provider (google, github, facebook)"
// @Param code query string true "Authorization code"
// @Param state query string true "State parameter"
// @Success 200 {object} map[string]interface{} "{"access_token": "jwt_token", "token_type": "Bearer", "user": {...}}"
// @Failure 400 {object} app.Response
// @Failure 401 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth/{provider}/callback [get]
func OAuthCallback(c *gin.Context) {
	appG := app.Gin{C: c}
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	if oauthService == nil {
		InitOAuthService()
	}

	// Verify state parameter (CSRF protection)
	storedState, err := c.Cookie("oauth_state")
	if err != nil || storedState != state {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{
			"error": "Invalid state parameter",
		})
		return
	}

	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Exchange code for token
	token, err := oauthService.ExchangeCodeForToken(provider, code)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_AUTH_TOKEN, map[string]string{
			"error": "Failed to exchange code for token",
		})
		return
	}

	// Get user info from OAuth provider
	userInfo, err := oauthService.GetUserInfo(provider, token)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, map[string]string{
			"error": "Failed to get user info",
		})
		return
	}

	// Create or get user
	socialAuth := &auth_service.SocialAuth{
		Email:      userInfo.Email,
		FirstName:  userInfo.FirstName,
		LastName:   userInfo.LastName,
		AvatarURL:  userInfo.AvatarURL,
		Provider:   userInfo.Provider,
		ProviderID: userInfo.ID,
	}

	user, err := auth_service.AuthenticateOrCreateSocialUser(socialAuth)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, map[string]string{
			"error": "Failed to authenticate user",
		})
		return
	}

	// Generate JWT token
	jwtToken, err := util.GenerateToken(user.GetUserIdentifier(), "")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"access_token": jwtToken,
		"token_type":   "Bearer",
		"expires_in":   3600,
		"user": map[string]interface{}{
			"id":           user.ID,
			"email":        user.Email,
			"username":     user.Username,
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"avatar_url":   user.AvatarURL,
			"provider":     user.Provider,
			"display_name": user.GetDisplayName(),
		},
	})
}

// generateRandomState generates a random state string for OAuth CSRF protection
func generateRandomState() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
