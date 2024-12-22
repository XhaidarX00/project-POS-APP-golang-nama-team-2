package middleware

import (
	"project_pos_app/service"

	"go.uber.org/zap"
)

type AllHandler struct {
	Access AccessController
}

func NewMiddleware(service *service.AllService, Log *zap.Logger) *AllHandler {
	return &AllHandler{
		Access: *NewAccessController(service, Log),
	}
}
