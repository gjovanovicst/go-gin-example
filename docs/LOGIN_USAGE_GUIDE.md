# Login Usage Guide

This guide shows you how to use the modernized authentication system with username/email and password login.

## Traditional Login (Enhanced)

The `/auth` endpoint now supports both username and email login.

### Endpoint
**POST** `/auth`

### Request Format
Send a JSON request with either username or email + password:

#### Option 1: Login with Email
```json
{
  "email": "user@example.com",
  "password": "your_password"
}
```

#### Option 2: Login with Username
```json
{
  "username": "johndoe",
  "password": "your_password"
}
```

### Example cURL Commands

#### Login with Email:
```bash
curl -X POST http://localhost:8000/auth \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "your_password"
  }'
```

#### Login with Username:
```bash
curl -X POST http://localhost:8000/auth \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "your_password"
  }'
```

### Response Format
```json
{
  "code": 200,
  "msg": "ok",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "johndoe",
      "first_name": "John",
      "last_name": "Doe",
      "avatar_url": "https://example.com/avatar.jpg",
      "provider": "local",
      "display_name": "John Doe"
    }
  }
}
```

## User Registration

Before you can login, you need to register a user with the new registration endpoint.

### Endpoint
**POST** `/auth/register`

### Request Format
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}
```

### Example cURL Command:
```bash
curl -X POST http://localhost:8000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

## JavaScript/Frontend Examples

### Register a New User
```javascript
async function registerUser(username, email, password) {
  const response = await fetch('/auth/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      username: username,
      email: email,
      password: password
    })
  });
  
  const data = await response.json();
  
  if (data.code === 200) {
    // Registration successful
    localStorage.setItem('token', data.data.access_token);
    return data.data.user;
  } else {
    throw new Error('Registration failed');
  }
}
```

### Login with Email
```javascript
async function loginWithEmail(email, password) {
  const response = await fetch('/auth', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email: email,
      password: password
    })
  });
  
  const data = await response.json();
  
  if (data.code === 200) {
    // Login successful
    localStorage.setItem('token', data.data.access_token);
    return data.data.user;
  } else {
    throw new Error('Login failed');
  }
}
```

### Login with Username
```javascript
async function loginWithUsername(username, password) {
  const response = await fetch('/auth', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      username: username,
      password: password
    })
  });
  
  const data = await response.json();
  
  if (data.code === 200) {
    // Login successful
    localStorage.setItem('token', data.data.access_token);
    return data.data.user;
  } else {
    throw new Error('Login failed');
  }
}
```

### Using the Token
After login, use the token for authenticated requests:

```javascript
async function makeAuthenticatedRequest(url) {
  const token = localStorage.getItem('token');
  
  const response = await fetch(url, {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  
  return await response.json();
}
```

## Testing the System

### 1. Start the Application
```bash
./go-gin-example.exe
```

### 2. Register a Test User
```bash
curl -X POST http://localhost:8000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Login with Email
```bash
curl -X POST http://localhost:8000/auth \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 4. Login with Username
```bash
curl -X POST http://localhost:8000/auth \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

## Key Changes from Original System

1. **JSON Instead of Form Data**: The login now accepts JSON instead of form data
2. **Email Support**: You can login with email address instead of just username
3. **Enhanced Response**: The response now includes detailed user information
4. **Registration Endpoint**: New endpoint to register users
5. **Backward Compatibility**: Existing username/password login still works

## Error Responses

### Invalid Credentials
```json
{
  "code": 401,
  "msg": "unauthorized",
  "data": null
}
```

### Missing Parameters
```json
{
  "code": 400,
  "msg": "invalid params",
  "data": null
}
```

### User Already Exists (Registration)
```json
{
  "code": 409,
  "msg": "tag already exists",
  "data": {
    "error": "User with this email already exists"
  }
}
```

This enhanced authentication system provides a modern, flexible login experience while maintaining backward compatibility with existing systems.