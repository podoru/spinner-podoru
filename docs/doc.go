// Package docs Podoru API Documentation
//
// This is the API documentation for Podoru - a Docker container management platform.
//
//	@title						Podoru API
//	@version					1.0.0
//	@description				Podoru is a Docker container management platform similar to EasyPanel.
//	@description				It provides APIs for managing containers, services, networks, and Docker Swarm clusters.
//
//	@contact.name				Podoru Team
//	@contact.url				https://github.com/podoru/spinner-podoru
//
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//
//	@host						localhost:8080
//	@BasePath					/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"
//
//	@tag.name					auth
//	@tag.description			Authentication endpoints
//
//	@tag.name					users
//	@tag.description			User management endpoints
//
//	@tag.name					teams
//	@tag.description			Team/Organization management endpoints
//
//	@tag.name					projects
//	@tag.description			Project management endpoints
//
//	@tag.name					services
//	@tag.description			Service/Container management endpoints
//
//	@tag.name					networks
//	@tag.description			Network management endpoints
//
//	@tag.name					volumes
//	@tag.description			Volume management endpoints
//
//	@tag.name					swarm
//	@tag.description			Docker Swarm management endpoints
package docs
