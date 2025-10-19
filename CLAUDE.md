# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

JOOSUM is a link management backend service built with Go and Gin framework. The application allows users to save, organize, and manage web links in link books with tags and notifications.

## Running the Application

### Local Development
```bash
cd backend
go run main.go              # Development environment
go run main.go -env=prod    # Production environment
```

### Docker
```bash
cd backend
make docker.dev    # Development environment
make docker.prod   # Production environment
make docker.stop   # Stop containers
```

### View Logs
```bash
docker logs joosum_dev -f              # Development logs
docker logs server 2>&1 -f | jq        # Production logs (requires jq)
```

## Building & Documentation

### Swagger Documentation
Swagger docs must be regenerated after API changes:
```bash
cd backend
swag init    # Regenerates docs/docs.go, docs/swagger.json, docs/swagger.yaml
```
Access Swagger UI at `/swagger/index.html` when server is running.

## Project Structure

The codebase follows a layered architecture pattern:

### Core Application (`backend/`)
- `main.go` - Application entry point, initializes MongoDB, loads Apple public keys, sets up routes
- `cmd/` - Additional executables:
  - `scheduler/main.go` - Notification scheduler (unread/unclassified links)
  - `jwt_generator/main.go` - JWT token generation utility

### Application Modules (`backend/app/`)
Each domain module follows Handler → Usecase → Model pattern:
- `auth/` - OAuth authentication (Google, Apple, Naver), JWT token issuance
- `user/` - User management, withdrawal, context utilities
- `link/` - Link CRUD operations (create, read, update, delete)
- `linkbook/` - Link book (folder) management
- `tag/` - Tag creation and management
- `notif/` - Push notifications
- `setting/` - User settings (device ID, notification preferences)
- `banner/` - Banner management
- `page/` - Page rendering

### Infrastructure (`backend/pkg/`)
- `config/` - Environment config via Viper (config.yml), JWT and Gin setup
- `db/` - MongoDB singleton client, collection initialization and indexing
- `middleware/` - JWT authentication, user context, internal API key, logging
- `routes/` - Route registration (public, private, internal, swagger)
- `util/` - Common utilities (JWT generation, validators, error handling, Apple public key loader)

### Background Jobs (`backend/job/`)
- `notification/` - Scheduled notification handlers for unread/unclassified links

## Architecture Patterns

### Authentication Flow
1. Public routes (`/auth/*`) handle OAuth provider verification (Google/Apple/Naver)
2. OAuth token is verified with the provider
3. JWT is issued using custom JWT secret (config.yml)
4. Private routes use `SetUserData()` middleware to verify JWT and load user into context
5. User is accessible via `c.Get("user")` in handlers

### Route Organization
- **Public routes** - No authentication required (OAuth callbacks, login)
- **Private routes** - JWT authentication via `SetUserData()` middleware
- **Internal routes** - Internal API key authentication via `InternalAPIKeyMiddleware()`
- **Swagger routes** - API documentation

Note: PrivateRoutes must be registered AFTER SwaggerRoutes, otherwise Swagger UI won't be accessible.

### Database Access
- MongoDB client is a singleton initialized in `main.go` via `util.StartMongoDB()`
- Collections are initialized with indexes in `db/collection.go`
- Global collection variables: `UserCollection`, `LinkCollection`, `LinkBookCollection`, `InactiveUserCollection`, `TagCollection`, `NotificationCollection`, `NotificationAgreeCollection`, `BannerCollection`

### Configuration
- Environment config loaded from `config.yml` using Viper
- Scheduler uses separate `scheduler-config.yml`
- No `.env` files - all config in `config.yml`
- Access config values: `config.GetEnvConfig("key")`

### Error Handling
Error codes are standardized (see README.md):
- 1000 series - Request validation errors
- 2000 series - Authentication errors
- 3000 series - Business logic errors

Use `util.APIError` for consistent error responses in Swagger docs.

## Testing

No test files currently exist in the codebase. When adding tests, use Go's standard testing package:
```bash
cd backend
go test ./...           # Run all tests
go test ./app/link      # Run tests for specific package
```

## Important Notes

- The main Go module path is `joosum-backend`, all imports use this prefix
- Server runs on port 5001
- Graceful shutdown is implemented (30s timeout)
- Production environment uses custom logging middleware
- Apple public keys are loaded at startup via `util.LoadApplePublicKeys()`
- Validator is initialized globally: `util.Validate = validator.New()`
