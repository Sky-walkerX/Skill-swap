#!/bin/bash

# Skill Swap Application Startup Script
# This script starts both the backend and frontend servers

set -e

echo "🚀 Starting Skill Swap Application..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default ports
BACKEND_PORT=${BACKEND_PORT:-8080}
FRONTEND_PORT=${FRONTEND_PORT:-3000}

echo -e "${BLUE}Configuration:${NC}"
echo -e "  Backend Port: ${GREEN}$BACKEND_PORT${NC}"
echo -e "  Frontend Port: ${GREEN}$FRONTEND_PORT${NC}"

# Function to cleanup processes on exit
cleanup() {
    echo -e "\n${YELLOW}Shutting down servers...${NC}"
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi
    exit 0
}

# Trap cleanup function on script exit
trap cleanup EXIT INT TERM

# Check if required directories exist
if [ ! -d "backend/skillswap" ] || [ ! -d "frontend" ]; then
    echo -e "${RED}Error: Please run this script from the root directory of the Skill-swap project${NC}"
    exit 1
fi

# Start backend server
echo -e "\n${BLUE}Starting Backend Server...${NC}"
cd backend/skillswap
export PORT=$BACKEND_PORT
export FRONTEND_URL="http://localhost:$FRONTEND_PORT"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Build and start the backend
go mod tidy
go run cmd/server/main.go &
BACKEND_PID=$!
cd ../..

echo -e "${GREEN}Backend server starting on port $BACKEND_PORT (PID: $BACKEND_PID)${NC}"

# Wait a moment for backend to start
sleep 3

# Start frontend server
echo -e "\n${BLUE}Starting Frontend Server...${NC}"
cd frontend

# Check if Node.js and npm are installed
if ! command -v node &> /dev/null; then
    echo -e "${RED}Error: Node.js is not installed${NC}"
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo -e "${RED}Error: npm is not installed${NC}"
    exit 1
fi

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}Installing frontend dependencies...${NC}"
    npm install
fi

# Set environment variables for frontend
export FRONTEND_PORT=$FRONTEND_PORT
export NEXT_PUBLIC_API_URL="http://localhost:$BACKEND_PORT"
export GO_BACKEND_URL="http://localhost:$BACKEND_PORT"

# Start the frontend development server
npm run dev -- -p $FRONTEND_PORT &
FRONTEND_PID=$!
cd ..

echo -e "${GREEN}Frontend server starting on port $FRONTEND_PORT (PID: $FRONTEND_PID)${NC}"

# Wait for servers to be ready
echo -e "\n${YELLOW}Waiting for servers to be ready...${NC}"
sleep 5

echo -e "\n${GREEN}✅ Skill Swap Application is running!${NC}"
echo -e "  🌐 Frontend: ${BLUE}http://localhost:$FRONTEND_PORT${NC}"
echo -e "  🔧 Backend API: ${BLUE}http://localhost:$BACKEND_PORT${NC}"
echo -e "  📚 API Documentation: ${BLUE}http://localhost:$BACKEND_PORT/api/v1${NC}"
echo -e "\n${YELLOW}Press Ctrl+C to stop both servers${NC}"

# Wait for both processes
wait $BACKEND_PID $FRONTEND_PID
