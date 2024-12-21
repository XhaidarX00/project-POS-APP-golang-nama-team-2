package superadminservice

import (
	"project_pos_app/model"
	"project_pos_app/repository"

	"go.uber.org/zap"
)

type SuperadminService interface {
	ListDataAdmin() ([]*model.ResponseEmployee, error)
}

type superadminService struct {
	Repo *repository.AllRepository
	Log  *zap.Logger
}

func NewSuperadminService(Repo *repository.AllRepository, Log *zap.Logger) SuperadminService {
	return &superadminService{Repo, Log}
}

func (ss *superadminService) ListDataAdmin() ([]*model.ResponseEmployee, error) {

	admins, err := ss.Repo.Superadmin.ListDataAdmin()
	if err != nil {
		return nil, err
	}

	return admins, err
}
