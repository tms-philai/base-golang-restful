package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	DisplayName string         `gorm:"type:varchar(150);not null" json:"display_name"`
	Description string         `gorm:"type:text" json:"description"`
	Resource    string         `gorm:"type:varchar(50);not null;index" json:"resource"`
	Action      string         `gorm:"type:varchar(50);not null;index" json:"action"`
	IsSystem    bool           `gorm:"default:false;not null" json:"is_system"`
	Roles       []Role         `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionList   = "list"
	ActionManage = "manage"
)

const (
	ResourceUser       = "user"
	ResourceRole       = "role"
	ResourcePermission = "permission"
	ResourceProduct    = "product"
	ResourceOrder      = "order"
	ResourceCategory   = "category"
	ResourceAll        = "*"
)

func MakePermissionName(resource, action string) string {
	return resource + "." + action
}
