# Development Scripts

This directory contains scripts to help manage the TemplateStore development environment. These scripts automate the process of starting and stopping the database, backend, and frontend servers.

## Scripts Overview

### 🚀 `start-dev.sh` (Primary Script)
**Full-featured development environment startup script for Unix/Linux/macOS**

- Automatically detects and starts PostgreSQL
- Creates the database if it doesn't exist
- Starts backend API server on port 8080
- Starts frontend web server on port 3000
- Provides comprehensive error handling and logging
- Shows detailed status and URLs

**Usage:**
```bash
./scripts/start-dev.sh
```

**Features:**
- ✅ Dependency checking (Go, PostgreSQL)
- ✅ Service status validation
- ✅ Automatic database setup
- ✅ Process management with PID files
- ✅ Comprehensive logging
- ✅ Graceful error handling

### 🛑 `stop-dev.sh`
**Graceful shutdown script for all services**

- Stops backend and frontend servers gracefully
- Optionally stops PostgreSQL
- Cleans up PID files and old logs
- Force kills processes if graceful shutdown fails

**Usage:**
```bash
./scripts/stop-dev.sh                    # Stop frontend and backend only
./scripts/stop-dev.sh --with-postgres    # Stop all services including PostgreSQL
./scripts/stop-dev.sh -h                 # Show help
```

### ⚡ `quick-start.sh`
**Simple, lightweight startup script**

- Minimal script for quick development setup
- Less error checking, faster startup
- Good for experienced developers who know their environment is set up correctly

**Usage:**
```bash
./scripts/quick-start.sh
```

### 🪟 `start-dev.bat`
**Windows batch script for Windows developers**

- Full-featured startup script for Windows
- Handles Windows-specific PostgreSQL service management
- Uses PowerShell for web requests and colored output
- Creates Windows command windows for each service

**Usage:**
```cmd
scripts\start-dev.bat
```

## Prerequisites

### All Platforms
- Go 1.21 or higher
- PostgreSQL 13 or higher
- Git (for version control)

### macOS
- Homebrew (recommended for PostgreSQL)
- Xcode Command Line Tools

### Linux
- PostgreSQL client tools (`postgresql-client`)
- systemd (for service management)

### Windows
- PostgreSQL Windows service
- PowerShell (for enhanced functionality)

## Quick Start Guide

1. **First Time Setup:**
   ```bash
   # Make scripts executable (Unix/Linux/macOS)
   chmod +x scripts/*.sh
   
   # Install dependencies
   brew install postgresql@13  # macOS
   # or
   sudo apt-get install postgresql postgresql-contrib  # Ubuntu/Debian
   ```

2. **Start Development Environment:**
   ```bash
   # Full startup with error checking
   ./scripts/start-dev.sh
   
   # Or quick start
   ./scripts/quick-start.sh
   ```

3. **Access Your Application:**
   - Main Application: http://localhost:3000
   - Test Page: http://localhost:3000/test.html
   - Backend API: http://localhost:8080/api/v1/
   - Health Check: http://localhost:8080/health

4. **Stop Services:**
   ```bash
   ./scripts/stop-dev.sh
   ```

## Environment Variables

The scripts automatically set these environment variables:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=alex              # Default user (configurable)
DB_NAME=template_store
DB_PASSWORD=""            # Empty for local development
DB_SSLMODE=disable
PORT=8080                # Backend port
GIN_MODE=debug           # Development mode
```

## Logging and Debugging

### Log Files Location
- **Backend:** `tmp/logs/backend.log`
- **Frontend:** `tmp/logs/frontend.log`

### Viewing Logs
```bash
# Watch backend logs in real-time
tail -f tmp/logs/backend.log

# Watch frontend logs in real-time
tail -f tmp/logs/frontend.log

# View recent errors
grep -i error tmp/logs/backend.log
```

### Process Management
- **PID Files:** `tmp/pids/backend.pid`, `tmp/pids/frontend.pid`
- **Manual Process Check:** `ps aux | grep "go run"`
- **Manual Kill:** `kill $(cat tmp/pids/backend.pid)`

## Troubleshooting

### Common Issues

#### 1. Port Already in Use
**Error:** `Port 8080 is already in use`

**Solution:**
```bash
# Find what's using the port
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use a different port
export PORT=8081
```

#### 2. PostgreSQL Not Starting
**Error:** `PostgreSQL failed to start`

**Solutions:**
```bash
# macOS with Homebrew
brew services restart postgresql@13

# Linux with systemd
sudo systemctl restart postgresql

# Check PostgreSQL status
pg_isready -h localhost -p 5432
```

#### 3. Database Connection Failed
**Error:** `Failed to connect to database`

**Solutions:**
```bash
# Check if database exists
psql -h localhost -l

# Create database manually
createdb template_store

# Check PostgreSQL logs
# macOS: /usr/local/var/log/postgres.log
# Linux: /var/log/postgresql/
```

#### 4. Go Dependencies Issues
**Error:** `go: module not found`

**Solutions:**
```bash
# Download dependencies
go mod download

# Clean module cache
go clean -modcache

# Tidy modules
go mod tidy
```

#### 5. Permission Denied (Unix/Linux)
**Error:** `Permission denied`

**Solution:**
```bash
# Make scripts executable
chmod +x scripts/*.sh

# Fix ownership if needed
sudo chown -R $USER:$USER .
```

### Advanced Debugging

#### Enable Verbose Logging
```bash
# Add to start-dev.sh
export GIN_MODE=debug
export DB_DEBUG=true
```

#### Database Debugging
```bash
# Connect to database directly
psql -h localhost -U alex -d template_store

# Check database tables
\dt

# View recent database activity
SELECT * FROM pg_stat_activity;
```

#### Network Debugging
```bash
# Check listening ports
netstat -tulnp | grep LISTEN

# Test API connectivity
curl -v http://localhost:8080/health

# Test frontend connectivity
curl -v http://localhost:3000
```

## Customization

### Changing Default Ports
Edit the configuration variables at the top of each script:

```bash
BACKEND_PORT=8080   # Change to your preferred port
FRONTEND_PORT=3000  # Change to your preferred port
```

### Custom Database Configuration
Modify the environment variables in the scripts:

```bash
export DB_USER=your_username
export DB_NAME=your_database
export DB_PASSWORD=your_password
```

### Adding Custom Services
To add additional services to the startup process:

1. Add service configuration to the script
2. Create start/stop functions following existing patterns
3. Add service to the main execution flow
4. Update status reporting

## Contributing

When modifying these scripts:

1. **Test on multiple platforms** (if possible)
2. **Maintain backward compatibility**
3. **Add appropriate error handling**
4. **Update this README**
5. **Follow existing code style**

## Platform-Specific Notes

### macOS
- Uses Homebrew for PostgreSQL management
- Requires Xcode Command Line Tools
- Uses `brew services` for service management

### Linux
- Uses systemd for service management
- Requires sudo for PostgreSQL operations
- Package manager varies by distribution

### Windows
- Uses Windows services for PostgreSQL
- Requires PowerShell for advanced features
- Different path separators and commands

## Performance Tips

1. **Use `quick-start.sh` for faster startup** when you know your environment is configured
2. **Keep PostgreSQL running** between sessions to avoid startup delays
3. **Use SSD storage** for better database performance
4. **Close unnecessary applications** to free up ports and resources

## Security Considerations

- These scripts are for **development only**
- Never use empty passwords in production
- Database runs without SSL for local development
- Services bind to localhost only for security

---

## Support

If you encounter issues:

1. Check the troubleshooting section above
2. Review the log files in `tmp/logs/`
3. Ensure all prerequisites are installed
4. Try the manual startup commands
5. Check the main project README.md

**Happy coding! 🚀**