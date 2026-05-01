package user

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEntityDisplayName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		email     string
		want      string
	}{
		{"first and last", "Alice", "Wonderland", "alice@example.com", "Alice Wonderland"},
		{"empty names falls back to email", "", "", "alice@example.com", "alice@example.com"},
		{"only first name", "Bob", "", "bob@example.com", "Bob "},
		{"only last name", "", "Builder", "bob@example.com", " Builder"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entity{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
				Email:     tt.email,
			}
			got := e.DisplayName()
			if got != tt.want {
				t.Errorf("DisplayName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEntityJSONMarshaling(t *testing.T) {
	e := Entity{
		ID:            "usr_1",
		TenantID:      "tnt_1",
		Email:         "user@example.com",
		PasswordHash:  "$2b$12$hashed",
		FirstName:     "Test",
		LastName:      "User",
		IsActive:      true,
		EmailVerified: false,
		CreatedAt:     "2025-01-01",
		UpdatedAt:     "2025-01-01",
	}
	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	s := string(data)

	if strings.Contains(s, "password_hash") {
		t.Error("password_hash should be hidden in JSON")
	}
	for _, field := range []string{"id", "tenant_id", "email", "first_name", "last_name", "is_active", "email_verified"} {
		if !strings.Contains(s, field) {
			t.Errorf("expected field %q in Entity JSON", field)
		}
	}
}

func TestCreateRequestJSON(t *testing.T) {
	cr := CreateRequest{
		Email:     "new@example.com",
		Password:  "securepass123",
		FirstName: "New",
		LastName:  "User",
	}
	data, err := json.Marshal(cr)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	s := string(data)
	for _, field := range []string{"email", "password", "first_name", "last_name"} {
		if !strings.Contains(s, field) {
			t.Errorf("expected field %q in CreateRequest JSON", field)
		}
	}
}

func TestUpdateRequestJSON_PointerFields(t *testing.T) {
	t.Run("nil pointers omitted", func(t *testing.T) {
		ur := UpdateRequest{}
		data, err := json.Marshal(ur)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		s := string(data)
		// nil pointers should be omitted (zero value for *string is nil, marshals to null)
		if !strings.Contains(s, "null") {
			t.Errorf("nil pointer fields should produce null, got: %s", s)
		}
	})
	t.Run("non-nil pointers present", func(t *testing.T) {
		first := "Updated"
		active := true
		ur := UpdateRequest{
			FirstName: &first,
			IsActive:  &active,
		}
		data, err := json.Marshal(ur)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		s := string(data)
		if !strings.Contains(s, "Updated") {
			t.Error("expected FirstName value in JSON")
		}
		if !strings.Contains(s, "first_name") {
			t.Error("expected first_name field in JSON")
		}
	})
}

func TestListRequestJSON(t *testing.T) {
	lr := ListRequest{Page: 2, PerPage: 50, Search: "john"}
	data, err := json.Marshal(lr)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	s := string(data)
	for _, field := range []string{"page", "per_page", "search"} {
		if !strings.Contains(s, field) {
			t.Errorf("expected field %q in ListRequest JSON", field)
		}
	}
}

func TestListRequestJSON_IsActivePointer(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		lr := ListRequest{Page: 1, PerPage: 10}
		data, _ := json.Marshal(lr)
		s := string(data)
		if !strings.Contains(s, "null") {
			t.Errorf("nil is_active should produce null, got: %s", s)
		}
	})
	t.Run("non-nil pointer", func(t *testing.T) {
		active := false
		lr := ListRequest{IsActive: &active}
		data, _ := json.Marshal(lr)
		s := string(data)
		if !strings.Contains(s, "is_active") {
			t.Error("expected is_active field")
		}
	})
}
