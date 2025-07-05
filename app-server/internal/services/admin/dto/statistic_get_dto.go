package dto

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

	UserOrdersOverview []int64 `json:"user_orders_overview"`
	MerchantsOverview  []int64 `json:"merchants_overview"`

	Message string `json:"message"`
}
