package models

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Laundry struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OwnerID            uuid.UUID  `gorm:"type:uuid;not null;index" json:"owner_id"`
	Owner              User       `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Name               string     `gorm:"type:varchar(255);not null" json:"name"`
	Description        string     `gorm:"type:text" json:"description"`
	Address            string     `gorm:"type:text;not null" json:"address"`
	Latitude         *float64    `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude          *float64    `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	ImageURL           string      `gorm:"type:text" json:"image_url,omitempty"`
	Rating             float64     `gorm:"type:decimal(3,2);default:0.0" json:"rating"`
	ReviewCount        int         `gorm:"default:0" json:"review_count"`
	IsOpen             bool        `gorm:"default:true" json:"is_open"`
	OperatingHoursOpen TimeOnly    `gorm:"type:time;not null" json:"operating_hours_open"`
	OperatingHoursClose TimeOnly   `gorm:"type:time;not null" json:"operating_hours_close"`
	Services           []Service   `gorm:"foreignKey:LaundryID" json:"services,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}

func (l *Laundry) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

// TimeOnly is a custom type for TIME in PostgreSQL
type TimeOnly string

func (t *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		*t = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*t = TimeOnly(v)
	case string:
		*t = TimeOnly(v)
	}
	return nil
}

func (t TimeOnly) Value() (driver.Value, error) {
	return string(t), nil
}

