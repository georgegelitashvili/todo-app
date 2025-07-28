#!/bin/bash

# Todo App Development Startup Script
echo "🚀 Starting Todo App Development Environment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to check if Docker is running
check_docker() {
    echo -e "${BLUE}📋 Checking Docker...${NC}"
    if ! docker --version >/dev/null 2>&1; then
        echo -e "${RED}❌ Docker is not installed or not running${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Docker is available${NC}"
}

# Function to start Cassandra
start_cassandra() {
    echo -e "${BLUE}🗄️  Setting up Cassandra database...${NC}"
    
    # Check if Cassandra container exists
    if docker ps -a --format "table {{.Names}}" | grep -q "^cassandra$"; then
        echo -e "${YELLOW}📦 Cassandra container exists${NC}"
        
        # Check if it's running
        if docker ps --format "table {{.Names}}" | grep -q "^cassandra$"; then
            echo -e "${GREEN}✅ Cassandra is already running${NC}"
        else
            echo -e "${YELLOW}🔄 Starting existing Cassandra container...${NC}"
            docker start cassandra
        fi
    else
        echo -e "${YELLOW}📥 Creating new Cassandra container...${NC}"
        docker run --name cassandra -p 9042:9042 -d cassandra:latest
    fi
    
    # Wait for Cassandra to be ready
    echo -e "${YELLOW}⏳ Waiting for Cassandra to initialize...${NC}"
    max_attempts=30
    attempts=0
    
    while [ $attempts -lt $max_attempts ]; do
        if docker exec cassandra cqlsh -e "DESCRIBE KEYSPACES;" >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Cassandra is ready!${NC}"
            break
        fi
        
        attempts=$((attempts + 1))
        echo -e "${YELLOW}   Attempt $attempts/$max_attempts - waiting...${NC}"
        sleep 2
    done
    
    if [ $attempts -eq $max_attempts ]; then
        echo -e "${RED}❌ Cassandra failed to start within expected time${NC}"
        exit 1
    fi
}

# Function to start the Go application
start_app() {
    echo -e "${BLUE}🚀 Starting Todo App...${NC}"
    
    # Navigate to project root
    cd "$(dirname "$0")/.."
    
    # Start the Go application in background
    go run . &
    APP_PID=$!
    
    echo -e "${GREEN}✅ Todo App started with PID: $APP_PID${NC}"
    
    # Wait a moment for the server to start
    sleep 3
    
    # Check if the server is responding
    max_attempts=10
    attempts=0
    
    while [ $attempts -lt $max_attempts ]; do
        if curl -s http://localhost:8080 >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Server is responding on http://localhost:8080${NC}"
            break
        fi
        
        attempts=$((attempts + 1))
        echo -e "${YELLOW}   Waiting for server... ($attempts/$max_attempts)${NC}"
        sleep 1
    done
    
    if [ $attempts -eq $max_attempts ]; then
        echo -e "${RED}❌ Server failed to respond${NC}"
        kill $APP_PID 2>/dev/null
        exit 1
    fi
}

# Function to open browser
open_browser() {
    echo -e "${BLUE}🌐 Opening browser...${NC}"
    
    # Detect OS and open browser accordingly
    case "$OSTYPE" in
        darwin*)  # macOS
            open http://localhost:8080
            ;;
        linux*)   # Linux
            if command -v xdg-open > /dev/null; then
                xdg-open http://localhost:8080
            elif command -v gnome-open > /dev/null; then
                gnome-open http://localhost:8080
            fi
            ;;
        msys*|cygwin*|mingw*)    # Windows Git Bash
            start http://localhost:8080
            ;;
        *)
            echo -e "${YELLOW}⚠️  Please open http://localhost:8080 in your browser${NC}"
            ;;
    esac
    
    echo -e "${GREEN}✅ Browser opened${NC}"
}

# Function to handle cleanup on exit
cleanup() {
    echo -e "\n${YELLOW}🛑 Shutting down...${NC}"
    if [ ! -z "$APP_PID" ]; then
        kill $APP_PID 2>/dev/null
        echo -e "${GREEN}✅ Todo App stopped${NC}"
    fi
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Main execution
main() {
    echo -e "${GREEN}🎯 Todo App Development Startup${NC}"
    echo -e "${GREEN}================================${NC}"
    
    check_docker
    start_cassandra
    start_app
    open_browser
    
    echo -e "\n${GREEN}🎉 Development environment is ready!${NC}"
    echo -e "${GREEN}📋 Application: http://localhost:8080${NC}"
    echo -e "${GREEN}🗄️  Database: Cassandra on localhost:9042${NC}"
    echo -e "\n${YELLOW}Press Ctrl+C to stop the development environment${NC}"
    
    # Keep the script running
    wait $APP_PID
}

# Run main function
main "$@"
