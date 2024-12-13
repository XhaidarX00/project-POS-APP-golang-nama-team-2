package database

import "gorm.io/gorm"

func Migration(db *gorm.DB) error {
	err := db.AutoMigrate()

	return err
}
