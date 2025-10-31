# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GOFLY LIVE CHAT is an open-source live chat support system built with Go. It provides real-time messaging between customers and support teams with a modern web-based interface.

### Key Technologies
- **Backend**: Gin (web framework), JWT (authentication), WebSocket (real-time communication), GORM (ORM), Cobra (CLI)
- **Frontend**: Vue.js, ElementUI
- **Database**: MySQL
- **Architecture**: MVC pattern with middleware, controllers, models, and WebSocket handlers

## Development Commands

### Installation & Setup
```bash
# Initialize database (run once)
go run main.go install

# Start development server
go run main.go server

# Start with custom port
go run main.go server -p 8082

# Run as daemon process
go run main.go server -d

# Build executable
go build -o gochat

# Run built binary
./gochat server
```

### Database Operations
- Database configuration: `config/mysql.json`
- Initial data import: `import.sql` (automatically processed by install command)
- Installation lock file: `install.lock` (created after successful installation)

### Testing
```bash
# Run tests
go test ./...

# Run specific test file
go test ./tools/geo_test.go

# Run tests with verbose output
go test -v ./...
```

## Code Architecture

### Project Structure
```
├── cmd/           # CLI commands (server, install)
├── controller/    # HTTP request handlers
├── middleware/    # Request middleware (auth, CORS, logging, etc.)
├── models/        # Database models and operations
├── router/        # Route definitions (API and view routes)
├── static/        # Static assets and templates
├── tools/         # Utility functions and helpers
├── ws/           # WebSocket handlers for real-time chat
├── common/       # Common configuration and utilities
└── config/       # Configuration files
```

### Key Components

#### CLI Interface (`cmd/`)
- **root.go**: Cobra CLI setup with command routing
- **server.go**: HTTP server startup, middleware configuration, background services
- **install.go**: Database initialization and lock file management

#### Web Layer (`controller/`, `router/`)
- **router/api.go**: REST API endpoints with middleware chains
- **router/view.go**: HTML page routes for web interface
- **controller/**: Request handlers for chat, users, authentication, etc.

#### Business Logic (`models/`, `ws/`)
- **models/**: GORM database models and queries
- **ws/ws.go**: WebSocket server, client management, message broadcasting
- **ws/user.go**, **ws/visitor.go**: User connection handling

#### Middleware Stack (`middleware/`)
- **session.go**: Cookie-based session management
- **jwt.go**: JWT authentication for API endpoints
- **rbac.go**: Role-based access control
- **ipblack.go**: IP blacklist filtering
- **logger.go**: Request logging
- **cross.go**: CORS handling

### Database Models
- Users, Roles, UserRoles (authentication & authorization)
- Messages, Replys (chat functionality)
- Visitors, Welcomes (visitor management)
-Configs, Abouts, Ipblacks (system configuration)

### WebSocket Communication
- **ClientList**: Active visitor connections
- **KefuList**: Active support agent connections
- **Message types**: ping/pong, inputing status, chat messages
- Real-time message broadcasting between visitors and support agents

## Configuration

### Database Configuration (`config/mysql.json`)
```json
{
  "Server": "localhost",
  "Port": "3306",
  "Database": "goflychat",
  "Username": "root",
  "Password": "root"
}
```

### Server Options
- **Default port**: 8081
- **Templates**: `static/templates/*`
- **Static files**: `/static` path
- **Logs**: `logs/` directory (in daemon mode)

## Development Guidelines

### Code Organization
- Follow MVC pattern with clear separation of concerns
- Use middleware for cross-cutting concerns (auth, logging, etc.)
- Implement proper error handling and logging
- Use descriptive variable names following Go conventions

### Authentication & Security
- JWT tokens for API authentication
- Role-based access control (RBAC)
- IP blacklist filtering
- Session management for web interface
- CORS enabled for cross-origin requests

### WebSocket Implementation
- Concurrent connection management with mutexes
- Message broadcasting with type-based routing
- Heartbeat/ping-pong for connection health
- Visitor status tracking and cleanup

### Database Operations
- Use GORM for database operations
- Connection pooling configured (10 idle, 100 max connections)
- Singular table naming convention
- Proper error handling for database failures

## Important Notes

- Installation creates `install.lock` file - remove to reinstall
- Server listens on all interfaces (0.0.0.0) by default
- WebSocket connections support cross-origin requests
- Daemon mode creates logs in `logs/gofly.log`
- Database uses utf8mb4 charset for full Unicode support