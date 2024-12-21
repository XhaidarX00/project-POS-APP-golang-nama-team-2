package categorycontroller

import (
	"net/http"
	"project_pos_app/helper"
	"project_pos_app/model"
	"project_pos_app/service"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryController struct {
	service *service.AllService
	log     *zap.Logger
}

func NewCategoryController(service *service.AllService, log *zap.Logger) *CategoryController {
	return &CategoryController{
		service: service,
		log:     log,
	}
}

// GetAllCategory godoc
// @Summary Get all categories
// @Description Get a list of categories with optional pagination
// @Tags Categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} model.SuccessResponse{data=[]model.Category} "List of categories retrieved successfully"
// @Failure 500 {object} model.ErrorResponse "Failed to fetch categories"
// @Router /api/categories [get]
func (cc *CategoryController) GetAllCategory(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	categories, total, totalPages, err := cc.service.Category.ShowAllCategory(page, limit)
	if err != nil {
		cc.log.Error("Failed to fetch categories", zap.Error(err))
		helper.Responses(c, http.StatusInternalServerError, "Failed to fetch categories", nil)
		return
	}

	response := gin.H{
		"categories":  categories,
		"total":       total,
		"totalPages":  totalPages,
		"currentPage": page,
	}
	helper.Responses(c, http.StatusOK, "Categories retrieved successfully", response)
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get a single category by its ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} model.SuccessResponse{data=model.Category} "Category retrieved successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid category ID"
// @Failure 404 {object} model.ErrorResponse "Category not found"
// @Router /api/categories/{id} [get]
func (cc *CategoryController) GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		cc.log.Error("Invalid category ID", zap.Error(err))
		helper.Responses(c, http.StatusBadRequest, "Invalid category ID format", nil)
		return
	}

	category, err := cc.service.Category.GetCategoryByID(id)
	if err != nil {
		cc.log.Error("Category not found", zap.Error(err))
		helper.Responses(c, http.StatusNotFound, "Category not found", nil)
		return
	}

	helper.Responses(c, http.StatusOK, "Category retrieved successfully", category)
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Add a new category to the database
// @Tags Categories
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Category Name"
// @Param description formData string true "Category Description"
// @Param icon formData file true "Category Icon"
// @Success 201 {object} model.SuccessResponse{data=model.Category} "Category created successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid category data"
// @Failure 500 {object} model.ErrorResponse "Failed to create category"
// @Router /api/categories [post]
func (cc *CategoryController) CreateCategory(c *gin.Context) {
	cc.log.Info("Starting category creation")

	// Membaca form data dari request
	form, err := c.MultipartForm()
	if err != nil {
		cc.log.Error("Error reading form data", zap.Error(err))
		helper.Responses(c, http.StatusBadRequest, "Invalid form data: "+err.Error(), nil)
		return
	}

	// Mendapatkan file gambar dari form
	files := form.File["icon_url"]
	if len(files) == 0 {
		cc.log.Error("No icon file provided")
		helper.Responses(c, http.StatusBadRequest, "Icon file is required", nil)
		return
	}

	// Menggunakan goroutine untuk upload gambar
	var wg sync.WaitGroup
	responses, err := helper.Upload(&wg, files)
	if err != nil {
		cc.log.Error("Failed to upload icon", zap.Error(err))
		helper.Responses(c, http.StatusInternalServerError, "Failed to upload icon", nil)
		return
	}

	// Menangani data lainnya dari form
	name := c.PostForm("name")
	description := c.PostForm("description")
	category := model.Category{
		IconURL:     responses[0].Data.Url,
		Name:        name,
		Description: description,
	}

	// Membuat category di database
	if err := cc.service.Category.CreateCategory(&category); err != nil {
		cc.log.Error("Failed to create category", zap.Error(err))
		helper.Responses(c, http.StatusInternalServerError, "Failed to create category", nil)
		return
	}

	// Menampilkan pesan sukses setelah category berhasil dibuat
	cc.log.Info("Category created successfully", zap.String("categoryName", category.Name))
	helper.Responses(c, http.StatusCreated, "Category created successfully", category)
}

// UpdateCategory godoc
// @Summary Update an existing category
// @Description Update the details of a category by its ID
// @Tags Categories
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Category ID"
// @Param name formData string false "Category Name"
// @Param description formData string false "Category Description"
// @Param icon formData file false "Category Icon"
// @Success 200 {object} model.SuccessResponse{data=model.Category} "Category updated successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid category ID or data"
// @Failure 404 {object} model.ErrorResponse "Category not found"
// @Failure 500 {object} model.ErrorResponse "Failed to update category"
// @Router /api/categories/{id} [put]
func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		cc.log.Error("Invalid category ID", zap.Error(err))
		helper.Responses(c, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	// Fetch existing category
	category, err := cc.service.Category.GetCategoryByID(id)
	if err != nil {
		cc.log.Error("Category not found", zap.Error(err))
		helper.Responses(c, http.StatusNotFound, "Category not found", nil)
		return
	}

	// Read form data
	form, err := c.MultipartForm()
	if err != nil {
		cc.log.Error("Error reading form data", zap.Error(err))
		helper.Responses(c, http.StatusBadRequest, "Invalid form data: "+err.Error(), nil)
		return
	}

	// Update file if provided
	files := form.File["icon"]
	if len(files) > 0 {
		var wg sync.WaitGroup
		responses, err := helper.Upload(&wg, files)
		if err != nil {
			cc.log.Error("Failed to upload icon", zap.Error(err))
			helper.Responses(c, http.StatusInternalServerError, "Failed to upload icon", nil)
			return
		}
		category.IconURL = responses[0].Data.Url
		cc.log.Info("Icon updated successfully", zap.String("iconURL", responses[0].Data.Url))
	}

	// Update other fields if provided
	if name := c.PostForm("name"); name != "" {
		category.Name = name
	}
	if description := c.PostForm("description"); description != "" {
		category.Description = description
	}

	// Save updates to database
	if err := cc.service.Category.UpdateCategory(category.ID, category); err != nil {
		cc.log.Error("Failed to update category", zap.Error(err))
		helper.Responses(c, http.StatusInternalServerError, "Failed to update category", nil)
		return
	}

	cc.log.Info("Category updated successfully", zap.String("categoryName", category.Name))
	helper.Responses(c, http.StatusOK, "Category updated successfully", category)
}
