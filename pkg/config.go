package mikrosaas

import (
	"fmt"
	"strings"
	"time"
)

// Config holds all configuration for the SaaS application.
// Values are loaded from environment variables or config file via stdlib.
type Config struct {
	// Server
	Host         string        // SERVER_HOST, default "0.0.0.0"
	Port         int           // SERVER_PORT, default 8080
	ReadTimeout  time.Duration // SERVER_READ_TIMEOUT, default 15s
	WriteTimeout time.Duration // SERVER_WRITE_TIMEOUT, default 15s
	IdleTimeout  time.Duration // SERVER_IDLE_TIMEOUT, default 60s

	// Database
	DatabaseURL           string // DATABASE_URL, required
	MaxOpenConns          int    // DB_MAX_OPEN_CONNS, default 25
	MaxIdleConns          int    // DB_MAX_IDLE_CONNS, default 5
	ConnMaxLifetime       time.Duration // DB_CONN_MAX_LIFETIME, default 5m

	// Auth
	JWTSecret          string        // JWT_SECRET, required
	JWTAccessTokenTTL  time.Duration // JWT_ACCESS_TOKEN_TTL, default 15m
	JWTRefreshTokenTTL time.Duration // JWT_REFRESH_TOKEN_TTL, default 7d
	BcryptCost         int           // BCRYPT_COST, default 12

	// App
	AppName string // APP_NAME, default "mikrosaas"
	Env     string // APP_ENV, default "development" (development|staging|production)
}

// SetDefaults fills zero values with sensible defaults.
func (c *Config) SetDefaults() {
	if c.Host == "" {
		c.Host = "0.0.0.0"
	}
	if c.Port == 0 {
		c.Port = 8080
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 15 * time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 15 * time.Second
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = 60 * time.Second
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 25
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 5
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = 5 * time.Minute
	}
	if c.JWTAccessTokenTTL == 0 {
		c.JWTAccessTokenTTL = 15 * time.Minute
	}
	if c.JWTRefreshTokenTTL == 0 {
		c.JWTRefreshTokenTTL = 7 * 24 * time.Hour
	}
	if c.BcryptCost == 0 {
		c.BcryptCost = 12
	}
	if c.AppName == "" {
		c.AppName = "mikrosaas"
	}
	if c.Env == "" {
		c.Env = "development"
	}
	c.Env = strings.ToLower(c.Env)
}

// Validate checks that required fields are set.
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("config: DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("config: JWT_SECRET is required")
	}
	if c.BcryptCost < 4 || c.BcryptCost > 31 {
		return fmt.Errorf("config: BCRYPT_COST must be between 4 and 31")
	}
	return nil
}

// Addr returns host:port for http.Server.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsDevelopment returns true in non-production environments.
func (c *Config) IsDevelopment() bool {
	return c.Env == "development" || c.Env == "staging"
}

// IsProduction returns true when APP_ENV=production.
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}
