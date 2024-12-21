package dashboardservice

import (
	"project_pos_app/model"
	"project_pos_app/repository"

	"go.uber.org/zap"
)

type ServiceDashboard interface {
	GetPopularProduct() ([]model.Product, error)
	GetNewProduct() ([]model.Product, error)
	GetSummary(summary *model.Summary) error
	GetReport(report *[]model.ReportExcel) error
}

type serviceDashboard struct {
	Repo *repository.AllRepository
	Log  *zap.Logger
}

func NewRevenueService(repo *repository.AllRepository, log *zap.Logger) ServiceDashboard {
	return &serviceDashboard{
		Repo: repo,
		Log:  log,
	}
}

func (s *serviceDashboard) GetPopularProduct() ([]model.Product, error) {
	return s.Repo.Dashboard.FindPopularProduct()
}
func (s *serviceDashboard) GetNewProduct() ([]model.Product, error) {
	return s.Repo.Dashboard.FindNewProduct()
}
func (s *serviceDashboard) GetSummary(summary *model.Summary) error {
	err := s.Repo.Dashboard.FindSummary(summary)
	if err != nil {
		return err
	}
	return nil
}
func (s *serviceDashboard) GetReport(report *[]model.ReportExcel) error {
	err := s.Repo.Dashboard.FindReport(report)
	if err != nil {
		return err
	}

	return nil
}
