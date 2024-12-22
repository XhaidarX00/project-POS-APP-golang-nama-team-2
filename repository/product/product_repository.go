package productrepository

import (
	"fmt"
	"math"
	"project_pos_app/model"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProductRepo defines the interface for product operations.
type ProductRepo interface {
	ShowAllProducts(page, limit int) (*[]model.Product, int, int, error)
	GetProductByID(id uint) (*model.Product, error)
	CreateProduct(product *model.Product) error
	UpdateProduct(productID uint, product *model.Product) error
	DeleteProduct(id uint) error
}

// productRepo implements the ProductRepo interface.
type productRepo struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewProductRepo creates a new instance of productRepo.
func NewProductRepo(db *gorm.DB, log *zap.Logger) ProductRepo {
	return &productRepo{db, log}
}

// ShowAllProducts fetches all products with pagination.
func (pr *productRepo) ShowAllProducts(page, limit int) (*[]model.Product, int, int, error) {
	pr.log.Info("Fetching all products", zap.Int("page", page), zap.Int("limit", limit))

	var products []model.Product
	var totalRecords int64

	// Count total records
	if err := pr.db.Model(&model.Product{}).Count(&totalRecords).Error; err != nil {
		pr.log.Error("Error counting products", zap.Error(err))
		return nil, 0, 0, err
	}

	// Fetch paginated results
	offset := (page - 1) * limit
	if err := pr.db.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		pr.log.Error("Error fetching products", zap.Error(err))
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	pr.log.Info("Successfully fetched products", zap.Int("totalRecords", int(totalRecords)), zap.Int("totalPages", totalPages))
	return &products, int(totalRecords), totalPages, nil
}

// GetProductByID fetches a product by its ID.
func (pr *productRepo) GetProductByID(id uint) (*model.Product, error) {
	pr.log.Info("Fetching product by ID", zap.Uint("id", id))

	var product model.Product
	if err := pr.db.First(&product, id).Error; err != nil {
		pr.log.Error("Error fetching product", zap.Uint("id", id), zap.Error(err))
		return nil, fmt.Errorf("product not found")
	}

	pr.log.Info("Successfully fetched product", zap.Uint("id", id))
	return &product, nil
}

// CreateProduct creates a new product.
func (pr *productRepo) CreateProduct(product *model.Product) error {
	pr.log.Info("Creating product", zap.String("name", product.Name))

	var wg sync.WaitGroup
	var err error

	err = pr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			pr.log.Error("Failed to create product", zap.String("name", product.Name), zap.Error(err))
			return fmt.Errorf("failed to create product: %w", err)
		}
		wg.Wait()
		return nil
	})

	if err != nil {
		pr.log.Error("Transaction failed", zap.String("name", product.Name), zap.Error(err))
		return err
	}

	pr.log.Info("Successfully created product", zap.String("name", product.Name))
	return nil
}

// UpdateProduct updates an existing product by its ID.
func (pr *productRepo) UpdateProduct(productID uint, product *model.Product) error {
	pr.log.Info("Updating product", zap.Uint("productID", productID))

	result := pr.db.Model(&model.Product{}).Where("id = ?", productID).Updates(product)
	if result.Error != nil {
		pr.log.Error("Failed to update product", zap.Uint("productID", productID), zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		pr.log.Warn("No product found to update", zap.Uint("productID", productID))
		return fmt.Errorf("no product found with id %d", productID)
	}

	pr.log.Info("Successfully updated product", zap.Uint("productID", productID))
	return nil
}

// DeleteProduct deletes a product by its ID.
func (pr *productRepo) DeleteProduct(id uint) error {
	pr.log.Info("Deleting product", zap.Uint("id", id))

	result := pr.db.Delete(&model.Product{}, id)
	if result.Error != nil {
		pr.log.Error("Failed to delete product", zap.Uint("id", id), zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		pr.log.Warn("No product found to delete", zap.Uint("id", id))
		return fmt.Errorf("product not found with id %d", id)
	}

	pr.log.Info("Successfully deleted product", zap.Uint("id", id))
	return nil
}
