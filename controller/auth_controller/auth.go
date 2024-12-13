package authcontroller

import (
	"net/http"
	"project_pos_app/database"
	"project_pos_app/helper"
	"project_pos_app/model"
	"project_pos_app/service"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

type AuthHadler struct {
	Service *service.AllService
	Log     *zap.Logger
	Cacher  *database.Cache
}

func NewUserHandler(service *service.AllService, log *zap.Logger, rdb *database.Cache) AuthHadler {
	return AuthHadler{
		Service: service,
		Log:     log,
		Cacher:  rdb,
	}
}

func (auth *AuthHadler) Login(c *gin.Context) {
	login := model.Login{}
	ipAddress := c.ClientIP()

	err := c.ShouldBindJSON(&login)
	if err != nil {
		auth.Log.Error("Invalid payload", zap.Error(err))
		helper.Responses(c, http.StatusInternalServerError, "Invalid Payload: "+err.Error(), nil)
		return
	}

	session, idKey, err := auth.Service.Auth.Login(&login, ipAddress)
	if err != nil {
		auth.Log.Error("Failed to Login"+err.Error(), zap.Error(err))
		helper.Responses(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token := session.Token
	IDKEY := idKey

	auth.Log.Info("Saving token to Redis", zap.String("IDKEY", IDKEY), zap.String("token", token))

	err = auth.Cacher.Set(IDKEY, token)
	if err != nil {
		helper.Responses(c, http.StatusInternalServerError, err.Error(), nil)
	}

	helper.Responses(c, http.StatusOK, "successfully login", session)

}
