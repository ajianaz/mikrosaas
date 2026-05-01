package tenant

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEntityJSONMarshaling(t *testing.T) {
	e := Entity{
		ID:        "tnt_1",
		Name:      "Acme Corp",
		Slug:      "acme-corp",
		Domain:    "acme.example.com",
		IsActive:  true,
		CreatedAt: "2025-01-01",
		UpdatedAt: "2025-01-01",
	}
	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	s := string(data)
	for _, field := range []string{"id", "name", "slug", "domain", "is_active", "created_at", "updated_at"} {
		if !strings.Contains(s, field) {
			t.Errorf("expected field %q in Entity JSON", field)
		}
	}
}

func TestCreateRequestJSON(t *testing.T) {
	cr := CreateRequest{
		Name:   "New Tenant",
		Slug:   "new-tenant",
		Domain: "new.example.com",
	}
	data, err := json.Marshal(cr)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	s := string(data)
	for _, field := range []string{"name", "slug", "domain"} {
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
		if !strings.Contains(s, "null") {
			t.Errorf("nil pointer fields should produce null, got: %s", s)
		}
	})
	t.Run("non-nil pointers present", func(t *testing.T) {
		name := "Updated Name"
		active := false
		ur := UpdateRequest{
			Name:     &name,
			IsActive: &active,
		}
		data, err := json.Marshal(ur)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		s := string(data)
		if !strings.Contains(s, "Updated Name") {
			t.Error("expected Name value in JSON")
		}
		if !strings.Contains(s, "name") {
			t.Error("expected name field in JSON")
		}
		if !strings.Contains(s, "is_active") {
			t.Error("expected is_active field in JSON")
		}
	})
}

func TestCreateRequestJSONRoundTrip(t *testing.T) {
	original := CreateRequest{
		Name:   "Test Tenant",
		Slug:   "test-tenant",
		Domain: "test.example.com",
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	var decoded CreateRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if decoded != original {
		t.Errorf("round-trip mismatch: got %+v, want %+v", decoded, original)
	}
}
