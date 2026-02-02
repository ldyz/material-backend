package db

import (
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New creates a DB connection. Only supports PostgreSQL and MySQL as requested.
// No fallback database will be used if DSN is invalid.
func New(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, errors.New("invalid DSN: empty connection string")
	}

	// Check for PostgreSQL connection string
	if (len(dsn) >= 8 && dsn[:8] == "postgres://") || 
	   (len(dsn) >= 9 && dsn[:9] == "postgresql://") ||
	   (contains(dsn, "host=") && contains(dsn, "dbname=") && contains(dsn, "user=")) {
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	// Check for MySQL connection string
	if contains(dsn, "tcp(") || contains(dsn, "@tcp") || contains(dsn, ":") {
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	// If none of the above, return invalid DSN error
	return nil, errors.New("invalid DSN: unsupported database type or format")
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}