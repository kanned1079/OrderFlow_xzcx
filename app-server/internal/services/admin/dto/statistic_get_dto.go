package dto

type OverviewItem struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

type FetchAdminDashboardStatisticResponseDto struct {
	RegisteredMerchants struct {
		Total     int64 `json:"total"`
		Yesterday int64 `json:"yesterday"`
		ThisMonth int64 `json:"this_month"`
	} `json:"registered_merchants"`

	RegisteredUsers struct {
		Total     int64 `json:"total"`
		Yesterday int64 `json:"yesterday"`
		ThisMonth int64 `json:"this_month"`
	} `json:"registered_users"`

	UserOrdersOverview []OverviewItem `json:"user_orders_overview"`
	MerchantsOverview  []OverviewItem `json:"merchants_overview"`

	Message string `json:"message"`
	P1      string `json:"p1"`
}
