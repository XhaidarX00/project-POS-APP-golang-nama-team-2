package helper

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status     int
	Message    string
	Data       any `json:"data,omitempty"`
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

func Responses(c *gin.Context, status int, massage string, data any) {
	c.JSON(status, Response{
		Status:  status,
		Message: massage,
		Data:    data,
	})
}

func ResponsePagination(c *gin.Context, data interface{}, message string, page, limit, totalItems, totalPages, httpStatusCode int) {
	c.JSON(httpStatusCode, Response{
		Status:     httpStatusCode,
		Message:    message,
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Data:       data,
	})
}
