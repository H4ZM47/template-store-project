#!/bin/bash

# TemplateStore Quick Start Script
# Simple script to quickly start all services for development

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}🚀 Starting TemplateStore Development Environment...${NC}"

# Navigate to project root
cd "$(dirname "$0")/.."

# Start PostgreSQL (macOS with Homebrew)
echo -e "${BLUE}📊 Starting PostgreSQL...${NC}"
if command -v brew >/dev/null 2>&1; then
    brew services start postgresql@13 2>/dev/null || brew services start postgresql
else
    echo -e "${RED}⚠️  Please start PostgreSQL manually${NC}"
fi

# Wait a moment for PostgreSQL to start
sleep 2

# Create database if it doesn't exist
echo -e "${BLUE}🗄️  Setting up database...${NC}"
createdb template_store 2>/dev/null || echo "Database already exists"

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=alex
export DB_NAME=template_store
export DB_PASSWORD=""
export DB_SSLMODE=disable
export GIN_MODE=debug

# Start backend server in background
echo -e "${BLUE}⚙️  Starting backend server...${NC}"
go run cmd/server/main.go &
BACKEND_PID=$!

# Wait for backend to start
sleep 3

# Start frontend server in background
echo -e "${BLUE}🌐 Starting frontend server...${NC}"
go run cmd/web/main.go &
FRONTEND_PID=$!

# Wait for frontend to start
sleep 2

echo
echo -e "${GREEN}✅ TemplateStore is now running!${NC}"
echo
echo "🌐 URLs:"
echo "  • Main Application: http://localhost:3000"
echo "  • Test Page: http://localhost:3000/test.html"
echo "  • Backend API: http://localhost:8080/api/v1/"
echo "  • Health Check: http://localhost:8080/health"
echo
echo "📝 Backend PID: $BACKEND_PID"
echo "📝 Frontend PID: $FRONTEND_PID"
echo
echo "🛑 To stop servers:"
echo "  kill $BACKEND_PID $FRONTEND_PID"
echo "  or use: ./scripts/stop-dev.sh"
echo

# Keep script running and show logs
echo -e "${BLUE}📋 Press Ctrl+C to stop all services${NC}"
echo

# Function to cleanup on exit
cleanup() {
    echo
    echo -e "${BLUE}🛑 Stopping services...${NC}"
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
    echo -e "${GREEN}✅ Services stopped${NC}"
    exit 0
}

trap cleanup INT

# Wait for user interrupt
wait
