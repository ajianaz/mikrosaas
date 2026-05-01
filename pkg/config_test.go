package mikrosaas

import (
	"testing"
	"time"
)

func TestSetDefaults_ZeroConfig(t *testing.T) {
	var c Config
	c.SetDefaults()

	assertEqual(t, "Host", "0.0.0.0", c.Host)
	assertEqual(t, "Port", 8080, c.Port)
	assertEqual(t, "ReadTimeout", 15*time.Second, c.ReadTimeout)
	assertEqual(t, "WriteTimeout", 15*time.Second, c.WriteTimeout)
	assertEqual(t, "IdleTimeout", 60*time.Second, c.IdleTimeout)
	assertEqual(t, "MaxOpenConns", 25, c.MaxOpenConns)
	assertEqual(t, "MaxIdleConns", 5, c.MaxIdleConns)
	assertEqual(t, "ConnMaxLifetime", 5*time.Minute, c.ConnMaxLifetime)
	assertEqual(t, "JWTAccessTokenTTL", 15*time.Minute, c.JWTAccessTokenTTL)
	assertEqual(t, "JWTRefreshTokenTTL", 7*24*time.Hour, c.JWTRefreshTokenTTL)
	assertEqual(t, "BcryptCost", 12, c.BcryptCost)
	assertEqual(t, "AppName", "mikrosaas", c.AppName)
	assertEqual(t, "Env", "development", c.Env)
}

func TestSetDefaults_PartialConfig(t *testing.T) {
	c := Config{
		Host:       "127.0.0.1",
		Port:       3000,
		BcryptCost: 10,
	}
	c.SetDefaults()

	// Pre-set values should be preserved
	assertEqual(t, "Host preserved", "127.0.0.1", c.Host)
	assertEqual(t, "Port preserved", 3000, c.Port)
	assertEqual(t, "BcryptCost preserved", 10, c.BcryptCost)

	// Missing values should get defaults
	assertEqual(t, "ReadTimeout default", 15*time.Second, c.ReadTimeout)
	assertEqual(t, "AppName default", "mikrosaas", c.AppName)
}

func TestSetDefaults_EnvLowerCased(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want string
	}{
		{"uppercase", "PRODUCTION", "production"},
		{"mixed case", "Staging", "staging"},
		{"already lowercase", "development", "development"},
		{"empty gets default", "", "development"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{Env: tt.env}
			c.SetDefaults()
			assertEqual(t, "Env", tt.want, c.Env)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errSub  string
	}{
		{
			name:    "missing DatabaseURL",
			config:  Config{JWTSecret: "secret", BcryptCost: 12},
			wantErr: true,
			errSub:  "DATABASE_URL is required",
		},
		{
			name:    "missing JWTSecret",
			config:  Config{DatabaseURL: "postgres://localhost", BcryptCost: 12},
			wantErr: true,
			errSub:  "JWT_SECRET is required",
		},
		{
			name:    "bcrypt cost too low",
			config:  Config{DatabaseURL: "postgres://localhost", JWTSecret: "secret", BcryptCost: 3},
			wantErr: true,
			errSub:  "BCRYPT_COST must be between",
		},
		{
			name:    "bcrypt cost too high",
			config:  Config{DatabaseURL: "postgres://localhost", JWTSecret: "secret", BcryptCost: 32},
			wantErr: true,
			errSub:  "BCRYPT_COST must be between",
		},
		{
			name: "bcrypt cost min boundary (4)",
			config: Config{
				DatabaseURL: "postgres://localhost",
				JWTSecret:   "secret",
				BcryptCost:  4,
			},
			wantErr: false,
		},
		{
			name: "bcrypt cost max boundary (31)",
			config: Config{
				DatabaseURL: "postgres://localhost",
				JWTSecret:   "secret",
				BcryptCost:  31,
			},
			wantErr: false,
		},
		{
			name: "valid config with bcrypt 10",
			config: Config{
				DatabaseURL: "postgres://localhost/db",
				JWTSecret:   "mysecret",
				BcryptCost:  10,
			},
			wantErr: false,
		},
		{
			name: "fully valid config",
			config: Config{
				DatabaseURL: "postgres://localhost/db",
				JWTSecret:   "mysecret",
				BcryptCost:  12,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errSub != "" {
					assertContains(t, "error message", err.Error(), tt.errSub)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
			}
		})
	}
}

func TestAddr(t *testing.T) {
	tests := []struct {
		name string
		host string
		port int
		want string
	}{
		{"default", "0.0.0.0", 8080, "0.0.0.0:8080"},
		{"custom", "127.0.0.1", 3000, "127.0.0.1:3000"},
		{"localhost", "localhost", 443, "localhost:443"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{Host: tt.host, Port: tt.port}
			assertEqual(t, "Addr", tt.want, c.Addr())
		})
	}
}

func TestIsDevelopment(t *testing.T) {
	tests := []struct {
		env  string
		want bool
	}{
		{"development", true},
		{"staging", true},
		{"production", false},
	}
	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			c := Config{Env: tt.env}
			assertEqual(t, "IsDevelopment", tt.want, c.IsDevelopment())
		})
	}
}

func TestIsProduction(t *testing.T) {
	tests := []struct {
		env  string
		want bool
	}{
		{"production", true},
		{"development", false},
		{"staging", false},
	}
	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			c := Config{Env: tt.env}
			assertEqual(t, "IsProduction", tt.want, c.IsProduction())
		})
	}
}

// --- helpers ---

func assertEqual[T comparable](t *testing.T, field string, want, got T) {
	t.Helper()
	if want != got {
		t.Errorf("%s: want %v, got %v", field, want, got)
	}
}

func assertContains(t *testing.T, field, s, sub string) {
	t.Helper()
	if sub != "" {
		// simple substring check
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return
			}
		}
		t.Errorf("%s: expected to contain %q, got %q", field, sub, s)
	}
}
