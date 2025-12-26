package repository

import (
	"laundry-go/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(service *models.Service) error
	FindByID(id uuid.UUID) (*models.Service, error)
	FindByLaundryID(laundryID uuid.UUID) ([]models.Service, error)
	Update(service *models.Service) error
	Delete(id uuid.UUID) error
	GetPriceRange(laundryID uuid.UUID) (minPrice, maxPrice float64, err error)
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(service *models.Service) error {
	return r.db.Create(service).Error
}

func (r *serviceRepository) FindByID(id uuid.UUID) (*models.Service, error) {
	var service models.Service
	err := r.db.Where("id = ?", id).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) FindByLaundryID(laundryID uuid.UUID) ([]models.Service, error) {
	var services []models.Service
	err := r.db.Where("laundry_id = ? AND is_active = ?", laundryID, true).Find(&services).Error
	return services, err
}

func (r *serviceRepository) Update(service *models.Service) error {
	return r.db.Save(service).Error
}

func (r *serviceRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Service{}, id).Error
}

func (r *serviceRepository) GetPriceRange(laundryID uuid.UUID) (minPrice, maxPrice float64, err error) {
	var result struct {
		MinPrice float64
		MaxPrice float64
	}

	err = r.db.Model(&models.Service{}).
		Where("laundry_id = ? AND is_active = ?", laundryID, true).
		Select("MIN(price) as min_price, MAX(price) as max_price").
		Scan(&result).Error

	if err != nil {
		return 0, 0, err
	}

	return result.MinPrice, result.MaxPrice, nil
}

