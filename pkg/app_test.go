package mikrosaas

import (
	"strings"
	"testing"
)

func validConfig() Config {
	return Config{
		DatabaseURL: "postgres://localhost/test",
		JWTSecret:   "test-secret-key",
		Env:         "test", // non-prod, non-dev to avoid logger quirks
	}
}

func TestNewApp_ValidConfig(t *testing.T) {
	cfg := validConfig()
	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	if app == nil {
		t.Fatal("NewApp() returned nil app")
	}
	if app.Handler() == nil {
		t.Error("Handler() returned nil")
	}
	got := app.Config()
	if got.DatabaseURL != cfg.DatabaseURL {
		t.Errorf("Config().DatabaseURL = %q, want %q", got.DatabaseURL, cfg.DatabaseURL)
	}
	if got.JWTSecret != cfg.JWTSecret {
		t.Errorf("Config().JWTSecret = %q, want %q", got.JWTSecret, cfg.JWTSecret)
	}
}

func TestNewApp_MissingDatabaseURL(t *testing.T) {
	cfg := Config{
		JWTSecret: "test-secret-key",
	}
	_, err := NewApp(cfg)
	if err == nil {
		t.Fatal("expected error for missing DatabaseURL, got nil")
	}
	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("error should mention DATABASE_URL, got: %v", err)
	}
}

func TestNewApp_MissingJWTSecret(t *testing.T) {
	cfg := Config{
		DatabaseURL: "postgres://localhost/test",
	}
	_, err := NewApp(cfg)
	if err == nil {
		t.Fatal("expected error for missing JWTSecret, got nil")
	}
	if !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Errorf("error should mention JWT_SECRET, got: %v", err)
	}
}
