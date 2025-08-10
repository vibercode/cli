package config

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// Validator handles configuration validation
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new configuration validator
func NewValidator() *Validator {
	v := validator.New()
	
	// Register custom validation functions
	v.RegisterValidation("port", validatePort)
	v.RegisterValidation("host", validateHost)
	v.RegisterValidation("url_http", validateHTTPURL)
	v.RegisterValidation("file_path", validateFilePath)
	v.RegisterValidation("database_provider", validateDatabaseProvider)
	v.RegisterValidation("log_level", validateLogLevel)
	v.RegisterValidation("storage_provider", validateStorageProvider)
	v.RegisterValidation("cache_provider", validateCacheProvider)
	v.RegisterValidation("auth_provider", validateAuthProvider)
	
	// Register custom tag name function for better error messages
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	
	return &Validator{validator: v}
}

// ValidateConfig validates the entire configuration
func (v *Validator) ValidateConfig(config *Config) error {
	if err := v.validator.Struct(config); err != nil {
		return v.formatValidationErrors(err)
	}
	
	// Additional custom validations
	if err := v.validateCustomRules(config); err != nil {
		return err
	}
	
	return nil
}

// ValidateSection validates a specific configuration section
func (v *Validator) ValidateSection(section interface{}) error {
	if err := v.validator.Struct(section); err != nil {
		return v.formatValidationErrors(err)
	}
	return nil
}

// formatValidationErrors converts validator errors to custom format
func (v *Validator) formatValidationErrors(err error) error {
	var validationErrors ValidationErrors
	
	if validatorErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validatorError := range validatorErrors {
			validationError := ValidationError{
				Field:   validatorError.Field(),
				Tag:     validatorError.Tag(),
				Value:   fmt.Sprintf("%v", validatorError.Value()),
				Message: v.getErrorMessage(validatorError),
			}
			validationErrors = append(validationErrors, validationError)
		}
	}
	
	return validationErrors
}

// getErrorMessage returns a human-readable error message
func (v *Validator) getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "must be a valid email address"
	case "url":
		return "must be a valid URL"
	case "url_http":
		return "must be a valid HTTP/HTTPS URL"
	case "min":
		return fmt.Sprintf("must be at least %s", err.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", err.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s characters", err.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", strings.ReplaceAll(err.Param(), " ", ", "))
	case "port":
		return "must be a valid port number (1-65535)"
	case "host":
		return "must be a valid hostname or IP address"
	case "file_path":
		return "must be a valid file path"
	case "database_provider":
		return "must be a supported database provider (postgres, mysql, sqlite, mongodb, supabase, redis)"
	case "log_level":
		return "must be a valid log level (debug, info, warn, error)"
	case "storage_provider":
		return "must be a valid storage provider (local, s3, gcs, supabase)"
	case "cache_provider":
		return "must be a valid cache provider (memory, redis)"
	case "auth_provider":
		return "must be a valid auth provider (jwt, oauth2, supabase)"
	default:
		return fmt.Sprintf("validation failed for tag '%s'", err.Tag())
	}
}

// validateCustomRules performs additional validation logic
func (v *Validator) validateCustomRules(config *Config) error {
	var errors ValidationErrors
	
	// Validate server configuration
	if config.Server.EnableHTTPS {
		if config.Server.TLSCertFile == "" {
			errors = append(errors, ValidationError{
				Field:   "server.tls_cert_file",
				Tag:     "required_with_https",
				Message: "TLS certificate file is required when HTTPS is enabled",
			})
		}
		if config.Server.TLSKeyFile == "" {
			errors = append(errors, ValidationError{
				Field:   "server.tls_key_file",
				Tag:     "required_with_https",
				Message: "TLS key file is required when HTTPS is enabled",
			})
		}
	}
	
	// Validate database configuration based on provider
	switch config.Database.Provider {
	case "postgres", "mysql":
		if config.Database.Host == "" {
			errors = append(errors, ValidationError{
				Field:   "database.host",
				Tag:     "required_for_network_db",
				Message: "database host is required for network database providers",
			})
		}
		if config.Database.Username == "" {
			errors = append(errors, ValidationError{
				Field:   "database.username",
				Tag:     "required_for_network_db",
				Message: "database username is required for network database providers",
			})
		}
	case "supabase":
		if config.Database.Supabase.ProjectURL == "" {
			errors = append(errors, ValidationError{
				Field:   "database.supabase.project_url",
				Tag:     "required_for_supabase",
				Message: "Supabase project URL is required",
			})
		}
		if config.Database.Supabase.APIKey == "" {
			errors = append(errors, ValidationError{
				Field:   "database.supabase.api_key",
				Tag:     "required_for_supabase",
				Message: "Supabase API key is required",
			})
		}
	}
	
	// Validate auth configuration
	if config.Auth.Provider == "jwt" && config.Auth.JWTSecret == "" {
		errors = append(errors, ValidationError{
			Field:   "auth.jwt_secret",
			Tag:     "required_for_jwt",
			Message: "JWT secret is required when using JWT authentication",
		})
	}
	
	// Validate storage configuration
	if config.Storage.Provider == "s3" {
		if config.Storage.S3.Bucket == "" {
			errors = append(errors, ValidationError{
				Field:   "storage.s3.bucket",
				Tag:     "required_for_s3",
				Message: "S3 bucket is required when using S3 storage",
			})
		}
		if config.Storage.S3.Region == "" {
			errors = append(errors, ValidationError{
				Field:   "storage.s3.region",
				Tag:     "required_for_s3",
				Message: "S3 region is required when using S3 storage",
			})
		}
	}
	
	// Validate cache configuration
	if config.Cache.Provider == "redis" {
		if config.Cache.Redis.Host == "" {
			errors = append(errors, ValidationError{
				Field:   "cache.redis.host",
				Tag:     "required_for_redis",
				Message: "Redis host is required when using Redis cache",
			})
		}
	}
	
	if len(errors) > 0 {
		return errors
	}
	
	return nil
}

// Custom validation functions

func validatePort(fl validator.FieldLevel) bool {
	port := fl.Field().Int()
	return port >= 1 && port <= 65535
}

func validateHost(fl validator.FieldLevel) bool {
	host := fl.Field().String()
	if host == "" {
		return false
	}
	
	// Check if it's a valid IP address
	if ip := net.ParseIP(host); ip != nil {
		return true
	}
	
	// Check if it's a valid hostname
	if len(host) > 253 {
		return false
	}
	
	// Basic hostname validation
	for _, part := range strings.Split(host, ".") {
		if len(part) == 0 || len(part) > 63 {
			return false
		}
		for _, char := range part {
			if !((char >= 'a' && char <= 'z') || 
				 (char >= 'A' && char <= 'Z') || 
				 (char >= '0' && char <= '9') || 
				 char == '-') {
				return false
			}
		}
	}
	
	return true
}

func validateHTTPURL(fl validator.FieldLevel) bool {
	urlStr := fl.Field().String()
	if urlStr == "" {
		return true // Allow empty for optional fields
	}
	
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}

func validateFilePath(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if path == "" {
		return true // Allow empty for optional fields
	}
	
	// Basic file path validation (more comprehensive validation would require OS-specific logic)
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			return false
		}
	}
	
	return true
}

func validateDatabaseProvider(fl validator.FieldLevel) bool {
	provider := fl.Field().String()
	validProviders := []string{"postgres", "mysql", "sqlite", "mongodb", "supabase", "redis"}
	
	for _, valid := range validProviders {
		if provider == valid {
			return true
		}
	}
	
	return false
}

func validateLogLevel(fl validator.FieldLevel) bool {
	level := fl.Field().String()
	validLevels := []string{"debug", "info", "warn", "error"}
	
	for _, valid := range validLevels {
		if level == valid {
			return true
		}
	}
	
	return false
}

func validateStorageProvider(fl validator.FieldLevel) bool {
	provider := fl.Field().String()
	validProviders := []string{"local", "s3", "gcs", "supabase"}
	
	for _, valid := range validProviders {
		if provider == valid {
			return true
		}
	}
	
	return false
}

func validateCacheProvider(fl validator.FieldLevel) bool {
	provider := fl.Field().String()
	validProviders := []string{"memory", "redis"}
	
	for _, valid := range validProviders {
		if provider == valid {
			return true
		}
	}
	
	return false
}

func validateAuthProvider(fl validator.FieldLevel) bool {
	provider := fl.Field().String()
	validProviders := []string{"jwt", "oauth2", "supabase"}
	
	for _, valid := range validProviders {
		if provider == valid {
			return true
		}
	}
	
	return false
}

// ValidateField validates a single field value against a validation tag
func (v *Validator) ValidateField(value interface{}, tag string) error {
	return v.validator.Var(value, tag)
}

// GetValidationRules returns the validation rules for a given struct
func (v *Validator) GetValidationRules(structType reflect.Type) map[string]string {
	rules := make(map[string]string)
	
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if validateTag := field.Tag.Get("validate"); validateTag != "" {
			jsonName := field.Tag.Get("json")
			if jsonName != "" {
				jsonName = strings.Split(jsonName, ",")[0]
				rules[jsonName] = validateTag
			}
		}
	}
	
	return rules
}