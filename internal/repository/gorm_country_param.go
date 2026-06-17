package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm/internal/cache"
	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

type GormCountryParamRepo struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewGormCountryParamRepository(db *gorm.DB, c *cache.Cache) domain.CountryParamRepo {
	return &GormCountryParamRepo{
		db:    db,
		cache: c,
	}
}

func (r *GormCountryParamRepo) Create(ctx context.Context, cp *domain.CountryParam) error {
	if cp == nil {
		return errors.New("country param cannot be nil")
	}
	if r.cache != nil {
		r.cache.Del("country_param:" + cp.CountryCode)
	}
	return r.db.WithContext(ctx).Create(cp).Error
}

func (r *GormCountryParamRepo) GetByID(ctx context.Context, id uint) (*domain.CountryParam, error) {
	if id == 0 {
		return nil, errors.New("id cannot be zero")
	}
	var cp domain.CountryParam
	err := r.db.WithContext(ctx).First(&cp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCountryParamNotFound
		}
		return nil, err
	}
	return &cp, nil
}

func (r *GormCountryParamRepo) GetByCountryCode(ctx context.Context, countryCode string) (*domain.CountryParam, error) {
	if countryCode == "" {
		return nil, errors.New("country code is required")
	}
	if r.cache != nil {
		if v, ok := r.cache.Get("country_param:" + countryCode); ok {
			return v.(*domain.CountryParam), nil
		}
	}
	var cp domain.CountryParam
	err := r.db.WithContext(ctx).Where("country_code = ?", countryCode).First(&cp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCountryParamNotFound
		}
		return nil, err
	}
	if r.cache != nil {
		r.cache.Set("country_param:"+countryCode, &cp)
	}
	return &cp, nil
}

func (r *GormCountryParamRepo) List(ctx context.Context, page, limit int) ([]domain.CountryParam, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	var params []domain.CountryParam
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&domain.CountryParam{}).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&params).Error; err != nil {
		return nil, 0, err
	}
	return params, total, nil
}

func (r *GormCountryParamRepo) Update(ctx context.Context, cp *domain.CountryParam) error {
	if cp == nil || cp.ID == 0 {
		return errors.New("country param cannot be nil or have zero id")
	}
	if r.cache != nil {
		r.cache.Del("country_param:" + cp.CountryCode)
	}
	return r.db.WithContext(ctx).Save(cp).Error
}

func (r *GormCountryParamRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid country param id")
	}
	var old domain.CountryParam
	if err := r.db.WithContext(ctx).First(&old, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrCountryParamNotFound
		}
		return err
	}
	result := r.db.WithContext(ctx).Delete(&domain.CountryParam{}, id)
	if result.Error != nil {
		return result.Error
	}
	if r.cache != nil {
		r.cache.Del("country_param:" + old.CountryCode)
	}
	return nil
}
