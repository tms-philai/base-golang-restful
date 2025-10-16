package setup

import (
	"base-gin/internal/pkg/database"
	"base-gin/test/config"
	"fmt"
	"log"
	"os"
	"testing"

	"gorm.io/gorm"
)

// TestDatabaseSetup handles test database setup and teardown
type TestDatabaseSetup struct {
	Config *config.TestConfig
}

// NewTestDatabaseSetup creates a new test database setup
func NewTestDatabaseSetup() (*TestDatabaseSetup, error) {
	testConfig, err := config.LoadTestConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load test config: %w", err)
	}

	return &TestDatabaseSetup{
		Config: testConfig,
	}, nil
}

// SetupTestDatabase sets up the test database
func (tds *TestDatabaseSetup) SetupTestDatabase() error {
	// Connect to test database
	dbConfig := database.Config{
		Host:     tds.Config.Database.Host,
		Port:     tds.Config.Database.Port,
		User:     tds.Config.Database.User,
		Password: tds.Config.Database.Password,
		DBName:   tds.Config.Database.Name,
		SSLMode:  tds.Config.Database.SSLMode,
		TimeZone: tds.Config.Database.TimeZone,
	}

	if err := database.ConnectWithRetry(dbConfig, 5, 2); err != nil {
		return fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Run migrations
	if err := database.AutoMigrate(database.GetDB()); err != nil {
		return fmt.Errorf("failed to run test database migrations: %w", err)
	}

	// Seed default roles and permissions
	if err := database.SeedDefaultRoles(database.GetDB()); err != nil {
		log.Printf("Warning: Failed to seed default roles: %v", err)
	}

	return nil
}

// CleanupTestDatabase cleans up the test database
func (tds *TestDatabaseSetup) CleanupTestDatabase() error {
	// Close database connection
	database.Close()
	return nil
}

// TruncateTestTables truncates all test tables
func (tds *TestDatabaseSetup) TruncateTestTables() error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// List of tables to truncate (in order to respect foreign key constraints)
	tables := []string{
		"user_roles",
		"role_permissions",
		"refresh_tokens",
		"notifications",
		"products",
		"files",
		"users",
		"roles",
		"permissions",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			log.Printf("Warning: Failed to truncate table %s: %v", table, err)
		}
	}

	return nil
}

// SetupTestData sets up test data
func (tds *TestDatabaseSetup) SetupTestData() error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Create test users, products, etc.
	// This would contain logic to create test data
	// For now, it's a placeholder

	return nil
}

// CleanupTestData cleans up test data
func (tds *TestDatabaseSetup) CleanupTestData() error {
	return tds.TruncateTestTables()
}

// TestMain is the main test function that sets up and tears down the test environment
func TestMain(m *testing.M) {
	// Setup test database
	setup, err := NewTestDatabaseSetup()
	if err != nil {
		log.Fatalf("Failed to create test database setup: %v", err)
	}

	// Setup test database
	if err := setup.SetupTestDatabase(); err != nil {
		log.Fatalf("Failed to setup test database: %v", err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if err := setup.CleanupTestDatabase(); err != nil {
		log.Printf("Failed to cleanup test database: %v", err)
	}

	// Cleanup test config
	if err := config.CleanupTestConfig(); err != nil {
		log.Printf("Failed to cleanup test config: %v", err)
	}

	os.Exit(code)
}

// SetupTestEnvironment sets up the test environment for a specific test
func SetupTestEnvironment(t *testing.T) *TestDatabaseSetup {
	setup, err := NewTestDatabaseSetup()
	if err != nil {
		t.Fatalf("Failed to create test database setup: %v", err)
	}

	if err := setup.SetupTestDatabase(); err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Clean up tables before each test
	if err := setup.TruncateTestTables(); err != nil {
		t.Fatalf("Failed to truncate test tables: %v", err)
	}

	return setup
}

// CleanupTestEnvironment cleans up the test environment
func CleanupTestEnvironment(t *testing.T, setup *TestDatabaseSetup) {
	if setup != nil {
		if err := setup.CleanupTestDatabase(); err != nil {
			t.Logf("Failed to cleanup test database: %v", err)
		}
	}
}

// GetTestDB returns the test database connection
func GetTestDB() *gorm.DB {
	return database.GetDB()
}

// IsTestEnvironment checks if we're running in a test environment
func IsTestEnvironment() bool {
	return testing.Testing()
}
