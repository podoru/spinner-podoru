# Podoru Documentation

Welcome to Podoru, a self-hosted Docker container management platform with automatic SSL and domain routing.

## What is Podoru?

Podoru is a lightweight alternative to platforms like EasyPanel, Coolify, or CapRover. It provides a REST API to manage Docker containers with:

- **Automatic Traefik Integration**: Deploy containers and automatically configure reverse proxy routing
- **Domain Management**: Bind custom domains to your services with automatic SSL via Let's Encrypt
- **Multi-tenant Architecture**: Organize resources into teams and projects
- **Simple API**: RESTful interface for all operations

## Quick Links

- [Installation Guide](getting-started/installation.md) - Get Podoru running in production
- [Quick Start](getting-started/quick-start.md) - Deploy your first container
- [Configuration](getting-started/configuration.md) - Environment variables and settings
- [API Reference](api/README.md) - Complete API documentation

## Core Concepts

### Teams
Teams are the top-level organizational unit. Each team can have multiple members with different roles (owner, admin, member).

### Projects
Projects belong to teams and group related services together. A project might represent an application with multiple microservices.

### Services
Services are the actual containers you deploy. Each service can have:
- Docker image configuration
- Environment variables
- Domain bindings
- Resource limits (CPU, memory)
- Health checks

### Domains
Domains are bound to services and automatically configured in Traefik for routing. SSL certificates are automatically provisioned via Let's Encrypt.

## Requirements

- Linux server (Ubuntu 20.04+ recommended)
- Docker Engine 20.10+
- Docker Compose 2.0+
- 2GB RAM minimum
- Ports 80, 443, 8080 available

## Getting Help

- [GitHub Issues](https://github.com/podoru/spinner-podoru/issues) - Bug reports and feature requests
- [API Documentation](api/README.md) - Endpoint reference
