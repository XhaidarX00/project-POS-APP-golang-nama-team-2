package database

import (
	"fmt"
	"log"
	"project_pos_app/model"

	"gorm.io/gorm"
)

func Migration(db *gorm.DB) error {
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`).Error; err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Define migrations
	allModel := []struct {
		name  string
		model interface{}
	}{
		{"notification", model.Notification{}},
		{"revenue_product", model.ProductRevenue{}},
		{"revenue_order", model.OrderRevenue{}},
		{"category", model.Category{}},
		{"order_product", model.OrderProduct{}},
		{"order", model.Order{}},
		{"product", model.Product{}},
		{"table", model.Table{}},
		{"payment", model.Payment{}},
		{"user", model.User{}},
		{"reservation", model.Reservation{}},
		{"superadmin", model.Superadmin{}},
		{"employes", model.Employee{}},
		{"permission", model.Permission{}},
		{"access_permission", model.AccessPermission{}},
		{"session", model.Session{}},
		{"employee", model.Employee{}},
	}

	for _, migration := range allModel {
		var count int64
		err := db.Raw("SELECT COUNT(1) FROM migrations WHERE name = ?", migration.name).Scan(&count).Error
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", migration.name, err)
		}

		if count > 0 {
			log.Printf("Migration '%s' already applied, skipping.", migration.name)
			continue
		}

		// Run migration
		if err := db.AutoMigrate(migration.model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", migration.model, err)
		}

		// Record migration as applied
		if err := db.Exec("INSERT INTO migrations (name) VALUES (?)", migration.name).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.name, err)
		}

		log.Printf("Migration '%s' applied successfully.", migration.name)
	}

	return nil
}
