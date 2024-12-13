package examplerepository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ExampleRepository interface {
}

type exampleRepository struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewExampleRepo(DB *gorm.DB, Log *zap.Logger) ExampleRepository {
	return &exampleRepository{DB, Log}
}
