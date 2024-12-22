package accessrepository

import (
	"project_pos_app/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AccessRepository interface {
	GetAccessRepo(token string) ([]*model.ResponseAccess, error)
}

type accessRepository struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewAccessRepository(DB *gorm.DB, Log *zap.Logger) AccessRepository {
	return &accessRepository{DB, Log}
}

func (ar *accessRepository) GetAccessRepo(token string) ([]*model.ResponseAccess, error) {

	access := []*model.ResponseAccess{}

	err := ar.DB.Table("access_permissions AS ap").Select("ap.user_id, p.name AS permission, ap.status, u.role").
		Joins("JOIN permissions AS p ON p.id = ap.permission_id").
		Joins("JOIN users AS u ON ap.user_id = u.id").
		Joins("JOIN sessions AS s ON u.id = s.user_id").
		Where("s.token = ? AND ap.status = ?", token, true).Find(&access).Error

	if err != nil {
		return nil, err
	}

	return access, nil
}
