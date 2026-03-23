package models

import "time"

type Subscription struct {
	ID           string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID       string     `gorm:"type:uuid;not null" json:"user_id"`
	PlanID       string     `gorm:"type:uuid;not null" json:"plan_id"`
	Status       string     `gorm:"type:varchar(30);not null" json:"status"` // trial / active / expired / cancelled
	StartAt      time.Time  `gorm:"type:timestamp;not null" json:"start_at"`
	EndAt        time.Time  `gorm:"type:timestamp;not null" json:"end_at"`
	TrialStartAt *time.Time `gorm:"type:timestamp" json:"trial_start_at"`
	TrialEndAt   *time.Time `gorm:"type:timestamp" json:"trial_end_at"`
	ActivatedAt  *time.Time `gorm:"type:timestamp" json:"activated_at"`
	ExpiredAt    *time.Time `gorm:"type:timestamp" json:"expired_at"`
	CancelledAt  *time.Time `gorm:"type:timestamp" json:"cancelled_at"`
	CancelReason string     `gorm:"type:text" json:"cancel_reason"`
	AutoRenew    bool       `gorm:"type:boolean;not null;default:false" json:"auto_renew"`
	CreatedAt    time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`

	Plan Plan `gorm:"foreignKey:PlanID" json:"plan"`
}
