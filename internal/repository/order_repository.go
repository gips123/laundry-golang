package repository

import (
	"laundry-go/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id uuid.UUID) (*models.Order, error)
	FindByUserID(userID uuid.UUID, status string, page, limit int) ([]models.Order, int64, error)
	FindByLaundryID(laundryID uuid.UUID, status string, page, limit int) ([]models.Order, int64, error)
	Update(order *models.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderServices").Preload("Laundry").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID uuid.UUID, status string, page, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Preload("OrderServices").Preload("Laundry").
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) FindByLaundryID(laundryID uuid.UUID, status string, page, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{}).Where("laundry_id = ?", laundryID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Preload("OrderServices").Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

