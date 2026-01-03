# Projects API

Projects group related services within a team.

## List Projects

Get all projects in a team.

```http
GET /api/v1/teams/:teamId/projects
Authorization: Bearer {access_token}
```

### Response

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "team_id": "team-uuid",
      "name": "My App",
      "slug": "my-app",
      "github_repo": null,
      "github_branch": "main",
      "auto_deploy": false,
      "created_at": "2026-01-03T10:00:00Z",
      "updated_at": "2026-01-03T10:00:00Z"
    }
  ]
}
```

## Create Project

```http
POST /api/v1/teams/:teamId/projects
Authorization: Bearer {access_token}
```

### Request

```json
{
  "name": "My App",
  "slug": "my-app",
  "github_repo": "username/repo",
  "github_branch": "main",
  "auto_deploy": true
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Display name |
| `slug` | string | Yes | URL-safe identifier |
| `github_repo` | string | No | GitHub repository (owner/repo) |
| `github_branch` | string | No | Branch for auto-deploy |
| `auto_deploy` | bool | No | Enable auto-deploy on push |

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "team_id": "team-uuid",
    "name": "My App",
    "slug": "my-app",
    "github_branch": "main",
    "auto_deploy": false,
    "created_at": "2026-01-03T10:00:00Z"
  }
}
```

## Get Project

```http
GET /api/v1/projects/:projectId
Authorization: Bearer {access_token}
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "team_id": "team-uuid",
    "name": "My App",
    "slug": "my-app",
    "github_repo": null,
    "github_branch": "main",
    "auto_deploy": false,
    "created_at": "2026-01-03T10:00:00Z",
    "updated_at": "2026-01-03T10:00:00Z"
  }
}
```

## Update Project

```http
PUT /api/v1/projects/:projectId
Authorization: Bearer {access_token}
```

### Request

```json
{
  "name": "Updated App Name",
  "github_repo": "username/new-repo",
  "auto_deploy": true
}
```

## Delete Project

```http
DELETE /api/v1/projects/:projectId
Authorization: Bearer {access_token}
```

Deletes the project and all associated services.

### Response

```json
{
  "success": true,
  "data": {
    "message": "Project deleted successfully"
  }
}
```

## Errors

| Code | Description |
|------|-------------|
| `NOT_FOUND` | Project not found |
| `FORBIDDEN` | Not a team member |
| `CONFLICT` | Slug already exists in team |
