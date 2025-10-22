# GonsGarage Development Agent

This agent file defines the coding standards, conventions, and guidelines for the GonsGarage project. All contributors and AI assistants must follow these rules when generating, refactoring, or testing code.

## 🧩 Project Overview

**GonsGarage** is a comprehensive auto repair shop management system built with modern technologies and clean architecture principles.

### Tech Stack
- **Backend**: Go (Golang) with Clean Architecture
- **Frontend**: Next.js 15 (React + TypeScript)
- **Database**: PostgreSQL (primary), Redis (caching/sessions)
- **Development**: Docker-based local environment
- **Testing**: Test-Driven Development (TDD)

---

## 1️⃣ Coding Standards

### Naming Conventions
```go
// Go - Correct
type UserService struct {}        // PascalCase for structs
func (s *UserService) CreateUser() // PascalCase for exported methods
var userName string               // camelCase for variables
const maxRetries = 3             // camelCase for constants

// Go - Incorrect
type userService struct {}        // ❌ should be PascalCase
func (s *UserService) create_user() // ❌ should be camelCase
var UserName string               // ❌ unexported should be camelCase
```

```typescript
// TypeScript - Correct
interface UserProps {             // PascalCase for interfaces
  firstName: string;              // camelCase for properties
}

function UserComponent({ firstName }: UserProps) { // PascalCase for components
  const [isLoading, setIsLoading] = useState(false); // camelCase for variables
  
  const handleSubmit = () => {};  // camelCase for functions
}

// TypeScript - Incorrect
interface userProps {             // ❌ should be PascalCase
  first_name: string;             // ❌ should be camelCase
}
```

### Code Quality Rules
- ✅ Keep functions small (max 20-30 lines)
- ✅ Single responsibility principle
- ✅ Explicit return types
- ✅ No magic numbers/strings - use constants
- ✅ Consistent error handling
- ✅ Meaningful variable names

```go
// Good
const (
    defaultTimeout = 30 * time.Second
    maxRetryAttempts = 3
)

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    // implementation
}

// Bad
func CreateUser(req interface{}) interface{} { // ❌ no types, magic interface{}
    time.Sleep(30000000000) // ❌ magic number
    // implementation
}
```

---

## 2️⃣ Testing and TDD Rules

### Testing Frameworks
- **Go**: `testing` + `testify/assert` + `testify/mock`
- **Next.js**: `Jest` + `React Testing Library`

### TDD Workflow
1. **Red**: Write a failing test
2. **Green**: Write minimal code to pass
3. **Refactor**: Improve code while keeping tests green

### Test Structure (AAA Pattern)
```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    req := CreateUserRequest{
        Email: "test@example.com",
        FirstName: "John",
        LastName: "Doe",
    }
    
    // Act
    user, err := service.CreateUser(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, req.Email, user.Email)
    mockRepo.AssertExpectations(t)
}
```

```typescript
describe('UserComponent', () => {
  it('should display user name when provided', () => {
    // Arrange
    const props: UserProps = {
      firstName: 'John',
      lastName: 'Doe'
    };
    
    // Act
    render(<UserComponent {...props} />);
    
    // Assert
    expect(screen.getByText('John Doe')).toBeInTheDocument();
  });
});
```

### Test File Naming
- Go: `filename_test.go`
- TypeScript: `filename.test.tsx` or `filename.test.ts`

---

## 3️⃣ Backend Conventions (Go)

### Project Structure
```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── core/
│   │   ├── domain/           # Entities, value objects
│   │   ├── ports/            # Interfaces/contracts
│   │   └── services/         # Business logic
│   ├── adapters/
│   │   ├── http/            # HTTP handlers
│   │   ├── repository/      # Database implementations
│   │   └── thirdparty/      # External service clients
│   └── pkg/                 # Shared utilities
├── migrations/              # Database migrations
├── docker/
└── .env.example
```

### Clean Architecture Layers

#### Domain Layer
```go
// internal/core/domain/user.go
type User struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
    if u.Email == "" {
        return errors.New("email is required")
    }
    return nil
}
```

#### Ports Layer
```go
// internal/core/ports/repositories.go
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    GetByID(ctx context.Context, id string) (*domain.User, error)
    GetByEmail(ctx context.Context, email string) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id string) error
}

// internal/core/ports/services.go
type UserService interface {
    CreateUser(ctx context.Context, req CreateUserRequest) (*domain.User, error)
    GetUser(ctx context.Context, id string) (*domain.User, error)
    AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error)
}
```

#### Services Layer
```go
// internal/core/services/user_service.go
type userService struct {
    userRepo ports.UserRepository
    logger   *slog.Logger
}

func NewUserService(userRepo ports.UserRepository, logger *slog.Logger) ports.UserService {
    return &userService{
        userRepo: userRepo,
        logger:   logger,
    }
}

func (s *userService) CreateUser(ctx context.Context, req CreateUserRequest) (*domain.User, error) {
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    user := &domain.User{
        ID:        generateID(),
        Email:     req.Email,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Role:      req.Role,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        s.logger.Error("failed to create user", "error", err, "email", req.Email)
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}
```

### Go Best Practices
- ✅ Always pass `context.Context` for I/O operations
- ✅ Use dependency injection
- ✅ Structured logging (slog)
- ✅ Explicit error wrapping
- ✅ Interface segregation

```go
// Good
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid JSON")
        return
    }
    
    user, err := h.userService.CreateUser(ctx, req)
    if err != nil {
        h.logger.Error("failed to create user", "error", err)
        h.respondError(w, http.StatusInternalServerError, "failed to create user")
        return
    }
    
    h.respondJSON(w, http.StatusCreated, user)
}
```

---

## 4️⃣ Frontend Conventions (Next.js 15)

### Project Structure
```
frontend/
├── src/
│   ├── app/                 # App Router pages
│   │   ├── (auth)/
│   │   ├── dashboard/
│   │   └── layout.tsx
│   ├── components/          # Reusable components
│   │   ├── ui/             # Basic UI components
│   │   └── features/       # Feature-specific components
│   ├── contexts/           # React contexts
│   ├── hooks/              # Custom hooks
│   ├── lib/                # Utilities and API clients
│   ├── types/              # TypeScript type definitions
│   └── styles/             # Global styles
├── public/
└── __tests__/              # Test setup files
```

### Component Conventions
```typescript
// components/ui/Button/Button.tsx
interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  children: React.ReactNode;
  onClick?: () => void;
}

export default function Button({
  variant = 'primary',
  size = 'md',
  disabled = false,
  children,
  onClick
}: ButtonProps) {
  const className = `btn btn--${variant} btn--${size}`;
  
  return (
    <button
      className={className}
      disabled={disabled}
      onClick={onClick}
      type="button"
    >
      {children}
    </button>
  );
}
```

### API Client Pattern
```typescript
// lib/api/client.ts
class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  setToken(token: string) {
    this.token = token;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<{ data: T | null; error: string | null }> {
    try {
      const url = `${this.baseURL}${endpoint}`;
      const headers = {
        'Content-Type': 'application/json',
        ...(this.token && { Authorization: `Bearer ${this.token}` }),
        ...options.headers,
      };

      const response = await fetch(url, { ...options, headers });
      
      if (!response.ok) {
        const errorText = await response.text();
        return { data: null, error: errorText };
      }

      const data = await response.json();
      return { data, error: null };
    } catch (error) {
      return { 
        data: null, 
        error: error instanceof Error ? error.message : 'Unknown error' 
      };
    }
  }

  async getUsers(): Promise<{ data: User[] | null; error: string | null }> {
    return this.request<User[]>('/api/v1/users');
  }

  async createUser(user: CreateUserRequest): Promise<{ data: User | null; error: string | null }> {
    return this.request<User>('/api/v1/users', {
      method: 'POST',
      body: JSON.stringify(user),
    });
  }
}

export const apiClient = new ApiClient(process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080');
```

### State Management
```typescript
// contexts/AuthContext.tsx
interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>;
  logout: () => void;
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  // Implementation...

  const value = {
    user,
    isAuthenticated,
    isLoading,
    login,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
```

---

## 5️⃣ Database Layer

### Environment Variables
```bash
# .env.local (frontend)
NEXT_PUBLIC_API_URL=http://localhost:8080

# .env (backend)
DATABASE_URL=postgres://admindb:password@localhost:5432/gonsgarage?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
PORT=8080
```

### Migration Files
```sql
-- migrations/001_create_users_table.up.sql
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

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
```

### Repository Implementation
```go
// internal/adapters/repository/user_repository.go
type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (id, email, password_hash, first_name, last_name, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
    
    _, err := r.db.ExecContext(ctx, query,
        user.ID, user.Email, user.PasswordHash, user.FirstName, 
        user.LastName, user.Role, user.CreatedAt, user.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}
```

---

## 6️⃣ Error Handling and Logging

### Go Error Handling
```go
// Good - Wrapped errors with context
func (s *userService) GetUser(ctx context.Context, id string) (*domain.User, error) {
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user not found: %w", err)
        }
        return nil, fmt.Errorf("failed to get user %s: %w", id, err)
    }
    return user, nil
}

// Structured logging
s.logger.Error("database operation failed",
    "operation", "GetUser",
    "user_id", id,
    "error", err)
```

### TypeScript Error Handling
```typescript
// Error boundary component
class ErrorBoundary extends React.Component<
  { children: React.ReactNode },
  { hasError: boolean }
> {
  constructor(props: { children: React.ReactNode }) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(_: Error) {
    return { hasError: true };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error boundary caught an error:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return <div>Something went wrong. Please refresh the page.</div>;
    }

    return this.props.children;
  }
}

// API error handling
async function handleApiCall<T>(apiCall: () => Promise<{ data: T | null; error: string | null }>) {
  try {
    const result = await apiCall();
    if (result.error) {
      throw new Error(result.error);
    }
    return result.data;
  } catch (error) {
    console.error('API call failed:', error);
    throw new Error('Operation failed. Please try again.');
  }
}
```

---

## 7️⃣ Commit and Workflow Rules

### Conventional Commits
```bash
# Format: <type>(<scope>): <description>

# Examples:
feat(auth): add user registration endpoint
fix(ui): resolve button styling issue
refactor(database): optimize user queries
test(services): add user service unit tests
docs(api): update authentication documentation
```

### Git Workflow
1. Create feature branch: `git checkout -b feat/user-authentication`
2. Write tests first (TDD)
3. Implement feature
4. Run tests: `go test ./...` and `npm test`
5. Lint code: `golangci-lint run` and `npm run lint`
6. Create PR with:
   - Clear description
   - Test results
   - Screenshots (if UI changes)

---

## 8️⃣ Agent Behaviors

### Code Generation Rules
When generating code, AI assistants must:

1. **Follow TDD**: Always write tests before implementation
2. **Respect architecture**: No business logic in adapters layer
3. **Use proper naming**: camelCase/PascalCase as specified
4. **Add documentation**: All exported functions need comments
5. **Handle errors**: Proper error wrapping and user-friendly messages
6. **Validate inputs**: Check data at API boundaries

### Example Generated Code
```go
// CreateUser creates a new user in the system
// It validates the input, checks for duplicates, and stores the user
func (s *userService) CreateUser(ctx context.Context, req CreateUserRequest) (*domain.User, error) {
    // Validate input
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid user data: %w", err)
    }
    
    // Check for existing user
    existing, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil && !errors.Is(err, sql.ErrNoRows) {
        return nil, fmt.Errorf("failed to check existing user: %w", err)
    }
    if existing != nil {
        return nil, fmt.Errorf("user with email %s already exists", req.Email)
    }
    
    // Create user domain object
    user := &domain.User{
        ID:        generateID(),
        Email:     req.Email,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Role:      req.Role,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Store in database
    if err := s.userRepo.Create(ctx, user); err != nil {
        s.logger.Error("failed to create user", "email", req.Email, "error", err)
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    s.logger.Info("user created successfully", "user_id", user.ID, "email", user.Email)
    return user, nil
}
```

---

## 🧠 Additional Guidelines

### Configuration Management
- Use environment variables for all configuration
- Provide `.env.example` files
- Never commit secrets or credentials
- Use different configs for different environments

### Performance Considerations
- Use database indexes appropriately
- Implement caching with Redis for expensive operations
- Use React.memo for expensive components
- Implement pagination for large data sets

### Security Best Practices
- Validate all inputs
- Use prepared statements for SQL queries
- Implement proper authentication and authorization
- Use HTTPS in production
- Sanitize user inputs to prevent XSS

### Documentation Requirements
- All public APIs must have OpenAPI/Swagger documentation
- Complex business logic needs inline comments
- README files for each major component
- Architecture decision records (ADRs) for significant changes

---

This agent file serves as the single source of truth for all development practices in the GonsGarage project. All contributors and AI assistants must adhere to these guidelines to ensure code quality, maintainability, and consistency across the entire codebase.