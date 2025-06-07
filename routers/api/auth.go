package api

import (
	"net/http"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary Login
// @Accept application/x-www-form-urlencoded
// @Produce  json
// @Param username formData string true "userName"
// @Param password formData string true "password"
// @Param grant_type formData string false "Grant Type" default(password)
// @Success 200 {object} map[string]interface{} "{"access_token": "jwt_token", "token_type": "Bearer"}"
// @Failure 400 {object} app.Response
// @Failure 401 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [post]
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	// Handle both form data and JSON
	var username, password string
	
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		username = c.PostForm("username")
		password = c.PostForm("password")
	} else {
		username = c.PostForm("username")
		password = c.PostForm("password")
		// Also try JSON format as fallback
		if username == "" || password == "" {
			username = c.Query("username")
			password = c.Query("password")
		}
	}

	a := auth{Username: username, Password: password}
	
	ok, _ := valid.Valid(&a)

	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"access_token": token,
		"token_type": "Bearer",
		"expires_in": 3600,
	})
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
