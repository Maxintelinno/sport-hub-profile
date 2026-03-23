package models

type DashboardResponse struct {
	Owner            DashboardOwner     `json:"owner"`
	Plan             DashboardPlan      `json:"plan"`
	Summary          DashboardSummary   `json:"summary"`
	BookingCount     int64              `json:"booking_count"`
	FieldCount       int64              `json:"field_count"`
	RevenueGrowthPct float64            `json:"revenue_growth_pct"`
	TotalRevenue     float64            `json:"total_revenue"`
	Alerts           []DashboardAlert   `json:"alerts"`
	NextActions      []DashboardAction  `json:"next_actions"`
	RevenueTrend7d   []RevenueTrendItem `json:"revenue_trend_7d"`
	Upsell           *DashboardUpsell   `json:"upsell,omitempty"`
}

type DashboardOwner struct {
	ID            string `json:"id"`
	Fullname      string `json:"fullname"`
	Phone         string `json:"phone"`
	AvatarInitial string `json:"avatar_initial"`
}

type DashboardPlan struct {
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Status          string  `json:"status"`
	TrialDaysLeft   int     `json:"trial_days_left"`
	IsTrial         bool    `json:"is_trial"`
	PriceAfterTrial float64 `json:"price_after_trial"`
}

type DashboardSummary struct {
	TotalRevenue     float64 `json:"total_revenue"`
	RevenueGrowthPct float64 `json:"revenue_growth_pct"`
	BookingCount     int64   `json:"booking_count"`
	FieldCount       int64   `json:"field_count"`
}

type DashboardAlert struct {
	Type       string `json:"type"` // warning, info, success, error
	Title      string `json:"title"`
	Message    string `json:"message"`
	ActionText string `json:"action_text"`
	ActionType string `json:"action_type"`
}

type DashboardAction struct {
	Title      string `json:"title"`
	Message    string `json:"message"`
	ActionText string `json:"action_text"`
	ActionType string `json:"action_type"`
}

type RevenueTrendItem struct {
	Label  string  `json:"label"`
	Amount float64 `json:"amount"`
}

type DashboardUpsell struct {
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	ActionText string `json:"action_text"`
	ActionType string `json:"action_type"`
}
