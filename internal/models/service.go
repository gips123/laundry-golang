package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	LaundryID        uuid.UUID `gorm:"type:uuid;not null;index" json:"laundry_id"`
	Laundry          Laundry   `gorm:"foreignKey:LaundryID" json:"laundry,omitempty"`
	Name             string    `gorm:"type:varchar(255);not null" json:"name"`
	Description      string    `gorm:"type:text" json:"description"`
	Price            float64   `gorm:"type:decimal(12,2);not null" json:"price"`
	Unit             string    `gorm:"type:varchar(20);not null" json:"unit"`
	EstimatedTimeHours int    `gorm:"type:integer;not null" json:"estimated_time_hours"`
	Category         string    `gorm:"type:varchar(50);not null" json:"category"`
	IsActive         bool      `gorm:"default:true" json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (s *Service) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

