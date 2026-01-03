# Quick Start

This guide walks you through deploying your first service with Podoru.

## Prerequisites

- Podoru installed and running
- Access to the API endpoint
- A domain configured (optional, but recommended)

## Step 1: Create an Account

Register your admin account (first user becomes superadmin):

```bash
curl -X POST https://your-domain.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "your-secure-password",
    "name": "Admin"
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid-here",
      "email": "admin@example.com",
      "name": "Admin",
      "role": "superadmin"
    },
    "tokens": {
      "access_token": "eyJhbG...",
      "refresh_token": "abc123...",
      "expires_in": 900
    }
  }
}
```

## Step 2: Login

Get your access token:

```bash
curl -X POST https://your-domain.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "your-secure-password"
  }'
```

Save the `access_token` for subsequent requests:

```bash
export TOKEN="eyJhbG..."
```

## Step 3: Create a Team

Teams organize your projects and members:

```bash
curl -X POST https://your-domain.com/api/v1/teams \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Team",
    "slug": "my-team"
  }'
```

Save the team ID:
```bash
export TEAM_ID="team-uuid-here"
```

## Step 4: Create a Project

Projects group related services:

```bash
curl -X POST https://your-domain.com/api/v1/teams/$TEAM_ID/projects \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My App",
    "slug": "my-app"
  }'
```

Save the project ID:
```bash
export PROJECT_ID="project-uuid-here"
```

## Step 5: Create a Service

Create a service with a Docker image:

```bash
curl -X POST https://your-domain.com/api/v1/projects/$PROJECT_ID/services \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Web Server",
    "slug": "web",
    "deploy_type": "image",
    "image": "nginx:alpine"
  }'
```

Save the service ID:
```bash
export SERVICE_ID="service-uuid-here"
```

## Step 6: Add a Domain

Bind a domain to your service:

```bash
curl -X POST https://your-domain.com/api/v1/services/$SERVICE_ID/domains \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "web.example.com",
    "ssl_enabled": true,
    "ssl_auto": true
  }'
```

## Step 7: Deploy

Deploy your service:

```bash
curl -X POST https://your-domain.com/api/v1/services/$SERVICE_ID/deploy \
  -H "Authorization: Bearer $TOKEN"
```

Response:
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

## Step 8: Verify

Check the service status:

```bash
curl https://your-domain.com/api/v1/services/$SERVICE_ID \
  -H "Authorization: Bearer $TOKEN"
```

Access your deployed service:

```bash
curl https://web.example.com
```

## Complete Example Script

Here's a complete script to deploy nginx:

```bash
#!/bin/bash
API_URL="https://your-domain.com/api/v1"

# Login
TOKEN=$(curl -s -X POST $API_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}' \
  | jq -r '.data.tokens.access_token')

# Create team
TEAM_ID=$(curl -s -X POST $API_URL/teams \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Demo","slug":"demo"}' \
  | jq -r '.data.id')

# Create project
PROJECT_ID=$(curl -s -X POST $API_URL/teams/$TEAM_ID/projects \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Demo App","slug":"demo-app"}' \
  | jq -r '.data.id')

# Create service
SERVICE_ID=$(curl -s -X POST $API_URL/projects/$PROJECT_ID/services \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Nginx","slug":"nginx","deploy_type":"image","image":"nginx:alpine"}' \
  | jq -r '.data.id')

# Add domain
curl -s -X POST $API_URL/services/$SERVICE_ID/domains \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"domain":"demo.example.com","ssl_enabled":true}'

# Deploy
curl -s -X POST $API_URL/services/$SERVICE_ID/deploy \
  -H "Authorization: Bearer $TOKEN"

echo "Deployed! Access at https://demo.example.com"
```

## Next Steps

- [Configure environment variables](configuration.md)
- [Learn about domain management](../guides/domains.md)
- [Explore the full API](../api/README.md)
