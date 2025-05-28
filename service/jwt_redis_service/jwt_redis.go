package jwt_redis_service

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
)

const (
	JWT_TOKEN_PREFIX = "jwt_token:"
	JWT_TOKEN_TTL    = 3 * 60 * 60 // 3 hours in seconds
)

// StoreToken stores JWT token in Redis with expiration
func StoreToken(username string, token string) error {
	key := JWT_TOKEN_PREFIX + username
	return gredis.Set(key, token, JWT_TOKEN_TTL)
}

// GetToken retrieves JWT token from Redis
func GetToken(username string) (string, error) {
	key := JWT_TOKEN_PREFIX + username
	data, err := gredis.Get(key)
	if err != nil {
		return "", err
	}
	
	var token string
	err = json.Unmarshal(data, &token)
	if err != nil {
		return "", err
	}
	
	return token, nil
}

// DeleteToken removes JWT token from Redis
func DeleteToken(username string) error {
	key := JWT_TOKEN_PREFIX + username
	_, err := gredis.Delete(key)
	return err
}

// IsTokenValid checks if token exists and matches in Redis
func IsTokenValid(username string, token string) (bool, error) {
	storedToken, err := GetToken(username)
	if err != nil {
		if err == redis.ErrNil {
			return false, nil // Token not found
		}
		return false, err
	}
	return storedToken == token, nil
}

// RefreshToken updates the token and resets its expiration
func RefreshToken(username string, newToken string) error {
	return StoreToken(username, newToken)
}

// GetTokenTTL returns the remaining time-to-live for a token
func GetTokenTTL(username string) (time.Duration, error) {
	conn := gredis.RedisConn.Get()
	defer conn.Close()

	key := JWT_TOKEN_PREFIX + username
	ttl, err := redis.Int(conn.Do("TTL", key))
	if err != nil {
		return 0, err
	}

	if ttl == -1 {
		return 0, nil // Key exists but has no expiration
	} else if ttl == -2 {
		return 0, redis.ErrNil // Key doesn't exist
	}

	return time.Duration(ttl) * time.Second, nil
}