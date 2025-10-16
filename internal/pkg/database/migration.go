package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) AutoMigrate(models ...interface{}) error {
	log.Println("Starting database migration...")

	if err := m.db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

func (m *Migrator) DropTables(models ...interface{}) error {
	log.Println("Dropping tables...")

	for _, model := range models {
		if err := m.db.Migrator().DropTable(model); err != nil {
			return fmt.Errorf("failed to drop table: %w", err)
		}
	}

	log.Println("Tables dropped successfully")
	return nil
}

func (m *Migrator) CreateIndexes(tableName string, indexes map[string][]string) error {
	log.Printf("Creating indexes for table %s...", tableName)

	for indexName, columns := range indexes {
		sql := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s (%s)",
			indexName, tableName, joinColumns(columns))

		if err := m.db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to create index %s: %w", indexName, err)
		}
	}

	log.Printf("Indexes created successfully for table %s", tableName)
	return nil
}

func joinColumns(columns []string) string {
	result := ""
	for i, col := range columns {
		if i > 0 {
			result += ", "
		}
		result += col
	}
	return result
}

func (m *Migrator) HasTable(tableName string) bool {
	return m.db.Migrator().HasTable(tableName)
}

func (m *Migrator) HasColumn(tableName, columnName string) bool {
	return m.db.Migrator().HasColumn(tableName, columnName)
}
