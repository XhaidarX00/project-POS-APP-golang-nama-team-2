package accessservice

import (
	"project_pos_app/model"
	"project_pos_app/repository"

	"go.uber.org/zap"
)

type AccessService interface {
	GetAccessRepo(token string) ([]*model.ResponseAccess, error)
}

type accessService struct {
	Repo *repository.AllRepository
	Log  *zap.Logger
}

func NewAccessService(Repo *repository.AllRepository, Log *zap.Logger) AccessService {
	return &accessService{Repo, Log}
}

func (as *accessService) GetAccessRepo(token string) ([]*model.ResponseAccess, error) {

	access, err := as.Repo.Access.GetAccessRepo(token)
	if err != nil {
		return nil, err
	}

	return access, nil
}
