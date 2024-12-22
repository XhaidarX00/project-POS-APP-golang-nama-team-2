package superadminservice

import (
	"project_pos_app/model"
	"project_pos_app/repository"

	"go.uber.org/zap"
)

type SuperadminService interface {
	ListDataAdmin() ([]*model.ResponseEmployee, error)
	UpdateSuperadmin(id int, admin *model.Superadmin) error
	UpdateAccessUser(id int, input *model.AccessPermission) error
	Logout(token string) error
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

func (ss *superadminService) UpdateSuperadmin(id int, admin *model.Superadmin) error {

	if err := ss.Repo.Superadmin.UpdateSuperadmin(id, admin); err != nil {
		return err
	}
	return nil
}

func (as *superadminService) UpdateAccessUser(id int, input *model.AccessPermission) error {

	err := as.Repo.Superadmin.UpdateAccessUser(id, input)
	if err != nil {
		return err
	}

	return nil
}

func (as *superadminService) Logout(token string) error {

	if err := as.Repo.Superadmin.Logout(token); err != nil {
		return err
	}

	return nil
}
