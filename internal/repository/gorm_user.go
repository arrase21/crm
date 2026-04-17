package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

type GormUserRepo struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) domain.UserRepo {
	return &GormUserRepo{
		db: db,
	}
}

func (r *GormUserRepo) Create(ctx context.Context, usr *domain.User) error {
	if usr == nil {
		return errors.New("user cannot be nil")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return nil
	}
	usr.TenantID = tenantID
	err = r.db.WithContext(ctx).Create(usr).Error
	if err != nil {
		if isDuplicateError(err) {
			if strings.Contains(err.Error(), "dni") {
				return domain.ErrDniAlreadyExist
			}
			if strings.Contains(err.Error(), "email") {
				return domain.ErrEmailAlreadyExist
			}
			if strings.Contains(err.Error(), "phone") {
				return domain.ErrPhoneAlreadyExist
			}
			return errors.New("duplicate entry")
		}
		return err
	}
	return nil
}

func (r *GormUserRepo) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user id")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var user domain.User
	err = r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepo) GetByDNI(ctx context.Context, dni string) (*domain.User, error) {
	if dni == "" {
		return nil, errors.New("dni cannot be empty")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var user domain.User
	err = r.db.WithContext(ctx).Where("tenant_id = ? AND dni = ?", tenantID, dni).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepo) List(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	var users []domain.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("tenant_id = ?", tenantID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("id DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *GormUserRepo) Update(ctx context.Context, usr *domain.User) error {
	if usr == nil || usr.ID == 0 {
		return errors.New("user cannot be nil or have zero id")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return err
	}
	existing, err := r.GetByID(ctx, usr.ID)
	if err != nil {
		return err
	}
	usr.TenantID = existing.TenantID

	err = r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ? AND tenant_id = ?", usr.ID, tenantID).Updates(usr).Error
	if err != nil {
		if isDuplicateError(err) {
			if strings.Contains(err.Error(), "dni") {
				return domain.ErrDniAlreadyExist
			}
			if strings.Contains(err.Error(), "email") {
				return domain.ErrEmailAlreadyExist
			}
			if strings.Contains(err.Error(), "phone") {
				return domain.ErrPhoneAlreadyExist
			}
			return errors.New("duplicate entry")
		}
		return err
	}
	return nil
}

func (r *GormUserRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid user id")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
