package productcontroller

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

type ProductController struct {
	service *service.AllService
	log     *zap.Logger
}

func NewProductController(service *service.AllService, log *zap.Logger) *ProductController {
	return &ProductController{
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
// @Success 200 {object} model.SuccessResponse{data=[]model.Product} "List of categories retrieved successfully"
// @Failure 500 {object} model.ErrorResponse "Failed to fetch categories"
// @Router /product [get]
func (pc *ProductController) GetAllProducts(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	products, total, totalPages, err := pc.service.Product.ShowAllProduct(page, limit)
	if err != nil {
		pc.log.Error("Failed to fetch products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	pc.log.Info("Fetched products successfully", zap.Int("total", total), zap.Int("pages", totalPages))
	c.JSON(http.StatusOK, gin.H{
		"data":        products,
		"total":       total,
		"totalPages":  totalPages,
		"currentPage": page,
	})
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Retrieve a single product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} model.SuccessResponse{data=model.Product} "Product retrieved successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid product ID"
// @Failure 404 {object} model.ErrorResponse "Product not found"
// @Router /product/{id} [get]
func (pc *ProductController) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pc.log.Error("Invalid product ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := pc.service.Product.GetProductByID((id))
	if err != nil {
		pc.log.Error("Product not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	pc.log.Info("Fetched product successfully", zap.Uint("id", uint(id)))
	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Add a new product to the inventory
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Product Name"
// @Param description formData string true "Product Description"
// @Param price formData float64 true "Product Price"
// @Param stock formData int true "Product Stock"
// @Param image formData file true "Product Image"
// @Success 201 {object} model.SuccessResponse{data=model.Product} "Product created successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid product data"
// @Failure 500 {object} model.ErrorResponse "Failed to create product"
// @Router /product [post]
func (pc *ProductController) CreateProduct(c *gin.Context) {
	pc.log.Info("Starting product creation")

	// Membaca form data dari request
	form, err := c.MultipartForm()
	if err != nil {
		pc.log.Error("Error reading form data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data: " + err.Error()})
		return
	}

	// Mendapatkan file gambar dari form
	files := form.File["image_url"]
	if len(files) == 0 {
		pc.log.Error("No image file provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	// Menggunakan goroutine untuk upload gambar
	var wg sync.WaitGroup
	responses, err := helper.Upload(&wg, files)
	if err != nil {
		pc.log.Error("Failed to upload image", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	// Menangani jika ada respons dari upload gambar
	if len(responses) > 0 {
		pc.log.Info("Image uploaded successfully", zap.String("imageURL", responses[0].Data.Url))
	}

	// Menangani data lainnya dari form
	name := c.PostForm("name")
	itemID := c.PostForm("item_id")
	stock := c.PostForm("stock")
	categoryID, err := strconv.Atoi(c.DefaultPostForm("category_id", "0"))
	if err != nil {
		pc.log.Error("Invalid category_id", zap.String("category_id", c.PostForm("category_id")), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
		return
	}
	price, err := strconv.ParseFloat(c.PostForm("price"), 64)
	if err != nil {
		pc.log.Error("Invalid price value", zap.String("price", c.PostForm("price")), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price value"})
		return
	}
	status := c.DefaultPostForm("status", "available")

	product := model.Product{
		Name:       name,
		ItemID:     itemID,
		Stock:      stock,
		CategoryID: uint(categoryID),
		Qty:        1,
		Price:      price,
		Status:     status,
		ImageURL:   responses[0].Data.Url,
	}

	// Membuat produk di database
	if err := pc.service.Product.CreateProduct(&product); err != nil {
		pc.log.Error("Failed to create product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// Menampilkan pesan sukses setelah produk berhasil dibuat
	pc.log.Info("Product created successfully", zap.String("productName", product.Name))
	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product": product})
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update the details of a product by its ID
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param name formData string false "Product Name"
// @Param description formData string false "Product Description"
// @Param price formData float64 false "Product Price"
// @Param stock formData int false "Product Stock"
// @Param image formData file false "Product Image"
// @Success 200 {object} model.SuccessResponse{data=model.Product} "Product updated successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid product ID or data"
// @Failure 404 {object} model.ErrorResponse "Product not found"
// @Failure 500 {object} model.ErrorResponse "Failed to update product"
// @Router /product/{id} [put]
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pc.log.Error("Invalid product ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product model.Product
	if err := c.ShouldBind(&product); err != nil {
		pc.log.Error("Invalid product data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
		return
	}

	files := c.Request.MultipartForm.File["image_url"]
	if len(files) > 0 {
		var wg sync.WaitGroup
		responses, err := helper.Upload(&wg, files)
		if err != nil {
			pc.log.Error("Failed to upload image", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}
		if len(responses) > 0 {
			product.ImageURL = responses[0].Data.Url
		}
	}

	if err := pc.service.Product.UpdateProduct(uint(id), &product); err != nil {
		pc.log.Error("Failed to update product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	pc.log.Info("Product updated successfully", zap.Uint("id", uint(id)))
	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Remove a product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID" example(1)
// @Success 200 {object} model.SuccessResponse{data=map[string]string} "Product deleted successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid product ID"
// @Failure 404 {object} model.ErrorResponse "Product not found"
// @Failure 500 {object} model.ErrorResponse "Failed to delete product"
// @Router /product/{id} [delete]
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pc.log.Error("Invalid product ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := pc.service.Product.DeleteProduct((id)); err != nil {
		pc.log.Error("Failed to delete product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	pc.log.Info("Product deleted successfully", zap.Uint("id", uint(id)))
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
