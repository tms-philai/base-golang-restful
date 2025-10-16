# Test Documentation

This directory contains comprehensive unit tests and integration tests for the base-gin project.

## 📁 Directory Structure

```
test/
├── unit/                    # Unit tests for individual components
│   ├── auth_handler_test.go
│   ├── user_handler_test.go
│   ├── product_handler_test.go
│   ├── email_handler_test.go
│   ├── file_handler_test.go
│   └── monitoring_handler_test.go
├── integration/             # Integration tests for API endpoints
│   ├── auth_integration_test.go
│   ├── user_integration_test.go
│   └── product_integration_test.go
├── helpers/                 # Test helper functions and utilities
│   └── test_helpers.go
├── config/                  # Test configuration
│   └── test_config.go
├── setup/                   # Test setup and teardown
│   └── database_setup.go
└── README.md               # This file
```

## 🧪 Test Types

### Unit Tests
Unit tests focus on testing individual components in isolation using mocks and stubs. They are located in the `unit/` directory and test:

- **AuthHandler**: Authentication endpoints (register, login, refresh, profile, change-password)
- **UserHandler**: User management endpoints (CRUD operations)
- **ProductHandler**: Product management endpoints (CRUD operations, stock management)
- **EmailHandler**: Email service endpoints (send, bulk send, status, connection test)
- **FileHandler**: File management endpoints (upload, download, list, delete)
- **MonitoringHandler**: Health check and metrics endpoints

### Integration Tests
Integration tests test the complete API endpoints with real database connections. They are located in the `integration/` directory and test:

- **Auth Integration**: Complete authentication flow
- **User Integration**: Complete user management flow
- **Product Integration**: Complete product management flow

## 🚀 Running Tests

### Prerequisites

1. **PostgreSQL Database**: Ensure PostgreSQL is running
2. **Test Database**: Create a test database:
   ```sql
   CREATE DATABASE test_base_gin;
   CREATE USER test_user WITH PASSWORD 'test_password';
   GRANT ALL PRIVILEGES ON DATABASE test_base_gin TO test_user;
   ```

### Using Makefile

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run tests with coverage
make test-coverage

# Run unit tests with coverage
make test-unit-coverage

# Run integration tests with coverage
make test-integration-coverage

# Run tests with race detection
make test-race

# Run benchmarks
make benchmark

# Setup test environment
make test-setup

# Clean test environment
make test-clean

# Show help
make help
```

### Using Go Commands

```bash
# Run all tests
go test ./...

# Run unit tests only
go test ./test/unit/...

# Run integration tests only
go test ./test/integration/...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific test package
go test -v ./test/unit

# Run specific test function
go test -v -run TestAuthHandler_Register ./test/unit

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...
```

## 📊 Test Coverage

The test suite aims for high coverage across all components:

- **Unit Tests**: Test individual handlers with mocked dependencies
- **Integration Tests**: Test complete API flows with real database
- **Coverage Reports**: Generate HTML coverage reports for detailed analysis

### Coverage Targets

- **Handlers**: >90% coverage
- **Services**: >85% coverage
- **Repositories**: >80% coverage
- **Overall**: >85% coverage

## 🔧 Test Configuration

### Environment Variables

Test configuration is automatically generated in `configs/.env.test` with the following defaults:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=test_user
DB_PASSWORD=test_password
DB_NAME=test_base_gin
DB_SSLMODE=disable
DB_TIMEZONE=UTC

# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8001
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30

# JWT Configuration
JWT_SECRET_KEY=test-secret-key-for-testing-only
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=7d

# Email Configuration
SMTP_HOST=localhost
SMTP_PORT=587
SMTP_USERNAME=test@example.com
SMTP_PASSWORD=test_password
SMTP_FROM_EMAIL=test@example.com
SMTP_FROM_NAME=Test Sender

# Storage Configuration
STORAGE_TYPE=local
LOCAL_STORAGE_PATH=/tmp/test_files
S3_BUCKET=test-bucket
S3_REGION=us-east-1
S3_ACCESS_KEY=test-access-key
S3_SECRET_KEY=test-secret-key
```

## 🛠️ Test Helpers

The `helpers/` directory contains utility functions for testing:

### TestUser and TestProduct
- `CreateTestUser()`: Creates a test user with default values
- `CreateTestAdmin()`: Creates a test admin user
- `CreateTestProduct()`: Creates a test product with default values

### Request Helpers
- `CreateJSONRequest()`: Creates JSON HTTP requests
- `CreateMultipartRequest()`: Creates multipart form requests

### Assertion Helpers
- `AssertJSONResponse()`: Asserts JSON response format
- `AssertErrorResponse()`: Asserts error response format
- `AssertSuccessResponse()`: Asserts success response format
- `AssertAuthResponse()`: Asserts authentication response format
- `AssertUserResponse()`: Asserts user response format
- `AssertProductResponse()`: Asserts product response format
- `AssertFileResponse()`: Asserts file response format
- `AssertPaginationResponse()`: Asserts pagination response format

### Middleware Helpers
- `MockAuthMiddleware()`: Creates mock authentication middleware
- `MockOptionalAuthMiddleware()`: Creates mock optional authentication middleware
- `MockAdminMiddleware()`: Creates mock admin middleware

## 🗄️ Database Setup

### Test Database Setup

The `setup/` directory contains database setup and teardown utilities:

- **TestDatabaseSetup**: Handles test database connection and cleanup
- **SetupTestDatabase()**: Sets up test database with migrations
- **CleanupTestDatabase()**: Cleans up test database connection
- **TruncateTestTables()**: Truncates all test tables
- **SetupTestData()**: Sets up test data
- **CleanupTestData()**: Cleans up test data

### Test Data Management

- **SetupTest()**: Sets up each test with clean data
- **TearDownSuite()**: Tears down the test suite
- **TruncateTestTables()**: Ensures clean state between tests

## 📝 Writing Tests

### Unit Test Guidelines

1. **Use Mocks**: Mock all external dependencies
2. **Test Edge Cases**: Test validation errors, unauthorized access, etc.
3. **Assert Response Format**: Use helper functions for consistent assertions
4. **Clean Naming**: Use descriptive test names
5. **Single Responsibility**: Each test should test one specific behavior

### Integration Test Guidelines

1. **Real Database**: Use real database connections
2. **Complete Flows**: Test complete API flows
3. **Authentication**: Test with real JWT tokens
4. **Data Cleanup**: Ensure clean state between tests
5. **Error Scenarios**: Test error handling and edge cases

### Example Unit Test

```go
func TestAuthHandler_Register(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    tests := []struct {
        name           string
        request        models.RegisterRequest
        mockSetup      func(*MockUserService, *MockRoleService, *MockJWTManager)
        expectedStatus int
        expectedError  string
    }{
        {
            name: "successful registration",
            request: models.RegisterRequest{
                Username:  "testuser",
                Email:     "test@example.com",
                Password:  "password123",
                FirstName: "Test",
                LastName:  "User",
            },
            mockSetup: func(userService *MockUserService, roleService *MockRoleService, jwtManager *MockJWTManager) {
                // Setup mocks
            },
            expectedStatus: http.StatusCreated,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Example Integration Test

```go
func (suite *AuthIntegrationTestSuite) TestRegister() {
    tests := []struct {
        name           string
        request        helpers.TestUser
        expectedStatus int
        expectError    bool
    }{
        {
            name: "successful registration",
            request: helpers.TestUser{
                Username:  "newuser",
                Email:     "newuser@example.com",
                Password:  "password123",
                FirstName: "New",
                LastName:  "User",
            },
            expectedStatus: http.StatusCreated,
            expectError:    false,
        },
    }
    
    for _, tt := range tests {
        suite.Run(tt.name, func() {
            // Test implementation
        })
    }
}
```

## 🐛 Debugging Tests

### Common Issues

1. **Database Connection**: Ensure PostgreSQL is running and test database exists
2. **Port Conflicts**: Ensure test port (8001) is available
3. **File Permissions**: Ensure test file directory has proper permissions
4. **Environment Variables**: Check test configuration is properly loaded

### Debug Commands

```bash
# Run specific test with verbose output
go test -v -run TestAuthHandler_Register ./test/unit

# Run tests with race detection
go test -race ./test/integration

# Run tests with timeout
go test -timeout=30s ./test/integration

# Run tests in parallel
go test -parallel=4 ./test/unit
```

## 📈 Continuous Integration

### CI/CD Pipeline

The test suite is designed to work with CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: test_password
          POSTGRES_DB: test_base_gin
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Run tests
        run: make ci-test
```

## 📚 Additional Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Gin Testing Documentation](https://gin-gonic.com/docs/testing/)
- [Go Mock Documentation](https://github.com/golang/mock)

## 🤝 Contributing

When adding new tests:

1. Follow the existing test structure
2. Use appropriate test helpers
3. Ensure high test coverage
4. Update this documentation if needed
5. Run all tests before submitting

## 📞 Support

For questions about testing:

1. Check this documentation
2. Review existing test examples
3. Check the Makefile for available commands
4. Ensure all prerequisites are met
