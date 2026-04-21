package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

type GormPositionRepo struct {
	db *gorm.DB
}

func NewGormPositionRepository(db *gorm.DB) domain.PositionRepo {
	return &GormPositionRepo{
		db: db,
	}
}

func (r *GormPositionRepo) Create(ctx context.Context, pstn *domain.Position) error {
	if pstn == nil {
		return errors.New("position cannot be nil")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return err
	}
	pstn.TenantID = tenantID
	err = r.db.WithContext(ctx).Create(pstn).Error
	if err != nil {
		if isDuplicateError(err) {
			if strings.Contains(err.Error(), "name") {
				return domain.ErrPositionNameExists
			}
		}
		return err
	}
	return nil
}

func (r *GormPositionRepo) GetByID(ctx context.Context, id uint) (*domain.Position, error) {
	if id == 0 {
		return nil, errors.New("invalid position id")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var position domain.Position
	err = r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&position).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPositionNotFound
		}
		return nil, err
	}
	return &position, nil
}

func (r *GormPositionRepo) GetByName(ctx context.Context, name string) (*domain.Position, error) {
	if name == "" {
		return nil, errors.New("position name cannot be empty")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var pstn domain.Position
	err = r.db.WithContext(ctx).Where("tenant_id = ? AND name = ?", tenantID, name).First(&pstn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPositionNotFound
		}
		return nil, err
	}
	return &pstn, nil
}

func (r *GormPositionRepo) List(ctx context.Context, page, limit int) ([]domain.Position, int64, error) {
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

	var positions []domain.Position
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&domain.Position{}).
		Where("tenant_id = ?", tenantID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&positions).Error; err != nil {
		return nil, 0, err
	}
	return positions, total, nil
}

func (r *GormPositionRepo) Update(ctx context.Context, pstn *domain.Position) error {
	if pstn == nil || pstn.ID == 0 {
		return errors.New("position cannot be nil or have zero id")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return err
	}
	existing, err := r.GetByID(ctx, pstn.ID)
	if err != nil {
		return err
	}
	pstn.TenantID = existing.TenantID
	err = r.db.WithContext(ctx).Model(&domain.Position{}).Where("id = ? AND tenant_id = ?", pstn.ID, tenantID).Updates(&pstn).Error
	if err != nil {
		if isDuplicateError(err) {
			if strings.Contains(err.Error(), "name") {
				return domain.ErrPositionNameExists
			}
		}
		return err
	}
	return nil
}

func (r *GormPositionRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid position id")
	}
	tenantID, err := tenatFromctx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Position{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrPositionNotFound
	}
	return nil
}
