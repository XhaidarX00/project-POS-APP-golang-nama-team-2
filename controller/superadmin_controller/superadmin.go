package superadmincontroller

import (
	"net/http"
	"project_pos_app/helper"
	"project_pos_app/model"
	"project_pos_app/service"
	"project_pos_app/utils"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SuperadminController interface {
	ListDataAdmin(c *gin.Context)
	UpdateSuperadmin(c *gin.Context)
	UpdateAccessUser(c *gin.Context)
	Logout(c *gin.Context)
}

type superadminController struct {
	service *service.AllService
	log     *zap.Logger
}

func NewSuperadminController(service *service.AllService, log *zap.Logger) SuperadminController {
	return &superadminController{service, log}
}

// ListDataAdmin godoc
// @Summary Retrieve list of admins
// @Description Get a list of all admins with their names and emails
// @Tags Superadmin
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse{data=[]ResponseEmployee} "Successfully retrieved admin data"
// @Failure 404 {object} ErrorResponse "Admin data not found"
// @Router /superadmin [get]
func (sc *superadminController) ListDataAdmin(c *gin.Context) {

	admins, err := sc.service.Superadmin.ListDataAdmin()
	if err != nil {
		helper.Responses(c, http.StatusNotFound, "Error: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusOK, "Succesfully Retrieved data", admins)

}

// UpdateSuperadmin godoc
// @Summary Update superadmin details
// @Description Update the details of a superadmin, including email, full name, address, password, and profile image
// @Tags Superadmin
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "Email of the superadmin" minlength(3) maxlength(100)
// @Param full_name formData string true "Full name of the superadmin" minlength(3) maxlength(100)
// @Param address formData string false "Address of the superadmin" maxlength(255)
// @Param new_password formData string false "New password for the superadmin" minlength(8)
// @Param confirm_password formData string false "Confirm password (must match new_password)"
// @Param image formData file true "Profile image (maximum size 5MB)"
// @Success 200 {object} SuccessResponse "Successfully updated superadmin"
// @Failure 400 {object} ErrorResponse "Invalid input data or validation error"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /superadmin [put]
func (sc *superadminController) UpdateSuperadmin(c *gin.Context) {
	var SuperadminInput struct {
		Email           string `form:"email" binding:"required,email"`
		FullName        string `form:"full_name" binding:"required,min=3,max=100"`
		Address         string `form:"address" binding:"omitempty,max=255"`
		NewPassword     string `form:"new_password" binding:"omitempty,min=8"`
		ConfirmPassword string `form:"confirm_password" binding:"omitempty,eqfield=NewPassword"`
	}

	if err := c.ShouldBind(&SuperadminInput); err != nil {
		data := map[string]string{
			"error": err.Error(),
		}
		sc.log.Error("Error binding form data", zap.Error(err))
		helper.Responses(c, http.StatusBadRequest, "Invalid form data: ", data)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		sc.log.Error("Error reading form data", zap.Error(err))
		helper.Responses(c, http.StatusBadRequest, "Invalid form data: "+err.Error(), nil)
		return
	}

	files, ok := form.File["image"]
	if !ok || len(files) == 0 {
		helper.Responses(c, http.StatusBadRequest, "Image is required", nil)
		return
	}

	for _, file := range files {
		if file.Size > 5*1024*1024 {
			helper.Responses(c, http.StatusBadRequest, "Image size exceeds 5MB", nil)
			return
		}
	}

	var wg sync.WaitGroup
	responses, err := helper.Upload(&wg, files)
	if err != nil {
		sc.log.Error("Failed to upload images", zap.Error(err))
		helper.Responses(c, http.StatusInternalServerError, "Failed to upload images: "+err.Error(), nil)
		return
	}

	if SuperadminInput.NewPassword != "" {
		hashedPassword, err := utils.HashPassword(SuperadminInput.ConfirmPassword)
		if err != nil {
			sc.log.Error("Error hashing password", zap.Error(err))
			helper.Responses(c, http.StatusInternalServerError, "Failed to process password", nil)
			return
		}
		SuperadminInput.ConfirmPassword = hashedPassword
	}

	admin := model.Superadmin{
		User: model.User{
			Email:    SuperadminInput.Email,
			Password: SuperadminInput.ConfirmPassword,
		},
		FullName: SuperadminInput.FullName,
		Address:  SuperadminInput.Address,
		Image:    responses[0].Data.Url,
	}

	// Update Superadmin di database
	if err := sc.service.Superadmin.UpdateSuperadmin(1, &admin); err != nil {
		helper.Responses(c, http.StatusBadRequest, "Failed to update superadmin: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusOK, "Successfully updated Superadmin", nil)
}

// UpdateAccessUser godoc
// @Summary Update access permissions for a user
// @Description Update the access permissions of a specific user based on their ID
// @Tags Superadmin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param input body model.AccessPermission true "Access Permission Payload"
// @Success 200 {object} SuccessResponse "Successfully updated access permissions"
// @Failure 400 {object} ErrorResponse "Invalid input payload"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /superadmin/{id} [put]
func (ac *superadminController) UpdateAccessUser(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	input := model.AccessPermission{}

	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Responses(c, http.StatusBadRequest, "Invalid payload request"+err.Error(), nil)
		return
	}

	err := ac.service.Superadmin.UpdateAccessUser(id, &input)
	if err != nil {
		helper.Responses(c, http.StatusInternalServerError, "Error: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusOK, "Successfully Update Access", nil)
}

// Logout godoc
// @Summary Logout user
// @Description Log the user out by invalidating their authorization token
// @Tags Superadmin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} SuccessResponse "Successfully logged out"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /logout [post]
func (ac *superadminController) Logout(c *gin.Context) {

	token := c.GetHeader("Authorization")
	if token == "" {
		helper.Responses(c, http.StatusUnauthorized, "Unauthorized", nil)
		c.Abort()
		return
	}

	err := ac.service.Superadmin.Logout(token)
	if err != nil {
		helper.Responses(c, http.StatusInternalServerError, "Error: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusOK, "Successfully Logout", nil)
}
