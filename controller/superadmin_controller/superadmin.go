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
}

type superadminController struct {
	service *service.AllService
	log     *zap.Logger
}

func NewSuperadminController(service *service.AllService, log *zap.Logger) SuperadminController {
	return &superadminController{service, log}
}

func (sc *superadminController) ListDataAdmin(c *gin.Context) {

	admins, err := sc.service.Superadmin.ListDataAdmin()
	if err != nil {
		helper.Responses(c, http.StatusNotFound, "Error: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusOK, "Succesfully Retrieved data", admins)

}

func (sc *superadminController) UpdateSuperadmin(c *gin.Context) {
	// Struct untuk validasi input
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
