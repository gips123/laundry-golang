package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string     `gorm:"type:varchar(255);not null" json:"name"`
	Email        string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	Phone        string     `gorm:"type:varchar(20);not null" json:"phone"`
	Address      string     `gorm:"type:text;not null" json:"address"`
	Latitude     *float64   `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude    *float64   `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	Role         string     `gorm:"type:varchar(20);default:'customer'" json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

