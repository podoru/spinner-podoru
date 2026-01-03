# Services API

Services represent deployable containers within a project.

## List Services

```http
GET /api/v1/projects/:projectId/services
Authorization: Bearer {access_token}
```

### Response

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "project_id": "project-uuid",
      "name": "Web Server",
      "slug": "web",
      "deploy_type": "image",
      "image": "nginx:alpine",
      "status": "running",
      "container_id": "abc123...",
      "created_at": "2026-01-03T10:00:00Z"
    }
  ]
}
```

## Create Service

```http
POST /api/v1/projects/:projectId/services
Authorization: Bearer {access_token}
```

### Request

```json
{
  "name": "Web Server",
  "slug": "web",
  "deploy_type": "image",
  "image": "nginx:alpine",
  "replicas": 1,
  "cpu_limit": 0.5,
  "memory_limit": 256,
  "restart_policy": "unless-stopped",
  "health_check_path": "/health",
  "health_check_interval": 30
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Display name |
| `slug` | string | Yes | URL-safe identifier |
| `deploy_type` | string | Yes | `image`, `dockerfile`, `compose` |
| `image` | string | For image | Docker image name |
| `dockerfile_path` | string | For dockerfile | Path to Dockerfile |
| `build_context` | string | For dockerfile | Build context path |
| `replicas` | int | No | Number of instances (default: 1) |
| `cpu_limit` | float | No | CPU limit (0.5 = 50%) |
| `memory_limit` | int | No | Memory limit in MB |
| `restart_policy` | string | No | Restart policy |
| `health_check_path` | string | No | HTTP health check path |
| `health_check_interval` | int | No | Interval in seconds |

## Get Service

```http
GET /api/v1/services/:serviceId
Authorization: Bearer {access_token}
```

## Update Service

```http
PUT /api/v1/services/:serviceId
Authorization: Bearer {access_token}
```

### Request

```json
{
  "image": "nginx:1.25",
  "memory_limit": 512
}
```

## Delete Service

```http
DELETE /api/v1/services/:serviceId
Authorization: Bearer {access_token}
```

Stops and removes the container before deleting.

## Deploy Service

Trigger a deployment.

```http
POST /api/v1/services/:serviceId/deploy
Authorization: Bearer {access_token}
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "deployment-uuid",
    "service_id": "service-uuid",
    "status": "pending",
    "started_at": "2026-01-03T10:00:00Z"
  }
}
```

## Start Service

Start a stopped container.

```http
POST /api/v1/services/:serviceId/start
Authorization: Bearer {access_token}
```

## Stop Service

Stop a running container.

```http
POST /api/v1/services/:serviceId/stop
Authorization: Bearer {access_token}
```

## Restart Service

Restart the container.

```http
POST /api/v1/services/:serviceId/restart
Authorization: Bearer {access_token}
```

## Get Logs

Retrieve container logs.

```http
GET /api/v1/services/:serviceId/logs
Authorization: Bearer {access_token}
```

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `tail` | int | 100 | Number of lines |
| `since` | string | - | Timestamp filter |

### Response

```json
{
  "success": true,
  "data": {
    "logs": "2026-01-03T10:00:00Z nginx started...\n..."
  }
}
```

## Service Status

| Status | Description |
|--------|-------------|
| `stopped` | Not deployed or stopped |
| `deploying` | Deployment in progress |
| `running` | Container is running |
| `failed` | Deployment or container failed |

## Restart Policies

| Policy | Description |
|--------|-------------|
| `no` | Never restart |
| `always` | Always restart |
| `on-failure` | Restart on non-zero exit |
| `unless-stopped` | Restart unless manually stopped |
