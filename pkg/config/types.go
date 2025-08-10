package config

import "time"

// Config represents the complete application configuration
type Config struct {
	Environment string           `json:"environment" yaml:"environment" validate:"required,oneof=development staging production"`
	Server      ServerConfig     `json:"server" yaml:"server" validate:"required"`
	Database    DatabaseConfig   `json:"database" yaml:"database" validate:"required"`
	Auth        AuthConfig       `json:"auth" yaml:"auth"`
	Storage     StorageConfig    `json:"storage" yaml:"storage"`
	Cache       CacheConfig      `json:"cache" yaml:"cache"`
	Logging     LoggingConfig    `json:"logging" yaml:"logging"`
	Monitoring  MonitoringConfig `json:"monitoring" yaml:"monitoring"`
	Security    SecurityConfig   `json:"security" yaml:"security"`
	Features    FeatureConfig    `json:"features" yaml:"features"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Host           string        `json:"host" yaml:"host" validate:"required"`
	Port           int           `json:"port" yaml:"port" validate:"required,min=1,max=65535"`
	ReadTimeout    time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout    time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	MaxHeaderBytes int           `json:"max_header_bytes" yaml:"max_header_bytes"`
	EnableHTTPS    bool          `json:"enable_https" yaml:"enable_https"`
	TLSCertFile    string        `json:"tls_cert_file" yaml:"tls_cert_file"`
	TLSKeyFile     string        `json:"tls_key_file" yaml:"tls_key_file"`
	CORS           CORSConfig    `json:"cors" yaml:"cors"`
}

// CORSConfig contains CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string `json:"allowed_origins" yaml:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods" yaml:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers" yaml:"allowed_headers"`
	ExposedHeaders   []string `json:"exposed_headers" yaml:"exposed_headers"`
	AllowCredentials bool     `json:"allow_credentials" yaml:"allow_credentials"`
	MaxAge           int      `json:"max_age" yaml:"max_age"`
}

// DatabaseConfig contains database connection configuration
type DatabaseConfig struct {
	Provider         string        `json:"provider" yaml:"provider" validate:"required,oneof=postgres mysql sqlite mongodb supabase redis"`
	Host             string        `json:"host" yaml:"host"`
	Port             int           `json:"port" yaml:"port"`
	Username         string        `json:"username" yaml:"username"`
	Password         string        `json:"password" yaml:"password"`
	Database         string        `json:"database" yaml:"database"`
	Schema           string        `json:"schema" yaml:"schema"`
	SSLMode          string        `json:"ssl_mode" yaml:"ssl_mode"`
	MaxOpenConns     int           `json:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns     int           `json:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime  time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	ConnMaxIdleTime  time.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	Migrations       MigrationConfig `json:"migrations" yaml:"migrations"`
	Supabase         SupabaseConfig  `json:"supabase" yaml:"supabase"`
}

// MigrationConfig contains database migration settings
type MigrationConfig struct {
	AutoMigrate    bool   `json:"auto_migrate" yaml:"auto_migrate"`
	BackupBefore   bool   `json:"backup_before" yaml:"backup_before"`
	Versioned      bool   `json:"versioned" yaml:"versioned"`
	MigrationsPath string `json:"migrations_path" yaml:"migrations_path"`
}

// SupabaseConfig contains Supabase-specific configuration
type SupabaseConfig struct {
	ProjectURL    string `json:"project_url" yaml:"project_url"`
	APIKey        string `json:"api_key" yaml:"api_key"`
	ServiceKey    string `json:"service_key" yaml:"service_key"`
	JWTSecret     string `json:"jwt_secret" yaml:"jwt_secret"`
	EnableAuth    bool   `json:"enable_auth" yaml:"enable_auth"`
	EnableStorage bool   `json:"enable_storage" yaml:"enable_storage"`
	EnableRealtime bool   `json:"enable_realtime" yaml:"enable_realtime"`
}

// AuthConfig contains authentication and authorization settings
type AuthConfig struct {
	Provider       string        `json:"provider" yaml:"provider" validate:"oneof=jwt oauth2 supabase"`
	JWTSecret      string        `json:"jwt_secret" yaml:"jwt_secret"`
	TokenExpiry    time.Duration `json:"token_expiry" yaml:"token_expiry"`
	RefreshExpiry  time.Duration `json:"refresh_expiry" yaml:"refresh_expiry"`
	PasswordMinLen int           `json:"password_min_length" yaml:"password_min_length"`
	EnableMFA      bool          `json:"enable_mfa" yaml:"enable_mfa"`
	OAuth2         OAuth2Config  `json:"oauth2" yaml:"oauth2"`
}

// OAuth2Config contains OAuth2 provider settings
type OAuth2Config struct {
	GoogleClientID     string `json:"google_client_id" yaml:"google_client_id"`
	GoogleClientSecret string `json:"google_client_secret" yaml:"google_client_secret"`
	GitHubClientID     string `json:"github_client_id" yaml:"github_client_id"`
	GitHubClientSecret string `json:"github_client_secret" yaml:"github_client_secret"`
	RedirectURL        string `json:"redirect_url" yaml:"redirect_url"`
}

// StorageConfig contains file storage configuration
type StorageConfig struct {
	Provider    string      `json:"provider" yaml:"provider" validate:"oneof=local s3 gcs supabase"`
	LocalPath   string      `json:"local_path" yaml:"local_path"`
	MaxFileSize int64       `json:"max_file_size" yaml:"max_file_size"`
	AllowedTypes []string    `json:"allowed_types" yaml:"allowed_types"`
	S3          S3Config    `json:"s3" yaml:"s3"`
	Supabase    StorageSupabaseConfig `json:"supabase" yaml:"supabase"`
}

// S3Config contains AWS S3 configuration
type S3Config struct {
	Bucket    string `json:"bucket" yaml:"bucket"`
	Region    string `json:"region" yaml:"region"`
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	Endpoint  string `json:"endpoint" yaml:"endpoint"`
}

// StorageSupabaseConfig contains Supabase storage configuration
type StorageSupabaseConfig struct {
	BucketName string `json:"bucket_name" yaml:"bucket_name"`
	PublicURL  string `json:"public_url" yaml:"public_url"`
}

// CacheConfig contains caching configuration
type CacheConfig struct {
	Provider string      `json:"provider" yaml:"provider" validate:"oneof=memory redis"`
	TTL      time.Duration `json:"ttl" yaml:"ttl"`
	Redis    RedisConfig `json:"redis" yaml:"redis"`
}

// RedisConfig contains Redis-specific configuration
type RedisConfig struct {
	Host         string        `json:"host" yaml:"host"`
	Port         int           `json:"port" yaml:"port"`
	Password     string        `json:"password" yaml:"password"`
	Database     int           `json:"database" yaml:"database"`
	PoolSize     int           `json:"pool_size" yaml:"pool_size"`
	DialTimeout  time.Duration `json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level      string `json:"level" yaml:"level" validate:"oneof=debug info warn error"`
	Format     string `json:"format" yaml:"format" validate:"oneof=json text"`
	Output     string `json:"output" yaml:"output" validate:"oneof=stdout stderr file"`
	Filename   string `json:"filename" yaml:"filename"`
	MaxSize    int    `json:"max_size" yaml:"max_size"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups"`
	MaxAge     int    `json:"max_age" yaml:"max_age"`
	Compress   bool   `json:"compress" yaml:"compress"`
}

// MonitoringConfig contains monitoring and metrics configuration
type MonitoringConfig struct {
	EnableMetrics bool         `json:"enable_metrics" yaml:"enable_metrics"`
	EnableTracing bool         `json:"enable_tracing" yaml:"enable_tracing"`
	EnableHealth  bool         `json:"enable_health" yaml:"enable_health"`
	MetricsPort   int          `json:"metrics_port" yaml:"metrics_port"`
	Prometheus    PrometheusConfig `json:"prometheus" yaml:"prometheus"`
	Jaeger        JaegerConfig `json:"jaeger" yaml:"jaeger"`
}

// PrometheusConfig contains Prometheus metrics configuration
type PrometheusConfig struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Path     string `json:"path" yaml:"path"`
}

// JaegerConfig contains Jaeger tracing configuration
type JaegerConfig struct {
	Endpoint    string  `json:"endpoint" yaml:"endpoint"`
	ServiceName string  `json:"service_name" yaml:"service_name"`
	SampleRate  float64 `json:"sample_rate" yaml:"sample_rate"`
}

// SecurityConfig contains security settings
type SecurityConfig struct {
	EnableRateLimiting bool              `json:"enable_rate_limiting" yaml:"enable_rate_limiting"`
	RateLimit          RateLimitConfig   `json:"rate_limit" yaml:"rate_limit"`
	EnableCSRF         bool              `json:"enable_csrf" yaml:"enable_csrf"`
	CSRFSecret         string            `json:"csrf_secret" yaml:"csrf_secret"`
	TrustedProxies     []string          `json:"trusted_proxies" yaml:"trusted_proxies"`
	SecretKeys         SecretKeysConfig  `json:"secret_keys" yaml:"secret_keys"`
}

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond int           `json:"requests_per_second" yaml:"requests_per_second"`
	BurstSize         int           `json:"burst_size" yaml:"burst_size"`
	CleanupInterval   time.Duration `json:"cleanup_interval" yaml:"cleanup_interval"`
}

// SecretKeysConfig contains encryption keys
type SecretKeysConfig struct {
	EncryptionKey string `json:"encryption_key" yaml:"encryption_key"`
	SigningKey    string `json:"signing_key" yaml:"signing_key"`
}

// FeatureConfig contains feature flags
type FeatureConfig struct {
	EnableGraphQL     bool `json:"enable_graphql" yaml:"enable_graphql"`
	EnableWebSocket   bool `json:"enable_websocket" yaml:"enable_websocket"`
	EnableFileUpload  bool `json:"enable_file_upload" yaml:"enable_file_upload"`
	EnableNotifications bool `json:"enable_notifications" yaml:"enable_notifications"`
	EnableSearch      bool `json:"enable_search" yaml:"enable_search"`
	EnableCaching     bool `json:"enable_caching" yaml:"enable_caching"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Environment: "development",
		Server: ServerConfig{
			Host:           "localhost",
			Port:           8080,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1MB
			CORS: CORSConfig{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders: []string{"*"},
				MaxAge:         86400,
			},
		},
		Database: DatabaseConfig{
			Provider:        "postgres",
			Host:           "localhost",
			Port:           5432,
			MaxOpenConns:   25,
			MaxIdleConns:   10,
			ConnMaxLifetime: 5 * time.Minute,
			ConnMaxIdleTime: 5 * time.Minute,
			Migrations: MigrationConfig{
				AutoMigrate:    true,
				BackupBefore:   false,
				Versioned:      true,
				MigrationsPath: "./migrations",
			},
		},
		Auth: AuthConfig{
			Provider:       "jwt",
			TokenExpiry:    24 * time.Hour,
			RefreshExpiry:  7 * 24 * time.Hour,
			PasswordMinLen: 8,
		},
		Storage: StorageConfig{
			Provider:    "local",
			LocalPath:   "./uploads",
			MaxFileSize: 10 << 20, // 10MB
			AllowedTypes: []string{
				"image/jpeg", "image/png", "image/gif",
				"application/pdf", "text/plain",
			},
		},
		Cache: CacheConfig{
			Provider: "memory",
			TTL:      1 * time.Hour,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		},
		Monitoring: MonitoringConfig{
			EnableMetrics: true,
			EnableHealth:  true,
			MetricsPort:   9090,
			Prometheus: PrometheusConfig{
				Path: "/metrics",
			},
		},
		Security: SecurityConfig{
			EnableRateLimiting: true,
			RateLimit: RateLimitConfig{
				RequestsPerSecond: 100,
				BurstSize:         200,
				CleanupInterval:   1 * time.Minute,
			},
			TrustedProxies: []string{"127.0.0.1"},
		},
		Features: FeatureConfig{
			EnableGraphQL:    false,
			EnableWebSocket:  false,
			EnableFileUpload: true,
			EnableCaching:    true,
		},
	}
}