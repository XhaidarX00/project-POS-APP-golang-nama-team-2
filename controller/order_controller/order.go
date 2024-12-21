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

// GetAllOrder godoc
// @Summary Retrieve all orders
// @Description Retrieve a list of orders with optional search and status filtering.
// @Tags Orders
// @Accept json
// @Produce json
// @Param search query string false "Search keyword to filter orders"
// @Param status query string false "Filter orders by status"
// @Success 200 {object} model.SuccessResponse{data=[]model.OrderResponse} "Order successfully retrieved"
// @Failure 404 {object} model.ErrorResponse "Order not found"
// @Security Authentication
// @Router /order [get]
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

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with associated products
// @Tags Orders
// @Accept json
// @Produce json
// @Security Authentication
// @Param order body model.Order true "Order payload"
// @Success 201 {object} model.SuccessResponse "Order successfully created"
// @Failure 400 {object} model.ErrorResponse "Failed to create order"
// @Failure 500 {object} model.ErrorResponse "Invalid input"
// @Router /order [post]
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

	helper.Responses(c, http.StatusCreated, "Order Succesfully Created", nil)
}

// UpdateOrder godoc
// @Summary Update an existing order
// @Description Update the details of an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security Authentication
// @Param order body model.Order true "Updated order payload"
// @Success 200 {object} model.SuccessResponse "Order successfully updated"
// @Failure 400 {object} model.ErrorResponse "Failed to update order"
// @Failure 500 {object} model.ErrorResponse "Invalid input"
// @Router /order/{id} [put]
func (oc *orderController) UpdateOrder(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	order := model.Order{}

	if err := c.ShouldBindJSON(&order); err != nil {
		helper.Responses(c, http.StatusInternalServerError, "Invalid Input: "+err.Error(), nil)
		return
	}

	if err := oc.service.Order.UpdateOrder(id, &order); err != nil {
		helper.Responses(c, http.StatusBadRequest, "failed to update order: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusOK, "Order Succesfully Updated", nil)
}

// DeleteOrder godoc
// @Summary Delete an order
// @Description Delete an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Security Authentication
// @Param id path int true "Order ID"
// @Success 200 {object} model.SuccessResponse{data=map[string]int} "Successfully deleted order"
// @Failure 404 {object} model.ErrorResponse "Order not found"
// @Router /order/{id} [delete]
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
