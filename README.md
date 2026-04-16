# GonsGarage - Auto Repair Shop Management System

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)]()
[![Next.js](https://img.shields.io/badge/Next.js-15.0+-black.svg)]()
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-blue.svg)]()
[![Swagger](https://img.shields.io/badge/Swagger-OpenAPI%203.0-brightgreen.svg)]()
[![License](https://img.shields.io/badge/license-MIT-green.svg)]()

A comprehensive auto repair shop management system built with modern technologies and clean architecture principles. GonsGarage provides a complete solution for managing clients, vehicles, repairs, appointments, and employees in an auto repair business.

## Documentation hub

Technical documentation (architecture analysis, dev setup, roadmap, and external spec placeholders such as Arnela) lives in **[docs/](docs/)**. Start at [docs/README.md](docs/README.md).

- **TDD / tests:** [docs/testing-tdd.md](docs/testing-tdd.md) · **MVP phases:** [docs/mvp-minimum-phases.md](docs/mvp-minimum-phases.md) · **API versionado:** [docs/api/versioning.md](docs/api/versioning.md) · **Observabilidad:** [docs/observability.md](docs/observability.md) · **Changelog:** [CHANGELOG.md](CHANGELOG.md) · **Contributing:** [CONTRIBUTING.md](CONTRIBUTING.md)

## 🚀 Features

### For Clients (Users with Client Role)
- ✅ **Vehicle Management**: Register and manage personal vehicles
- ✅ **Appointment Scheduling**: Book repair services online
- ✅ **Repair History**: View complete repair history for each vehicle
- ✅ **Real-time Updates**: Track repair progress and status
- ✅ **Profile Management**: Update personal information and preferences

### For Employees (Users with Employee Role)
- ✅ **Work Order Management**: Manage assigned repair tasks
- ✅ **Client Communication**: Update clients on repair progress
- ✅ **Inventory Tracking**: Track parts and materials usage
- ✅ **Time Tracking**: Log work hours for accurate billing
- ✅ **Employee Profile**: Manage position, hourly rate, and schedule

### For Administrators (Users with Admin Role)
- ✅ **User Management**: Manage all users (clients, employees, admins)
- ✅ **Employee Profiles**: Create and manage employee profiles
- ✅ **Business Analytics**: View performance metrics and reports
- ✅ **System Configuration**: Manage application settings and permissions
- ✅ **Financial Reporting**: Generate financial reports and insights

### For Developers
- ✅ **API Documentation**: Interactive Swagger/OpenAPI documentation
- ✅ **Type Safety**: Full TypeScript support with auto-generated types
- ✅ **Testing**: Comprehensive test suite with TDD approach
- ✅ **Development Tools**: Hot reload, linting, and code formatting

## 🏗️ Architecture

GonsGarage follows **Clean Architecture** principles with unified domain entities:

```
┌─────────────────────────────────────────────────────────────┐
│           Frontend (Next.js + TypeScript + Zustand)         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │    Pages    │  │ Components  │  │   Zustand Stores    │  │
│  │ (App Router)│  │   (TSX)     │  │  (State Mgmt)       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────┬───────────────────────────────────────┘
                      │ HTTP/REST API (camelCase JSON)
                      │ Swagger/OpenAPI Documentation
┌─────────────────────▼───────────────────────────────────────┐
│              Backend (Go + Gin Framework)                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Handlers  │  │   Services  │  │   Repositories      │  │
│  │ (Gin HTTP)  │  │ (Business)  │  │   (Data Access)     │  │
│  │ + Middleware│  │             │  │                     │  │
│  │ + Swagger   │  │             │  │                     │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                    Database Layer                           │
│  ┌─────────────┐              ┌─────────────────────────┐   │
│  │ PostgreSQL  │              │        Redis            │   │
│  │ (Primary)   │              │      (Cache)            │   │
│  └─────────────┘              └─────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin HTTP framework with native middleware support
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Authentication**: JWT tokens with Gin middleware
- **API Documentation**: Swagger/OpenAPI 3.0 with swaggo
- **Architecture**: Clean Architecture with unified domain entities
- **Logging**: Structured logging with slog
- **Testing**: Go testing + testify (TDD approach)

### Frontend
- **Framework**: Next.js 15 (App Router)
- **Language**: TypeScript 5+
- **State Management**: Zustand for global state management
- **Styling**: Tailwind CSS + CSS Modules
- **Testing**: Jest + React Testing Library
- **HTTP Client**: Fetch API with TypeScript wrappers
- **API Integration**: Auto-generated types from OpenAPI specs

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **Database**: PostgreSQL with persistent volumes
- **Cache**: Redis for session management and caching
- **Development**: Hot reload for both frontend and backend
- **Documentation**: Automated API docs with Swagger UI

## 📋 Prerequisites

Make sure you have the following installed on your system:

- **Node.js** 22+ ([Download](https://nodejs.org/))
- **pnpm** 9+ (`corepack enable` recommended)
- **Go** 1.21+ ([Download](https://golang.org/dl/))
- **Docker** and **Docker Compose** ([Download](https://www.docker.com/))
- **Git** ([Download](https://git-scm.com/))

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/gonsgarage.git
cd gonsgarage
```

### 2. Environment Setup

Create environment files from examples:

```bash
# Backend environment
cp backend/.env.example backend/.env

# Frontend environment
cp frontend/.env.local.example frontend/.env.local
```

### 3. Start with Docker (Recommended)

From the **repository root**, start PostgreSQL and Redis (credentials match `backend/.env.example` and the defaults in `backend/cmd/api/main.go`):

```bash
docker compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

Services: **`postgres`** and **`redis`**. Run the Go API and Next.js locally (see manual setup below) unless you add your own compose profiles for app containers.

### 4. Manual Setup (Alternative)

#### Start database services
```bash
# From repository root
docker compose up -d postgres redis
```

#### Backend Setup
```bash
cd backend

# Install dependencies
go mod tidy

# Generar Swagger (no hace falta instalar swag globalmente)
go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g main.go -o docs -d ./cmd/api,./internal/adapters/http/handlers,./internal/core/ports --parseInternal

# Schema: GORM AutoMigrate runs on startup (no separate migrate command required)

# Start the server
go run ./cmd/api
```

#### Frontend Setup
```bash
cd frontend

# Install dependencies (package manager: pnpm)
pnpm install

# Generate API types from Swagger (optional)
pnpm run generate-types

# Start development server
pnpm dev
```

## 🌐 Access the Application

After successful setup, access the application at:

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **📋 Swagger UI Documentation**: http://localhost:8080/swagger/index.html
- **📄 OpenAPI JSON Spec**: http://localhost:8080/swagger/doc.json
- **🔍 API Health Check**: http://localhost:8080/health

### Default Admin User
```
Email: admin@gonsgarage.com
Password: admin123
Role: admin
```

## 📚 API Documentation (Swagger/OpenAPI)

GonsGarage provides comprehensive API documentation through Swagger/OpenAPI 3.0:

### 🔍 **Interactive Documentation**
Visit [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) for:
- **Interactive API Explorer**: Test endpoints directly in the browser
- **Request/Response Examples**: See real data structures
- **Authentication Testing**: Test protected endpoints with JWT tokens
- **Schema Validation**: Validate request payloads before sending

### 📋 **API Specification**
- **OpenAPI 3.0 Compliant**: Industry-standard specification format
- **JSON Export**: Available at `/swagger/doc.json`
- **Code Generation**: Use specs to generate client SDKs
- **Postman Integration**: Import OpenAPI spec directly into Postman

### 🛠️ **Documentation Features**
```go
// Example of documented endpoint
// @Summary Create a new user
// @Description Create a new user in the system (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} domain.User "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Router /api/v1/users [post]
```

### 📊 **Available API Groups**
- **🔐 Authentication** (`/auth/*`) - Login, register, token management
- **👥 Users** (`/users/*`) - User management (unified clients/employees/admins)  
- **🚗 Cars** (`/cars/*`) - Vehicle management
- **👔 Employees** (`/employees/*`) - Employee profile management
- **🔧 Repairs** (`/repairs/*`) - Repair order management
- **📅 Appointments** (`/appointments/*`) - Appointment scheduling

## 📊 Unified Database Schema

### Core Entities (Simplified Architecture)

```sql
-- Unified Users table (replaces separate Client table)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'client',  -- 'client', 'employee', 'admin'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Employee profiles (only for users with role='employee')
CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    position VARCHAR(100) NOT NULL,
    hourly_rate DECIMAL(10,2),
    hire_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Cars (owned by users with role='client')
CREATE TABLE cars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER NOT NULL,
    license_plate VARCHAR(20) UNIQUE NOT NULL,
    vin VARCHAR(17),
    color VARCHAR(50) NOT NULL,
    mileage INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Repair orders
CREATE TABLE repairs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    car_id UUID REFERENCES cars(id) ON DELETE CASCADE,
    technician_id UUID REFERENCES employees(id),
    description TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    estimated_cost DECIMAL(10,2),
    actual_cost DECIMAL(10,2),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Appointment scheduling
CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID REFERENCES users(id) ON DELETE CASCADE,
    car_id UUID REFERENCES cars(id) ON DELETE CASCADE,
    service_type VARCHAR(100) NOT NULL,
    scheduled_at TIMESTAMP NOT NULL,
    status VARCHAR(50) DEFAULT 'scheduled',
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
```

## 🧪 Testing

### Backend Testing (TDD Approach)

```bash
cd backend

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific service tests
go test ./internal/core/services/...

# Run tests with verbose output
go test -v ./...

# Generate and view coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Frontend Testing (TypeScript + Jest)

```bash
cd frontend

# Run all tests
pnpm test

# Run tests in watch mode
pnpm test:watch

# Run tests with coverage
pnpm test:coverage

# Type checking
pnpm typecheck

# Generate API types from Swagger
pnpm run generate-types
```

### Test Structure Examples

```go
// Backend service test example
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewMockUserRepository(t)
    logger := slog.Default()
    service := NewUserService(mockRepo, logger)
    
    user := &domain.User{
        FirstName: "John",
        LastName:  "Doe",
        Email:     "john@example.com", 
        Role:      domain.RoleClient,
    }
    
    mockRepo.EXPECT().CreateUser(mock.Anything, user).Return(user, nil)
    
    // Act
    result, err := service.CreateUser(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, user.Email, result.Email)
    assert.Equal(t, domain.RoleClient, result.Role)
    mockRepo.AssertExpectations(t)
}
```

```typescript
// Frontend Zustand store test example
describe('UserStore', () => {
  beforeEach(() => {
    useUserStore.getState().users = [];
    useUserStore.getState().error = null;
  });

  it('should create user successfully', async () => {
    // Arrange
    const mockUser: User = {
      id: '1', 
      firstName: 'John', 
      lastName: 'Doe',
      email: 'john@example.com',
      role: 'client'
    };
    
    jest.spyOn(userApi, 'createUser').mockResolvedValue(mockUser);
    
    // Act
    await useUserStore.getState().createUser({
      firstName: 'John',
      lastName: 'Doe', 
      email: 'john@example.com',
      role: 'client'
    });
    
    // Assert
    const state = useUserStore.getState();
    expect(state.users).toContain(mockUser);
    expect(state.loading).toBe(false);
    expect(state.error).toBeNull();
  });
});
```

## 🔧 API Documentation Examples

### Authentication Endpoints

```http
POST /api/v1/auth/register    # User registration (roles: admin, manager, employee, client; default client if role omitted)
POST /api/v1/auth/login       # User authentication
GET /api/v1/auth/me           # Current user (requires Bearer token; protected route)
POST /api/v1/auth/refresh     # Token refresh (if implemented in client)
POST /api/v1/auth/logout      # User logout (if implemented)
```

### User Management Endpoints (Unified)

```http
GET    /api/v1/users          # List users (admin only)
POST   /api/v1/users          # Create user (admin only)
GET    /api/v1/users/:id      # Get user details
PUT    /api/v1/users/:id      # Update user
DELETE /api/v1/users/:id      # Delete user (admin only)
```

### Car Management Endpoints

```http
GET    /api/v1/cars           # List user's cars
POST   /api/v1/cars           # Create new car
GET    /api/v1/cars/:id       # Get car details
PUT    /api/v1/cars/:id       # Update car
DELETE /api/v1/cars/:id       # Delete car
```

### Employee Management Endpoints

```http
GET    /api/v1/employees      # List employees (admin only)
POST   /api/v1/employees      # Create employee profile (admin only)
GET    /api/v1/employees/:id  # Get employee details
PUT    /api/v1/employees/:id  # Update employee
DELETE /api/v1/employees/:id  # Delete employee (admin only)
```

### Example API Request/Response (camelCase JSON)

```bash
# Create a new user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Doe", 
    "email": "john@example.com",
    "password": "securepassword",
    "role": "client"
  }'
```

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "firstName": "John",
  "lastName": "Doe",
  "email": "john@example.com", 
  "role": "client",
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

## 🎯 Frontend State Management (Zustand)

### User Store Example

```typescript
// stores/user.store.ts
interface UserStore {
  // State
  users: User[];
  currentUser: User | null;
  loading: boolean;
  error: string | null;
  
  // Actions  
  fetchUsers: () => Promise<void>;
  createUser: (userData: CreateUserRequest) => Promise<void>;
  updateUser: (id: string, userData: UpdateUserRequest) => Promise<void>;
  deleteUser: (id: string) => Promise<void>;
  setCurrentUser: (user: User | null) => void;
  clearError: () => void;
}

export const useUserStore = create<UserStore>((set, get) => ({
  // Initial state
  users: [],
  currentUser: null,
  loading: false,
  error: null,

  // Actions implementation
  fetchUsers: async () => {
    set({ loading: true, error: null });
    try {
      const users = await userApi.getUsers();
      set({ users, loading: false });
    } catch (error) {
      set({ error: (error as Error).message, loading: false });
    }
  },

  createUser: async (userData) => {
    set({ loading: true, error: null });
    try {
      const newUser = await userApi.createUser(userData);
      set(state => ({ 
        users: [...state.users, newUser], 
        loading: false 
      }));
    } catch (error) {
      set({ error: (error as Error).message, loading: false });
    }
  },

  // ... other actions
}));
```

### Component Usage

```typescript
// components/UserList.tsx
import { useUserStore } from '@/stores/user.store';

export default function UserList() {
  const { 
    users, 
    loading, 
    error, 
    fetchUsers, 
    clearError 
  } = useUserStore();

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  if (loading) return <div>Loading users...</div>;
  if (error) {
    return (
      <div className="error">
        <p>Error: {error}</p>
        <button onClick={clearError}>Dismiss</button>
      </div>
    );
  }

  return (
    <div className="user-list">
      <h2>Users ({users.length})</h2>
      {users.map(user => (
        <div key={user.id} className="user-card">
          <h3>{user.firstName} {user.lastName}</h3>
          <p>Email: {user.email}</p>
          <p>Role: {user.role}</p>
        </div>
      ))}
    </div>
  );
}
```

## 🔧 Gin Middleware Configuration

### Authentication Middleware

```go
// middleware/auth_middleware.go
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "missing authorization header"
            })
            c.Abort()
            return
        }

        // JWT validation...
        userID, err := m.validateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "invalid token"
            })
            c.Abort()
            return
        }

        // Set user context with standard key
        c.Set("userID", userID)
        c.Next()
    }
}
```

### CORS Middleware

```go
// middleware/cors_middleware.go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

## 📋 Swagger Documentation Setup

### 1. Install Swagger Dependencies

```bash
# Backend - Install swaggo
cd backend
go install github.com/swaggo/swag/cmd/swag@latest
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/files
```

### 2. Generate Documentation

```bash
cd backend
go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g main.go -o docs -d ./cmd/api,./internal/adapters/http/handlers,./internal/core/ports --parseInternal
```

### 3. Integration Example

```go
// cmd/api/main.go
package main

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "github.com/your-org/gonsgarage/docs" // Import generated docs
)

// @title GonsGarage API
// @version 1.0
// @description Auto repair shop management system API
// @contact.name API Support
// @contact.email support@gonsgarage.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
    r := gin.Default()
    
    // Swagger endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // Your API routes...
    setupRoutes(r)
    
    r.Run(":8080")
}
```

## 🚀 Deployment

### Environment Variables

```bash
# Backend (.env)
DATABASE_URL=postgres://admindb:secure_password@localhost:5432/gonsgarage?sslmode=require
# REDIS_URL is host:port (e.g. localhost:6379), not a redis:// URL
REDIS_URL=localhost:6379
JWT_SECRET=your-very-secure-jwt-secret-key-here
GIN_MODE=release
SERVER_PORT=8080
SWAGGER_ENABLED=true

# Frontend (.env.local)
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
NEXT_PUBLIC_SWAGGER_URL=https://api.yourdomain.com/swagger
NODE_ENV=production
```

### Docker Compose Production

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: gonsgarage
      POSTGRES_USER: admindb
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

  backend:
    build: 
      context: ./backend
      dockerfile: Dockerfile.prod
    environment:
      DATABASE_URL: postgres://admindb:${DB_PASSWORD}@postgres:5432/gonsgarage?sslmode=disable
      REDIS_URL: redis:6379
      JWT_SECRET: ${JWT_SECRET}
      GIN_MODE: release
      SWAGGER_ENABLED: true
    depends_on:
      - postgres
      - redis

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.prod
    environment:
      NEXT_PUBLIC_API_URL: http://backend:8080
      NEXT_PUBLIC_SWAGGER_URL: http://backend:8080/swagger
      NODE_ENV: production
    depends_on:
      - backend

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - frontend
      - backend

volumes:
  postgres_data:
  redis_data:
```

## 🤝 Contributing

We welcome contributions! See **[CONTRIBUTING.md](CONTRIBUTING.md)** for setup (pnpm, Docker) and the **TDD** policy. Summary:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feat/amazing-feature`
3. **Follow TDD**: Write tests first, then implementation ([docs/testing-tdd.md](docs/testing-tdd.md))
4. **Follow coding standards**: See [Agent.md](Agent.md) for detailed guidelines
5. **Use TypeScript**: Ensure type safety in frontend code
6. **Document APIs**: Add Swagger annotations for all new endpoints
7. **Test Zustand stores**: Write tests for state management logic
8. **Update documentation**: Keep Swagger docs up to date
9. **Commit with conventional commits**: `feat(auth): add user registration`
10. **Push to branch**: `git push origin feat/amazing-feature`
11. **Create Pull Request**

### Code Review Checklist

- [ ] Tests pass locally (both Go and TypeScript)
- [ ] Code follows naming conventions (PascalCase/camelCase)
- [ ] Business logic is in the correct layer (Clean Architecture)
- [ ] Error handling is implemented properly
- [ ] API endpoints follow camelCase JSON convention
- [ ] **Swagger documentation is updated for new/changed endpoints**
- [ ] TypeScript types are properly defined
- [ ] Zustand stores are tested
- [ ] Gin middleware is properly configured
- [ ] **API examples work in Swagger UI**

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

- **Documentation**: [Wiki](https://github.com/your-username/gonsgarage/wiki)
- **API Documentation**: [Swagger UI](http://localhost:8080/swagger/index.html)
- **Issues**: [GitHub Issues](https://github.com/your-username/gonsgarage/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/gonsgarage/discussions)
- **Email**: support@gonsgarage.com

## 🙏 Acknowledgments

- Clean Architecture principles by Robert C. Martin
- Go community for excellent libraries and tools  
- Next.js team for the fantastic React framework
- Gin framework for excellent Go HTTP middleware support
- Zustand community for simple and effective state management
- Swagger/OpenAPI community for API documentation standards
- All contributors who help improve this project

---

**Built with ❤️ for the auto repair industry using Go + Gin + Next.js + TypeScript + Zustand + Swagger**