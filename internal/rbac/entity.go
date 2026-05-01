package rbac

// Role represents a role that can be assigned to users.
type Role struct {
	ID          string `json:"id" db:"id"`
	TenantID    string `json:"tenant_id" db:"tenant_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	IsSystem    bool   `json:"is_system" db:"is_system"` // System roles cannot be deleted
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}

// Permission represents a specific action on a resource.
type Permission struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`           // e.g. "users:read", "users:write"
	Resource    string `json:"resource" db:"resource"`   // e.g. "users", "tenants"
	Action      string `json:"action" db:"action"`       // e.g. "read", "write", "delete"
	Description string `json:"description" db:"description"`
	CreatedAt   string `json:"created_at" db:"created_at"`
}

// RolePermission maps permissions to roles.
type RolePermission struct {
	RoleID       string `json:"role_id" db:"role_id"`
	PermissionID string `json:"permission_id" db:"permission_id"`
}

// UserRole assigns a role to a user within a tenant.
type UserRole struct {
	UserID string `json:"user_id" db:"user_id"`
	RoleID string `json:"role_id" db:"role_id"`
}
