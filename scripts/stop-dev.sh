#!/bin/bash

# TemplateStore Development Shutdown Script
# This script stops the backend and frontend servers

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
BACKEND_PORT="${BACKEND_PORT:-8080}"
FRONTEND_PORT="${FRONTEND_PORT:-3000}"

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

# Function to check if a process is running
is_process_running() {
    kill -0 "$1" 2>/dev/null
}

# Function to stop a service by PID file
stop_service() {
    local service_name=$1
    local pid_file="$PID_DIR/${service_name}.pid"

    if [ ! -f "$pid_file" ]; then
        print_warning "No PID file found for $service_name"
        return 0
    fi

    local pid=$(cat "$pid_file")

    if [ -z "$pid" ]; then
        print_warning "Empty PID file for $service_name"
        rm -f "$pid_file"
        return 0
    fi

    if is_process_running "$pid"; then
        print_status "Stopping $service_name (PID: $pid)..."

        # Try graceful shutdown first
        kill "$pid" 2>/dev/null

        # Wait up to 10 seconds for graceful shutdown
        for i in {1..10}; do
            if ! is_process_running "$pid"; then
                print_success "$service_name stopped gracefully"
                rm -f "$pid_file"
                return 0
            fi
            sleep 1
        done

        # Force kill if still running
        print_warning "Force killing $service_name..."
        kill -9 "$pid" 2>/dev/null

        # Check if force kill worked
        if ! is_process_running "$pid"; then
            print_success "$service_name force stopped"
            rm -f "$pid_file"
        else
            print_error "Failed to stop $service_name"
            return 1
        fi
    else
        print_warning "$service_name was not running (stale PID file)"
        rm -f "$pid_file"
    fi
}

# Function to stop services by port (fallback method)
stop_by_port() {
    local service_name=$1
    local port=$2

    local pid=$(lsof -ti :$port 2>/dev/null)

    if [ -n "$pid" ]; then
        print_status "Found $service_name running on port $port (PID: $pid)"

        # Try graceful shutdown
        kill "$pid" 2>/dev/null

        # Wait for graceful shutdown
        for i in {1..10}; do
            if ! kill -0 "$pid" 2>/dev/null; then
                print_success "$service_name stopped"
                return 0
            fi
            sleep 1
        done

        # Force kill
        print_warning "Force killing $service_name on port $port..."
        kill -9 "$pid" 2>/dev/null
        print_success "$service_name force stopped"
    else
        print_success "$service_name is not running on port $port"
    fi
}

# Function to stop PostgreSQL (optional)
stop_postgres() {
    if [ "$1" = "--with-postgres" ] || [ "$1" = "-p" ]; then
        print_status "Stopping PostgreSQL..."

        if command -v brew >/dev/null 2>&1; then
            # macOS with Homebrew
            brew services stop postgresql@13 2>/dev/null || brew services stop postgresql
            if [ $? -eq 0 ]; then
                print_success "PostgreSQL stopped via Homebrew"
            else
                print_warning "Could not stop PostgreSQL via Homebrew (may not be running)"
            fi
        elif command -v systemctl >/dev/null 2>&1; then
            # Linux with systemd
            sudo systemctl stop postgresql
            if [ $? -eq 0 ]; then
                print_success "PostgreSQL stopped via systemctl"
            else
                print_warning "Could not stop PostgreSQL via systemctl"
            fi
        elif command -v pg_ctl >/dev/null 2>&1; then
            # Direct pg_ctl command
            pg_ctl stop -D /usr/local/var/postgres
            if [ $? -eq 0 ]; then
                print_success "PostgreSQL stopped via pg_ctl"
            else
                print_warning "Could not stop PostgreSQL via pg_ctl"
            fi
        else
            print_warning "Could not find a way to stop PostgreSQL"
        fi
    else
        print_status "PostgreSQL left running (use --with-postgres to stop it)"
    fi
}

# Function to clean up temporary files
cleanup_temp_files() {
    print_status "Cleaning up temporary files..."

    # Clean up old log files (keep last 5)
    if [ -d "$LOG_DIR" ]; then
        find "$LOG_DIR" -name "*.log" -type f -mtime +7 -delete 2>/dev/null || true
    fi

    # Clean up any stale PID files
    if [ -d "$PID_DIR" ]; then
        for pid_file in "$PID_DIR"/*.pid; do
            if [ -f "$pid_file" ]; then
                local pid=$(cat "$pid_file" 2>/dev/null)
                if [ -n "$pid" ] && ! is_process_running "$pid"; then
                    rm -f "$pid_file"
                fi
            fi
        done
    fi

    print_success "Cleanup completed"
}

# Function to show final status
show_final_status() {
    echo
    print_success "🛑 TemplateStore Development Environment Stopped"
    echo
    echo "📊 Final Status:"

    # Check if services are still running
    local backend_running=$(lsof -ti :$BACKEND_PORT 2>/dev/null)
    local frontend_running=$(lsof -ti :$FRONTEND_PORT 2>/dev/null)

    echo "  • Backend API: $([ -n "$backend_running" ] && echo "❌ Still running on port $BACKEND_PORT" || echo "✅ Stopped")"
    echo "  • Frontend Web: $([ -n "$frontend_running" ] && echo "❌ Still running on port $FRONTEND_PORT" || echo "✅ Stopped")"

    if [ -n "$backend_running" ] || [ -n "$frontend_running" ]; then
        echo
        print_warning "Some services are still running. You may need to stop them manually:"
        [ -n "$backend_running" ] && echo "  Backend: kill $backend_running"
        [ -n "$frontend_running" ] && echo "  Frontend: kill $frontend_running"
    fi

    echo
    echo "📝 Logs are preserved in: $LOG_DIR"
    echo "🚀 To start again: ./scripts/start-dev.sh"
    echo
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -p, --with-postgres    Also stop PostgreSQL server"
    echo "  -h, --help            Show this help message"
    echo
    echo "Examples:"
    echo "  $0                    # Stop frontend and backend only"
    echo "  $0 --with-postgres    # Stop all services including PostgreSQL"
    echo
}

# Main execution
main() {
    local stop_postgres_flag=""

    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -p|--with-postgres)
                stop_postgres_flag="--with-postgres"
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done

    print_status "Stopping TemplateStore Development Environment..."
    echo

    # Stop services in reverse order
    stop_service "frontend"
    stop_service "backend"

    # Fallback: stop by port if PID files didn't work
    stop_by_port "frontend" $FRONTEND_PORT
    stop_by_port "backend" $BACKEND_PORT

    # Stop PostgreSQL if requested
    stop_postgres "$stop_postgres_flag"

    # Clean up
    cleanup_temp_files

    show_final_status
}

# Check if script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
