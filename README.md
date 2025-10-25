# GonsGarage - Auto Repair Shop Management System

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)]()
[![Next.js](https://img.shields.io/badge/Next.js-15.0+-black.svg)]()
[![License](https://img.shields.io/badge/license-MIT-green.svg)]()

A comprehensive auto repair shop management system built with modern technologies and clean architecture principles. GonsGarage provides a complete solution for managing clients, vehicles, repairs, appointments, and employees in an auto repair business.

## 🚀 Features

### For Clients
- ✅ **Vehicle Management**: Register and manage personal vehicles
- ✅ **Appointment Scheduling**: Book repair services online
- ✅ **Repair History**: View complete repair history for each vehicle
- ✅ **Real-time Updates**: Track repair progress and status
- ✅ **Invoicing**: View and download repair invoices

### For Employees
- ✅ **Work Order Management**: Manage assigned repair tasks
- ✅ **Client Communication**: Update clients on repair progress
- ✅ **Inventory Tracking**: Track parts and materials usage
- ✅ **Time Tracking**: Log work hours for accurate billing

### For Administrators
- ✅ **Employee Management**: Manage staff accounts and permissions
- ✅ **Business Analytics**: View performance metrics and reports
- ✅ **Inventory Management**: Manage parts inventory and suppliers
- ✅ **Financial Reporting**: Generate financial reports and insights

## 🏗️ Architecture

GonsGarage follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend (Next.js)                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │    Pages    │  │ Components  │  │    API Client       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────┬───────────────────────────────────────┘
                      │ HTTP/REST API
┌─────────────────────▼───────────────────────────────────────┐
│                     Backend (Go)                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Handlers   │  │   Services  │  │    Repositories     │  │
│  │  (HTTP)     │  │ (Business)  │  │   (Data Access)     │  │
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
- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Authentication**: JWT tokens
- **Migration**: Custom SQL migrations
- **Testing**: Go testing + testify
- **Documentation**: Swagger/OpenAPI

### Frontend
- **Framework**: Next.js 15 (App Router)
- **Language**: TypeScript 5+
- **Styling**: Tailwind CSS + CSS Modules
- **State Management**: React Context + Custom Hooks
- **Testing**: Jest + React Testing Library
- **HTTP Client**: Fetch API with custom wrapper

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **Database**: PostgreSQL with persistent volumes
- **Cache**: Redis for session management and caching
- **Development**: Hot reload for both frontend and backend

## 📋 Prerequisites

Make sure you have the following installed on your system:

- **Node.js** 18+ ([Download](https://nodejs.org/))
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

```bash
# Start all services (database, redis, backend, frontend)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### 4. Manual Setup (Alternative)

#### Start Database Services
```bash
# Start PostgreSQL and Redis
docker-compose up -d postgres redis
```

#### Backend Setup
```bash
cd backend

# Install dependencies
go mod tidy

# Run database migrations
go run cmd/migrate/main.go up

# Start the server
go run cmd/server/main.go
```

#### Frontend Setup
```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

## 🌐 Access the Application

After successful setup, access the application at:

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/index.html

### Default Admin Credentials
```
Email: admin@gonsgarage.com
Password: admin123
```

## 📊 Database Schema

### Core Entities

```sql
-- Users (base for all user types)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'client',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Employee profiles
CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    position VARCHAR(100) NOT NULL,
    hourly_rate DECIMAL(10,2),
    hire_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Client vehicles
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
    updated_at TIMESTAMP DEFAULT NOW()
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
    updated_at TIMESTAMP DEFAULT NOW()
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
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## 🧪 Testing

### Backend Testing

```bash
cd backend

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./internal/core/services/...

# Run tests with verbose output
go test -v ./...
```

### Frontend Testing

```bash
cd frontend

# Run all tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage

# Run E2E tests (if configured)
npm run test:e2e
```

### Test Structure Example

```go
// Backend test example
func TestCarService_CreateCar(t *testing.T) {
    // Arrange
    mockRepo := &MockCarRepository{}
    service := NewCarService(mockRepo, nil)
    
    req := CreateCarRequest{
        Make:         "Toyota",
        Model:        "Camry",
        Year:         2023,
        LicensePlate: "ABC-123",
        Color:        "Blue",
    }
    
    // Act
    car, err := service.CreateCar(context.Background(), req, userID)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, req.Make, car.Make)
    mockRepo.AssertExpectations(t)
}
```

```typescript
// Frontend test example
describe('CarForm', () => {
  it('should create car when form is submitted', async () => {
    // Arrange
    const mockCreateCar = jest.fn().mockResolvedValue(true);
    render(<CarForm onCreate={mockCreateCar} />);
    
    // Act
    await user.type(screen.getByLabelText('Make'), 'Toyota');
    await user.type(screen.getByLabelText('Model'), 'Camry');
    await user.click(screen.getByText('Create Car'));
    
    // Assert
    await waitFor(() => {
      expect(mockCreateCar).toHaveBeenCalledWith({
        make: 'Toyota',
        model: 'Camry'
      });
    });
  });
});
```

## 🔧 API Documentation

### Authentication Endpoints

```http
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
```

### Car Management Endpoints

```http
GET    /api/v1/cars              # List user's cars
POST   /api/v1/cars              # Create new car
GET    /api/v1/cars/:id          # Get car details
PUT    /api/v1/cars/:id          # Update car
DELETE /api/v1/cars/:id          # Delete car
```

### Example API Request/Response

```bash
# Create a new car
curl -X POST http://localhost:8080/api/v1/cars \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "make": "Toyota",
    "model": "Camry",
    "year": 2023,
    "licensePlate": "ABC-123",
    "color": "Blue",
    "mileage": 15000
  }'
```

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "make": "Toyota",
  "model": "Camry",
  "year": 2023,
  "licensePlate": "ABC-123",
  "color": "Blue",
  "mileage": 15000,
  "ownerId": "user-uuid",
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

## 🐳 Docker Configuration

### Development Environment

```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: gonsgarage
      POSTGRES_USER: admindb
      POSTGRES_PASSWORD: gonsgarage123
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  backend:
    build: 
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://admindb:gonsgarage123@postgres:5432/gonsgarage?sslmode=disable
      REDIS_URL: redis://redis:6379
      JWT_SECRET: your-super-secret-jwt-key
    depends_on:
      - postgres
      - redis
    volumes:
      - ./backend:/app
    command: air # Hot reload in development

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    environment:
      NEXT_PUBLIC_API_URL: http://localhost:8080
    volumes:
      - ./frontend:/app
      - /app/node_modules
    command: npm run dev

volumes:
  postgres_data:
  redis_data:
```

### Production Dockerfile Examples

```dockerfile
# Backend Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

```dockerfile
# Frontend Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

EXPOSE 3000
CMD ["npm", "start"]
```

## 🚀 Deployment

### Production Environment Setup

1. **Server Requirements**:
   - 2+ CPU cores
   - 4GB+ RAM
   - 50GB+ storage
   - Ubuntu 20.04+ or similar

2. **Install Dependencies**:
```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

3. **Deploy Application**:
```bash
# Clone repository
git clone https://github.com/your-username/gonsgarage.git
cd gonsgarage

# Set production environment variables
cp .env.production.example .env.production
# Edit .env.production with your production values

# Start production services
docker-compose -f docker-compose.prod.yml up -d

# Run database migrations
docker-compose exec backend go run cmd/migrate/main.go up

# Set up SSL with Let's Encrypt (optional)
sudo apt install certbot
sudo certbot --nginx -d yourdomain.com
```

### Environment Variables for Production

```bash
# Backend (.env)
DATABASE_URL=postgres://admindb:secure_password@localhost:5432/gonsgarage?sslmode=require
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-very-secure-jwt-secret-key-here
GIN_MODE=release
PORT=8080

# Frontend (.env.local)
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
NEXTAUTH_SECRET=your-nextauth-secret
NEXTAUTH_URL=https://yourdomain.com
```

## 🔧 Development Workflow

### Setting up Development Environment

1. **Clone and Setup**:
```bash
git clone https://github.com/your-username/gonsgarage.git
cd gonsgarage
make setup # Runs all setup commands
```

2. **Database Operations**:
```bash
# Create new migration
make migration name=add_new_table

# Run migrations
make migrate-up

# Rollback migration
make migrate-down

# Reset database
make db-reset
```

3. **Development Commands**:
```bash
# Start development environment
make dev

# Run tests
make test

# Run linting
make lint

# Build for production
make build
```

### Code Quality Tools

```bash
# Backend linting
golangci-lint run

# Frontend linting
npm run lint
npm run type-check

# Format code
gofmt -w .
npm run format
```

## 📈 Monitoring and Logging

### Application Metrics

The application includes built-in monitoring endpoints:

- **Health Check**: `GET /health`
- **Metrics**: `GET /metrics` (Prometheus format)
- **Debug Info**: `GET /debug/pprof/` (development only)

### Logging Configuration

```go
// Structured logging example
logger.Info("car created successfully",
    "user_id", userID,
    "car_id", car.ID,
    "make", car.Make,
    "model", car.Model)

logger.Error("database operation failed",
    "operation", "CreateCar",
    "error", err,
    "user_id", userID)
```

### Production Monitoring Stack

For production, consider setting up:
- **Prometheus** for metrics collection
- **Grafana** for visualization
- **Jaeger** for distributed tracing
- **ELK Stack** for log aggregation

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feat/amazing-feature`
3. **Follow TDD**: Write tests first, then implementation
4. **Follow coding standards**: See [Agent.md](Agent.md) for detailed guidelines
5. **Commit with conventional commits**: `feat(auth): add user registration`
6. **Push to branch**: `git push origin feat/amazing-feature`
7. **Create Pull Request**

### Code Review Checklist

- [ ] Tests pass locally
- [ ] Code follows naming conventions (camelCase/PascalCase)
- [ ] Business logic is in the correct layer (Clean Architecture)
- [ ] Error handling is implemented properly
- [ ] API endpoints are documented
- [ ] Database migrations are provided (if needed)

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

- **Documentation**: [Wiki](https://github.com/your-username/gonsgarage/wiki)
- **Issues**: [GitHub Issues](https://github.com/your-username/gonsgarage/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/gonsgarage/discussions)
- **Email**: support@gonsgarage.com

## 🙏 Acknowledgments

- Clean Architecture principles by Robert C. Martin
- Go community for excellent libraries and tools
- Next.js team for the fantastic React framework
- All contributors who help improve this project

---

**Built with ❤️ for the auto repair industry**