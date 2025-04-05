package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Config holds database connection parameters
type Config struct {
	Host        string
	Port        string
	User        string
	Password    string
	Database    string
	SSLMode     string
	SSLRootCert string

	// Pool settings
	MaxOpenConns    int // Maximum number of open connections
	MaxIdleConns    int // Maximum number of idle connections
	ConnMaxLifetime int // Maximum lifetime of a connection in seconds
}

// DB wraps the database connection pool and sqlc queries
type DB struct {
	Conn    *sql.DB
	Queries *Queries
}

// NewDB creates a new database connection pool and initializes it
func NewDB() (*DB, error) {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	// Create the configuration from environment variables
	cfg := Config{
		Host:            os.Getenv("POSTGRES_HOST"),
		Port:            os.Getenv("POSTGRES_PORT"),
		User:            os.Getenv("POSTGRES_USER"),
		Password:        os.Getenv("POSTGRES_PASSWORD"),
		Database:        os.Getenv("POSTGRES_DATABASE"),
		SSLMode:         os.Getenv("POSTGRES_SSLMODE"),
		SSLRootCert:     os.Getenv("POSTGRES_CERTIFICATE_PATH"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 10),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 300),
	}

	testing_mode := os.Getenv("TESTING_MODE")
	if testing_mode == "true" {
		fmt.Println("Running in testing mode, targeting local database.")
		cfg.Host = "127.0.0.1"
		cfg.Port = "31415"
		cfg.User = "postgres"
		cfg.Password = "postgres"
		cfg.Database = "testdb"
		cfg.SSLMode = ""
		cfg.SSLRootCert = ""
	}

	// Validate that all required configuration fields are set
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.Password == "" || cfg.Database == "" {
		log.Fatalf("One or more required configuration fields are missing")
	}

	// Check if sslmode is set, then CA should be set
	if cfg.SSLMode == "verify-full" && cfg.SSLRootCert == "" {
		log.Fatalf("SSL mode is set to verify-full, but SSLRootCert is not provided")
	}

	// Create the connection string with SSL parameters
	var dsn string
	if cfg.SSLMode != "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode, cfg.SSLRootCert,
		)
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database,
		)
	}

	// Open a new database connection
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(cfg.MaxOpenConns)
	conn.SetMaxIdleConns(cfg.MaxIdleConns)
	conn.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	if err = conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to the database with SSL!")
	fmt.Printf("Connection pool configured: MaxOpen=%d, MaxIdle=%d, ConnMaxLifetime=%d seconds\n",
		cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.ConnMaxLifetime)

	return &DB{
		Conn:    conn,
		Queries: New(conn),
	}, nil
}

// getEnvAsInt reads an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s, using default %d: %v", key, defaultValue, err)
		return defaultValue
	}
	return value
}

// Close closes the database connection pool
func (d *DB) Close() {
	if err := d.Conn.Close(); err != nil {
		log.Printf("Error closing database connection pool: %v", err)
	}
}
