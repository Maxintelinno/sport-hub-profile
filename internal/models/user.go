package models

import "time"

type User struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Phone        string    `gorm:"type:text;unique;not null" json:"phone"`
	Username     string    `gorm:"type:text;unique;not null" json:"username"`
	PasswordHash string    `gorm:"type:text;not null" json:"password_hash"`
	Fullname     string    `gorm:"type:text;not null" json:"fullname"`
	Province     string    `gorm:"type:varchar(100);not null" json:"province"`
	District     string    `gorm:"type:varchar(100);not null" json:"district"`
	Role         string    `gorm:"type:varchar(100);default:'user';not null" json:"role"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`
}
