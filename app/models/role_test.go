package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRole_HasPermission(t *testing.T) {
	role := &Role{
		ID:   uuid.New(),
		Name: "admin",
		Permissions: []Permission{
			{Name: "user.create"},
			{Name: "user.read"},
			{Name: "user.update"},
		},
	}

	tests := []struct {
		name           string
		permissionName string
		expected       bool
	}{
		{
			name:           "has permission",
			permissionName: "user.create",
			expected:       true,
		},
		{
			name:           "does not have permission",
			permissionName: "user.delete",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := role.HasPermission(tt.permissionName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRole_HasAnyPermission(t *testing.T) {
	role := &Role{
		ID:   uuid.New(),
		Name: "admin",
		Permissions: []Permission{
			{Name: "user.create"},
			{Name: "user.read"},
		},
	}

	tests := []struct {
		name            string
		permissionNames []string
		expected        bool
	}{
		{
			name:            "has one permission",
			permissionNames: []string{"user.create", "user.delete"},
			expected:        true,
		},
		{
			name:            "has no permissions",
			permissionNames: []string{"user.delete", "user.manage"},
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := role.HasAnyPermission(tt.permissionNames)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRole_HasAllPermissions(t *testing.T) {
	role := &Role{
		ID:   uuid.New(),
		Name: "admin",
		Permissions: []Permission{
			{Name: "user.create"},
			{Name: "user.read"},
			{Name: "user.update"},
		},
	}

	tests := []struct {
		name            string
		permissionNames []string
		expected        bool
	}{
		{
			name:            "has all permissions",
			permissionNames: []string{"user.create", "user.read"},
			expected:        true,
		},
		{
			name:            "missing one permission",
			permissionNames: []string{"user.create", "user.delete"},
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := role.HasAllPermissions(tt.permissionNames)
			assert.Equal(t, tt.expected, result)
		})
	}
}
