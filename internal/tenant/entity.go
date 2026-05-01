package tenant

// Entity represents a tenant (organization/customer) in the system.
type Entity struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Slug      string `json:"slug" db:"slug"`
	Domain    string `json:"domain" db:"domain"`
	IsActive  bool   `json:"is_active" db:"is_active"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

// CreateRequest is the payload for creating a new tenant.
type CreateRequest struct {
	Name   string `json:"name"   validate:"required,min=1,max=200"`
	Slug   string `json:"slug"   validate:"required,min=1,max=100,slug"`
	Domain string `json:"domain" validate:"omitempty,max=253"`
}

// UpdateRequest is the payload for updating a tenant.
type UpdateRequest struct {
	Name     *string `json:"name"   validate:"omitempty,min=1,max=200"`
	Slug     *string `json:"slug"   validate:"omitempty,min=1,max=100,slug"`
	Domain   *string `json:"domain" validate:"omitempty,max=253"`
	IsActive *bool   `json:"is_active"`
}
