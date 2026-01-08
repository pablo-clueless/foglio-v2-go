# Testing Guide

This project includes comprehensive testing setup with unit tests, e2e tests, and CI/CD integration.

## Test Structure

```
tests/
├── e2e/           # End-to-end integration tests
│   └── api_test.go
└── utils/         # Test utilities and helpers
    └── test_helper.go

src/
└── handlers/
    └── auth_test.go  # Unit tests for handlers
```

## Running Tests

### All Tests
```bash
make test
```

### Unit Tests Only
```bash
make test-unit
```

### E2E Tests Only
```bash
make test-e2e
```

### Test Coverage
```bash
make test-coverage
```

### CI Tests (with test environment)
```bash
make test-ci
```

### Docker Tests
```bash
make docker-test
```

## Test Environment

Tests use a separate test database configuration defined in `.env.test`. Make sure to:

1. Create a test database: `foglio_test`
2. Update `.env.test` with your test database credentials
3. Run `make test-ci` to use the test environment

## CI/CD

The project includes GitHub Actions workflow (`.github/workflows/ci.yml`) that:

- Runs on push/PR to main/develop branches
- Sets up PostgreSQL service
- Runs unit tests, e2e tests, linting, and security scans
- Builds the application

## Test Database Setup

For local testing, create a test database:

```sql
CREATE DATABASE foglio_test;
```

## Writing Tests

### Unit Tests
Place unit tests alongside source files with `_test.go` suffix:
```go
// src/handlers/auth_test.go
func TestRegisterHandler(t *testing.T) {
    // Test implementation
}
```

### E2E Tests
Add e2e tests in `tests/e2e/` directory using testify suite:
```go
func (suite *E2ETestSuite) TestNewEndpoint() {
    // Test implementation
}
```

## Test Utilities

Use the test utilities in `tests/utils/test_helper.go`:
- `SetupTestServer()` - Initialize test server
- `MakeRequest()` - Make HTTP requests
- `MakeAuthenticatedRequest()` - Make authenticated requests
- `AssertJSONResponse()` - Assert JSON response format