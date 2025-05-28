package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/EDDYCJY/go-gin-example/service/jwt_redis_service"
)

var jwtSecret []byte

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth and store in Redis
func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		EncodeMD5(username),
		EncodeMD5(password),
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	// Store token in Redis
	err = jwt_redis_service.StoreToken(EncodeMD5(username), token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ParseToken parsing token and validate against Redis
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			// Validate token against Redis store
			isValid, err := jwt_redis_service.IsTokenValid(claims.Username, token)
			if err != nil || !isValid {
				return nil, jwt.NewValidationError("token not found in store or invalid", jwt.ValidationErrorClaimsInvalid)
			}
			return claims, nil
		}
	}

	return nil, err
}

// InvalidateToken removes token from Redis
func InvalidateToken(username string) error {
	return jwt_redis_service.DeleteToken(username)
}
