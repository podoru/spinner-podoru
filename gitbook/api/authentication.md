# Authentication API

Podoru uses JWT tokens for authentication with access and refresh token flow.

## Register

Create a new user account. The first registered user becomes superadmin.

```http
POST /api/v1/auth/register
```

### Request

```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe"
}
```

### Response

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "user",
      "is_active": true,
      "created_at": "2026-01-03T10:00:00Z"
    },
    "tokens": {
      "access_token": "eyJhbG...",
      "refresh_token": "abc123...",
      "expires_in": 900,
      "token_type": "Bearer"
    }
  }
}
```

### Errors

| Code | Description |
|------|-------------|
| `FORBIDDEN` | Registration is disabled |
| `CONFLICT` | Email already registered |
| `VALIDATION_ERROR` | Invalid email or password |

## Login

Authenticate and receive tokens.

```http
POST /api/v1/auth/login
```

### Request

```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

### Response

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "user"
    },
    "tokens": {
      "access_token": "eyJhbG...",
      "refresh_token": "abc123...",
      "expires_in": 900,
      "token_type": "Bearer"
    }
  }
}
```

### Errors

| Code | Description |
|------|-------------|
| `UNAUTHORIZED` | Invalid email or password |

## Refresh Token

Get a new access token using a refresh token.

```http
POST /api/v1/auth/refresh
```

### Request

```json
{
  "refresh_token": "abc123..."
}
```

### Response

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbG...",
    "refresh_token": "xyz789...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

### Errors

| Code | Description |
|------|-------------|
| `UNAUTHORIZED` | Invalid or expired refresh token |

## Logout

Invalidate the current refresh token.

```http
POST /api/v1/auth/logout
Authorization: Bearer {access_token}
```

### Request

```json
{
  "refresh_token": "abc123..."
}
```

### Response

```json
{
  "success": true,
  "data": {
    "message": "Logged out successfully"
  }
}
```

## Token Usage

Include the access token in all authenticated requests:

```bash
curl -H "Authorization: Bearer eyJhbG..." https://api.example.com/api/v1/teams
```

## Token Expiry

| Token | Default Expiry | Configurable |
|-------|---------------|--------------|
| Access Token | 15 minutes | `JWT_ACCESS_EXPIRY` |
| Refresh Token | 7 days | `JWT_REFRESH_EXPIRY` |

## Get Current User

```http
GET /api/v1/users/me
Authorization: Bearer {access_token}
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user",
    "is_active": true,
    "created_at": "2026-01-03T10:00:00Z",
    "updated_at": "2026-01-03T10:00:00Z"
  }
}
```

## Update Current User

```http
PUT /api/v1/users/me
Authorization: Bearer {access_token}
```

### Request

```json
{
  "name": "Jane Doe"
}
```

## Change Password

```http
PUT /api/v1/users/me/password
Authorization: Bearer {access_token}
```

### Request

```json
{
  "current_password": "oldpassword",
  "new_password": "newpassword"
}
```
