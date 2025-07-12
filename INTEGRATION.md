# Skill Swap - Frontend & Backend Integration

This document describes how the frontend and backend are integrated and how to run the application with configurable ports.

## 🏗️ Architecture Overview

- **Backend**: Go (Gin framework) - API server
- **Frontend**: Next.js (React) - Web application
- **Communication**: RESTful API with CORS enabled
- **Ports**: Fully configurable via environment variables

## 🚀 Quick Start

### Option 1: Using the Startup Script (Recommended)

```bash
# Clone and navigate to the repository
git clone <your-repo-url>
cd Skill-swap

# Run with default ports (Backend: 8080, Frontend: 3000)
./start.sh

# Run with custom ports
BACKEND_PORT=8081 FRONTEND_PORT=3001 ./start.sh
```

### Option 2: Manual Setup

1. **Start Backend**:
```bash
cd backend/skillswap
export PORT=8080
export FRONTEND_URL=http://localhost:3000
go run cmd/server/main.go
```

2. **Start Frontend**:
```bash
cd frontend
export NEXT_PUBLIC_API_URL=http://localhost:8080
npm install
npm run dev
```

## ⚙️ Configuration

### Environment Variables

#### Backend (.env in backend/skillswap/)
```bash
# Server Configuration
PORT=8080
BASE_URL=http://localhost:8080
FRONTEND_URL=http://localhost:3000

# Database
DB_URL=postgresql://username:password@localhost:5432/skillswap?sslmode=disable

# Security
JWT_SECRET=your-super-secret-jwt-key

# File Upload
UPLOAD_DIR=./uploads
```

#### Frontend (.env.local in frontend/)
```bash
# Frontend Configuration
FRONTEND_PORT=3000
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-secret-key

# Backend Configuration
BACKEND_PORT=8080
GO_BACKEND_URL=http://localhost:8080
NEXT_PUBLIC_API_URL=http://localhost:8080

# OAuth (optional)
GITHUB_ID=your-github-id
GITHUB_SECRET=your-github-secret
```

## 🔗 API Integration

### API Base URL Configuration

The frontend automatically detects the backend URL using:

1. `NEXT_PUBLIC_API_URL` for client-side requests
2. `GO_BACKEND_URL` for server-side requests
3. Fallback to `http://localhost:8080`

### API Utilities

The frontend includes a centralized API utility (`src/lib/api.ts`) with:

- `apiGet()` - GET requests
- `apiPost()` - POST requests
- `apiPut()` - PUT requests
- `apiDelete()` - DELETE requests
- Automatic CORS and authentication handling

### Example Usage

```typescript
import { apiGet, apiPost } from '@/lib/api';

// GET request
const users = await apiGet('/api/v1/users');

// POST request
const newUser = await apiPost('/api/v1/users', { name: 'John', email: 'john@example.com' });
```

## 🛡️ Security & CORS

### CORS Configuration

The backend automatically configures CORS to allow:
- The configured frontend URL
- `localhost:3000` and `localhost:3001` for development
- Credentials and common headers

### Authentication

- JWT-based authentication
- NextAuth.js integration
- Automatic token handling in API requests

## 🧪 Testing

### Run Integration Tests

```bash
# Test the integration between frontend and backend
./test-integration.sh

# Test with custom ports
BACKEND_PORT=8081 FRONTEND_PORT=3001 ./test-integration.sh
```

### Health Checks

- **Backend Health**: `GET /health`
- **Backend Ready**: `GET /ready` 
- **Backend Live**: `GET /live`

Example health response:
```json
{
  "status": "healthy",
  "service": "skillswap-api",
  "version": "1.0.0",
  "backend_url": "http://localhost:8080",
  "frontend_url": "http://localhost:3000"
}
```

## 📝 Development Notes

### Port Configuration

Ports are configurable but must be consistent:

1. Set `BACKEND_PORT` and `FRONTEND_PORT` environment variables
2. Update `NEXT_PUBLIC_API_URL` to match backend port
3. Update `FRONTEND_URL` in backend config to match frontend port

### Development Workflow

1. Make changes to backend or frontend
2. Both servers will auto-reload on file changes
3. API changes are immediately available to frontend
4. CORS is configured to allow local development

### Production Deployment

For production, ensure:

1. Set production URLs in environment variables
2. Configure CORS for production frontend domain
3. Use HTTPS in production
4. Set secure JWT secrets

## 🔧 Troubleshooting

### Common Issues

1. **CORS Errors**: Check that `FRONTEND_URL` is set correctly in backend config
2. **API Not Found**: Verify `NEXT_PUBLIC_API_URL` matches backend port
3. **Port Conflicts**: Use different ports if default ones are in use

### Logs

- Backend logs show startup port and configuration
- Frontend logs show API base URL being used
- Both services log CORS and request information

## 📚 API Documentation

The backend provides comprehensive API documentation. Key endpoints:

- `GET /health` - Health check
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/users` - List users
- `POST /api/v1/skills` - Create skill

For complete API documentation, see [backend/backend_api_endpoints.md](backend/backend_api_endpoints.md).
