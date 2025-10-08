package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	DisplayName string         `gorm:"type:varchar(100);not null" json:"display_name"`
	Description string         `gorm:"type:text" json:"description"`
	IsSystem    bool           `gorm:"default:false;not null" json:"is_system"`
	IsActive    bool           `gorm:"default:true;not null" json:"is_active"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User         `gorm:"many2many:user_roles;" json:"users,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (r *Role) HasPermission(permissionName string) bool {
	for _, perm := range r.Permissions {
		if perm.Name == permissionName {
			return true
		}
	}
	return false
}

func (r *Role) HasAnyPermission(permissionNames []string) bool {
	permMap := make(map[string]bool)
	for _, perm := range r.Permissions {
		permMap[perm.Name] = true
	}

	for _, name := range permissionNames {
		if permMap[name] {
			return true
		}
	}
	return false
}

func (r *Role) HasAllPermissions(permissionNames []string) bool {
	permMap := make(map[string]bool)
	for _, perm := range r.Permissions {
		permMap[perm.Name] = true
	}

	for _, name := range permissionNames {
		if !permMap[name] {
			return false
		}
	}
	return true
}
