package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

type GormPayrollRecordRepo struct {
	db *gorm.DB
}

func NewGormPayrollRecordRepository(db *gorm.DB) domain.PayrollRecordRepo {
	return &GormPayrollRecordRepo{db: db}
}

func (r *GormPayrollRecordRepo) Create(ctx context.Context, pr *domain.PayrollRecord) error {
	if pr == nil {
		return errors.New("payroll record cannot be nil")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	pr.TenantID = tenantID
	return r.db.WithContext(ctx).Create(pr).Error
}

func (r *GormPayrollRecordRepo) GetByID(ctx context.Context, id uint) (*domain.PayrollRecord, error) {
	if id == 0 {
		return nil, errors.New("id cannot be zero")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var pr domain.PayrollRecord
	err = r.db.WithContext(ctx).
		Preload("EmployeeContract").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&pr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPayrollRecordNotFound
		}
		return nil, err
	}
	return &pr, nil
}

func (r *GormPayrollRecordRepo) GetByContractID(ctx context.Context, contractID uint, page, limit int) ([]domain.PayrollRecord, int64, error) {
	if contractID == 0 {
		return nil, 0, errors.New("contract id is required")
	}
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
	var records []domain.PayrollRecord
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&domain.PayrollRecord{}).
		Where("tenant_id = ? AND employee_contract_id = ?", tenantID, contractID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).
		Preload("EmployeeContract.Employee.User").
		Where("tenant_id = ? AND employee_contract_id = ?", tenantID, contractID).
		Order("period_start DESC").
		Offset(offset).
		Limit(limit).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

func (r *GormPayrollRecordRepo) List(ctx context.Context, page, limit int) ([]domain.PayrollRecord, int64, error) {
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
	var records []domain.PayrollRecord
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&domain.PayrollRecord{}).
		Where("tenant_id = ?", tenantID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).
		Preload("EmployeeContract.Employee.User").
		Where("tenant_id = ?", tenantID).
		Order("period_start DESC").
		Offset(offset).
		Limit(limit).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}
