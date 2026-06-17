package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

type GormAttendanceRepo struct {
	db *gorm.DB
}

func NewGormAttendanceRepository(db *gorm.DB) domain.AttendanceRepo {
	return &GormAttendanceRepo{db: db}
}

func (r *GormAttendanceRepo) Create(ctx context.Context, a *domain.Attendance) error {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	a.TenantID = tenantID
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *GormAttendanceRepo) GetByID(ctx context.Context, id uint) (*domain.Attendance, error) {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var a domain.Attendance
	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("attendance not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *GormAttendanceRepo) ListByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]domain.Attendance, int64, error) {
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
	var records []domain.Attendance
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Attendance{}).
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

func (r *GormAttendanceRepo) ListByDepartment(ctx context.Context, departmentID uint, page, limit int) ([]domain.Attendance, int64, error) {
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
	var records []domain.Attendance
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Attendance{}).
		Joins("JOIN employees ON employees.id = attendances.employee_id").
		Where("attendances.tenant_id = ? AND employees.department_id = ?", tenantID, departmentID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Joins("JOIN employees ON employees.id = attendances.employee_id").
		Where("attendances.tenant_id = ? AND employees.department_id = ?", tenantID, departmentID).
		Order("date DESC").
		Offset(offset).Limit(limit).
		Find(&records).Error
	return records, total, err
}

func (r *GormAttendanceRepo) List(ctx context.Context, page, limit int) ([]domain.Attendance, int64, error) {
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
	var records []domain.Attendance
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Attendance{}).
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

func (r *GormAttendanceRepo) Update(ctx context.Context, a *domain.Attendance) error {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, a.ID).
		Updates(a).Error
}

func (r *GormAttendanceRepo) Delete(ctx context.Context, id uint) error {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&domain.Attendance{}).Error
}
