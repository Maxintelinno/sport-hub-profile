package models

import (
	"time"
)

type Field struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OwnerID       string    `gorm:"type:uuid;not null" json:"owner_id"`
	Name          string    `gorm:"type:varchar(150);not null" json:"name"`
	SportType     string    `gorm:"type:varchar(50);not null" json:"sport_type"`
	PricePerHour  float64   `gorm:"type:numeric(10,2);not null" json:"price_per_hour"`
	OpenTime      string    `gorm:"type:time;not null" json:"open_time"`
	CloseTime     string    `gorm:"type:time;not null" json:"close_time"`
	Province      string    `gorm:"type:varchar(100);not null" json:"province"`
	District      string    `gorm:"type:varchar(100);not null" json:"district"`
	AddressLine   string    `gorm:"type:text;not null" json:"address_line"`
	Description   string    `gorm:"type:text" json:"description"`
	Status        string    `gorm:"type:varchar(20);default:'draft';not null" json:"status"`
	ThumbnailURL  string    `gorm:"type:varchar(255)" json:"thumbnail_url"`
	Latitude      float64   `gorm:"type:numeric(10,7)" json:"latitude"`
	Longitude     float64   `gorm:"type:numeric(10,7)" json:"longitude"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Courts        []FieldCourt `gorm:"foreignKey:FieldID" json:"courts"`
}

type FieldCourt struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FieldID      string    `gorm:"type:uuid;not null" json:"field_id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	PricePerHour float64   `gorm:"type:numeric(10,2);not null" json:"price_per_hour"`
	Status       string    `gorm:"type:varchar(20);default:'active';not null" json:"status"`
	Capacity     int       `gorm:"type:int4" json:"capacity"`
	CourtType    string    `gorm:"type:varchar(50)" json:"court_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
