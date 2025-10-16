package database

import (
	"base-gin/internal/domain/models"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Seeder provides database seeding functionality
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// AutoMigrate runs all database migrations
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running auto migrations...")

	err := db.AutoMigrate(
		&models.Permission{},
		&models.Role{},
		&models.User{},
		&models.RefreshToken{},
	)

	if err != nil {
		return err
	}

	log.Println("Auto migrations completed successfully!")
	return nil
}

// Seed executes multiple seeders in sequence
func (s *Seeder) Seed(seeders ...func(*gorm.DB) error) error {
	log.Println("Starting database seeding...")

	for _, seeder := range seeders {
		if err := seeder(s.db); err != nil {
			return fmt.Errorf("failed to seed database: %w", err)
		}
	}

	log.Println("Database seeding completed successfully")
	return nil
}

// Truncate truncates specified tables
func (s *Seeder) Truncate(tables ...string) error {
	log.Println("Truncating tables...")

	for _, table := range tables {
		sql := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		if err := s.db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	log.Println("Tables truncated successfully")
	return nil
}

// Clear removes all data from specified models
func (s *Seeder) Clear(models ...interface{}) error {
	log.Println("Clearing data...")

	for _, model := range models {
		if err := s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model).Error; err != nil {
			return fmt.Errorf("failed to clear data: %w", err)
		}
	}

	log.Println("Data cleared successfully")
	return nil
}

// SeedDefaultRoles seeds the default roles and permissions
func SeedDefaultRoles(db *gorm.DB) error {
	// Check if roles already exist
	var count int64
	if err := db.Model(&models.Role{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Roles already seeded, skipping...")
		return nil
	}

	log.Println("Seeding default roles and permissions...")

	permissions := createDefaultPermissions()
	if err := seedPermissions(db, permissions); err != nil {
		return err
	}

	roles := createDefaultRoles(permissions)
	if err := seedRoles(db, roles); err != nil {
		return err
	}

	log.Println("✓ Default roles and permissions seeded successfully!")
	return nil
}

// createDefaultPermissions creates the default permissions
func createDefaultPermissions() []models.Permission {
	return []models.Permission{
		// User permissions
		{
			ID:          uuid.New(),
			Name:        models.MakePermissionName(models.ResourceUser, models.ActionCreate),
			DisplayName: "Create User",
			Description: "Create new users",
			Resource:    models.ResourceUser,
			Action:      models.ActionCreate,
			IsSystem:    true,
		},
		{
			ID:          uuid.New(),
			Name:        models.MakePermissionName(models.ResourceUser, models.ActionRead),
			DisplayName: "Read User",
			Description: "View user details",
			Resource:    models.ResourceUser,
			Action:      models.ActionRead,
			IsSystem:    true,
		},
		{
			ID:          uuid.New(),
			Name:        models.MakePermissionName(models.ResourceUser, models.ActionUpdate),
			DisplayName: "Update User",
			Description: "Update user information",
			Resource:    models.ResourceUser,
			Action:      models.ActionUpdate,
			IsSystem:    true,
		},
		{
			ID:          uuid.New(),
			Name:        models.MakePermissionName(models.ResourceUser, models.ActionDelete),
			DisplayName: "Delete User",
			Description: "Delete users",
			Resource:    models.ResourceUser,
			Action:      models.ActionDelete,
			IsSystem:    true,
		},
		{
			ID:          uuid.New(),
			Name:        models.MakePermissionName(models.ResourceUser, models.ActionList),
			DisplayName: "List Users",
			Description: "List all users",
			Resource:    models.ResourceUser,
			Action:      models.ActionList,
			IsSystem:    true,
		},
		// Role permissions
		{
			ID:          uuid.New(),
			Name:        models.MakePermissionName(models.ResourceRole, models.ActionCreate),
			DisplayName: "Create Role",
			Description: "Create new roles",
			Resource:    models.ResourceRole,
			Action:      models.ActionCreate,
			IsSystem:    true,
		},
		{
			ID:          uuid.New(),
			Name:        models.MakePermissionName(models.ResourceRole, models.ActionRead),
			DisplayName: "Read Role",
			Description: "View role details",
			Resource:    models.ResourceRole,
			Action:      models.ActionRead,
			IsSystem:    true,
		},
		{
			ID:          uuid.New(),
			Name:        models.MakePermissionName(models.ResourceRole, models.ActionUpdate),
			DisplayName: "Update Role",
			Description: "Update role information",
			Resource:    models.ResourceRole,
			Action:      models.ActionUpdate,
			IsSystem:    true,
		},
		{
			ID:          uuid.New(),
			Name:        models.MakePermissionName(models.ResourceRole, models.ActionDelete),
			DisplayName: "Delete Role",
			Description: "Delete roles",
			Resource:    models.ResourceRole,
			Action:      models.ActionDelete,
			IsSystem:    true,
		},
		{
			ID:          uuid.New(),
			Name:        models.MakePermissionName(models.ResourceRole, models.ActionList),
			DisplayName: "List Roles",
			Description: "List all roles",
			Resource:    models.ResourceRole,
			Action:      models.ActionList,
			IsSystem:    true,
		},
		// Permission permissions
		{
			ID:          uuid.New(),
			Name:        models.MakePermissionName(models.ResourcePermission, models.ActionRead),
			DisplayName: "Read Permission",
			Description: "View permission details",
			Resource:    models.ResourcePermission,
			Action:      models.ActionRead,
			IsSystem:    true,
		},
		{
			ID:          uuid.New(),
			Name:        models.MakePermissionName(models.ResourcePermission, models.ActionList),
			DisplayName: "List Permissions",
			Description: "List all permissions",
			Resource:    models.ResourcePermission,
			Action:      models.ActionList,
			IsSystem:    true,
		},
	}
}

// createDefaultRoles creates the default roles with their permissions
func createDefaultRoles(permissions []models.Permission) []models.Role {
	return []models.Role{
		{
			ID:          uuid.New(),
			Name:        "admin",
			DisplayName: "Administrator",
			Description: "Full system access with all permissions",
			IsSystem:    true,
			IsActive:    true,
			Permissions: permissions, // Admin has all permissions
		},
		{
			ID:          uuid.New(),
			Name:        "user",
			DisplayName: "User",
			Description: "Basic user role with limited permissions",
			IsSystem:    true,
			IsActive:    true,
			Permissions: []models.Permission{permissions[1]}, // user.read
		},
		{
			ID:          uuid.New(),
			Name:        "moderator",
			DisplayName: "Moderator",
			Description: "Moderator role with some administrative permissions",
			IsSystem:    true,
			IsActive:    true,
			Permissions: []models.Permission{
				permissions[1], // user.read
				permissions[4], // user.list
				permissions[6], // role.read
				permissions[9], // role.list
			},
		},
	}
}

// seedPermissions seeds permissions to database
func seedPermissions(db *gorm.DB, permissions []models.Permission) error {
	for _, perm := range permissions {
		if err := db.Create(&perm).Error; err != nil {
			log.Printf("Warning: Failed to create permission %s: %v", perm.Name, err)
		}
	}
	return nil
}

// seedRoles seeds roles to database
func seedRoles(db *gorm.DB, roles []models.Role) error {
	for _, role := range roles {
		if err := db.Create(&role).Error; err != nil {
			log.Printf("Warning: Failed to create role %s: %v", role.Name, err)
		} else {
			log.Printf("✓ Created role: %s with %d permissions", role.Name, len(role.Permissions))
		}
	}
	return nil
}
