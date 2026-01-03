-- Drop indexes
DROP INDEX IF EXISTS idx_refresh_tokens_expires;
DROP INDEX IF EXISTS idx_refresh_tokens_user;
DROP INDEX IF EXISTS idx_domains_service;
DROP INDEX IF EXISTS idx_deployments_service;
DROP INDEX IF EXISTS idx_services_project;
DROP INDEX IF EXISTS idx_projects_team;
DROP INDEX IF EXISTS idx_team_members_user;

-- Drop tables in reverse order of creation (due to foreign key dependencies)
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS deployments;
DROP TABLE IF EXISTS service_networks;
DROP TABLE IF EXISTS networks;
DROP TABLE IF EXISTS volumes;
DROP TABLE IF EXISTS port_mappings;
DROP TABLE IF EXISTS domains;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
