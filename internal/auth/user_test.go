package auth

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDisplayName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		email     string
		want      string
	}{
		{"first and last", "John", "Doe", "john@example.com", "John Doe"},
		{"empty names falls back to email", "", "", "john@example.com", "john@example.com"},
		{"only first name", "John", "", "john@example.com", "John "},
		{"only last name", "", "Doe", "john@example.com", " Doe"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
				Email:     tt.email,
			}
			got := u.DisplayName()
			if got != tt.want {
				t.Errorf("DisplayName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUserJSONMarshaling(t *testing.T) {
	u := User{
		ID:            "usr_1",
		TenantID:      "tnt_1",
		Email:         "john@example.com",
		PasswordHash:  "$2b$12$hashed",
		FirstName:     "John",
		LastName:      "Doe",
		IsActive:      true,
		EmailVerified: true,
		CreatedAt:     "2025-01-01T00:00:00Z",
		UpdatedAt:     "2025-01-01T00:00:00Z",
	}
	data, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}
	s := string(data)

	// password_hash must NOT appear in JSON (json:"-")
	if strings.Contains(s, "password_hash") {
		t.Error("password_hash should be hidden in JSON output")
	}
	if strings.Contains(s, "$2b$12$hashed") {
		t.Error("password hash value should not appear in JSON")
	}

	// Other fields must be present
	for _, field := range []string{"id", "tenant_id", "email", "first_name", "last_name", "is_active", "email_verified", "created_at", "updated_at"} {
		if !strings.Contains(s, field) {
			t.Errorf("expected field %q in JSON, got: %s", field, s)
		}
	}
}

func TestUserJSONRoundTrip(t *testing.T) {
	original := User{
		ID:            "usr_123",
		TenantID:      "tnt_456",
		Email:         "test@test.com",
		PasswordHash:  "should-be-hidden",
		FirstName:     "Jane",
		LastName:      "Smith",
		IsActive:      false,
		EmailVerified: false,
		CreatedAt:     "2025-06-01",
		UpdatedAt:     "2025-06-01",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded User
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	// password_hash should be empty after round-trip (json:"-")
	if decoded.PasswordHash != "" {
		t.Errorf("PasswordHash should be empty after round-trip, got %q", decoded.PasswordHash)
	}
	if decoded.ID != original.ID {
		t.Errorf("ID mismatch: got %q, want %q", decoded.ID, original.ID)
	}
	if decoded.Email != original.Email {
		t.Errorf("Email mismatch: got %q, want %q", decoded.Email, original.Email)
	}
}

func TestAuthEntityJSONMarshaling(t *testing.T) {
	t.Run("TokenPair", func(t *testing.T) {
		tp := TokenPair{
			AccessToken:  "access-abc",
			RefreshToken: "refresh-xyz",
			ExpiresIn:    900,
			TokenType:    "Bearer",
		}
		data, err := json.Marshal(tp)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		s := string(data)
		for _, field := range []string{"access_token", "refresh_token", "expires_in", "token_type"} {
			if !strings.Contains(s, field) {
				t.Errorf("expected field %q in JSON", field)
			}
		}
	})
	t.Run("LoginRequest", func(t *testing.T) {
		lr := LoginRequest{Email: "user@test.com", Password: "password123"}
		data, err := json.Marshal(lr)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		if !strings.Contains(string(data), "email") || !strings.Contains(string(data), "password") {
			t.Error("expected email and password fields")
		}
	})
	t.Run("RefreshRequest", func(t *testing.T) {
		rr := RefreshRequest{RefreshToken: "ref-token"}
		data, err := json.Marshal(rr)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		if !strings.Contains(string(data), "refresh_token") {
			t.Error("expected refresh_token field")
		}
	})
	t.Run("LoginResponse", func(t *testing.T) {
		lr := LoginResponse{
			User:  User{ID: "usr_1", Email: "test@test.com"},
			Token: TokenPair{AccessToken: "acc", RefreshToken: "ref", ExpiresIn: 900, TokenType: "Bearer"},
		}
		data, err := json.Marshal(lr)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		s := string(data)
		for _, field := range []string{"user", "token", "access_token"} {
			if !strings.Contains(s, field) {
				t.Errorf("expected field %q in LoginResponse JSON", field)
			}
		}
	})
}
