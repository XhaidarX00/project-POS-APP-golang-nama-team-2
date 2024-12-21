package categoryrepository

import (
	"fmt"
	"math"
	"project_pos_app/model"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	ShowAllCategory(page, limit int) (*[]model.Category, int, int, error)
	GetCategoryByID(id uint) (*model.Category, error)
	CreateCategory(category *model.Category) error
	UpdateCategory(categoryID uint, category *model.Category) error
}

type categoryRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewCategoryRepo(db *gorm.DB, log *zap.Logger) CategoryRepository {
	return &categoryRepository{db, log}
}

func (cr *categoryRepository) ShowAllCategory(page, limit int) (*[]model.Category, int, int, error) {
	cr.log.Info("Fetching all category", zap.Int("page", page), zap.Int("limit", limit))

	var category []model.Category
	var totalRecords int64

	// Count total records
	if err := cr.db.Model(&model.Category{}).Count(&totalRecords).Error; err != nil {
		cr.log.Error("Error counting category", zap.Error(err))
		return nil, 0, 0, err
	}

	// Fetch paginated results
	offset := (page - 1) * limit
	if err := cr.db.Offset(offset).Limit(limit).Find(&category).Error; err != nil {
		cr.log.Error("Error fetching products", zap.Error(err))
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	cr.log.Info("Successfully fetched products", zap.Int("totalRecords", int(totalRecords)), zap.Int("totalPages", totalPages))
	return &category, int(totalRecords), totalPages, nil
}

func (cr *categoryRepository) GetCategoryByID(id uint) (*model.Category, error) {
	cr.log.Info("Fetching category by ID", zap.Uint("id", id))

	var category model.Category
	if err := cr.db.First(&category, id).Error; err != nil {
		cr.log.Error("Error fetching category", zap.Uint("id", id), zap.Error(err))
		return nil, fmt.Errorf("category not found")
	}

	cr.log.Info("Successfully fetched category", zap.Uint("id", id))
	return &category, nil
}

func (cr *categoryRepository) CreateCategory(category *model.Category) error {
	cr.log.Info("Creating category", zap.String("name", category.Name))

	var wg sync.WaitGroup
	var err error

	err = cr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(category).Error; err != nil {
			cr.log.Error("Failed to create category", zap.String("name", category.Name), zap.Error(err))
			return fmt.Errorf("failed to create category: %w", err)
		}
		wg.Wait()
		return nil
	})

	if err != nil {
		cr.log.Error("Transaction failed", zap.String("name", category.Name), zap.Error(err))
		return err
	}

	cr.log.Info("Successfully created product", zap.String("name", category.Name))
	return nil
}

func (cr *categoryRepository) UpdateCategory(categoryID uint, category *model.Category) error {
	cr.log.Info("Updating category with data", zap.Any("category", category))
	result := cr.db.Model(&model.Category{}).Where("id = ?", categoryID).Updates(map[string]interface{}{
		"name":        category.Name,
		"description": category.Description,
		"icon_url":    category.IconURL,
	})
	cr.log.Info("Update result", zap.Int64("RowsAffected", result.RowsAffected))

	if result.Error != nil {
		cr.log.Error("Failed to update category", zap.Uint("categoryID", categoryID), zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		cr.log.Warn("No category found to update", zap.Uint("categoryID", categoryID))
		return fmt.Errorf("no category found with id %d", categoryID)
	}

	cr.log.Info("Successfully updated category", zap.Uint("categoryID", categoryID))
	return nil
}
