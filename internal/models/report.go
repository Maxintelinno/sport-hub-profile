package models

type RevenueReportResponse struct {
	TotalRevenue float64          `json:"total_revenue"`
	Period       string           `json:"period"`
	ByField      []RevenueByField `json:"by_field"`
}

type RevenueByField struct {
	FieldID      string  `json:"field_id"`
	FieldName    string  `json:"field_name"`
	Revenue      float64 `json:"revenue"`
	BookingCount int     `json:"booking_count"`
}
