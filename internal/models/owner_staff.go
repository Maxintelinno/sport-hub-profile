package models

import "time"

type OwnerStaff struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OwnerUserID string    `gorm:"type:uuid;not null" json:"owner_user_id"`
	StaffUserID string    `gorm:"type:uuid;not null" json:"staff_user_id"`
	RoleCode    string    `gorm:"type:varchar(30);default:'staff';not null" json:"role_code"`
	Status      string    `gorm:"type:varchar(20);default:'active';not null" json:"status"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`
}

type OwnerStaffResponse struct {
	ID          string    `json:"id"`
	OwnerUserID string    `json:"owner_user_id"`
	StaffUserID string    `json:"staff_user_id"`
	Username    string    `json:"username"`
	Fullname    string    `json:"fullname"`
	Phone       string    `json:"phone"`
	RoleCode    string    `json:"role_code"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
