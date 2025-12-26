package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID             uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User               User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LaundryID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"laundry_id"`
	Laundry            Laundry        `gorm:"foreignKey:LaundryID" json:"laundry,omitempty"`
	Status             string         `gorm:"type:varchar(50);not null;default:'pending';index" json:"status"`
	TotalPrice         float64        `gorm:"type:decimal(12,2);not null" json:"total_price"`
	DeliveryAddress    string         `gorm:"type:text;not null" json:"delivery_address"`
	Notes              string         `gorm:"type:text" json:"notes"`
	EstimatedPickupAt *time.Time     `json:"estimated_pickup_at,omitempty"`
	EstimatedDeliveryAt *time.Time    `json:"estimated_delivery_at,omitempty"`
	ActualPickupAt     *time.Time     `json:"actual_pickup_at,omitempty"`
	ActualDeliveryAt   *time.Time     `json:"actual_delivery_at,omitempty"`
	OrderServices      []OrderService `gorm:"foreignKey:OrderID" json:"order_services,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

type OrderService struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	Order       Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ServiceID   uuid.UUID `gorm:"type:uuid;not null" json:"service_id"`
	Service     Service   `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	ServiceName string    `gorm:"type:varchar(255);not null" json:"service_name"`
	Quantity    float64   `gorm:"type:decimal(10,2);not null" json:"quantity"`
	UnitPrice   float64   `gorm:"type:decimal(12,2);not null" json:"unit_price"`
	Unit        string    `gorm:"type:varchar(20);not null" json:"unit"`
	Subtotal    float64   `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	CreatedAt   time.Time `json:"created_at"`
}

func (os *OrderService) BeforeCreate(tx *gorm.DB) error {
	if os.ID == uuid.Nil {
		os.ID = uuid.New()
	}
	return nil
}

