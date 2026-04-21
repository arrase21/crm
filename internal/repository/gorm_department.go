package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

type GormDepartmentRepo struct {
	db *gorm.DB
}

func NewGormDepartmentRepository(db *gorm.DB) domain.DepartmentRepo {
	return &GormDepartmentRepo{
		db: db,
	}
}

func (r *GormDepartmentRepo) Create(ctx context.Context, dept *domain.Department) error {
	if dept == nil {
		return errors.New("department cannot be nil")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return err
	}
	dept.TenantID = tenantID
	err = r.db.WithContext(ctx).Create(dept).Error
	if err != nil {
		if isDuplicateError(err) {
			if strings.Contains(err.Error(), "name") {
				return domain.ErrDepartmentNameExists
			}
			if strings.Contains(err.Error(), "code") {
				return domain.ErrDepartmentCodeExists
			}
		}
		return err
	}
	return nil
}

func (r *GormDepartmentRepo) GetByID(ctx context.Context, id uint) (*domain.Department, error) {
	if id == 0 {
		return nil, errors.New("invalid department id")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var dept domain.Department
	err = r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&dept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDepartmentNotFound
		}
		return nil, err
	}
	return &dept, nil
}

func (r *GormDepartmentRepo) GetByCode(ctx context.Context, code string) (*domain.Department, error) {
	if code == "" {
		return nil, errors.New("department code cannot be nil")
	}
	tenanID, err := tenatFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var dept domain.Department
	err = r.db.WithContext(ctx).Where("tenant_id = ? AND code = ?", tenanID, code).First(&dept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDepartmentNotFound
		}
		return nil, err
	}
	return &dept, nil
}

func (r *GormDepartmentRepo) GetByName(ctx context.Context, name string) (*domain.Department, error) {
	if name == "" {
		return nil, errors.New("department name cannot be nil")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var dept domain.Department
	err = r.db.WithContext(ctx).Where("tenant_id = ? AND name =?", tenantID, name).First(&dept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDepartmentNotFound
		}
		return nil, err
	}
	return &dept, nil
}

func (r *GormDepartmentRepo) List(ctx context.Context, page, limit int) ([]domain.Department, int64, error) {
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

	var dept []domain.Department
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&domain.Department{}).
		Where("tenant_id = ?", tenantID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&dept).Error; err != nil {
		return nil, 0, err
	}
	return dept, total, nil
}

func (r *GormDepartmentRepo) Update(ctx context.Context, dept *domain.Department) error {
	if dept == nil || dept.ID == 0 {
		return errors.New("user cannot be nil or have zero")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return err
	}
	existing, err := r.GetByID(ctx, dept.ID)
	if err != nil {
		return err
	}
	dept.TenantID = existing.TenantID
	err = r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", dept.ID, tenantID).Updates(dept).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormDepartmentRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid department id")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Department{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrDepartmentNotFound
	}
	return nil
}
