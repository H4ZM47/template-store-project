# Template Store - Startup Guide

This guide provides step-by-step instructions for setting up and running the Template Store application locally.

## Quick Start (Automated)

We provide an automated startup tool that handles all the setup steps for you:

**Option 1: Using Make (Simplest)**
```bash
# First time setup
cp env.example .env
# Edit .env with your configuration

# Start everything
make auto-start
```

**Option 2: Direct Command**
```bash
# Build the startup tool
go build -o startup cmd/startup/main.go

# Run it
./startup
```

The automated tool will:
- ✓ Check all prerequisites (Docker, Go, etc.)
- ✓ Start PostgreSQL database
- ✓ Build the backend server
- ✓ Start the backend server
- ✓ Start the frontend server
- ✓ Open your browser automatically
- ✓ Handle graceful shutdown with Ctrl+C

**First time setup:**
1. Make sure you have Docker Desktop running
2. Copy and configure your `.env` file: `cp env.example .env`
3. Edit the `.env` file with your AWS, Stripe, and other credentials
4. Run: `make auto-start` or `./startup`

That's it! Continue reading for manual setup instructions and detailed configuration.

---

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go** (version 1.21 or higher)
  - Download from: https://golang.org/dl/
  - Verify installation: `go version`

- **Docker Desktop** (for PostgreSQL database)
  - Download from: https://www.docker.com/products/docker-desktop
  - Verify installation: `docker --version` and `docker-compose --version`

- **Python 3** (for serving frontend files)
  - Most systems have this pre-installed
  - Verify installation: `python3 --version`

- **Git** (for version control)
  - Download from: https://git-scm.com/downloads
  - Verify installation: `git --version`

## Initial Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd template-store-project
```

### 2. Configure Environment Variables

Copy the example environment file and configure it:

```bash
cp env.example .env
```

Edit `.env` file with your configuration:

**Required for local development:**
```env
# Server Configuration
PORT=8080
GIN_MODE=debug

# Database Configuration (matches docker-compose.yml)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=template_store
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-minimum-32-characters
JWT_ISSUER=template-store
```

**Required for AWS services (Cognito, S3):**
```env
# AWS Configuration
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
AWS_S3_BUCKET=your-s3-bucket-name
AWS_COGNITO_POOL_ID=your_cognito_pool_id
AWS_COGNITO_APP_CLIENT_ID=your_cognito_app_client_id
```

**Required for payment processing:**
```env
# Stripe Configuration
STRIPE_API_KEY=sk_test_your_stripe_secret_key
STRIPE_PUBLISHABLE_KEY=pk_test_your_stripe_publishable_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret
STRIPE_SUCCESS_URL=http://localhost:3000/payment/success
STRIPE_CANCEL_URL=http://localhost:3000/payment/cancel
```

**Optional (for email notifications):**
```env
# SendGrid Configuration
SENDGRID_API_KEY=your_sendgrid_api_key
SENDGRID_FROM_EMAIL=noreply@templatestore.com
```

### 3. Install Go Dependencies

```bash
go mod download
```

## Database Setup

### Option 1: Using Docker Compose (Recommended)

Start PostgreSQL database using Docker Compose:

```bash
# Start PostgreSQL in detached mode
docker-compose up -d postgres

# Verify database is running
docker ps

# Check database logs
docker-compose logs postgres
```

The database will be available at `localhost:5432` with:
- **Database**: `template_store`
- **User**: `postgres`
- **Password**: `postgres`

### Option 2: Local PostgreSQL Installation

If you prefer to use a local PostgreSQL installation:

1. Install PostgreSQL: https://www.postgresql.org/download/
2. Create database:
   ```bash
   psql -U postgres
   CREATE DATABASE template_store;
   \q
   ```
3. Update `.env` with your PostgreSQL credentials

## Backend Setup

### 1. Build the Backend

Build the Go application:

```bash
go build -o main cmd/server/main.go
```

This creates an executable named `main` in your project root.

### 2. Start the Backend Server

Run the backend server:

```bash
./main
```

The server will:
- Start on port `8080` (or the port specified in `.env`)
- Automatically run database migrations
- Connect to PostgreSQL
- Initialize AWS services (Cognito, S3)
- Set up Stripe payment processing

**Expected output:**
```
[GIN-debug] Listening and serving HTTP on :8080
```

**Health check:**
```bash
curl http://localhost:8080/health
```

## Frontend Setup

### 1. Navigate to Web Directory

The frontend files are located in the `web/` directory:
- `web/index.html` - Main HTML file
- `web/css/` - Stylesheets
- `web/js/` - JavaScript files

### 2. Start Frontend Server

Start a simple HTTP server to serve the frontend:

```bash
cd web
python3 -m http.server 3000
```

**Alternative (if Python 3 is not available):**
```bash
# Using Python 2
python -m SimpleHTTPServer 3000

# Using Node.js (if you have it installed)
npx http-server -p 3000
```

The frontend will be available at: http://localhost:3000

## Testing the Application

### 1. Open in Browser

Navigate to: http://localhost:3000

### 2. Verify Components

Check that the following work:
- **Homepage loads** with template listings
- **Template modals** open when clicking templates
- **Stripe buy buttons** appear for configured templates
- **API connectivity** (check browser console for errors)

### 3. Test API Endpoints

Test key endpoints:

```bash
# Health check
curl http://localhost:8080/health

# Get templates
curl http://localhost:8080/api/templates

# Get specific template (replace {id} with actual ID)
curl http://localhost:8080/api/templates/{id}
```

### 4. Check Database Connection

Verify database tables were created:

```bash
docker exec -it template-store-project-postgres-1 psql -U postgres -d template_store -c "\dt"
```

Expected tables:
- `users`
- `templates`
- `orders`
- `customizations`

## Startup with Docker Compose (Full Stack)

To start the entire application using Docker Compose:

```bash
# Build and start all services
docker-compose up --build

# Or run in detached mode
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes database data)
docker-compose down -v
```

This will start:
- PostgreSQL database on port `5432`
- Go backend server on port `8080`
- Frontend will still need to be served separately

## Troubleshooting

### Database Connection Issues

**Problem:** Backend cannot connect to database

**Solutions:**
1. Verify PostgreSQL is running:
   ```bash
   docker ps
   ```

2. Check database logs:
   ```bash
   docker-compose logs postgres
   ```

3. Test database connection:
   ```bash
   docker exec -it template-store-project-postgres-1 psql -U postgres -d template_store
   ```

4. Verify `.env` file has correct database credentials

### Port Already in Use

**Problem:** Port 8080 or 3000 is already in use

**Solutions:**
1. Find and kill process using the port:
   ```bash
   # On macOS/Linux
   lsof -ti:8080 | xargs kill -9
   lsof -ti:3000 | xargs kill -9

   # On Windows
   netstat -ano | findstr :8080
   taskkill /PID <pid> /F
   ```

2. Or change the port in `.env`:
   ```env
   PORT=8081
   ```

### Build Errors

**Problem:** Go build fails with missing dependencies

**Solutions:**
1. Clean and reinstall dependencies:
   ```bash
   go clean -modcache
   go mod download
   go mod tidy
   ```

2. Update dependencies:
   ```bash
   go get -u ./...
   go mod tidy
   ```

### Frontend Not Loading

**Problem:** Frontend shows blank page or errors

**Solutions:**
1. Check browser console for errors
2. Verify backend is running: `curl http://localhost:8080/health`
3. Check CORS configuration in backend
4. Ensure you're serving from the `web/` directory

### AWS Services Not Working

**Problem:** Cognito or S3 errors

**Solutions:**
1. Verify AWS credentials in `.env`
2. Check AWS region matches your resources
3. Ensure IAM permissions are correct
4. Test AWS credentials:
   ```bash
   aws sts get-caller-identity
   ```

### Stripe Payment Issues

**Problem:** Buy buttons not working

**Solutions:**
1. Verify Stripe keys in `.env`
2. Check browser console for Stripe errors
3. Ensure Stripe.js is loaded: view page source and check for script tag
4. Test with Stripe test mode cards: https://stripe.com/docs/testing

## Development Workflow

### Making Changes

1. **Backend changes**: Rebuild and restart
   ```bash
   go build -o main cmd/server/main.go
   ./main
   ```

2. **Frontend changes**: Just refresh browser (no rebuild needed)

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/handlers/...
```

### Viewing Logs

```bash
# Backend logs (if running in terminal, view directly)
# If running in background, check process output

# Database logs
docker-compose logs -f postgres

# Docker Compose logs (if using full stack)
docker-compose logs -f
```

## Shutting Down

### Stop Services

1. **Backend**: Press `Ctrl+C` in the terminal running `./main`

2. **Frontend**: Press `Ctrl+C` in the terminal running `python3 -m http.server`

3. **Database**:
   ```bash
   # Stop but keep data
   docker-compose stop postgres

   # Stop and remove container (keeps data volume)
   docker-compose down

   # Stop and remove everything including data
   docker-compose down -v
   ```

## Quick Start Summary

For experienced developers, here's the quick start:

```bash
# 1. Setup
cp env.example .env
# Edit .env with your configuration

# 2. Start database
docker-compose up -d postgres

# 3. Build and start backend
go build -o main cmd/server/main.go
./main &

# 4. Start frontend
cd web && python3 -m http.server 3000 &

# 5. Open browser
open http://localhost:3000
```

## Additional Resources

- **Project Structure**: See `README.md` (if available)
- **API Documentation**: `/docs` endpoint (if configured)
- **Stripe Documentation**: https://stripe.com/docs
- **AWS Cognito Documentation**: https://docs.aws.amazon.com/cognito/
- **Gin Framework Documentation**: https://gin-gonic.com/docs/

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review application logs
3. Check Docker container status
4. Verify environment configuration
5. Create an issue in the project repository
