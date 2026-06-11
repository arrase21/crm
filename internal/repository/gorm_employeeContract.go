package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

type GormEmployeeContractRepo struct {
	db *gorm.DB
}

func NewGormEmployeeContractRepository(db *gorm.DB) domain.EmployeeContractRepo {
	return &GormEmployeeContractRepo{
		db: db,
	}
}

func (r *GormEmployeeContractRepo) Create(ctx context.Context, empct *domain.EmployeeContract) error {
	if empct == nil {
		return errors.New("employee contract cannot be nil")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	empct.TenantID = tenantID
	return r.db.WithContext(ctx).Create(empct).Error
}

func (r *GormEmployeeContractRepo) GetByID(ctx context.Context, id uint) (*domain.EmployeeContract, error) {
	if id == 0 {
		return nil, errors.New("id cannot be zero")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var empct domain.EmployeeContract
	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Preload("ContractType").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&empct).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrContractNotFound
		}
		return nil, err
	}
	return &empct, nil
}

func (r *GormEmployeeContractRepo) GetByEmployeeID(ctx context.Context, employeeID uint) ([]domain.EmployeeContract, error) {
	if employeeID == 0 {
		return nil, errors.New("invalid employee id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var contracts []domain.EmployeeContract
	err = r.db.WithContext(ctx).
		Preload("ContractType").
		Where("tenant_id = ? AND employee_id = ?", tenantID, employeeID).
		Order("start_date DESC").
		Find(&contracts).Error
	if err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *GormEmployeeContractRepo) List(ctx context.Context, page, limit int) ([]domain.EmployeeContract, int64, error) {
	tenantID, err := tenantFromCtx(ctx)
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
	var contracts []domain.EmployeeContract
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&domain.EmployeeContract{}).
		Where("tenant_id = ?", tenantID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("Employee.User").
		Preload("ContractType").
		Where("tenant_id = ?", tenantID).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&contracts).Error; err != nil {
		return nil, 0, err
	}
	return contracts, total, nil
}

func (r *GormEmployeeContractRepo) Update(ctx context.Context, empct *domain.EmployeeContract) error {
	if empct == nil || empct.ID == 0 {
		return errors.New("contract cannot be nil or have zero id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	existing, err := r.GetByID(ctx, empct.ID)
	if err != nil {
		return err
	}
	empct.TenantID = existing.TenantID
	return r.db.WithContext(ctx).
		Model(&domain.EmployeeContract{}).
		Where("id = ? AND tenant_id = ?", empct.ID, tenantID).
		Updates(empct).Error
}

func (r *GormEmployeeContractRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid contract id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&domain.EmployeeContract{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrContractNotFound
	}
	return nil
}
