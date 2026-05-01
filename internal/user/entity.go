package user

// Entity represents a user in the system.
// This is the user domain package — auth/user.go contains the auth-specific view.
type Entity struct {
	ID            string `json:"id" db:"id"`
	TenantID      string `json:"tenant_id" db:"tenant_id"`
	Email         string `json:"email" db:"email"`
	PasswordHash  string `json:"-" db:"password_hash"`
	FirstName     string `json:"first_name" db:"first_name"`
	LastName      string `json:"last_name" db:"last_name"`
	IsActive      bool   `json:"is_active" db:"is_active"`
	EmailVerified bool   `json:"email_verified" db:"email_verified"`
	CreatedAt     string `json:"created_at" db:"created_at"`
	UpdatedAt     string `json:"updated_at" db:"updated_at"`
}

// DisplayName returns the user's full name or email fallback.
func (u *Entity) DisplayName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Email
	}
	return u.FirstName + " " + u.LastName
}

// CreateRequest is the payload for creating a new user.
type CreateRequest struct {
	Email     string `json:"email"     validate:"required,email"`
	Password  string `json:"password"  validate:"required,min=8,max=128"`
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name"  validate:"required,max=100"`
}

// UpdateRequest is the payload for updating a user.
type UpdateRequest struct {
	FirstName *string `json:"first_name" validate:"omitempty,max=100"`
	LastName  *string `json:"last_name"  validate:"omitempty,max=100"`
	IsActive  *bool   `json:"is_active"`
}

// ListRequest is the payload for listing users with pagination.
type ListRequest struct {
	Page     int    `json:"page"     validate:"omitempty,min=1"`
	PerPage  int    `json:"per_page" validate:"omitempty,min=1,max=100"`
	Search   string `json:"search"   validate:"omitempty"`
	IsActive *bool  `json:"is_active"`
}
