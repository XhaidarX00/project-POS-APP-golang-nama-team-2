package authservice

import (
	"project_pos_app/model"
	"project_pos_app/repository"

	"go.uber.org/zap"
)

type AuthService interface {
	Login(login *model.Login, ipAddress string) (*model.Session, string, error)
}

type authService struct {
	repo *repository.AllRepository
	log  *zap.Logger
}

func NewManagementVoucherService(repo *repository.AllRepository, log *zap.Logger) AuthService {
	return &authService{repo, log}
}

func (as *authService) Login(login *model.Login, ipAddress string) (*model.Session, string, error) {

	session, idKey, err := as.repo.Auth.Login(login, ipAddress)
	if err != nil {
		return nil, "", err
	}

	return session, idKey, nil
}
