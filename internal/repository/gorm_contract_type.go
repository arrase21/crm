package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm/internal/cache"
	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

type GormContractTypeRepo struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewGormContractTypeRepository(db *gorm.DB, c *cache.Cache) domain.ContractTypeRepo {
	return &GormContractTypeRepo{
		db:    db,
		cache: c,
	}
}

func (r *GormContractTypeRepo) Create(ctx context.Context, ct *domain.ContractType) error {
	if ct == nil {
		return errors.New("contract type cannot be nil")
	}
	if r.cache != nil {
		r.cache.Del("contract_type_country:" + ct.CountryCode)
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
	if r.cache != nil {
		if v, ok := r.cache.Get("contract_type_country:" + countryCode); ok {
			return v.([]domain.ContractType), nil
		}
	}
	var types []domain.ContractType
	err := r.db.WithContext(ctx).
		Where("country_code = ? AND is_active = ?", countryCode, true).
		Find(&types).Error
	if err != nil {
		return nil, err
	}
	if r.cache != nil {
		r.cache.Set("contract_type_country:"+countryCode, types)
	}
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
	if r.cache != nil {
		r.cache.Del("contract_type_country:" + ct.CountryCode)
	}
	return r.db.WithContext(ctx).Save(ct).Error
}

func (r *GormContractTypeRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid contract type id")
	}
	var old domain.ContractType
	if err := r.db.WithContext(ctx).First(&old, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrContractTypeNotFound
		}
		return err
	}
	result := r.db.WithContext(ctx).Delete(&domain.ContractType{}, id)
	if result.Error != nil {
		return result.Error
	}
	if r.cache != nil {
		r.cache.Del("contract_type_country:" + old.CountryCode)
	}
	return nil
}
