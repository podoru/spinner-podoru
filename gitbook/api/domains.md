# Domains API

Domains bind custom hostnames to services with automatic Traefik routing.

## List Domains

Get all domains for a service.

```http
GET /api/v1/services/:serviceId/domains
Authorization: Bearer {access_token}
```

### Response

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "service_id": "service-uuid",
      "domain": "myapp.example.com",
      "ssl_enabled": true,
      "ssl_auto": true,
      "created_at": "2026-01-03T10:00:00Z"
    }
  ]
}
```

## Add Domain

Bind a domain to a service.

```http
POST /api/v1/services/:serviceId/domains
Authorization: Bearer {access_token}
```

### Request

```json
{
  "domain": "myapp.example.com",
  "ssl_enabled": true,
  "ssl_auto": true
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `domain` | string | Yes | Full domain name |
| `ssl_enabled` | bool | No | Enable HTTPS (default: false) |
| `ssl_auto` | bool | No | Auto-provision Let's Encrypt cert |

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "service_id": "service-uuid",
    "domain": "myapp.example.com",
    "ssl_enabled": true,
    "ssl_auto": true,
    "created_at": "2026-01-03T10:00:00Z"
  }
}
```

## Delete Domain

Remove a domain binding.

```http
DELETE /api/v1/services/:serviceId/domains/:domainId
Authorization: Bearer {access_token}
```

### Response

```json
{
  "success": true,
  "data": {
    "message": "Domain deleted successfully"
  }
}
```

## Domain Validation

Domains must:

- Be a valid hostname format
- Not already be assigned to another service
- Have DNS configured before SSL can be provisioned

## SSL Modes

### No SSL

```json
{
  "domain": "internal.local",
  "ssl_enabled": false
}
```

### Let's Encrypt Auto SSL

```json
{
  "domain": "myapp.example.com",
  "ssl_enabled": true,
  "ssl_auto": true
}
```

Requirements:
- Valid domain with DNS pointing to server
- Port 80 accessible for ACME challenge
- `TRAEFIK_ACME_EMAIL` configured

### Custom Certificate (Coming Soon)

```json
{
  "domain": "myapp.example.com",
  "ssl_enabled": true,
  "ssl_auto": false,
  "ssl_cert": "-----BEGIN CERTIFICATE-----...",
  "ssl_key": "-----BEGIN PRIVATE KEY-----..."
}
```

## Traefik Integration

When a domain is added and the service is deployed, Traefik labels are automatically configured:

```yaml
labels:
  traefik.http.routers.service.rule: "Host(`myapp.example.com`)"
  traefik.http.routers.service-secure.tls: "true"
  traefik.http.routers.service-secure.tls.certresolver: "letsencrypt"
```

## Multiple Domains

A service can have multiple domains:

```bash
# Add primary domain
POST /services/:id/domains
{"domain": "myapp.com", "ssl_enabled": true}

# Add www alias
POST /services/:id/domains
{"domain": "www.myapp.com", "ssl_enabled": true}
```

Both domains route to the same container.

## Errors

| Code | Description |
|------|-------------|
| `CONFLICT` | Domain already exists |
| `VALIDATION_ERROR` | Invalid domain format |
| `NOT_FOUND` | Service or domain not found |
