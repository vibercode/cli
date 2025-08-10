# Task 01: Database Providers Enhancement

## Status: ✅ Completed

## Overview
Enhance the database provider support in Vibercode CLI to include modern cloud database services, with particular focus on Supabase integration.

## Current State
- Basic support for PostgreSQL, MySQL, and SQLite
- Database connection configuration in generated projects
- Basic GORM integration

## Requirements

### 1. Supabase Integration
- Add Supabase as a database provider option
- Include Supabase connection configuration
- Support for Supabase auth integration
- Real-time subscriptions setup
- Storage configuration for file uploads

### 2. Additional Database Providers
- MongoDB support with appropriate ORM
- Redis support for caching
- CockroachDB support
- PlanetScale support

### 3. Connection Management
- Connection pooling configuration
- Environment-specific database configurations
- Database URL parsing and validation
- SSL/TLS configuration options

## Implementation Details

### Files to Modify
- `internal/generator/api.go` - Add provider selection logic
- `internal/templates/` - Add database-specific templates
- `cmd/generate.go` - Update prompts for provider selection

### New Files to Create
- `internal/templates/supabase.go` - Supabase-specific templates
- `internal/models/database.go` - Database provider models
- `pkg/database/supabase.go` - Supabase connection utilities

### Configuration Structure
```go
type DatabaseProvider struct {
    Type        string // "postgres", "mysql", "sqlite", "supabase", etc.
    Host        string
    Port        int
    Database    string
    Username    string
    Password    string
    SSLMode     string
    URL         string // For cloud providers
    ProjectRef  string // For Supabase
    AnonKey     string // For Supabase
    ServiceKey  string // For Supabase
}
```

## Acceptance Criteria
- [x] Supabase can be selected as database provider
- [x] Generated projects connect successfully to Supabase
- [x] Database-specific configurations are properly generated
- [x] Connection utilities handle different provider types
- [x] Error handling for connection failures
- [x] Documentation updated with provider-specific instructions

## Implementation Completed
- ✅ Enhanced `DatabaseProvider` model with Supabase-specific fields
- ✅ Created Supabase database template (`internal/templates/supabase.go`)
- ✅ Updated API generator to use Supabase templates
- ✅ Added environment variable templates for Supabase
- ✅ Implemented connection validation and error handling

## Dependencies
- None (can be implemented immediately)

## Effort Estimate
- 2-3 days of development
- 1 day for testing and documentation

## Testing Requirements
- Unit tests for database provider logic
- Integration tests with actual database connections
- Template generation validation
- Error handling scenarios