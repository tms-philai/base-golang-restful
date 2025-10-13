package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	FirstName string         `gorm:"type:varchar(100)" json:"first_name"`
	LastName  string         `gorm:"type:varchar(100)" json:"last_name"`
	IsActive  bool           `gorm:"default:true;not null" json:"is_active"`
	Roles     []Role         `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

func (u *User) HasAnyRole(roleNames []string) bool {
	roleMap := make(map[string]bool)
	for _, role := range u.Roles {
		roleMap[role.Name] = true
	}

	for _, name := range roleNames {
		if roleMap[name] {
			return true
		}
	}
	return false
}

func (u *User) HasAllRoles(roleNames []string) bool {
	roleMap := make(map[string]bool)
	for _, role := range u.Roles {
		roleMap[role.Name] = true
	}

	for _, name := range roleNames {
		if !roleMap[name] {
			return false
		}
	}
	return true
}

func (u *User) HasPermission(permissionName string) bool {
	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}
		if role.HasPermission(permissionName) {
			return true
		}
	}
	return false
}

func (u *User) HasAnyPermission(permissionNames []string) bool {
	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}
		if role.HasAnyPermission(permissionNames) {
			return true
		}
	}
	return false
}

func (u *User) HasAllPermissions(permissionNames []string) bool {
	allPerms := make(map[string]bool)
	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}
		for _, perm := range role.Permissions {
			allPerms[perm.Name] = true
		}
	}

	for _, name := range permissionNames {
		if !allPerms[name] {
			return false
		}
	}
	return true
}

func (u *User) GetAllPermissions() []Permission {
	permMap := make(map[uuid.UUID]Permission)

	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}
		for _, perm := range role.Permissions {
			permMap[perm.ID] = perm
		}
	}

	permissions := make([]Permission, 0, len(permMap))
	for _, perm := range permMap {
		permissions = append(permissions, perm)
	}

	return permissions
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID.String(),
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
