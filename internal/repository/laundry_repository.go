package repository

import (
	"laundry-go/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LaundryRepository interface {
	Create(laundry *models.Laundry) error
	FindByID(id uuid.UUID) (*models.Laundry, error)
	FindAll(search string, isOpen *bool, page, limit int) ([]models.Laundry, int64, error)
	FindByOwnerID(ownerID uuid.UUID) ([]models.Laundry, error)
	Update(laundry *models.Laundry) error
	Delete(id uuid.UUID) error
}

type laundryRepository struct {
	db *gorm.DB
}

func NewLaundryRepository(db *gorm.DB) LaundryRepository {
	return &laundryRepository{db: db}
}

func (r *laundryRepository) Create(laundry *models.Laundry) error {
	return r.db.Create(laundry).Error
}

func (r *laundryRepository) FindByID(id uuid.UUID) (*models.Laundry, error) {
	var laundry models.Laundry
	err := r.db.Preload("Services", "is_active = ?", true).Where("id = ?", id).First(&laundry).Error
	if err != nil {
		return nil, err
	}
	return &laundry, nil
}

func (r *laundryRepository) FindAll(search string, isOpen *bool, page, limit int) ([]models.Laundry, int64, error) {
	var laundries []models.Laundry
	var total int64

	query := r.db.Model(&models.Laundry{})

	if search != "" {
		query = query.Where("name ILIKE ? OR address ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if isOpen != nil {
		query = query.Where("is_open = ?", *isOpen)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&laundries).Error
	if err != nil {
		return nil, 0, err
	}

	return laundries, total, nil
}

func (r *laundryRepository) FindByOwnerID(ownerID uuid.UUID) ([]models.Laundry, error) {
	var laundries []models.Laundry
	err := r.db.Where("owner_id = ?", ownerID).Find(&laundries).Error
	return laundries, err
}

func (r *laundryRepository) Update(laundry *models.Laundry) error {
	return r.db.Save(laundry).Error
}

func (r *laundryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Laundry{}, id).Error
}

