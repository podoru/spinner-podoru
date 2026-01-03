# Getting Started

This section covers everything you need to get Podoru up and running.

## Contents

1. [Installation](installation.md) - Install Podoru on your server
2. [Quick Start](quick-start.md) - Deploy your first service
3. [Configuration](configuration.md) - Configure environment and settings

## Prerequisites

Before installing Podoru, ensure your server has:

- **Operating System**: Linux (Ubuntu 20.04+, Debian 11+, or similar)
- **Docker**: Version 20.10.0 or higher
- **Docker Compose**: Version 2.0.0 or higher
- **Memory**: Minimum 2GB RAM
- **Ports**: 80, 443, and 8080 available

## Installation Methods

### Production (Recommended)

Use the automated install script for a complete production setup:

```bash
curl -fsSL https://raw.githubusercontent.com/podoru/spinner-podoru/master/scripts/install.sh -o install.sh
chmod +x install.sh
./install.sh install
```

### Development

For local development with hot reload:

```bash
git clone git@github.com:podoru/spinner-podoru.git
cd spinner-podoru
./scripts/dev-setup.sh setup
make dev
```

## Next Steps

After installation:

1. Access the API at your configured domain
2. Register your admin account (first user becomes superadmin)
3. Create a team and project
4. Deploy your first service

Continue to [Quick Start](quick-start.md) for a step-by-step guide.
