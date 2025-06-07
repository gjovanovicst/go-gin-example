# Social Authentication Implementation

This document describes the implementation of social network login support for the Go Gin example application.

## Overview

The authentication system has been modernized to support:
- Email-based login (in addition to username)
- Social network authentication (Google, GitHub, Facebook)
- Enhanced user profiles with additional fields
- OAuth 2.0 integration

## Database Changes

### New Auth Table Schema

The `blog_auth` table has been upgraded with the following new columns:

- `email` (varchar 255) - User email address (unique)
- `first_name` (varchar 100) - User first name
- `last_name` (varchar 100) - User last name
- `avatar_url` (varchar 500) - Profile picture URL
- `provider` (varchar 50) - Authentication provider (local, google, github, facebook)
- `provider_id` (varchar 255) - Provider-specific user ID
- `is_email_verified` (boolean) - Email verification status
- `last_login` (timestamp) - Last login timestamp
- `created_at` (timestamp) - Account creation time
- `updated_at` (timestamp) - Last update time

### Migration

Run the migration to upgrade your database:

```bash
# Apply the migration
./bin/migrate up
```

## API Endpoints

### Traditional Login

**POST** `/auth`

```json
{
  "username": "user@example.com",  // Can be username or email
  "password": "password123"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "ok",
  "data": {
    "access_token": "jwt_token_here",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "johndoe",
      "first_name": "John",
      "last_name": "Doe",
      "avatar_url": "https://...",
      "provider": "local",
      "display_name": "John Doe"
    }
  }
}
```

### User Registration

**POST** `/auth/register`

```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}
```

### Social Login URLs

**GET** `/auth/{provider}`

Get OAuth authorization URL for social login:
- `/auth/google`
- `/auth/github`
- `/auth/facebook`

**Response:**
```json
{
  "code": 200,
  "msg": "ok",
  "data": {
    "auth_url": "https://accounts.google.com/oauth/authorize?..."
  }
}
```

### Social Login Callback

**GET** `/auth/{provider}/callback`

OAuth callback endpoint that handles the authorization code and creates/authenticates users.

## Configuration

### Environment Variables

Create a `.env` file with your OAuth provider credentials:

```env
BASE_URL=http://localhost:8000

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# GitHub OAuth
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret

# Facebook OAuth
FACEBOOK_CLIENT_ID=your-facebook-client-id
FACEBOOK_CLIENT_SECRET=your-facebook-client-secret
```

### OAuth Provider Setup

#### Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URI: `http://localhost:8000/auth/google/callback`

#### GitHub OAuth Setup

1. Go to GitHub Settings > Developer settings > OAuth Apps
2. Create a new OAuth App
3. Set Authorization callback URL: `http://localhost:8000/auth/github/callback`

#### Facebook OAuth Setup

1. Go to [Facebook Developers](https://developers.facebook.com/)
2. Create a new app
3. Add Facebook Login product
4. Set Valid OAuth Redirect URIs: `http://localhost:8000/auth/facebook/callback`

## Code Structure

### Models

- `models/auth.go` - Enhanced Auth model with social login support
- New functions:
  - `CheckAuthByEmail()` - Email-based authentication
  - `GetAuthByProvider()` - Find user by OAuth provider
  - `CreateSocialAuth()` - Create social login user
  - `CreateLocalAuth()` - Create local user

### Services

- `service/auth_service/auth.go` - Enhanced authentication service
- `service/oauth_service/oauth.go` - OAuth provider management
- New functions:
  - `AuthenticateOrCreateSocialUser()` - Handle social login flow
  - `RegisterLocalUser()` - Local user registration

### API Routes

- `routers/api/auth.go` - Enhanced authentication endpoints
- New endpoints:
  - `Register()` - User registration
  - `GetOAuthURL()` - Get OAuth authorization URL
  - `OAuthCallback()` - Handle OAuth callback

### Configuration

- `config/oauth.go` - OAuth provider configuration management

## Usage Examples

### Frontend Integration

```javascript
// Get Google OAuth URL
const response = await fetch('/auth/google');
const data = await response.json();
window.location.href = data.data.auth_url;

// Traditional login
const loginResponse = await fetch('/auth', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'password123'
  })
});
```

### User Registration

```javascript
const registerResponse = await fetch('/auth/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'johndoe',
    email: 'john@example.com',
    password: 'password123'
  })
});
```

## Security Features

- CSRF protection with state parameter in OAuth flow
- Secure cookie storage for OAuth state
- Password hashing with bcrypt
- JWT token-based authentication
- Email verification status tracking
- Provider-specific user ID storage

## Migration Notes

- Existing users will continue to work with username/password authentication
- The `username` and `password` fields are now nullable to support social logins
- Social login users don't have passwords and are identified by provider + provider_id
- Email field is unique and required for all new users

## Testing

The system supports multiple authentication methods:
1. Username + password (existing)
2. Email + password (new)
3. Social OAuth providers (new)

Test with various scenarios:
- Local user registration and login
- Social login with Google/GitHub/Facebook
- Email verification workflows
- User profile updates