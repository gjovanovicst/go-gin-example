# JWT Token Storage in Redis and Bearer Token Authentication

## Overview
This implementation enhances the existing JWT authentication system by:
1. Storing JWT tokens in Redis with automatic expiration
2. Using Authorization header with Bearer token instead of query parameters
3. Adding logout functionality to invalidate tokens

## Changes Made

### 1. JWT Redis Service (`service/jwt_redis_service/jwt_redis.go`)
- **StoreToken**: Stores JWT tokens in Redis with 3-hour expiration
- **GetToken**: Retrieves JWT tokens from Redis
- **DeleteToken**: Removes JWT tokens from Redis (used for logout)
- **IsTokenValid**: Validates if a token exists and matches in Redis
- **RefreshToken**: Updates token and resets expiration
- **GetTokenTTL**: Returns remaining time-to-live for a token

### 2. Updated JWT Utilities (`pkg/util/jwt.go`)
- **GenerateToken**: Now stores generated tokens in Redis automatically
- **ParseToken**: Now validates tokens against Redis store in addition to JWT validation
- **InvalidateToken**: New function to remove tokens from Redis (logout)

### 3. Updated JWT Middleware (`middleware/jwt/jwt.go`)
- Changed from query parameter `?token=xxx` to Authorization header
- Now expects: `Authorization: Bearer <token>`
- More secure and follows standard practices

### 4. Enhanced Authentication API (`routers/api/auth.go`)
- **GetAuth**: Login endpoint (unchanged functionality, but now stores token in Redis)
- **Logout**: New endpoint to invalidate tokens
- Added proper imports for string manipulation

### 5. Router Updates (`routers/router.go`)
- Added logout route: `POST /auth/logout`

## API Usage

### Login
```bash
curl -X GET "http://localhost:8000/auth?username=test&password=test"
```
Response:
```json
{
  "code": 200,
  "msg": "ok",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Using Protected Endpoints
```bash
curl -X GET "http://localhost:8000/api/v1/tags" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Logout
```bash
curl -X POST "http://localhost:8000/auth/logout" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```
Response:
```json
{
  "code": 200,
  "msg": "ok",
  "data": {
    "message": "Successfully logged out"
  }
}
```

## Security Benefits

1. **Token Revocation**: Tokens can now be invalidated immediately via logout
2. **Redis Expiration**: Automatic token cleanup after 3 hours
3. **Header-based Authentication**: More secure than query parameters
4. **Token Validation**: Double validation (JWT + Redis existence)

## Redis Key Structure

Tokens are stored in Redis with the following key pattern:
- Key: `jwt_token:<hashed_username>`
- Value: The actual JWT token
- TTL: 3 hours (10800 seconds)

## Configuration

Ensure Redis is properly configured in `conf/app.ini`:
```ini
[redis]
Host = 127.0.0.1:6379
Password =
MaxIdle = 30
MaxActive = 30
IdleTimeout = 200
```

## Testing

The system maintains backward compatibility for the login endpoint but requires:
1. Redis to be running and properly configured
2. Authorization header format for protected endpoints
3. Proper error handling for token validation failures

All existing functionality remains intact while adding the new Redis-based token management and Bearer token authentication.