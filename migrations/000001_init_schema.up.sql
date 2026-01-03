-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Teams/Organizations
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Team memberships
CREATE TABLE team_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(team_id, user_id)
);

-- Projects
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    github_repo VARCHAR(500),
    github_branch VARCHAR(255) DEFAULT 'main',
    github_token_encrypted BYTEA,
    auto_deploy BOOLEAN DEFAULT false,
    webhook_secret VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(team_id, slug)
);

-- Services (apps within a project)
CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    deploy_type VARCHAR(50) NOT NULL,
    image VARCHAR(500),
    dockerfile_path VARCHAR(500) DEFAULT 'Dockerfile',
    build_context VARCHAR(500) DEFAULT '.',
    compose_file VARCHAR(500),
    env_vars_encrypted BYTEA,
    replicas INTEGER DEFAULT 1,
    cpu_limit DECIMAL(5,2),
    memory_limit INTEGER,
    health_check_path VARCHAR(255),
    health_check_interval INTEGER DEFAULT 30,
    restart_policy VARCHAR(50) DEFAULT 'unless-stopped',
    status VARCHAR(50) DEFAULT 'stopped',
    container_id VARCHAR(255),
    swarm_service_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(project_id, slug)
);

-- Domains for services (Traefik routing)
CREATE TABLE domains (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    domain VARCHAR(255) NOT NULL UNIQUE,
    ssl_enabled BOOLEAN DEFAULT true,
    ssl_auto BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Port mappings
CREATE TABLE port_mappings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    container_port INTEGER NOT NULL,
    host_port INTEGER,
    protocol VARCHAR(10) DEFAULT 'tcp',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Volumes
CREATE TABLE volumes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    mount_path VARCHAR(500) NOT NULL,
    host_path VARCHAR(500),
    driver VARCHAR(50) DEFAULT 'local',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Networks
CREATE TABLE networks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    docker_network_id VARCHAR(255),
    driver VARCHAR(50) DEFAULT 'bridge',
    subnet VARCHAR(50),
    gateway VARCHAR(50),
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Service-Network junction
CREATE TABLE service_networks (
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    network_id UUID NOT NULL REFERENCES networks(id) ON DELETE CASCADE,
    PRIMARY KEY (service_id, network_id)
);

-- Deployment history
CREATE TABLE deployments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    triggered_by UUID REFERENCES users(id),
    commit_sha VARCHAR(50),
    commit_message TEXT,
    status VARCHAR(50) NOT NULL,
    logs TEXT,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    finished_at TIMESTAMP WITH TIME ZONE
);

-- Refresh tokens
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_team_members_user ON team_members(user_id);
CREATE INDEX idx_projects_team ON projects(team_id);
CREATE INDEX idx_services_project ON services(project_id);
CREATE INDEX idx_deployments_service ON deployments(service_id);
CREATE INDEX idx_domains_service ON domains(service_id);
CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);
