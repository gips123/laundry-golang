package repository

import (
	"laundry-go/internal/models"

	"gorm.io/gorm"
)

type OrderServiceRepository interface {
	Create(orderService *models.OrderService) error
	CreateBatch(orderServices []models.OrderService) error
	FindByOrderID(orderID string) ([]models.OrderService, error)
}

type orderServiceRepository struct {
	db *gorm.DB
}

func NewOrderServiceRepository(db *gorm.DB) OrderServiceRepository {
	return &orderServiceRepository{db: db}
}

func (r *orderServiceRepository) Create(orderService *models.OrderService) error {
	return r.db.Create(orderService).Error
}

func (r *orderServiceRepository) CreateBatch(orderServices []models.OrderService) error {
	return r.db.Create(&orderServices).Error
}

func (r *orderServiceRepository) FindByOrderID(orderID string) ([]models.OrderService, error) {
	var orderServices []models.OrderService
	err := r.db.Where("order_id = ?", orderID).Find(&orderServices).Error
	return orderServices, err
}

