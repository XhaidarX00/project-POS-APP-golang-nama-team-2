package categoryservice

import (
	"project_pos_app/model"
	"project_pos_app/repository"

	"go.uber.org/zap"
)

type CategoryService interface {
	ShowAllCategory(page, limit int) (*[]model.Category, int, int, error)
	GetCategoryByID(id int) (*model.Category, error)
	CreateCategory(category *model.Category) error
	UpdateCategory(categoryID uint, category *model.Category) error
}

type categoryService struct {
	repo *repository.AllRepository
	log  *zap.Logger
}

func NewCategoryService(repo *repository.AllRepository, log *zap.Logger) CategoryService {
	return &categoryService{repo: repo, log: log}
}

func (cs *categoryService) ShowAllCategory(page, limit int) (*[]model.Category, int, int, error) {
	cs.log.Info("Fetching all category", zap.Int("page", page), zap.Int("limit", limit))

	category, count, totalPages, err := cs.repo.Category.ShowAllCategory(page, limit)
	if err != nil {
		cs.log.Error("Error fetching category", zap.Error(err))
		return nil, 0, 0, err
	}

	cs.log.Info("Successfully fetched category", zap.Int("count", count), zap.Int("totalPages", totalPages))
	return category, count, totalPages, nil
}

func (cs *categoryService) GetCategoryByID(id int) (*model.Category, error) {
	cs.log.Info("Fetching category by ID", zap.Int("id", id))

	category, err := cs.repo.Category.GetCategoryByID(uint(id))
	if err != nil {
		cs.log.Error("Error fetching category", zap.Error(err))
		return nil, err
	}

	cs.log.Info("Successfully fetched category", zap.Int("id", id))
	return category, nil
}

func (cs *categoryService) CreateCategory(category *model.Category) error {
	cs.log.Info("Creating category", zap.String("name", category.Name))

	err := cs.repo.Category.CreateCategory(category)
	if err != nil {
		cs.log.Error("Error creating category", zap.Error(err))
		return err
	}

	cs.log.Info("Successfully created category", zap.String("name", category.Name))
	return nil
}

func (ps *categoryService) UpdateCategory(categoryID uint, category *model.Category) error {
	ps.log.Info("Updating category", zap.Uint("categoryID", categoryID), zap.String("name", category.Name))

	if err := ps.repo.Category.UpdateCategory(categoryID, category); err != nil {
		ps.log.Error("Error updating category", zap.Error(err))
		return err
	}

	ps.log.Info("Successfully updated category", zap.Uint("categoryID", categoryID))
	return nil
}
