package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

type GormRoleRepo struct {
	db *gorm.DB
}

func NewGormRoleRepository(db *gorm.DB) domain.RoleRepo {
	return &GormRoleRepo{db: db}
}

func (r *GormRoleRepo) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *GormRoleRepo) List(ctx context.Context) ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.WithContext(ctx).Find(&roles).Error
	return roles, err
}

func (r *GormRoleRepo) GetUserRoles(ctx context.Context, userID uint) ([]string, error) {
	var roleNames []string
	err := r.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Pluck("name", &roleNames).Error
	if roleNames == nil {
		roleNames = []string{}
	}
	return roleNames, err
}

func (r *GormRoleRepo) AssignRole(ctx context.Context, userID, roleID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		FirstOrCreate(&domain.UserRole{UserID: userID, RoleID: roleID}).Error
}

func (r *GormRoleRepo) UnassignRole(ctx context.Context, userID, roleID uint) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&domain.UserRole{})
	if result.RowsAffected == 0 {
		return errors.New("role not assigned to user")
	}
	return result.Error
}
