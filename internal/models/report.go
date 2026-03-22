package models

type RevenueReportResponse struct {
	TotalRevenue float64               `json:"total_revenue"`
	Period       string                `json:"period"`
	ByField      []RevenueByField      `json:"by_field"`
	BySportType  []RevenueBySportType `json:"by_sport_type"`
}

type RevenueByField struct {
	FieldID      string  `json:"field_id"`
	FieldName    string  `json:"field_name"`
	Revenue      float64 `json:"revenue"`
	BookingCount int     `json:"booking_count"`
}

type RevenueBySportType struct {
	SportType    string  `json:"sport_type"`
	Revenue      float64 `json:"revenue"`
	BookingCount int     `json:"booking_count"`
}
