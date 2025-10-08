package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser_HasRole(t *testing.T) {
	user := &User{
		ID: uuid.New(),
		Roles: []Role{
			{Name: "admin"},
			{Name: "editor"},
		},
	}

	tests := []struct {
		name     string
		roleName string
		expected bool
	}{
		{
			name:     "has role",
			roleName: "admin",
			expected: true,
		},
		{
			name:     "does not have role",
			roleName: "viewer",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.HasRole(tt.roleName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_HasAnyRole(t *testing.T) {
	user := &User{
		ID: uuid.New(),
		Roles: []Role{
			{Name: "admin"},
			{Name: "editor"},
		},
	}

	tests := []struct {
		name      string
		roleNames []string
		expected  bool
	}{
		{
			name:      "has one role",
			roleNames: []string{"admin", "viewer"},
			expected:  true,
		},
		{
			name:      "has no roles",
			roleNames: []string{"viewer", "guest"},
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.HasAnyRole(tt.roleNames)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_HasAllRoles(t *testing.T) {
	user := &User{
		ID: uuid.New(),
		Roles: []Role{
			{Name: "admin"},
			{Name: "editor"},
		},
	}

	tests := []struct {
		name      string
		roleNames []string
		expected  bool
	}{
		{
			name:      "has all roles",
			roleNames: []string{"admin", "editor"},
			expected:  true,
		},
		{
			name:      "missing one role",
			roleNames: []string{"admin", "viewer"},
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.HasAllRoles(tt.roleNames)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_HasPermission(t *testing.T) {
	user := &User{
		ID: uuid.New(),
		Roles: []Role{
			{
				Name:     "admin",
				IsActive: true,
				Permissions: []Permission{
					{Name: "user.create"},
					{Name: "user.read"},
				},
			},
			{
				Name:     "editor",
				IsActive: false,
				Permissions: []Permission{
					{Name: "post.create"},
				},
			},
		},
	}

	tests := []struct {
		name           string
		permissionName string
		expected       bool
	}{
		{
			name:           "has permission from active role",
			permissionName: "user.create",
			expected:       true,
		},
		{
			name:           "does not have permission from inactive role",
			permissionName: "post.create",
			expected:       false,
		},
		{
			name:           "does not have permission",
			permissionName: "user.delete",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.HasPermission(tt.permissionName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_GetAllPermissions(t *testing.T) {
	perm1 := Permission{ID: uuid.New(), Name: "user.create"}
	perm2 := Permission{ID: uuid.New(), Name: "user.read"}
	perm3 := Permission{ID: uuid.New(), Name: "post.create"}

	user := &User{
		ID: uuid.New(),
		Roles: []Role{
			{
				Name:        "admin",
				IsActive:    true,
				Permissions: []Permission{perm1, perm2},
			},
			{
				Name:        "editor",
				IsActive:    true,
				Permissions: []Permission{perm2, perm3},
			},
			{
				Name:        "inactive",
				IsActive:    false,
				Permissions: []Permission{{ID: uuid.New(), Name: "inactive.perm"}},
			},
		},
	}

	permissions := user.GetAllPermissions()

	assert.Len(t, permissions, 3)

	permNames := make(map[string]bool)
	for _, perm := range permissions {
		permNames[perm.Name] = true
	}

	assert.True(t, permNames["user.create"])
	assert.True(t, permNames["user.read"])
	assert.True(t, permNames["post.create"])
	assert.False(t, permNames["inactive.perm"])
}
