package models

type ProfileResponse struct {
	User           UserSummary    `json:"user"`
	Stats          ProfileStats   `json:"stats"`
	Plan           PlanInfo       `json:"plan"`
	RevenueSummary RevenueSummary `json:"revenue_summary"`
}

type UserSummary struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatar_url"`
	Initials  string `json:"initials"`
}

type ProfileStats struct {
	FieldCount   int     `json:"field_count"`
	BookingCount int     `json:"booking_count"`
	TotalRevenue float64 `json:"total_revenue"`
}

type PlanInfo struct {
	Name       string `json:"name"`
	FieldUsage string `json:"field_usage"`
	CourtUsage string `json:"court_usage"`
	CanUpgrade bool   `json:"can_upgrade"`
}

type RevenueSummary struct {
	Total   float64 `json:"total"`
	Daily   float64 `json:"daily"`
	Monthly float64 `json:"monthly"`
	Yearly  float64 `json:"yearly"`
}
