package auth

// Entity types for the auth domain.

// TokenPair represents an access + refresh token pair returned after login.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	TokenType    string `json:"token_type"`  // "Bearer"
}

// LoginRequest is the payload for user login.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RefreshRequest is the payload for token refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LoginResponse is the response after successful login.
type LoginResponse struct {
	User  User       `json:"user"`
	Token TokenPair  `json:"token"`
}
