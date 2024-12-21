package dashboardcontroller

import (
	"bytes"
	"fmt"
	"net/http"
	"project_pos_app/helper"
	"project_pos_app/model"
	"project_pos_app/service"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type ControllerDashboard struct {
	Service *service.AllService
	Log     *zap.Logger
}

func NewControllerDashboard(service *service.AllService, log *zap.Logger) ControllerDashboard {
	return ControllerDashboard{Service: service, Log: log}
}

// @Summary Get Popular Product
// @Description Endpoint For Popular Product Dashboard
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response{data=[]model.Product} "Get Summary Success"
// @Failure 500 {object} helper.Response "server error"
// @Router  /api/dashboard/popular [get]
func (ctrl *ControllerDashboard) GetPopularProduct(ctx *gin.Context) {
	data, err := ctrl.Service.Dashboard.GetPopularProduct()
	if err != nil {
		helper.Responses(ctx, http.StatusInternalServerError, err.Error(), nil)
		ctx.Abort()
		return
	}
	helper.Responses(ctx, http.StatusOK, "Get Popular Product success", data)
}

// @Summary Get New Product
// @Description Endpoint For New Product Dashboard
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response{data=[]model.Product} "Get Summary Success"
// @Failure 500 {object} helper.Response "server error"
// @Router  /api/dashboard/new [get]
func (ctrl *ControllerDashboard) GetNewProduct(ctx *gin.Context) {
	data, err := ctrl.Service.Dashboard.GetNewProduct()
	if err != nil {
		helper.Responses(ctx, http.StatusInternalServerError, err.Error(), nil)
		ctx.Abort()
		return
	}
	helper.Responses(ctx, http.StatusOK, "Get New Product success", data)
}

// @Summary Get Summary
// @Description Endpoint For Summary Dashboard
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response{data=model.Summary} "Get Summary Success"
// @Failure 500 {object} helper.Response "server error"
// @Router  /api/dashboard/summary [get]
func (ctrl *ControllerDashboard) GetSummary(ctx *gin.Context) {
	var summary model.Summary
	err := ctrl.Service.Dashboard.GetSummary(&summary)
	if err != nil {
		helper.Responses(ctx, http.StatusInternalServerError, err.Error(), nil)
		ctx.Abort()
		return
	}
	helper.Responses(ctx, http.StatusOK, "Get Summary success", summary)
}

// generateExcelReport menggenerate file excel dan mengirimkannya sebagai download
// @Summary Generate an Excel report
// @Description Generate an Excel report and return it as an attachment
// @Tags Dashboard
// @Accept  json
// @Produce  application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success 200 {file} file "Order_Report.xlsx"
// @Failure 500 {object} helper.Response "server error"
// @Router /api/dashboard/report [get]
func (ctrl *ControllerDashboard) GetReport(ctx *gin.Context) {
	var reports []model.ReportExcel
	err := ctrl.Service.Dashboard.GetReport(&reports)
	if err != nil {
		helper.Responses(ctx, http.StatusInternalServerError, err.Error(), nil)
		ctx.Abort()
		return
	}
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	val := reflect.ValueOf(model.ReportExcel{})
	var data [][]interface{} = [][]interface{}{}
	var header []interface{}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		header = append(header, field.Name)
	}
	data = append(data, header)
	for i, report := range reports {
		var temp []interface{}
		report.No = i + 1
		val := reflect.ValueOf(report)
		// fmt.Println(val, "------", val.NumField(), val.String())
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			temp = append(temp, field)
		}
		data = append(data, temp)
	}
	for i, rowData := range data {
		cell, err := excelize.CoordinatesToCellName(1, i+1)
		if err != nil {
			ctrl.Log.Error("Failed to excelize", zap.Error(err))
			// return errors.New(" Internal Server Error")
			helper.Responses(ctx, http.StatusInternalServerError, err.Error(), nil)
		}
		f.SetSheetRow("Sheet1", cell, &rowData)
	}
	// if err := f.SaveAs("Order_Report.xlsx"); err != nil {
	// 	ctrl.Log.Error("Failed to save excel", zap.Error(err))
	// 	// return errors.New(" Internal Server Error")
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Excel file"})
	// }
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		helper.Responses(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=Order_Report.xlsx")
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	// helper.Responses(ctx, http.StatusOK, "Get Report success", nil)
}
