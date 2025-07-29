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

### 1. Clone the repository
```bash
git clone <repository-url>
cd template-store-project
```

### 2. Set up environment variables
```bash
cp env.example .env
# Edit .env with your configuration values
```

### 3. Install dependencies
```bash
go mod download
```

### 4. Set up the database
```bash
# Create the database
createdb template_store

# Run migrations (when implemented)
# go run cmd/migrate/main.go
```

### 5. Run the application
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

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