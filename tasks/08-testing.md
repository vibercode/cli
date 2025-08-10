# Task 08: Testing Framework Integration

## Overview
Implement comprehensive testing framework integration that generates test files and test utilities for Go web APIs. This includes unit tests, integration tests, API tests, and testing utilities with multiple testing frameworks support.

## Objectives
- Generate unit test files for all generated components
- Create integration test suites for API endpoints
- Implement testing utilities and helpers
- Support multiple testing frameworks (Testify, Ginkgo, GoConvey)
- Generate mock files and test data
- Create benchmark tests for performance testing

## Implementation Details

### Command Structure
```bash
# Generate tests for existing code
vibercode generate test --target resource --name User
vibercode generate test --target middleware --name auth
vibercode generate test --target api --full

# Generate test utilities
vibercode generate test --type utils
vibercode generate test --type mocks

# Generate specific test types
vibercode generate test --type unit --target handlers
vibercode generate test --type integration --api-tests
vibercode generate test --type benchmark --performance

# Generate test suite with framework
vibercode generate test --framework testify --suite api
vibercode generate test --framework ginkgo --bdd
```

### Test Types

#### 1. Unit Tests
- Handler unit tests with mocked dependencies
- Service layer tests with business logic validation
- Repository tests with database mocking
- Middleware tests with request/response simulation
- Model validation tests

#### 2. Integration Tests
- API endpoint testing with real database
- Database integration tests
- External service integration
- End-to-end workflow testing

#### 3. Test Utilities
- Database test helpers (setup/teardown)
- HTTP test clients
- Mock generators
- Test data factories
- Assertion helpers

#### 4. Performance Tests
- Benchmark tests for handlers
- Load testing utilities
- Memory usage profiling
- Concurrent request testing

### File Structure
```
test/
├── unit/
│   ├── handlers/
│   │   ├── user_handler_test.go
│   │   └── auth_handler_test.go
│   ├── services/
│   │   ├── user_service_test.go
│   │   └── auth_service_test.go
│   ├── repositories/
│   │   └── user_repository_test.go
│   └── middleware/
│       └── auth_middleware_test.go
├── integration/
│   ├── api/
│   │   ├── user_api_test.go
│   │   └── auth_api_test.go
│   ├── database/
│   │   └── migrations_test.go
│   └── e2e/
│       └── user_workflow_test.go
├── benchmark/
│   ├── handlers_bench_test.go
│   └── database_bench_test.go
├── mocks/
│   ├── user_service_mock.go
│   ├── user_repository_mock.go
│   └── external_api_mock.go
├── utils/
│   ├── test_database.go
│   ├── test_server.go
│   ├── test_client.go
│   └── test_factories.go
└── fixtures/
    ├── users.json
    └── test_data.sql
```

### Testing Frameworks Support

#### 1. Testify (Default)
- Assert and require packages
- Test suites with setup/teardown
- Mock generation with testify/mock
- HTTP testing utilities

#### 2. Ginkgo + Gomega (BDD)
- Behavior-driven development testing
- Descriptive test specifications
- Powerful matchers with Gomega
- Parallel test execution

#### 3. GoConvey
- Web UI for test results
- Nested test scenarios
- Real-time test monitoring

### Templates Required

#### Unit Test Template (Testify)
```go
package handlers

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/gin-gonic/gin"
)

func TestUserHandler_GetUser(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    mockService := &MockUserService{}
    handler := NewUserHandler(mockService)
    
    // Test cases
    tests := []struct {
        name           string
        userID         string
        mockSetup      func()
        expectedStatus int
        expectedBody   string
    }{
        {
            name:   "successful user retrieval",
            userID: "123",
            mockSetup: func() {
                mockService.On("GetUser", "123").Return(&User{ID: "123", Name: "John"}, nil)
            },
            expectedStatus: http.StatusOK,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mock
            tt.mockSetup()
            
            // Create request
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            c.Param("id", tt.userID)
            
            // Execute
            handler.GetUser(c)
            
            // Assert
            assert.Equal(t, tt.expectedStatus, w.Code)
            mockService.AssertExpectations(t)
        })
    }
}
```

#### Integration Test Template
```go
package integration

import (
    "testing"
    "net/http"
    "bytes"
    "encoding/json"
    "github.com/stretchr/testify/suite"
)

type UserAPITestSuite struct {
    suite.Suite
    server *TestServer
    client *TestClient
}

func (suite *UserAPITestSuite) SetupSuite() {
    suite.server = NewTestServer()
    suite.client = NewTestClient(suite.server.URL)
}

func (suite *UserAPITestSuite) TearDownSuite() {
    suite.server.Close()
}

func (suite *UserAPITestSuite) SetupTest() {
    suite.server.ResetDatabase()
}

func (suite *UserAPITestSuite) TestCreateUser() {
    // Test data
    userData := map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
    }
    
    // Make request
    resp, err := suite.client.POST("/api/users", userData)
    suite.NoError(err)
    suite.Equal(http.StatusCreated, resp.StatusCode)
    
    // Verify response
    var user User
    err = json.NewDecoder(resp.Body).Decode(&user)
    suite.NoError(err)
    suite.Equal("John Doe", user.Name)
    suite.NotEmpty(user.ID)
}

func TestUserAPITestSuite(t *testing.T) {
    suite.Run(t, new(UserAPITestSuite))
}
```

### Test Utilities

#### Database Test Helper
```go
package utils

import (
    "database/sql"
    "fmt"
    "testing"
)

type TestDatabase struct {
    DB *sql.DB
    originalDB *sql.DB
}

func NewTestDatabase(t *testing.T) *TestDatabase {
    // Create test database connection
    db := setupTestDB()
    
    return &TestDatabase{
        DB: db,
    }
}

func (td *TestDatabase) Reset() {
    // Clean all tables
    td.truncateAllTables()
    // Run migrations
    td.runMigrations()
}

func (td *TestDatabase) Close() {
    td.DB.Close()
}
```

### Mock Generation

#### Service Mock Template
```go
package mocks

import (
    "github.com/stretchr/testify/mock"
)

type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) GetUser(id string) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserService) CreateUser(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}
```

### Configuration Integration

#### Test Configuration
```go
type TestConfig struct {
    Database struct {
        TestDBName string `yaml:"test_db_name"`
        ResetDB    bool   `yaml:"reset_db"`
    } `yaml:"database"`
    
    Server struct {
        TestPort int `yaml:"test_port"`
    } `yaml:"server"`
    
    Logging struct {
        TestLevel string `yaml:"test_level"`
    } `yaml:"logging"`
}
```

## Dependencies
- Task 02: Template System Enhancement (for test templates)
- Task 07: Middleware Generator (for middleware tests)

## Deliverables
1. Test generator implementation
2. Unit test templates for all component types
3. Integration test templates and utilities
4. Mock generation system
5. Test utilities and helpers
6. Multiple testing framework support
7. Benchmark test templates
8. Test configuration management
9. Documentation and examples

## Acceptance Criteria
- [ ] Generate unit tests for handlers, services, repositories
- [ ] Create integration test suites with database setup/teardown
- [ ] Generate mock files automatically
- [ ] Support multiple testing frameworks (Testify, Ginkgo, GoConvey)
- [ ] Include test utilities and helpers
- [ ] Generate benchmark tests for performance testing
- [ ] Provide test data factories and fixtures
- [ ] Include comprehensive test documentation
- [ ] Support CI/CD integration with test reporting
- [ ] Pass all generated tests successfully

## Implementation Priority
Medium - Important for code quality and maintainability

## Estimated Effort
4-5 days

## Notes
- Focus on comprehensive test coverage
- Ensure tests are maintainable and readable
- Include performance testing capabilities
- Support both TDD and BDD approaches
- Consider test parallelization for large test suites