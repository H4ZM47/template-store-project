#!/bin/bash

# TemplateStore Development Startup Script
# This script starts the database, backend, and frontend servers

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PID_DIR="$PROJECT_ROOT/tmp/pids"
LOG_DIR="$PROJECT_ROOT/tmp/logs"
# Allow env override
BACKEND_PORT="${BACKEND_PORT:-8080}"
FRONTEND_PORT="${FRONTEND_PORT:-3000}"

# Create directories for PIDs and logs
mkdir -p "$PID_DIR" "$LOG_DIR"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[$(date '+%H:%M:%S')]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[$(date '+%H:%M:%S')] ✅ $1${NC}"
}

print_error() {
    echo -e "${RED}[$(date '+%H:%M:%S')] ❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}[$(date '+%H:%M:%S')] ⚠️  $1${NC}"
}

# Function to check if a service is running on a port
is_port_in_use() {
    lsof -Pi :$1 -sTCP:LISTEN -t >/dev/null 2>&1
}

# Function to check if PostgreSQL is running
check_postgres() {
    if command -v pg_isready >/dev/null 2>&1; then
        pg_isready -h localhost -p 5432 >/dev/null 2>&1
    else
        # Fallback: check if port 5432 is in use
        is_port_in_use 5432
    fi
}

# Function to start PostgreSQL
start_postgres() {
    print_status "Checking PostgreSQL status..."

    if check_postgres; then
        print_success "PostgreSQL is already running"
        return 0
    fi

    print_status "Starting PostgreSQL..."

    # Try different methods to start PostgreSQL based on the system
    if command -v brew >/dev/null 2>&1; then
        # macOS with Homebrew
        if brew services list | grep postgresql | grep started >/dev/null 2>&1; then
            print_success "PostgreSQL is already running via Homebrew"
        else
            brew services start postgresql@13 2>/dev/null || brew services start postgresql
            if [ $? -eq 0 ]; then
                print_success "PostgreSQL started via Homebrew"
                # Wait for PostgreSQL to be ready
                for i in {1..30}; do
                    if check_postgres; then
                        break
                    fi
                    sleep 1
                done
            else
                print_error "Failed to start PostgreSQL via Homebrew"
                return 1
            fi
        fi
    elif command -v systemctl >/dev/null 2>&1; then
        # Linux with systemd
        sudo systemctl start postgresql
        if [ $? -eq 0 ]; then
            print_success "PostgreSQL started via systemctl"
        else
            print_error "Failed to start PostgreSQL via systemctl"
            return 1
        fi
    elif command -v pg_ctl >/dev/null 2>&1; then
        # Direct pg_ctl command
        pg_ctl start -D /usr/local/var/postgres
        if [ $? -eq 0 ]; then
            print_success "PostgreSQL started via pg_ctl"
        else
            print_error "Failed to start PostgreSQL via pg_ctl"
            return 1
        fi
    else
        print_error "Could not find a way to start PostgreSQL"
        print_error "Please start PostgreSQL manually and run this script again"
        return 1
    fi

    # Final check
    if check_postgres; then
        print_success "PostgreSQL is running and ready"
    else
        print_error "PostgreSQL failed to start or is not responding"
        return 1
    fi
}

# Function to create database if it doesn't exist
setup_database() {
    print_status "Checking database setup..."

    # Check if database exists
    if psql -h localhost -U alex -lqt 2>/dev/null | cut -d \| -f 1 | grep -qw template_store; then
        print_success "Database 'template_store' already exists"
    else
        print_status "Creating database 'template_store'..."
        createdb -h localhost -U alex template_store 2>/dev/null || {
            print_warning "Could not create database as user 'alex', trying with current user..."
            createdb template_store || {
                print_error "Failed to create database 'template_store'"
                return 1
            }
        }
        print_success "Database 'template_store' created"
    fi
}

# Function to start the backend server
start_backend() {
    print_status "Starting backend server..."

    if is_port_in_use $BACKEND_PORT; then
        print_warning "Port $BACKEND_PORT is already in use. Checking if it's our backend..."
        if curl -s http://localhost:$BACKEND_PORT/health >/dev/null 2>&1; then
            print_success "Backend server is already running on port $BACKEND_PORT"
            return 0
        else
            print_error "Port $BACKEND_PORT is in use by another service"
            return 1
        fi
    fi

    # Set environment variables
    export DB_HOST=localhost
    export DB_PORT=5432
    export DB_USER=alex
    export DB_NAME=template_store
    export DB_PASSWORD=""
    export DB_SSLMODE=disable
    export PORT=$BACKEND_PORT
    export GIN_MODE=debug

    # Start backend server
    cd "$PROJECT_ROOT"
    nohup go run cmd/server/main.go > "$LOG_DIR/backend.log" 2>&1 &
    BACKEND_PID=$!
    echo $BACKEND_PID > "$PID_DIR/backend.pid"

    print_status "Backend server starting with PID $BACKEND_PID..."

    # Wait for backend to be ready
    for i in {1..30}; do
        if curl -s http://localhost:$BACKEND_PORT/health >/dev/null 2>&1; then
            print_success "Backend server is running on http://localhost:$BACKEND_PORT"
            return 0
        fi
        sleep 1
    done

    print_error "Backend server failed to start or is not responding"
    return 1
}

# Function to start the frontend server
start_frontend() {
    print_status "Starting frontend server..."

    if is_port_in_use $FRONTEND_PORT; then
        print_warning "Port $FRONTEND_PORT is already in use. Checking if it's our frontend..."
        if curl -s http://localhost:$FRONTEND_PORT >/dev/null 2>&1; then
            print_success "Frontend server is already running on port $FRONTEND_PORT"
            return 0
        else
            print_error "Port $FRONTEND_PORT is in use by another service"
            return 1
        fi
    fi

    # Start frontend server (pass through FRONTEND_PORT to app)
    cd "$PROJECT_ROOT"
    FRONTEND_PORT=$FRONTEND_PORT nohup go run cmd/web/main.go > "$LOG_DIR/frontend.log" 2>&1 &
    FRONTEND_PID=$!
    echo $FRONTEND_PID > "$PID_DIR/frontend.pid"

    print_status "Frontend server starting with PID $FRONTEND_PID..."

    # Wait for frontend to be ready
    for i in {1..15}; do
        if curl -s http://localhost:$FRONTEND_PORT >/dev/null 2>&1; then
            print_success "Frontend server is running on http://localhost:$FRONTEND_PORT"
            return 0
        fi
        sleep 1
    done

    print_error "Frontend server failed to start or is not responding"
    return 1
}

# Function to show status
show_status() {
    echo
    print_success "🚀 TemplateStore Development Environment Started!"
    echo
    echo "📊 Service Status:"
    echo "  • Database (PostgreSQL): $(check_postgres && echo "✅ Running" || echo "❌ Not Running")"
    echo "  • Backend API: $(is_port_in_use $BACKEND_PORT && echo "✅ Running on http://localhost:$BACKEND_PORT" || echo "❌ Not Running")"
    echo "  • Frontend Web: $(is_port_in_use $FRONTEND_PORT && echo "✅ Running on http://localhost:$FRONTEND_PORT" || echo "❌ Not Running")"
    echo
    echo "🌐 URLs:"
    echo "  • Main Application: http://localhost:$FRONTEND_PORT"
    echo "  • Test Page: http://localhost:$FRONTEND_PORT/test.html"
    echo "  • Backend API: http://localhost:$BACKEND_PORT/api/v1/"
    echo "  • Health Check: http://localhost:$BACKEND_PORT/health"
    echo
    echo "📝 Logs:"
    echo "  • Backend: tail -f $LOG_DIR/backend.log"
    echo "  • Frontend: tail -f $LOG_DIR/frontend.log"
    echo
    echo "🛑 To stop all services: ./scripts/stop-dev.sh"
    echo
}

# Function to cleanup on script exit
cleanup() {
    if [ $? -ne 0 ]; then
        print_error "Script failed. Check the logs for details."
        echo "Backend log: $LOG_DIR/backend.log"
        echo "Frontend log: $LOG_DIR/frontend.log"
    fi
}

trap cleanup EXIT

# Main execution
main() {
    print_status "Starting TemplateStore Development Environment..."
    echo

    # Check dependencies
    print_status "Checking dependencies..."

    if ! command -v go >/dev/null 2>&1; then
        print_error "Go is not installed. Please install Go 1.21 or higher."
        exit 1
    fi

    if ! command -v psql >/dev/null 2>&1; then
        print_error "PostgreSQL client tools are not installed."
        exit 1
    fi

    print_success "Dependencies check passed"
    echo

    # Start services in order
    start_postgres || exit 1
    setup_database || exit 1
    start_backend || exit 1
    start_frontend || exit 1

    show_status
}

# Check if script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
