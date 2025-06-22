# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Terraform provider for Lambda Labs GPU cloud services. It enables Terraform users to provision and manage GPU instances, SSH keys, and file systems through Infrastructure as Code.

## Common Development Commands

### Building and Testing
- **Build the provider**: `go build -v .`
- **Install locally**: `make install` or `go install .`
- **Run unit tests**: `go test ./...`
- **Run acceptance tests**: `make testacc` or `TF_ACC=1 go test ./... -cover -timeout 120m`
- **Generate documentation**: `make doc` or `go generate ./...`
- **Run linting**: `golangci-lint run`

### Environment Variables
- `LAMBDALABS_API_KEY`: Lambda Labs API key for authentication
- `LAMBDALABS_BASE_URL`: Override default API base URL (optional)
- `TF_ACC=1`: Required to run acceptance tests

## Code Architecture

### Directory Structure
- `main.go`: Provider entry point using Terraform Plugin Framework v2
- `internal/provider/`: Terraform provider implementation
  - `provider.go`: Main provider configuration and registration
  - `*_resource.go`: Resource implementations (CRUD operations)
  - `*_data.go`: Data source implementations (read-only)
  - `*_model.go`: Terraform state models for complex resources
- `pkg/lambdalabs/`: API client library
  - `client.go`: HTTP client wrapper
  - `transport.go`: Bearer token authentication
  - `schema.go`: Core data structures
  - `*.go`: API endpoint implementations

### Provider Resources
- `lambdalabs_instance`: GPU instances with regions, types, SSH keys
- `lambdalabs_ssh_key`: SSH key management
- `lambdalabs_filesystem`: Persistent storage

### Provider Data Sources  
- `lambdalabs_instance_types`: Available instance configurations
- `lambdalabs_image`: Available OS/software images
- `lambdalabs_ssh_key`: SSH key lookup
- `lambdalabs_filesystem`: File system lookup
- `lambdalabs_firewall`: Firewall configuration lookup

## Code Patterns and Conventions

### File Naming
- Resources: `{resource}_resource.go` (e.g., `instance_resource.go`)
- Data sources: `{resource}_data.go` (e.g., `instance_data.go`)
- Models: `{resource}_model.go` for complex state models
- Tests: `{component}_test.go`
- API client: `{api_entity}.go` in `pkg/lambdalabs/`

### Resource Implementation Pattern
Each resource implements standard methods:
- `Metadata()`: Sets Terraform resource type name
- `Schema()`: Defines attributes with tfsdk tags
- `Configure()`: Receives API client from provider
- `Create()`, `Read()`, `Update()`, `Delete()`: CRUD operations
- `ImportState()`: Enables `terraform import` support

### API Client Pattern
- Request/Response structs with JSON tags
- Context-aware methods accepting `context.Context`
- Consistent error handling with structured errors
- HTTP methods: `Get()`, `Post()`, `Put()`, `Delete()`

### Testing Conventions (from CONVENTIONS.md)
- Use standard library `testing` package
- Use `httptest` package for HTTP mocking
- Structure: `cases` variable with test scenarios
- Hard-coded expected outputs in test cases
- Each function has corresponding test (e.g., `TestCreate` for `Create`)

### Comments Policy
- Only document public functions and types using GoDoc format
- Explain **why** the code does something, not **what** it does
- Focus on business logic and design decisions

## Development Workflow

### Making Changes
1. Understand existing patterns by examining similar resources/data sources
2. Follow the established file naming and code structure conventions
3. Use the API client from `pkg/lambdalabs/` for external calls
4. Add comprehensive tests following the established patterns
5. Run `go generate ./...` to update documentation
6. Ensure acceptance tests pass with `make testacc`

### Provider Configuration
- Default API base URL: `https://cloud.lambdalabs.com/api/v1`
- Authentication via Bearer token in `Authorization` header
- Support for environment variable configuration
- Backward compatibility with deprecated `endpoint` field

### Release Process
- Uses GoReleaser for multi-platform builds
- Automated via GitHub Actions on version tags
- Supports Linux, Darwin, Windows, FreeBSD on multiple architectures
- Includes GPG signing and Terraform registry manifest generation