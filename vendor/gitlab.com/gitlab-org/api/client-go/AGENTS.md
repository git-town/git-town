# AI Agent Guidelines for GitLab client-go

This document provides comprehensive guidelines for AI agents working with the GitLab client-go repository. It covers development practices, testing requirements, code formatting, API alignment, and code generation procedures.

## Repository Overview

The GitLab client-go is a Go client library for the GitLab API, enabling Go programs to interact with GitLab in a simple and uniform way. The repository follows strict Go best practices and maintains close alignment with GitLab's official API documentation.

## Development Workflow

### Prerequisites

When asked to modify code, read CONTRIBUTING.md and README.md for examples and formatting instructions. Where the 
instructions in CONTRIBUTING.md and README.md conflict with information in AGENTS.md, prefer the instructions in 
CONTRIBUTING.md and README.md over the instructions in AGENTS.md

When asked to perform analysis on the codebase instead of changing code, skipping the read of CONTRIBUTING.md and README.md
is allowed and preferred, since understanding contributing guidelines is not required to perform analysis.

### Required Tools

- **Go** - Use the version specified in go.mod
- **gofumpt** - Code formatter
- **golangci-lint** - Linting tool
- **buf** - Protocol buffer tools for code generation
- **gomock** - Mock generation

### Running Tests

```bash
# Run all tests with race detection
make test

# Run the complete reviewable process (includes tests)
make reviewable
```

### Test Patterns

- All tests use the `testing` package with `testify/assert`
- Tests are parallelized using `t.Parallel()`
- Mock HTTP handlers are used for API testing
- Test data is stored in `testdata/` directory
- Each service method should have corresponding test coverage
  - **CRITICAL** - When fixing bugs or creating new features, ensure new test scenarios are added to cover the new logic.
- When writing a test, write Gherkin comments in-line with the test to make the tests easier to read. This means adding GIVEN/WHEN/THEN comments in tests.


### Test Structure Example

```go
func TestGetUser(t *testing.T) {
    t.Parallel()
    mux, client := setup(t)

    path := "/api/v4/users/1"
    mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
        testMethod(t, r, http.MethodGet)
        mustWriteHTTPResponse(t, w, "testdata/get_user.json")
    })

    user, _, err := client.Users.GetUser(1, GetUsersOptions{})
    assert.NoError(t, err)
    // ... assertions
}
```

## Code Formatting and Linting

### Formatting

The project uses `gofumpt` for code formatting:

```bash
# Format all Go files
make fmt
```

**Formatting Rules:**
- Line width for comments: < 80 characters
- Line width for code: < 100 characters (where sensible)
- Use `gofumpt` for consistent formatting
- Follow Go best practices

### Linting

```bash
# Run all linters
make lint
```

**Linting Configuration:**
- Uses `golangci-lint` with custom configuration in `.golangci.yml`
- Enabled linters: asciicheck, dogsled, dupword, errorlint, goconst, misspell, 
  nakedret, nolintlint, revive, staticcheck, testifylint, unconvert, 
  usestdlibvars, whitespace
- Excludes generated files and examples directory

## Mock generation

This repository uses gomock to generate testing structs, which are in the `testing/` folder. These need to be kept up-to-date with function signatures to that the Service implementations match the interfaces that have generated mocks.

### Available Generation Commands

```bash
# Generate all code (protobuf, mocks, testing client)
make generate

# Clean generated files
make clean
```

### Generation Scripts

1. **`scripts/generate_testing_client.sh`** - Generates testing client with mocks
2. **`scripts/generate_mock_api.sh`** - Generates mock interfaces for all services
3. **`scripts/generate_service_interface_map.sh`** - Generates service interface mapping

### When to Regenerate

- After adding new service interfaces
- After modifying existing interfaces
- Before committing changes
- When mock generation fails

## Function Comment Formatting

### Required Comment Structure

Every public function, type, and method must have properly formatted comments:

```go
// FunctionName performs a specific action with the given parameters.
//
// GitLab API docs: https://docs.gitlab.com/api/endpoint/
func (s *ServiceName) FunctionName(param Type, opt *Options, options ...RequestOptionFunc) (*ReturnType, *Response, error) {
    // Implementation
}
```

### Comment Guidelines

1. **Function Comments:**
   - Start with function name (no "The" or "This function")
   - Use present tense ("performs", "returns", "creates")
   - Keep under 80 characters per line
   - Include GitLab API documentation link

2. **Type Comments:**
   - Start with type name
   - Describe the purpose and usage
   - Include GitLab API documentation link

3. **Struct Field Comments:**
   - Use `json:"field_name"` tags
   - Include `url:"field_name,omitempty"` for query parameters
   - Document complex fields

### GitLab API Documentation Alignment

**CRITICAL: All code must align with GitLab's official API documentation.**

#### API Documentation References

Every function must reference the corresponding GitLab API documentation:

```go
// GitLab API docs: https://docs.gitlab.com/api/users/
// GitLab API docs: https://docs.gitlab.com/api/projects/#list-all-projects
// GitLab API docs: https://docs.gitlab.com/api/commits/#get-the-diff-of-a-commit
```

#### Field Ordering

Struct fields and methods should be ordered to match the GitLab API documentation:

```go
type CreateProjectOptions struct {
    Name                     *string `url:"name,omitempty" json:"name,omitempty"`
    Description              *string `url:"description,omitempty" json:"description,omitempty"`
    Visibility               *VisibilityValue `url:"visibility,omitempty" json:"visibility,omitempty"`
    // ... other fields in API documentation order
}
```

#### Parameter Validation

- Use `any` type for project/group IDs to support both int64 and string. 
- Implement proper parameter parsing with `parseID()` function
- Validate required parameters before making API calls

## Code Structure and Patterns

### Service Structure

Each GitLab API service follows this pattern:

```go
type (
    ServiceNameInterface interface {
        MethodName(opt *MethodOptions, options ...RequestOptionFunc) (*ReturnType, *Response, error)
        // ... other methods
    }

    // ServiceName handles communication with the service related methods
    // of the GitLab API.
    //
    // GitLab API docs: https://docs.gitlab.com/api/service/
    ServiceName struct {
        client *Client
    }
)

var _ ServiceNameInterface = (*ServiceName)(nil)
```

### Request Options Pattern

All API methods should accept `options ...RequestOptionFunc`:

```go
func (s *ServiceName) MethodName(opt *MethodOptions, options ...RequestOptionFunc) (*ReturnType, *Response, error) {
    project, err := parseID(pid)
    if err != nil {
        return nil, nil, err
    }
    
    u := fmt.Sprintf("projects/%s/endpoint", PathEscape(project))
    
    req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
    if err != nil {
        return nil, nil, err
    }
    
    var result *ReturnType
    resp, err := s.client.Do(req, &result)
    if err != nil {
        return nil, resp, err
    }
    
    return result, resp, nil
}
```

### Error Handling

- Always return `(*Type, *Response, error)` tuple
- Use `PathEscape()` for URL path parameters
- Use `url.PathEscape()` for query parameters
- Handle `parseID()` errors for project/group IDs

### Type Usage

- Do not use `interface{}`, use the `any` alias instead!
- Do not use `int`, use `int64` instead! This applies to both slices and maps.

## Pre-commit Checklist

**CRITICAL: Tests MUST be run for every build or code modification.**
**CRITICAL: Linting MUST pass for every build or code modification.**
**CRITICAL: Mock generation should be run any time function signatures change**

You can accomplish all three of these by running `make reviewable`, which will do:

1. `make setup` - Install dependencies
2. `make generate` - Generate required code
3. `make fmt` - Format code
4. `make lint` - Run linters
5. `make test` - Run tests

## Code Generation Guidelines

### When Adding New Services

1. Create the service file (e.g., `new_service.go`)
2. Define the interface and struct following the established pattern
3. Implement all methods with proper error handling
4. Add comprehensive tests in `new_service_test.go`
5. Run `make generate` to update mocks and testing client
6. Ensure all tests pass with `make test`

### Mock Generation

The repository uses `gomock` for generating mocks:

```bash
# Generate mocks for all interfaces
make generate
```

Mocks are automatically generated in the `testing/` package and should not be manually edited.

## File Organization

### Service Files

- One service per file (e.g., `users.go`, `projects.go`)
- Corresponding test file (e.g., `users_test.go`)
- Interface definition at the top of the file
- Service struct and implementation below

### Generated Files

- `testing/*_mock.go` - Generated mock files
- `testing/*_generated.go` - Generated testing client files
- `*_generated_test.go` - Generated test files

**Never edit generated files manually.**

## Common Patterns and Best Practices

### Pointer Usage

Use pointers for optional fields in structs:

```go
type CreateUserOptions struct {
    Name     *string `url:"name,omitempty" json:"name,omitempty"`
    Email    *string `url:"email,omitempty" json:"email,omitempty"`
    Username *string `url:"username,omitempty" json:"username,omitempty"`
}
```

### Time Handling

Use `*time.Time` for time fields and `ISOTime` for custom time types that only support year-month-day formatting:

```go
type User struct {
    CreatedAt *time.Time `json:"created_at"`
    LastActivityOn *ISOTime `json:"last_activity_on"`
}
```

### Response Handling

Always return the full response for pagination and metadata:

```go
users, resp, err := client.Users.ListUsers(&gitlab.ListUsersOptions{})
if err != nil {
    return err
}

// Access pagination info
fmt.Printf("Total pages: %d\n", resp.TotalPages)
```

## Troubleshooting

### Common Issues

1. **Tests failing after changes:**
   - Run `make generate` to update mocks
   - Check for linting errors with `make lint`
   - Ensure all imports are correct

2. **Linting errors:**
   - Run `make fmt` to fix formatting issues
   - Check `.golangci.yml` for specific rule configurations
   - Address any static analysis warnings

3. **Generation failures:**
   - Ensure all interfaces are properly defined
   - Check that service files follow the correct pattern
   - Verify that all required tools are installed

### Getting Help

- Check existing issues in the [issue tracker](https://gitlab.com/gitlab-org/api/client-go/-/issues)
- Review the [contributing guide](CONTRIBUTING.md)
- Examine similar implementations in the codebase
- Refer to [GitLab API documentation](https://docs.gitlab.com/ee/api/)

## Summary

When working with this repository:

1. **Always run tests** - `make test` is mandatory
2. **Follow formatting rules** - Use `gofumpt` and respect line limits
3. **Align with GitLab API docs** - Every function must reference official documentation
4. **Generate code when needed** - Run `make generate` after interface changes
5. **Use proper commenting** - Include GitLab API links and follow format guidelines
6. **Maintain consistency** - Follow established patterns and conventions