package rbac

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRoleJSONMarshaling(t *testing.T) {
	r := Role{
		ID:          "role_1",
		TenantID:    "tnt_1",
		Name:        "admin",
		Description: "Administrator role",
		IsSystem:    true,
		CreatedAt:   "2025-01-01",
		UpdatedAt:   "2025-01-01",
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	s := string(data)
	for _, field := range []string{"id", "tenant_id", "name", "description", "is_system", "created_at", "updated_at"} {
		if !strings.Contains(s, field) {
			t.Errorf("expected field %q in Role JSON", field)
		}
	}
}

func TestPermissionJSONMarshaling(t *testing.T) {
	p := Permission{
		ID:          "perm_1",
		Name:        "users:read",
		Resource:    "users",
		Action:      "read",
		Description: "Read users",
		CreatedAt:   "2025-01-01",
	}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	s := string(data)
	for _, field := range []string{"id", "name", "resource", "action", "description", "created_at"} {
		if !strings.Contains(s, field) {
			t.Errorf("expected field %q in Permission JSON", field)
		}
	}
}

func TestRolePermissionJSONMarshaling(t *testing.T) {
	rp := RolePermission{
		RoleID:       "role_1",
		PermissionID: "perm_1",
	}
	data, err := json.Marshal(rp)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	s := string(data)
	for _, field := range []string{"role_id", "permission_id"} {
		if !strings.Contains(s, field) {
			t.Errorf("expected field %q in RolePermission JSON", field)
		}
	}
}

func TestUserRoleJSONMarshaling(t *testing.T) {
	ur := UserRole{
		UserID: "usr_1",
		RoleID: "role_1",
	}
	data, err := json.Marshal(ur)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	s := string(data)
	for _, field := range []string{"user_id", "role_id"} {
		if !strings.Contains(s, field) {
			t.Errorf("expected field %q in UserRole JSON", field)
		}
	}
}

func TestRoleJSONRoundTrip(t *testing.T) {
	original := Role{
		ID:          "role_abc",
		TenantID:    "tnt_xyz",
		Name:        "editor",
		Description: "Can edit content",
		IsSystem:    false,
		CreatedAt:   "2025-06-01",
		UpdatedAt:   "2025-06-01",
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	var decoded Role
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if decoded != original {
		t.Errorf("round-trip mismatch: got %+v, want %+v", decoded, original)
	}
}

func TestPermissionJSONRoundTrip(t *testing.T) {
	original := Permission{
		ID:          "perm_abc",
		Name:        "tenants:write",
		Resource:    "tenants",
		Action:      "write",
		Description: "Write tenants",
		CreatedAt:   "2025-06-01",
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	var decoded Permission
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if decoded != original {
		t.Errorf("round-trip mismatch: got %+v, want %+v", decoded, original)
	}
}
