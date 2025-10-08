package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

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
