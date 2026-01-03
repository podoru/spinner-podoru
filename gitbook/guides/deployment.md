# Deploying Services

This guide covers deploying and managing services in Podoru.

## Service Types

Podoru supports deploying services from:

| Type | Description |
|------|-------------|
| `image` | Deploy from a Docker Hub or registry image |
| `dockerfile` | Build from a Dockerfile (coming soon) |
| `compose` | Deploy from docker-compose.yml (coming soon) |

## Creating a Service

### From Docker Image

```bash
curl -X POST https://api.example.com/api/v1/projects/$PROJECT_ID/services \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Web App",
    "slug": "my-web-app",
    "deploy_type": "image",
    "image": "nginx:alpine",
    "replicas": 1,
    "restart_policy": "unless-stopped"
  }'
```

### Service Options

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Display name |
| `slug` | string | URL-safe identifier |
| `deploy_type` | string | `image`, `dockerfile`, `compose` |
| `image` | string | Docker image (for `image` type) |
| `replicas` | int | Number of instances (default: 1) |
| `restart_policy` | string | `no`, `always`, `on-failure`, `unless-stopped` |
| `cpu_limit` | float | CPU limit (e.g., 0.5 = 50% of one core) |
| `memory_limit` | int | Memory limit in MB |
| `health_check_path` | string | HTTP path for health checks |
| `health_check_interval` | int | Health check interval in seconds |

## Deploying

Trigger a deployment:

```bash
curl -X POST https://api.example.com/api/v1/services/$SERVICE_ID/deploy \
  -H "Authorization: Bearer $TOKEN"
```

The deployment runs asynchronously. Check status:

```bash
curl https://api.example.com/api/v1/services/$SERVICE_ID \
  -H "Authorization: Bearer $TOKEN"
```

## Deployment Lifecycle

1. **pending** - Deployment created
2. **deploying** - Pulling image, creating container
3. **success** - Container running
4. **failed** - Deployment failed (check logs)

## Service Operations

### Start

```bash
curl -X POST https://api.example.com/api/v1/services/$SERVICE_ID/start \
  -H "Authorization: Bearer $TOKEN"
```

### Stop

```bash
curl -X POST https://api.example.com/api/v1/services/$SERVICE_ID/stop \
  -H "Authorization: Bearer $TOKEN"
```

### Restart

```bash
curl -X POST https://api.example.com/api/v1/services/$SERVICE_ID/restart \
  -H "Authorization: Bearer $TOKEN"
```

## Viewing Logs

Get container logs:

```bash
curl "https://api.example.com/api/v1/services/$SERVICE_ID/logs?tail=100" \
  -H "Authorization: Bearer $TOKEN"
```

Query parameters:
- `tail` - Number of lines (default: 100)
- `since` - Timestamp (e.g., `2024-01-01T00:00:00Z`)

## Updating a Service

Update service configuration:

```bash
curl -X PUT https://api.example.com/api/v1/services/$SERVICE_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "image": "nginx:1.25",
    "memory_limit": 256
  }'
```

After updating, redeploy to apply changes:

```bash
curl -X POST https://api.example.com/api/v1/services/$SERVICE_ID/deploy \
  -H "Authorization: Bearer $TOKEN"
```

## Deleting a Service

Delete a service (stops and removes container):

```bash
curl -X DELETE https://api.example.com/api/v1/services/$SERVICE_ID \
  -H "Authorization: Bearer $TOKEN"
```

## Resource Limits

Set CPU and memory limits:

```json
{
  "cpu_limit": 0.5,
  "memory_limit": 512
}
```

- `cpu_limit`: Fraction of CPU (0.5 = 50% of one core)
- `memory_limit`: Memory in MB

## Health Checks

Configure health checks:

```json
{
  "health_check_path": "/health",
  "health_check_interval": 30
}
```

The container is marked unhealthy if the health check fails.

## Private Registries

For private Docker registries, include credentials in the image URL:

```json
{
  "image": "registry.example.com/myapp:latest"
}
```

Registry authentication coming in a future release.
