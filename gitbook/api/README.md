# API Reference

Podoru provides a RESTful API for managing containers, services, and infrastructure.

## Base URL

```
https://your-domain.com/api/v1
```

## Authentication

All API requests (except registration and login) require a Bearer token:

```bash
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" https://api.example.com/api/v1/...
```

## Response Format

All responses follow this format:

### Success

```json
{
  "success": true,
  "data": { ... }
}
```

### Error

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Missing or invalid token |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `CONFLICT` | 409 | Resource already exists |
| `VALIDATION_ERROR` | 400 | Invalid request data |
| `INTERNAL_ERROR` | 500 | Server error |

## Pagination

List endpoints support pagination:

```bash
GET /api/v1/teams?page=1&limit=20
```

Response includes pagination metadata:

```json
{
  "success": true,
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

## Rate Limiting

API requests are rate limited:

- **Authenticated**: 1000 requests/hour
- **Unauthenticated**: 100 requests/hour

Rate limit headers:

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1704067200
```

## API Sections

- [Authentication](authentication.md) - Login, register, tokens
- [Teams](teams.md) - Team management
- [Projects](projects.md) - Project management
- [Services](services.md) - Service deployment
- [Domains](domains.md) - Domain management

## Interactive Documentation

Swagger/OpenAPI documentation is available at:

```
https://your-domain.com/api/v1/docs
```

OpenAPI JSON spec:

```
https://your-domain.com/api/v1/docs/openapi.json
```

## Quick Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register new user |
| POST | `/auth/login` | Login |
| POST | `/auth/refresh` | Refresh token |
| GET | `/users/me` | Get current user |
| GET | `/teams` | List teams |
| POST | `/teams` | Create team |
| GET | `/teams/:id/projects` | List projects |
| POST | `/teams/:id/projects` | Create project |
| GET | `/projects/:id/services` | List services |
| POST | `/projects/:id/services` | Create service |
| POST | `/services/:id/deploy` | Deploy service |
| GET | `/services/:id/domains` | List domains |
| POST | `/services/:id/domains` | Add domain |
