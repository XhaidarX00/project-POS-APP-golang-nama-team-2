package middleware

import (
	"net/http"
	"project_pos_app/helper"
	"project_pos_app/service"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AccessController struct {
	service *service.AllService
	log     *zap.Logger
}

func NewAccessController(service *service.AllService, log *zap.Logger) *AccessController {
	return &AccessController{service, log}
}

func (ac *AccessController) AccessMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			ctx.Abort()
			return
		}

		access, err := ac.service.Access.GetAccessRepo(token)
		if err != nil {
			helper.Responses(ctx, http.StatusForbidden, "Invalid or expired token", nil)
			ctx.Abort()
			return
		}

		isSuperAdmin := false
		for _, perm := range access {
			if perm.Role == "super_admin" {
				isSuperAdmin = true
				break
			}
		}

		if isSuperAdmin {
			ctx.Next()
			return
		}

		hasPermission := false

		route := ctx.Request.URL.Path
		corePath := extractCorePath(route)

		for _, perm := range access {

			if perm.Status && (perm.Permission == corePath) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			helper.Responses(ctx, http.StatusForbidden, "No valid access permissions for this route", nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func extractCorePath(route string) string {

	trimmed := strings.Trim(route, "/")

	parts := strings.Split(trimmed, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
