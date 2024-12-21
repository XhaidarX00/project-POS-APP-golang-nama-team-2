package model

import (
	"time"

	"gorm.io/gorm"
)

type Table struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `gorm:"type:varchar(255);not null" json:"name"`
	IsBook    bool            `json:"is_book"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func SeedTables() []Table {
	return []Table{
		{Name: "Book A", IsBook: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Table 1", IsBook: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Book B", IsBook: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Table 2", IsBook: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Book C", IsBook: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Table 3", IsBook: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Book D", IsBook: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Table 4", IsBook: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Book E", IsBook: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Table 5", IsBook: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Book F", IsBook: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Table 6", IsBook: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Book G", IsBook: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Table 7", IsBook: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Book H", IsBook: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Table 8", IsBook: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Book I", IsBook: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Table 9", IsBook: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Book J", IsBook: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Table 10", IsBook: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
}
