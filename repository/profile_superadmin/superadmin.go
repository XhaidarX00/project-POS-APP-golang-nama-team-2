package profilesuperadmin

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SuperadminRepo interface {
}

type superadminRepo struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewOrderRepo(DB *gorm.DB, Log *zap.Logger) SuperadminRepo {
	return &superadminRepo{DB, Log}
}

func (sr *superadminRepo) UpdateSuperAdmin() {

}
