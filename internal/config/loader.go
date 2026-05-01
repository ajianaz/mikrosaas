package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Loader reads configuration from environment variables.
// No Viper, no YAML, no TOML — just env vars. Simple and testable.
type Loader struct {
	prefix string // e.g. "MIKROSAAS_" — optional prefix for all vars
}

// NewLoader creates a config loader with an optional env var prefix.
func NewLoader(prefix ...string) *Loader {
	l := &Loader{}
	if len(prefix) > 0 {
		l.prefix = prefix[0]
	}
	return l
}

// String reads an environment variable or returns the default value.
func (l *Loader) String(key, fallback string) string {
	if v := os.Getenv(l.key(key)); v != "" {
		return v
	}
	return fallback
}

// Int reads an environment variable as integer or returns the default value.
func (l *Loader) Int(key string, fallback int) int {
	if v := os.Getenv(l.key(key)); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

// Duration reads an environment variable as a Go duration string
// or returns the default value. Supports formats: "15s", "5m", "2h".
func (l *Loader) Duration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(l.key(key)); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

// Bool reads an environment variable as boolean.
func (l *Loader) Bool(key string, fallback bool) bool {
	if v := os.Getenv(l.key(key)); v != "" {
		switch strings.ToLower(v) {
		case "1", "true", "yes", "on":
			return true
		case "0", "false", "no", "off":
			return false
		}
	}
	return fallback
}

// Required returns the env var value or an error if empty.
func (l *Loader) Required(key string) (string, error) {
	v := os.Getenv(l.key(key))
	if v == "" {
		return "", fmt.Errorf("required environment variable %s is not set", l.key(key))
	}
	return v, nil
}

// MustString panics if the env var is not set.
func (l *Loader) MustString(key string) string {
	v, err := l.Required(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (l *Loader) key(name string) string {
	return l.prefix + name
}
