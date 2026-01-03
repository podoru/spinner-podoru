# GitHub Integration

This guide covers integrating GitHub for automatic deployments.

> **Note**: GitHub integration is planned for a future release. This document outlines the intended functionality.

## Overview

GitHub integration enables:

- Auto-deploy on push to specified branches
- Webhook notifications for deployment status
- Build from Dockerfile in repository

## Planned Features

### Repository Connection

Connect a GitHub repository to a project:

```bash
curl -X POST https://api.example.com/api/v1/projects/$PROJECT_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "github_repo": "username/repo",
    "github_branch": "main",
    "auto_deploy": true
  }'
```

### Webhook Setup

Podoru will provide a webhook URL:

```
https://your-domain.com/api/v1/webhooks/github
```

Configure in GitHub:
1. Repository Settings → Webhooks → Add webhook
2. Payload URL: webhook URL above
3. Content type: `application/json`
4. Events: Push events

### Auto-Deploy Flow

1. Push to configured branch
2. GitHub sends webhook to Podoru
3. Podoru pulls latest code
4. Builds image from Dockerfile
5. Deploys new container
6. Sends status back to GitHub

### Build Configuration

Configure build settings in service:

```json
{
  "deploy_type": "dockerfile",
  "dockerfile_path": "Dockerfile",
  "build_context": ".",
  "build_args": {
    "NODE_ENV": "production"
  }
}
```

## Current Workaround

Until GitHub integration is complete, you can:

1. Build and push images to a registry
2. Update the service image
3. Trigger deployment via API

Example CI/CD script:

```yaml
# .github/workflows/deploy.yml
name: Deploy to Podoru
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build and push
        run: |
          docker build -t registry.example.com/myapp:${{ github.sha }} .
          docker push registry.example.com/myapp:${{ github.sha }}

      - name: Deploy to Podoru
        run: |
          # Update service image
          curl -X PUT ${{ secrets.PODORU_API }}/services/${{ secrets.SERVICE_ID }} \
            -H "Authorization: Bearer ${{ secrets.PODORU_TOKEN }}" \
            -H "Content-Type: application/json" \
            -d '{"image": "registry.example.com/myapp:${{ github.sha }}"}'

          # Trigger deploy
          curl -X POST ${{ secrets.PODORU_API }}/services/${{ secrets.SERVICE_ID }}/deploy \
            -H "Authorization: Bearer ${{ secrets.PODORU_TOKEN }}"
```

## Roadmap

- [ ] GitHub OAuth app integration
- [ ] Automatic webhook setup
- [ ] Dockerfile builds
- [ ] Build logs streaming
- [ ] Deployment status checks
- [ ] Branch protection integration
