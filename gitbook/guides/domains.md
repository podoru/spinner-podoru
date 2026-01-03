# Domain Management

This guide covers binding custom domains to your services with automatic SSL.

## Overview

Podoru integrates with Traefik to provide:

- Automatic domain routing to containers
- SSL certificates via Let's Encrypt
- HTTP to HTTPS redirection
- Multiple domains per service

## Prerequisites

- Traefik enabled (`TRAEFIK_ENABLED=true`)
- DNS configured to point to your server
- Ports 80 and 443 accessible

## Adding a Domain

Bind a domain to a service:

```bash
curl -X POST https://api.example.com/api/v1/services/$SERVICE_ID/domains \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "myapp.example.com",
    "ssl_enabled": true,
    "ssl_auto": true
  }'
```

### Domain Options

| Field | Type | Description |
|-------|------|-------------|
| `domain` | string | Full domain name |
| `ssl_enabled` | bool | Enable HTTPS |
| `ssl_auto` | bool | Auto-provision Let's Encrypt certificate |

## SSL Configuration

### Automatic SSL (Let's Encrypt)

```json
{
  "domain": "myapp.example.com",
  "ssl_enabled": true,
  "ssl_auto": true
}
```

Requirements:
- Domain DNS points to your server
- Port 80 accessible for HTTP challenge
- Valid email in `TRAEFIK_ACME_EMAIL`

### Without SSL

```json
{
  "domain": "internal.local",
  "ssl_enabled": false
}
```

## Listing Domains

Get all domains for a service:

```bash
curl https://api.example.com/api/v1/services/$SERVICE_ID/domains \
  -H "Authorization: Bearer $TOKEN"
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "id": "domain-uuid",
      "service_id": "service-uuid",
      "domain": "myapp.example.com",
      "ssl_enabled": true,
      "ssl_auto": true,
      "created_at": "2026-01-03T10:00:00Z"
    }
  ]
}
```

## Removing a Domain

Delete a domain binding:

```bash
curl -X DELETE https://api.example.com/api/v1/services/$SERVICE_ID/domains/$DOMAIN_ID \
  -H "Authorization: Bearer $TOKEN"
```

## Multiple Domains

A service can have multiple domains:

```bash
# Primary domain
curl -X POST .../domains -d '{"domain": "myapp.com", "ssl_enabled": true}'

# Alias domain
curl -X POST .../domains -d '{"domain": "www.myapp.com", "ssl_enabled": true}'
```

All domains route to the same container.

## DNS Configuration

### A Record

```
Type: A
Name: myapp
Value: YOUR_SERVER_IP
TTL: 300
```

### Wildcard (for subdomains)

```
Type: A
Name: *.myapp
Value: YOUR_SERVER_IP
TTL: 300
```

## Traefik Labels

When a service with domains is deployed, Podoru automatically adds Traefik labels:

```yaml
labels:
  traefik.enable: "true"
  traefik.http.routers.myapp.rule: "Host(`myapp.example.com`)"
  traefik.http.routers.myapp.entrypoints: "web"
  traefik.http.routers.myapp-secure.rule: "Host(`myapp.example.com`)"
  traefik.http.routers.myapp-secure.entrypoints: "websecure"
  traefik.http.routers.myapp-secure.tls: "true"
  traefik.http.routers.myapp-secure.tls.certresolver: "letsencrypt"
```

## Troubleshooting

### Domain Not Accessible

1. Verify DNS is pointing to your server:
   ```bash
   dig +short myapp.example.com
   ```

2. Check Traefik is running:
   ```bash
   docker ps | grep traefik
   ```

3. Verify container is on Traefik network:
   ```bash
   docker network inspect podoru_traefik
   ```

### SSL Certificate Not Issued

1. Check Traefik logs:
   ```bash
   docker logs podoru_traefik
   ```

2. Verify ACME email is set:
   ```bash
   echo $TRAEFIK_ACME_EMAIL
   ```

3. Ensure port 80 is accessible for HTTP challenge

### HTTP Not Redirecting to HTTPS

Redeploy the service to refresh Traefik labels:

```bash
curl -X POST .../services/$SERVICE_ID/deploy
```
