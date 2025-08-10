package models

import (
	"fmt"
	"strconv"
	"strings"
)

// DatabaseProvider represents database provider configuration
type DatabaseProvider struct {
	Type       string `json:"type"`        // "postgres", "mysql", "sqlite", "supabase", "mongodb", "redis"
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Database   string `json:"database"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	SSLMode    string `json:"ssl_mode"`
	URL        string `json:"url"`         // For cloud providers
	ProjectRef string `json:"project_ref"` // For Supabase
	AnonKey    string `json:"anon_key"`    // For Supabase
	ServiceKey string `json:"service_key"` // For Supabase
	JWTSecret  string `json:"jwt_secret"`  // For Supabase JWT
	
	// MongoDB specific fields
	AuthSource string `json:"auth_source"` // MongoDB auth source
	ReplicaSet string `json:"replica_set"` // MongoDB replica set
	
	// Redis specific fields
	MaxRetries int `json:"max_retries"` // Redis max retries
	PoolSize   int `json:"pool_size"`   // Redis connection pool size
}

// SupportedDatabaseTypes returns all supported database types
func SupportedDatabaseTypes() []string {
	return []string{
		"postgres",
		"mysql", 
		"sqlite",
		"supabase",
		"mongodb",
		"redis",
	}
}

// GetDSN returns the database connection string
func (db *DatabaseProvider) GetDSN() string {
	switch db.Type {
	case "postgres":
		return db.getPostgresDSN()
	case "mysql":
		return db.getMySQLDSN()
	case "sqlite":
		return db.getSQLiteDSN()
	case "supabase":
		return db.getSupabaseDSN()
	case "mongodb":
		return db.getMongoDSN()
	case "redis":
		return db.getRedisDSN()
	default:
		return ""
	}
}

// getPostgresDSN returns PostgreSQL connection string
func (db *DatabaseProvider) getPostgresDSN() string {
	if db.URL != "" {
		return db.URL
	}
	return "host=" + db.Host + " user=" + db.Username + " password=" + db.Password + " dbname=" + db.Database + " port=" + strconv.Itoa(db.Port) + " sslmode=" + db.SSLMode
}

// getMySQLDSN returns MySQL connection string
func (db *DatabaseProvider) getMySQLDSN() string {
	if db.URL != "" {
		return db.URL
	}
	return db.Username + ":" + db.Password + "@tcp(" + db.Host + ":" + strconv.Itoa(db.Port) + ")/" + db.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
}

// getSQLiteDSN returns SQLite connection string
func (db *DatabaseProvider) getSQLiteDSN() string {
	return db.Database
}

// getSupabaseDSN returns Supabase connection string
func (db *DatabaseProvider) getSupabaseDSN() string {
	if db.URL != "" {
		return db.URL
	}
	return "postgresql://" + db.Username + ":" + db.Password + "@" + db.Host + ":" + strconv.Itoa(db.Port) + "/" + db.Database + "?sslmode=require"
}

// getMongoDSN returns MongoDB connection string
func (db *DatabaseProvider) getMongoDSN() string {
	if db.URL != "" {
		return db.URL
	}
	dsn := "mongodb://"
	if db.Username != "" && db.Password != "" {
		dsn += db.Username + ":" + db.Password + "@"
	}
	dsn += db.Host + ":" + strconv.Itoa(db.Port) + "/" + db.Database
	if db.AuthSource != "" {
		dsn += "?authSource=" + db.AuthSource
	}
	if db.ReplicaSet != "" {
		separator := "?"
		if strings.Contains(dsn, "?") {
			separator = "&"
		}
		dsn += separator + "replicaSet=" + db.ReplicaSet
	}
	return dsn
}

// getRedisDSN returns Redis connection string
func (db *DatabaseProvider) getRedisDSN() string {
	if db.URL != "" {
		return db.URL
	}
	dsn := "redis://"
	if db.Password != "" {
		dsn += ":" + db.Password + "@"
	}
	dsn += db.Host + ":" + strconv.Itoa(db.Port)
	return dsn
}

// RequiredImports returns the imports needed for the database provider
func (db *DatabaseProvider) RequiredImports() []string {
	var imports []string
	
	switch db.Type {
	case "postgres":
		imports = []string{"gorm.io/gorm", "gorm.io/driver/postgres"}
	case "mysql":
		imports = []string{"gorm.io/gorm", "gorm.io/driver/mysql"}
	case "sqlite":
		imports = []string{"gorm.io/gorm", "gorm.io/driver/sqlite"}
	case "supabase":
		imports = []string{
			"gorm.io/gorm",
			"gorm.io/driver/postgres",
			"github.com/supabase-community/gotrue-go",
			"github.com/supabase-community/storage-go",
		}
	case "mongodb":
		imports = []string{
			"go.mongodb.org/mongo-driver/mongo",
			"go.mongodb.org/mongo-driver/mongo/options",
			"context",
		}
	case "redis":
		imports = []string{
			"github.com/go-redis/redis/v8",
			"context",
		}
	}
	
	return removeDuplicates(imports)
}

// GetDriverName returns the GORM driver name
func (db *DatabaseProvider) GetDriverName() string {
	switch db.Type {
	case "postgres", "supabase":
		return "postgres"
	case "mysql":
		return "mysql"
	case "sqlite":
		return "sqlite"
	case "mongodb":
		return "mongodb"
	case "redis":
		return "redis"
	default:
		return "postgres"
	}
}

// IsCloudProvider returns true if the provider is a cloud service
func (db *DatabaseProvider) IsCloudProvider() bool {
	return db.Type == "supabase"
}

// IsNoSQL returns true if the database is a NoSQL database
func (db *DatabaseProvider) IsNoSQL() bool {
	return db.Type == "mongodb" || db.Type == "redis"
}

// GetDisplayName returns a user-friendly display name for the database type
func (db *DatabaseProvider) GetDisplayName() string {
	switch db.Type {
	case "postgres":
		return "PostgreSQL"
	case "mysql":
		return "MySQL"
	case "sqlite":
		return "SQLite"
	case "supabase":
		return "Supabase (PostgreSQL + Auth + Storage)"
	case "mongodb":
		return "MongoDB"
	case "redis":
		return "Redis"
	default:
		return strings.Title(db.Type)
	}
}

// GetDescription returns a description of the database provider
func (db *DatabaseProvider) GetDescription() string {
	switch db.Type {
	case "postgres":
		return "Powerful, open source object-relational database system"
	case "mysql":
		return "The world's most popular open source database"
	case "sqlite":
		return "Lightweight, file-based SQL database"
	case "supabase":
		return "Open source Firebase alternative with PostgreSQL, Auth, and Storage"
	case "mongodb":
		return "Document-oriented NoSQL database"
	case "redis":
		return "In-memory data structure store for caching and sessions"
	default:
		return "Database provider"
	}
}

// GetEnvironmentVars returns environment variables for the database
func (db *DatabaseProvider) GetEnvironmentVars() map[string]string {
	vars := make(map[string]string)
	
	switch db.Type {
	case "postgres":
		vars["DB_HOST"] = db.Host
		vars["DB_PORT"] = strconv.Itoa(db.Port)
		vars["DB_USER"] = db.Username
		vars["DB_PASSWORD"] = db.Password
		vars["DB_NAME"] = db.Database
		vars["DB_SSLMODE"] = db.SSLMode
	case "mysql":
		vars["DB_HOST"] = db.Host
		vars["DB_PORT"] = strconv.Itoa(db.Port)
		vars["DB_USER"] = db.Username
		vars["DB_PASSWORD"] = db.Password
		vars["DB_NAME"] = db.Database
	case "sqlite":
		vars["DB_PATH"] = db.Database
	case "supabase":
		vars["SUPABASE_URL"] = "https://" + db.ProjectRef + ".supabase.co"
		vars["SUPABASE_ANON_KEY"] = db.AnonKey
		vars["SUPABASE_SERVICE_KEY"] = db.ServiceKey
		vars["SUPABASE_JWT_SECRET"] = db.JWTSecret
		vars["DATABASE_URL"] = db.GetDSN()
	case "mongodb":
		vars["MONGODB_HOST"] = db.Host
		vars["MONGODB_PORT"] = strconv.Itoa(db.Port)
		vars["MONGODB_DATABASE"] = db.Database
		if db.Username != "" {
			vars["MONGODB_USERNAME"] = db.Username
		}
		if db.Password != "" {
			vars["MONGODB_PASSWORD"] = db.Password
		}
		if db.AuthSource != "" {
			vars["MONGODB_AUTH_SOURCE"] = db.AuthSource
		}
		if db.ReplicaSet != "" {
			vars["MONGODB_REPLICA_SET"] = db.ReplicaSet
		}
	case "redis":
		vars["REDIS_HOST"] = db.Host
		vars["REDIS_PORT"] = strconv.Itoa(db.Port)
		if db.Password != "" {
			vars["REDIS_PASSWORD"] = db.Password
		}
		if db.MaxRetries > 0 {
			vars["REDIS_MAX_RETRIES"] = fmt.Sprintf("%d", db.MaxRetries)
		}
		if db.PoolSize > 0 {
			vars["REDIS_POOL_SIZE"] = fmt.Sprintf("%d", db.PoolSize)
		}
	}
	
	return vars
}

// ValidateConfiguration validates the database provider configuration
func (db *DatabaseProvider) ValidateConfiguration() error {
	if db.Type == "" {
		return fmt.Errorf("database type is required")
	}
	
	// Validate supported types
	supported := SupportedDatabaseTypes()
	isSupported := false
	for _, supportedType := range supported {
		if db.Type == supportedType {
			isSupported = true
			break
		}
	}
	if !isSupported {
		return fmt.Errorf("unsupported database type: %s", db.Type)
	}
	
	// Type-specific validation
	switch db.Type {
	case "supabase":
		if db.ProjectRef == "" && db.URL == "" {
			return fmt.Errorf("supabase requires either project_ref or url")
		}
		if db.AnonKey == "" {
			return fmt.Errorf("supabase requires anon_key")
		}
	case "redis":
		if db.Host == "" && db.URL == "" {
			return fmt.Errorf("redis requires either host or url")
		}
	}
	
	return nil
}

