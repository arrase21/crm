package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

type GormContractTypeRepo struct {
	db *gorm.DB
}

func NewGormContractTypeRepository(db *gorm.DB) domain.ContractTypeRepo {
	return &GormContractTypeRepo{db: db}
}

func (r *GormContractTypeRepo) Create(ctx context.Context, ct *domain.ContractType) error {
	if ct == nil {
		return errors.New("contract type cannot be nil")
	}
	return r.db.WithContext(ctx).Create(ct).Error
}

func (r *GormContractTypeRepo) GetByID(ctx context.Context, id uint) (*domain.ContractType, error) {
	if id == 0 {
		return nil, errors.New("id cannot be zero")
	}
	var ct domain.ContractType
	err := r.db.WithContext(ctx).First(&ct, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrContractTypeNotFound
		}
		return nil, err
	}
	return &ct, nil
}

func (r *GormContractTypeRepo) GetByCountry(ctx context.Context, countryCode string) ([]domain.ContractType, error) {
	if countryCode == "" {
		return nil, errors.New("country code is required")
	}
	var types []domain.ContractType
	err := r.db.WithContext(ctx).
		Where("country_code = ? AND is_active = ?", countryCode, true).
		Find(&types).Error
	return types, err
}

func (r *GormContractTypeRepo) List(ctx context.Context, page, limit int) ([]domain.ContractType, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	var types []domain.ContractType
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&domain.ContractType{}).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&types).Error; err != nil {
		return nil, 0, err
	}
	return types, total, nil
}

func (r *GormContractTypeRepo) Update(ctx context.Context, ct *domain.ContractType) error {
	if ct == nil || ct.ID == 0 {
		return errors.New("contract type cannot be nil or have zero id")
	}
	return r.db.WithContext(ctx).Save(ct).Error
}

func (r *GormContractTypeRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid contract type id")
	}
	result := r.db.WithContext(ctx).Delete(&domain.ContractType{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrContractTypeNotFound
	}
	return nil
}
