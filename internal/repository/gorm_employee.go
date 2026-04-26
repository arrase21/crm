package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

type GormEmployeeRepo struct {
	db *gorm.DB
}

func NewGormEmployeeRepository(db *gorm.DB) domain.EmployeeRepo {
	return &GormEmployeeRepo{
		db: db,
	}
}

func (r *GormEmployeeRepo) Create(ctx context.Context, emp *domain.Employee) error {
	if emp == nil {
		return errors.New("employee cannot be nil")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	emp.TenantID = tenantID
	err = r.db.WithContext(ctx).Create(emp).Error
	if err != nil {
		return err // validar si es nuevo o tiene info duplicada
	}
	return nil
}

func (r *GormEmployeeRepo) GetByID(ctx context.Context, id uint) (*domain.Employee, error) {
	if id == 0 {
		return nil, errors.New("id cannot be nil or zero")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var emp domain.Employee
	err = r.db.WithContext(ctx).Preload("User").Preload("Department").
		Preload("Position"). //anadir un preload para contratos
		Where("tenant_id = ? AND id = ?", tenantID, id).First(&emp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, err
	}
	return &emp, nil
}

func (r *GormEmployeeRepo) GetByUserID(ctx context.Context, userID uint) (*domain.Employee, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var emp domain.Employee
	err = r.db.WithContext(ctx).Preload("User").Preload("Department").
		Preload("Position"). //anadir un preload para contratos
		Where("tenant_id = ? and userID = ?", tenantID, userID).First(&emp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, err
	}
	return &emp, nil
}

func (r *GormEmployeeRepo) List(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	var emp []domain.Employee
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.Employee{}).Where("tenant_id = ?", tenantID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Preload("User").Preload("Department").Preload("Position").
		Where("tenant_id = ?", tenantID).Order("id DESC").Offset(offset).Limit(limit).Find(&emp).Error; err != nil {
		return nil, 0, err
	}
	return emp, total, nil
}

func (r *GormEmployeeRepo) ListActive(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	var emp []domain.Employee
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Employee{}).
		Where("tenant_id = ? AND is_active = ?)", tenantID, true).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Preload("User").
		// Preload("Contracts", "is_active = ?", true). // Solo contratos activos
		// Preload("Contracts.ContractType").
		Where("tenant_id = ? and is_active = ?", tenantID, true).Order("id ASC").Offset(offset).Limit(limit).Find(&emp).Error; err != nil {
		return nil, 0, err
	}
	return emp, total, nil
}

func (r *GormEmployeeRepo) Update(ctx context.Context, emp *domain.Employee) error {
	if emp == nil || emp.ID == 0 {
		return errors.New("employee cannot be nil or have 0")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	emp.TenantID = tenantID
	err = r.db.WithContext(ctx).Save(emp).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormEmployeeRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid employee id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&domain.Employee{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEmployeeNotFound
	}
	return nil
}
