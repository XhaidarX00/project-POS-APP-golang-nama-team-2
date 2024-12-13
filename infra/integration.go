package infra

import (
	"project_pos_app/config"
	"project_pos_app/controller"
	"project_pos_app/database"
	"project_pos_app/helper"
	"project_pos_app/repository"
	"project_pos_app/service"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IntegrationContext struct {
	Cfg   config.Config
	DB    *gorm.DB
	Log   *zap.Logger
	Ctl   controller.AllController
	Cache database.Cache
}

func NewIntegrateContext() (*IntegrationContext, error) {

	errorHandler := func(err error) (*IntegrationContext, error) {
		return nil, err
	}

	config, err := config.SetConfig()
	if err != nil {
		return errorHandler(err)
	}

	log, err := helper.InitLog(config)
	if err != nil {
		return errorHandler(err)
	}

	db, err := database.SetDatabase(config)
	if err != nil {
		return errorHandler(err)
	}

	rdb := database.NewCache(config, 60*60)

	repo := repository.NewAllRepo(db, log)

	service := service.NewAllService(repo, log)

	handler := controller.NewAllController(service, log, &rdb)

	return &IntegrationContext{
		Cfg:   config,
		DB:    db,
		Log:   log,
		Ctl:   handler,
		Cache: rdb,
	}, nil
}
