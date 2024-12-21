package reservationcontroller

import (
	"net/http"
	"project_pos_app/helper"
	"project_pos_app/model"
	"project_pos_app/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ControllerReservation struct {
	Service *service.AllService
	Log     *zap.Logger
}

func NewControllerReservation(service *service.AllService, log *zap.Logger) ControllerReservation {
	return ControllerReservation{Service: service, Log: log}
}

// @Summary Get All Reservation
// @Description Endpoint For All Reservation
// @Tags Reservation
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response{data=[]model.Reservation} "Get Summary Success"
// @Failure 500 {object} helper.Response "server error"
// @Router  /reservation [get]
func (ctrl *ControllerReservation) GetAll(ctx *gin.Context) {
	date := ctx.Query("date")
	data, err := ctrl.Service.Reservation.GetAll(date)
	if err != nil {
		helper.Responses(ctx, http.StatusInternalServerError, err.Error(), nil)
		ctx.Abort()
		return
	}
	helper.Responses(ctx, http.StatusOK, "Get Reservations success", data)
}

// @Summary Get Detail Reservation
// @Description Endpoint For Detail Reservation
// @Tags Reservation
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response{data=model.Reservation} "Get Summary Success"
// @Failure 500 {object} helper.Response "server error"
// @Router  /reservation/{id} [get]
func (ctrl *ControllerReservation) GetById(ctx *gin.Context) {
	var reservation model.Reservation
	id, _ := strconv.Atoi(ctx.Param("id"))
	reservation.ID = uint(id)
	err := ctrl.Service.Reservation.GetById(&reservation)
	if err != nil {
		helper.Responses(ctx, http.StatusInternalServerError, err.Error(), nil)
		ctx.Abort()
		return
	}
	// fmt.Println(reservation.TableID, "*****")
	helper.Responses(ctx, http.StatusOK, "Get Reservation success", reservation)
}

// @Summary Create a new Reservation
// @Description Create a new reservation.
// @Tags Reservation
// @Accept json
// @Produce json
// @Success 201 {string} helper.Response{data=model.Reservation} "Reservation successfully created"
// @Failure 400 {object} helper.Response "Invalid form data"
// @Failure 500 {object} helper.Response "Internal Server Error"
// @Router /reservation [post]
func (ctrl *ControllerReservation) Create(ctx *gin.Context) {
	var reservation model.Reservation
	err := ctx.ShouldBindJSON(&reservation)
	if err != nil {
		helper.Responses(ctx, http.StatusBadRequest, "Invalid Request", nil)
		ctx.Abort()
		return
	}
	err = ctrl.Service.Reservation.Create(&reservation)
	if err != nil {
		helper.Responses(ctx, http.StatusBadRequest, err.Error(), nil)
		ctx.Abort()
		return
	}
	helper.Responses(ctx, http.StatusCreated, "Create Reservation success", reservation)
}

// @Summary edit existing Reservation
// @Description edit existing reservation.
// @Tags Reservation
// @Accept json
// @Produce json
// @Success 201 {string} helper.Response{data=model.Reservation} "Reservation successfully updated"
// @Failure 400 {object} helper.Response "Invalid form data"
// @Failure 500 {object} helper.Response "Internal Server Error"
// @Router /reservation [put]
func (ctrl *ControllerReservation) Edit(ctx *gin.Context) {
	var reservation model.Reservation
	id, _ := strconv.Atoi(ctx.Param("id"))
	reservation.ID = uint(id)
	var form model.FormUpdate
	err := ctx.ShouldBindJSON(&form)
	if err != nil {
		helper.Responses(ctx, http.StatusBadRequest, "Invalid Request", nil)
		ctx.Abort()
		return
	}
	err = ctrl.Service.Reservation.Edit(&reservation, form)
	if err != nil {
		helper.Responses(ctx, http.StatusBadRequest, err.Error(), nil)
		ctx.Abort()
		return
	}
	helper.Responses(ctx, http.StatusOK, "Edit Reservation success", reservation)
}
