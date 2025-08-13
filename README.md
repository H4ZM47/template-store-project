# Template Store & Blog Platform

A modern web application for selling digital templates and hosting a company blog, built with Go, PostgreSQL, and AWS services.

## Architecture

This project follows a clean architecture pattern with:
- **Backend**: Go with Gin framework
- **Database**: PostgreSQL on AWS RDS
- **Storage**: AWS S3 with CloudFront CDN
- **Authentication**: AWS Cognito
- **Payments**: Stripe integration
- **Email**: SendGrid for notifications
- **Frontend**: Vanilla HTML/CSS/JS with Tailwind CSS

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 13 or higher
- Docker (optional, for containerized development)
- AWS CLI (for deployment)

## Quick Start

### Option 1: Using Development Scripts (Recommended)

The easiest way to get started is using our automated development scripts:

```bash
# Clone the repository
git clone <repository-url>
cd template-store-project

# Make scripts executable (Unix/Linux/macOS)
chmod +x scripts/*.sh

# Start all services (database, backend, frontend)
./scripts/start-dev.sh

# Or use make commands
make start
```

This will automatically:
- ✅ Start PostgreSQL database
- ✅ Create the `template_store` database
- ✅ Start backend API server on port 8080
- ✅ Start frontend web server on port 3000
- ✅ Open your browser to http://localhost:3000

### Option 2: Manual Setup

If you prefer manual setup or need to customize the configuration:

```bash
# 1. Clone and setup
git clone <repository-url>
cd template-store-project

# 2. Set up environment variables
cp env.example .env
# Edit .env with your configuration values

# 3. Install dependencies
go mod download

# 4. Set up PostgreSQL and database
createdb template_store

# 5. Start backend server
go run cmd/server/main.go &

# 6. Start frontend server (in another terminal)
go run cmd/web/main.go
```

### Access Your Application

Once started, you can access:
- **Main Application**: http://localhost:3000
- **Test Page**: http://localhost:3000/test.html  
- **Backend API**: http://localhost:8080/api/v1/
- **Health Check**: http://localhost:8080/health

## Development Scripts

We provide comprehensive scripts to manage your development environment:

### Available Scripts

- **`./scripts/start-dev.sh`** - Full development environment startup
- **`./scripts/stop-dev.sh`** - Graceful shutdown of all services
- **`./scripts/quick-start.sh`** - Lightweight startup for experienced developers
- **`./scripts/start-dev.bat`** - Windows batch script

### Make Commands

For even easier development, use our Makefile:

```bash
make help          # Show all available commands
make start         # Start all services
make stop          # Stop all services
make status        # Show service status
make logs          # View recent logs
make db-seed       # Seed database with sample data
make dev           # Full development setup + open browser
make clean         # Clean temporary files
```

### Development Workflow

```bash
# Fresh start
make fresh         # Stop, clean, start, seed database

# Daily development
make start         # Start all services
make logs          # Check logs if needed
make stop          # Stop when done

# Database operations
make db-seed       # Add sample data
make db-reset      # Reset database
make db-connect    # Connect with psql
```

### Script Features

- ✅ **Automatic PostgreSQL management** - Detects and starts database
- ✅ **Database setup** - Creates database if it doesn't exist
- ✅ **Process management** - Tracks PIDs, graceful shutdown
- ✅ **Comprehensive logging** - Logs stored in `tmp/logs/`
- ✅ **Cross-platform support** - Unix/Linux/macOS/Windows
- ✅ **Error handling** - Detailed error messages and troubleshooting
- ✅ **Service health checks** - Verifies all services are running

For detailed documentation, see [`scripts/README.md`](scripts/README.md).

## Development

### Project Structure
```
template-store-project/
├── cmd/                    # Application entry points
│   └── server/            # Main server application
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── models/           # Database models
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # Custom middleware
│   └── services/         # Business logic
├── pkg/                  # Public libraries
├── web/                  # Frontend assets
├── migrations/           # Database migrations
├── docs/                 # Documentation
└── scripts/              # Build and deployment scripts
```

### Available Endpoints

- `GET /health` - Health check
- `GET /api/v1/` - API information

### Environment Variables

See `env.example` for all available configuration options.

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/handlers
```

## Building

```bash
# Build for current platform
go build -o bin/server cmd/server/main.go

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o bin/server cmd/server/main.go
```

## Deployment

### Local Development
```bash
# Run with hot reload (requires air)
air

# Run with Docker
docker-compose up
```

### Production
The application is designed to be deployed on AWS ECS with:
- Application Load Balancer
- RDS PostgreSQL
- S3 for file storage
- CloudFront for CDN

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

[Add your license here]

## Support

For support and questions, please [create an issue](link-to-issues) or contact the development team. 