# Teams API

Teams are the top-level organizational unit for grouping projects and members.

## List Teams

Get all teams the current user is a member of.

```http
GET /api/v1/teams
Authorization: Bearer {access_token}
```

### Response

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "My Team",
      "slug": "my-team",
      "owner_id": "user-uuid",
      "created_at": "2026-01-03T10:00:00Z",
      "updated_at": "2026-01-03T10:00:00Z",
      "role": "owner"
    }
  ]
}
```

## Create Team

```http
POST /api/v1/teams
Authorization: Bearer {access_token}
```

### Request

```json
{
  "name": "My Team",
  "slug": "my-team"
}
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "My Team",
    "slug": "my-team",
    "owner_id": "user-uuid",
    "created_at": "2026-01-03T10:00:00Z"
  }
}
```

## Get Team

```http
GET /api/v1/teams/:teamId
Authorization: Bearer {access_token}
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "My Team",
    "slug": "my-team",
    "owner_id": "user-uuid",
    "created_at": "2026-01-03T10:00:00Z",
    "updated_at": "2026-01-03T10:00:00Z"
  }
}
```

## Update Team

```http
PUT /api/v1/teams/:teamId
Authorization: Bearer {access_token}
```

### Request

```json
{
  "name": "Updated Team Name"
}
```

## Delete Team

```http
DELETE /api/v1/teams/:teamId
Authorization: Bearer {access_token}
```

Requires owner role. Deletes all projects and services.

## Team Members

### List Members

```http
GET /api/v1/teams/:teamId/members
Authorization: Bearer {access_token}
```

### Response

```json
{
  "success": true,
  "data": [
    {
      "user_id": "uuid",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "owner",
      "joined_at": "2026-01-03T10:00:00Z"
    }
  ]
}
```

### Add Member

```http
POST /api/v1/teams/:teamId/members
Authorization: Bearer {access_token}
```

#### Request

```json
{
  "email": "newmember@example.com",
  "role": "member"
}
```

### Update Member Role

```http
PUT /api/v1/teams/:teamId/members/:userId
Authorization: Bearer {access_token}
```

#### Request

```json
{
  "role": "admin"
}
```

### Remove Member

```http
DELETE /api/v1/teams/:teamId/members/:userId
Authorization: Bearer {access_token}
```

## Roles

| Role | Permissions |
|------|-------------|
| `owner` | Full access, delete team, transfer ownership |
| `admin` | Manage projects, services, members |
| `member` | View and deploy services |
