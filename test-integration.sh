#!/bin/bash

# Integration Test Script for Skill Swap
# Tests the communication between frontend and backend

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BACKEND_PORT=${BACKEND_PORT:-8080}
FRONTEND_PORT=${FRONTEND_PORT:-3000}
BACKEND_URL="http://localhost:$BACKEND_PORT"
FRONTEND_URL="http://localhost:$FRONTEND_PORT"

echo -e "${BLUE}🧪 Starting Integration Tests...${NC}"

# Function to check if a service is running
check_service() {
    local url=$1
    local name=$2
    
    if curl -s -f "$url" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ $name is running${NC}"
        return 0
    else
        echo -e "${RED}❌ $name is not running${NC}"
        return 1
    fi
}

# Function to test API endpoint
test_api_endpoint() {
    local endpoint=$1
    local expected_status=$2
    local description=$3
    
    echo -e "${YELLOW}Testing: $description${NC}"
    
    local response=$(curl -s -w "%{http_code}" -o /dev/null "$BACKEND_URL$endpoint")
    
    if [ "$response" = "$expected_status" ]; then
        echo -e "${GREEN}✅ $description - Status: $response${NC}"
        return 0
    else
        echo -e "${RED}❌ $description - Expected: $expected_status, Got: $response${NC}"
        return 1
    fi
}

# Test 1: Check if backend is running
echo -e "\n${BLUE}Test 1: Backend Health Check${NC}"
if check_service "$BACKEND_URL/health" "Backend"; then
    # Get health info
    echo -e "${YELLOW}Backend Health Info:${NC}"
    curl -s "$BACKEND_URL/health" | jq . 2>/dev/null || curl -s "$BACKEND_URL/health"
else
    echo -e "${RED}Backend is not running. Please start the backend server first.${NC}"
    exit 1
fi

# Test 2: Check if frontend is running
echo -e "\n${BLUE}Test 2: Frontend Health Check${NC}"
if check_service "$FRONTEND_URL" "Frontend"; then
    echo -e "${GREEN}Frontend is accessible${NC}"
else
    echo -e "${YELLOW}Frontend may not be running or is starting up${NC}"
fi

# Test 3: Test API endpoints
echo -e "\n${BLUE}Test 3: API Endpoint Tests${NC}"
test_api_endpoint "/health" "200" "Health endpoint"
test_api_endpoint "/ready" "200" "Ready endpoint"
test_api_endpoint "/live" "200" "Live endpoint"

# Test 4: Test CORS headers
echo -e "\n${BLUE}Test 4: CORS Configuration Test${NC}"
echo -e "${YELLOW}Testing CORS headers from frontend origin...${NC}"

cors_response=$(curl -s -H "Origin: $FRONTEND_URL" -H "Access-Control-Request-Method: GET" -H "Access-Control-Request-Headers: Content-Type" -X OPTIONS "$BACKEND_URL/health" -D -)

if echo "$cors_response" | grep -q "Access-Control-Allow-Origin"; then
    echo -e "${GREEN}✅ CORS headers are present${NC}"
    echo -e "${YELLOW}CORS Headers:${NC}"
    echo "$cors_response" | grep -i "access-control" | head -5
else
    echo -e "${RED}❌ CORS headers are missing${NC}"
fi

# Test 5: Test API versioning
echo -e "\n${BLUE}Test 5: API Versioning Test${NC}"
test_api_endpoint "/api/v1" "404" "API v1 root (expected 404)"

# Summary
echo -e "\n${BLUE}🎯 Integration Test Summary${NC}"
echo -e "Backend URL: ${GREEN}$BACKEND_URL${NC}"
echo -e "Frontend URL: ${GREEN}$FRONTEND_URL${NC}"
echo -e "\n${GREEN}✅ Integration tests completed!${NC}"
echo -e "${YELLOW}Note: Some tests may fail if the services are still starting up.${NC}"
