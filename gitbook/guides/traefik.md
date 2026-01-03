# Traefik Integration

Podoru uses Traefik as a reverse proxy for automatic routing and SSL.

## Overview

Traefik provides:

- **Dynamic routing**: Automatic container discovery
- **SSL termination**: Let's Encrypt integration
- **Load balancing**: Distribute traffic across replicas
- **Dashboard**: Monitor routes and services

## Architecture

```
Internet → Traefik (80/443) → Docker Network → Containers
```

All deployed containers join the `podoru_traefik` network and receive Traefik labels for routing.

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TRAEFIK_ENABLED` | Enable Traefik integration | `true` |
| `TRAEFIK_NETWORK` | Docker network name | `podoru_traefik` |
| `TRAEFIK_HTTP_PORT` | HTTP port | `80` |
| `TRAEFIK_HTTPS_PORT` | HTTPS port | `443` |
| `TRAEFIK_DASHBOARD_PORT` | Dashboard port | `8081` |
| `TRAEFIK_ACME_EMAIL` | Let's Encrypt email | (required) |

### Docker Compose (Production)

```yaml
traefik:
  image: traefik:v3.2
  command:
    - "--api.dashboard=true"
    - "--providers.docker=true"
    - "--providers.docker.exposedbydefault=false"
    - "--providers.docker.network=podoru_traefik"
    - "--entrypoints.web.address=:80"
    - "--entrypoints.websecure.address=:443"
    - "--certificatesresolvers.letsencrypt.acme.email=${ACME_EMAIL}"
    - "--certificatesresolvers.letsencrypt.acme.storage=/etc/traefik/acme.json"
    - "--certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web"
  ports:
    - "80:80"
    - "443:443"
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock:ro
    - traefik_data:/etc/traefik
  networks:
    - podoru_traefik
```

## How It Works

### 1. Container Deployment

When you deploy a service with domains, Podoru:

1. Creates the container with Traefik labels
2. Connects container to `podoru_traefik` network
3. Traefik discovers the container automatically

### 2. Label Configuration

Podoru generates these labels automatically:

```yaml
labels:
  # Enable Traefik for this container
  traefik.enable: "true"

  # HTTP Router
  traefik.http.routers.myapp.rule: "Host(`myapp.example.com`)"
  traefik.http.routers.myapp.entrypoints: "web"

  # HTTPS Router (with SSL)
  traefik.http.routers.myapp-secure.rule: "Host(`myapp.example.com`)"
  traefik.http.routers.myapp-secure.entrypoints: "websecure"
  traefik.http.routers.myapp-secure.tls: "true"
  traefik.http.routers.myapp-secure.tls.certresolver: "letsencrypt"

  # Service port
  traefik.http.services.myapp.loadbalancer.server.port: "80"
```

### 3. Network Validation

Before deployment, Podoru validates the Traefik network exists:

```go
if err := containerManager.ValidateNetwork(ctx, networkID); err != nil {
    return fmt.Errorf("traefik network not found: %s", networkID)
}
```

## Dashboard

Access the Traefik dashboard:

- **Development**: `http://localhost:8081`
- **Production**: `https://traefik.your-domain.com/dashboard/`

### Dashboard Authentication (Production)

Generate password hash:

```bash
htpasswd -nb admin your-password
```

Add to `.env.prod`:

```bash
TRAEFIK_DASHBOARD_AUTH=admin:$apr1$...
```

## SSL Certificates

### Let's Encrypt

Certificates are automatically provisioned when:

1. Domain has `ssl_enabled: true` and `ssl_auto: true`
2. DNS points to your server
3. Port 80 is accessible

### Certificate Storage

Certificates are stored in the `traefik_data` Docker volume at `/etc/traefik/acme.json`.

### Backup Certificates

```bash
docker cp podoru_traefik:/etc/traefik/acme.json ./acme-backup.json
```

## Multiple Domains

For services with multiple domains, Traefik creates a combined rule:

```yaml
traefik.http.routers.myapp.rule: "Host(`myapp.com`) || Host(`www.myapp.com`)"
```

## Troubleshooting

### Container Not Accessible via Domain

1. Check container is on the network:
   ```bash
   docker network inspect podoru_traefik
   ```

2. Verify Traefik sees the container:
   ```bash
   curl http://localhost:8081/api/http/routers
   ```

3. Check container labels:
   ```bash
   docker inspect container_name --format '{{json .Config.Labels}}'
   ```

### Traefik Not Starting

Check logs:
```bash
docker logs podoru_traefik
```

Common issues:
- Port 80/443 already in use
- Docker socket permission denied
- Invalid configuration

### SSL Certificate Not Working

1. Check ACME logs in Traefik output
2. Verify domain DNS resolution
3. Ensure HTTP challenge port (80) is accessible
4. Check certificate limits (Let's Encrypt rate limits)
