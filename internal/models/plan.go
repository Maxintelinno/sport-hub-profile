package models

import "time"

type Plan struct {
	ID                 string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Code               string    `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name               string    `gorm:"type:varchar(100);not null" json:"name"`
	Description        string    `gorm:"type:text" json:"description"`
	Price              float64   `gorm:"type:numeric(10,2);not null;default:0" json:"price"`
	BillingCycle       string    `gorm:"type:varchar(20);not null;default:'monthly'" json:"billing_cycle"`
	TrialDays          int       `gorm:"type:int;not null;default:0" json:"trial_days"`
	MaxFields          *int      `gorm:"type:int" json:"max_fields"`
	MaxCourts          *int      `gorm:"type:int" json:"max_courts"`
	HasReports         bool      `gorm:"type:boolean;not null;default:false" json:"has_reports"`
	HasPromotion       bool      `gorm:"type:boolean;not null;default:false" json:"has_promotion"`
	HasHomepageFeature bool      `gorm:"type:boolean;not null;default:false" json:"has_homepage_feature"`
	HasPrioritySupport bool      `gorm:"type:boolean;not null;default:false" json:"has_priority_support"`
	IsActive           bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt          time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt          time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`
}
