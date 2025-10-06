package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username  string    `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Email     string    `json:"email" binding:"required,email" example:"john@example.com"`
	Password  string    `json:"password,omitempty" binding:"required,min=6" example:"password123"`
	FirstName string    `json:"first_name" binding:"required" example:"John"`
	LastName  string    `json:"last_name" binding:"required" example:"Doe"`
	Role      string    `json:"role" example:"user"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// UserCreateRequest represents the request payload for creating a user
type UserCreateRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Email     string `json:"email" binding:"required,email" example:"john@example.com"`
	Password  string `json:"password" binding:"required,min=6" example:"password123"`
	FirstName string `json:"first_name" binding:"required" example:"John"`
	LastName  string `json:"last_name" binding:"required" example:"Doe"`
}

// UserUpdateRequest represents the request payload for updating a user
type UserUpdateRequest struct {
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=50" example:"johndoe"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email" example:"john@example.com"`
	FirstName *string `json:"first_name,omitempty" example:"John"`
	LastName  *string `json:"last_name,omitempty" example:"Doe"`
	IsActive  *bool   `json:"is_active,omitempty" example:"true"`
}

// UserResponse represents the response payload for user operations (without password)
type UserResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john@example.com"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
	Role      string    `json:"role" example:"user"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// NewUser creates a new user with generated ID and timestamps
func NewUser(req UserCreateRequest) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password, // Will be hashed by service layer
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user", // Default role
		IsActive:  true,   // Default active status
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ToResponse converts User to UserResponse (removes sensitive data)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// Update applies updates from UserUpdateRequest to the user
func (u *User) Update(req UserUpdateRequest) {
	if req.Username != nil {
		u.Username = *req.Username
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.FirstName != nil {
		u.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		u.LastName = *req.LastName
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}
	u.UpdatedAt = time.Now()
}
