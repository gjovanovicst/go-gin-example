package jwt

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			code = e.INVALID_PARAMS
		} else {
			// Check if header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				code = e.INVALID_PARAMS
			} else {
				// Extract token from "Bearer <token>"
				token := strings.TrimPrefix(authHeader, "Bearer ")
				if token == "" {
					code = e.INVALID_PARAMS
				} else {
					_, err := util.ParseToken(token)
					if err != nil {
						// Check if it's a JWT validation error
						if ve, ok := err.(*jwt.ValidationError); ok {
							switch ve.Errors {
							case jwt.ValidationErrorExpired:
								code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
							default:
								code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
							}
						} else {
							// Handle other types of errors (Redis, etc.)
							code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
						}
					}
				}
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
