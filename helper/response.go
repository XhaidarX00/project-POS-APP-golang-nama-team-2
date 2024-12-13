package helper

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  bool
	Message string
	Data    any `json:"data,omitempty"`
}

func Responses(c *gin.Context, status int, massage string, data any) {
	c.JSON(status, Response{
		Status:  false,
		Message: massage,
		Data:    data,
	})
}
