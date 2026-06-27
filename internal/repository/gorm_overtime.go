package repository

import (
	"context"
	"errors"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

type GormOvertimeRepo struct {
	db *gorm.DB
}

func NewGormOvertimeRepository(db *gorm.DB) domain.OvertimeRepo {
	return &GormOvertimeRepo{db: db}
}

func (r *GormOvertimeRepo) Create(ctx context.Context, o *domain.Overtime) error {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	o.TenantID = tenantID
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *GormOvertimeRepo) GetByID(ctx context.Context, id uint) (*domain.Overtime, error) {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var o domain.Overtime
	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&o).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOvertimeNotFound
		}
		return nil, err
	}
	return &o, nil
}

func (r *GormOvertimeRepo) ListByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]domain.Overtime, int64, error) {
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
	var records []domain.Overtime
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Overtime{}).
		Where("tenant_id = ? AND employee_id = ?", tenantID, employeeID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Where("tenant_id = ? AND employee_id = ?", tenantID, employeeID).
		Order("date DESC").
		Offset(offset).Limit(limit).
		Find(&records).Error
	return records, total, err
}

func (r *GormOvertimeRepo) ListByEmployeeAndPeriod(ctx context.Context, employeeID uint, start, end time.Time) ([]domain.Overtime, error) {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var records []domain.Overtime
	err = r.db.WithContext(ctx).
		Where("tenant_id = ? AND employee_id = ? AND date >= ? AND date <= ? AND approved = ?", tenantID, employeeID, start, end, true).
		Find(&records).Error
	return records, err
}

func (r *GormOvertimeRepo) ListByDepartment(ctx context.Context, departmentID uint, page, limit int) ([]domain.Overtime, int64, error) {
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
	var records []domain.Overtime
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Overtime{}).
		Joins("JOIN employees ON employees.id = overtimes.employee_id").
		Where("overtimes.tenant_id = ? AND employees.department_id = ?", tenantID, departmentID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Joins("JOIN employees ON employees.id = overtimes.employee_id").
		Where("overtimes.tenant_id = ? AND employees.department_id = ?", tenantID, departmentID).
		Order("date DESC").
		Offset(offset).Limit(limit).
		Find(&records).Error
	return records, total, err
}

func (r *GormOvertimeRepo) List(ctx context.Context, page, limit int) ([]domain.Overtime, int64, error) {
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
	var records []domain.Overtime
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Overtime{}).
		Where("tenant_id = ?", tenantID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Where("tenant_id = ?", tenantID).
		Order("date DESC").
		Offset(offset).Limit(limit).
		Find(&records).Error
	return records, total, err
}

func (r *GormOvertimeRepo) Update(ctx context.Context, o *domain.Overtime) error {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, o.ID).
		Updates(o).Error
}

func (r *GormOvertimeRepo) Delete(ctx context.Context, id uint) error {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&domain.Overtime{}).Error
}
