package profilesuperadmin

import (
	"fmt"
	"project_pos_app/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SuperadminRepo interface {
	ListDataAdmin() ([]*model.ResponseEmployee, error)
}

type superadminRepo struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewSuperadmin(DB *gorm.DB, Log *zap.Logger) SuperadminRepo {
	return &superadminRepo{DB, Log}
}

func (sr *superadminRepo) ListDataAdmin() ([]*model.ResponseEmployee, error) {

	admin := []*model.ResponseEmployee{}
	result := sr.DB.Table("employees AS e").Select("e.name, a.email").Where("a.role = ?", "admin").Where("a.deleted_at is NULL").
		Joins("JOIN users AS a ON a.id = e.user_id").Scan(&admin)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("admin not found")
	}

	return admin, nil

}
