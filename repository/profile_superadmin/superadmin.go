package profilesuperadmin

import (
	"fmt"
	"project_pos_app/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SuperadminRepo interface {
	ListDataAdmin() ([]*model.ResponseEmployee, error)
	UpdateSuperadmin(id int, admin *model.Superadmin) error
	UpdateAccessUser(id int, input *model.AccessPermission) error
	Logout(token string) error
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

func (sr *superadminRepo) UpdateSuperadmin(id int, admin *model.Superadmin) error {

	tx := sr.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updateData := map[string]interface{}{
		"email": admin.User.Email,
	}

	if admin.User.Password != "" {
		updateData["password"] = admin.User.Password
	}

	if err := tx.Model(&model.User{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Superadmin{}).Where("user_id = ?", id).
		Updates(admin).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (ar *superadminRepo) UpdateAccessUser(id int, input *model.AccessPermission) error {

	var accessPermission model.AccessPermission
	err := ar.DB.Where("user_id = ? AND permission_id = ?", id, input.PermissionID).First(&accessPermission).Error
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("access permission not found for user_id %d and permission_id %d", input.UserID, input.PermissionID)
		}
		return err
	}

	err = ar.DB.Model(&model.AccessPermission{}).
		Where("user_id = ? AND permission_id = ?", id, input.PermissionID).
		Update("status", input.Status).Error

	if err != nil {
		return err
	}

	return nil
}

func (ar *superadminRepo) Logout(token string) error {

	updateResult := ar.DB.Model(&model.Session{}).Where("token = ?", token).Update("token", nil)

	if updateResult.Error != nil {
		return fmt.Errorf("failed to logout: %w", updateResult.Error)
	}

	return nil
}
