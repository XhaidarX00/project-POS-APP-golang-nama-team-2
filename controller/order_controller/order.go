package ordercontroller

import (
	"net/http"
	"project_pos_app/helper"
	"project_pos_app/model"
	"project_pos_app/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderController interface {
	GetAllOrder(c *gin.Context)
	CreateOrder(c *gin.Context)
	UpdateOrder(c *gin.Context)
	GetAllTable(c *gin.Context)
	GetAllPayment(c *gin.Context)
	DeleteOrder(c *gin.Context)
}

type orderController struct {
	service *service.AllService
	log     *zap.Logger
}

func NewOrderController(service *service.AllService, log *zap.Logger) OrderController {
	return &orderController{service, log}
}

func (oc *orderController) GetAllOrder(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")

	orders, err := oc.service.Order.GetAllOrder(search, status)
	if err != nil {
		helper.Responses(c, http.StatusNotFound, "order not found", nil)
		return
	}

	helper.Responses(c, http.StatusOK, "Order succesfully Retrived", orders)
}

func (oc *orderController) CreateOrder(c *gin.Context) {

	order := model.Order{}

	if err := c.ShouldBindJSON(&order); err != nil {
		helper.Responses(c, http.StatusInternalServerError, "Invalid Input: "+err.Error(), nil)
		return
	}

	if err := oc.service.Order.CreateOrder(&order); err != nil {
		helper.Responses(c, http.StatusBadRequest, "failed to create order: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusCreated, "Order Succesfully Created", order)
}

func (oc *orderController) UpdateOrder(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	order := model.Order{}

	if err := c.ShouldBindJSON(&order); err != nil {
		helper.Responses(c, http.StatusInternalServerError, "Invalid Input: "+err.Error(), nil)
		return
	}

	if order.PaymentMethod == "" {
		helper.Responses(c, http.StatusBadRequest, "payment is required", nil)
		return
	}

	if err := oc.service.Order.UpdateOrder(id, &order); err != nil {
		helper.Responses(c, http.StatusBadRequest, "failed to update order: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusOK, "Order Succesfully Updated", order)
}

func (oc *orderController) DeleteOrder(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	if err := oc.service.Order.DeleteOrder(id); err != nil {
		helper.Responses(c, http.StatusNotFound, "Error: "+err.Error(), nil)
		return
	}

	data := map[string]int{
		"id": id,
	}

	helper.Responses(c, http.StatusOK, "Successfully Deleted Order", data)
}
