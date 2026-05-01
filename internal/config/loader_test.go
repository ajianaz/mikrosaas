package config

import (
	"strings"
	"testing"
	"time"
)

func TestNewLoader(t *testing.T) {
	t.Run("no prefix", func(t *testing.T) {
		l := NewLoader()
		if l.prefix != "" {
			t.Errorf("expected empty prefix, got %q", l.prefix)
		}
	})
	t.Run("with prefix", func(t *testing.T) {
		l := NewLoader("PREFIX_")
		if l.prefix != "PREFIX_" {
			t.Errorf("expected %q, got %q", "PREFIX_", l.prefix)
		}
	})
	t.Run("with empty string prefix", func(t *testing.T) {
		l := NewLoader("")
		if l.prefix != "" {
			t.Errorf("expected empty prefix, got %q", l.prefix)
		}
	})
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		envKey   string
		envVal   string
		fallback string
		want     string
	}{
		{"env set", "HOST", "MY_HOST", "localhost", "0.0.0.0", "localhost"},
		{"env not set", "HOST", "MY_HOST", "", "0.0.0.0", "0.0.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoader("MY_")
			t.Setenv(tt.envKey, tt.envVal)
			got := l.String(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("String(%q, %q) = %q, want %q", tt.key, tt.fallback, got, tt.want)
			}
		})
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		name     string
		envVal   string
		fallback int
		want     int
	}{
		{"valid int", "42", 0, 42},
		{"invalid int", "notanumber", 10, 10},
		{"env not set", "", 99, 99},
		{"negative", "-5", 0, -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoader()
			t.Setenv("TEST_PORT", tt.envVal)
			got := l.Int("TEST_PORT", tt.fallback)
			if got != tt.want {
				t.Errorf("Int() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		name     string
		envVal   string
		fallback time.Duration
		want     time.Duration
	}{
		{"seconds", "30s", 0, 30 * time.Second},
		{"minutes", "5m", 0, 5 * time.Minute},
		{"hours", "2h", 0, 2 * time.Hour},
		{"invalid", "forever", 10*time.Second, 10 * time.Second},
		{"env not set", "", 5 * time.Minute, 5 * time.Minute},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoader()
			t.Setenv("TEST_TIMEOUT", tt.envVal)
			got := l.Duration("TEST_TIMEOUT", tt.fallback)
			if got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		name     string
		envVal   string
		fallback bool
		want     bool
	}{
		{"true", "true", false, true},
		{"True", "True", false, true},
		{"1", "1", false, true},
		{"yes", "yes", false, true},
		{"YES", "YES", false, true},
		{"on", "on", false, true},
		{"false", "false", true, false},
		{"0", "0", true, false},
		{"no", "no", true, false},
		{"off", "off", true, false},
		{"invalid returns fallback true", "maybe", true, true},
		{"invalid returns fallback false", "maybe", false, false},
		{"env not set fallback true", "", true, true},
		{"env not set fallback false", "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoader()
			t.Setenv("TEST_BOOL", tt.envVal)
			got := l.Bool("TEST_BOOL", tt.fallback)
			if got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequired(t *testing.T) {
	t.Run("env set", func(t *testing.T) {
		l := NewLoader()
		t.Setenv("TEST_REQUIRED", "somevalue")
		got, err := l.Required("TEST_REQUIRED")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "somevalue" {
			t.Errorf("Required() = %q, want %q", got, "somevalue")
		}
	})
	t.Run("env not set", func(t *testing.T) {
		l := NewLoader()
		_, err := l.Required("TEST_NOTSET_REQUIRED")
		if err == nil {
			t.Fatal("expected error for missing required var")
		}
		if !strings.Contains(err.Error(), "TEST_NOTSET_REQUIRED") {
			t.Errorf("error should mention var name, got: %v", err)
		}
	})
	t.Run("with prefix", func(t *testing.T) {
		l := NewLoader("APP_")
		t.Setenv("APP_SECRET", "s3cret")
		got, err := l.Required("SECRET")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "s3cret" {
			t.Errorf("Required() = %q, want %q", got, "s3cret")
		}
	})
}

func TestMustString(t *testing.T) {
	t.Run("env set", func(t *testing.T) {
		l := NewLoader()
		t.Setenv("TEST_MUST", "value")
		got := l.MustString("TEST_MUST")
		if got != "value" {
			t.Errorf("MustString() = %q, want %q", got, "value")
		}
	})
	t.Run("env not set panics", func(t *testing.T) {
		l := NewLoader()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected MustString to panic, but it didn't")
			}
		}()
		l.MustString("TEST_MUST_NOTSET")
	})
}
