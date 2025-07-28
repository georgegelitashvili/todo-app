# Development Scripts

This directory contains scripts to easily start the Todo App development environment.

## Scripts

### Windows
- `start.bat` - Windows batch script to start the development environment

### Unix/Linux/macOS
- `start.sh` - Bash script to start the development environment

## What the scripts do:

1. **Check Docker** - Verify Docker is installed and running
2. **Start Cassandra** - Create/start Cassandra container if needed
3. **Wait for Database** - Wait for Cassandra to be fully initialized
4. **Start Todo App** - Run the Go application
5. **Open Browser** - Automatically open http://localhost:8080
6. **Cleanup** - Handle graceful shutdown when stopped

## Usage

### Windows:
```cmd
cd C:\Users\User\Desktop\todo-app
dev\start.bat
```

### Unix/Linux/macOS:
```bash
cd /path/to/todo-app
chmod +x dev/start.sh
./dev/start.sh
```

## Requirements

- Docker installed and running
- Go installed and configured
- Port 8080 available for the web server
- Port 9042 available for Cassandra

## Stopping

Press Ctrl+C (Unix) or any key (Windows) to stop the development environment. The scripts will automatically clean up running processes.
