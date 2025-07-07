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

	UserOrdersOverview [7]int64 `json:"user_orders_overview"`
	MerchantsOverview  [7]int64 `json:"merchants_overview"`

	Message string `json:"message"`
}

type FetchMerchantStatisticRequest struct {
}

type FetchMerchantStatisticResponse struct {
	Goods struct {
		Total   int64 `json:"total"`
		Deleted int64 `json:"deleted"`
		OnSale  int64 `json:"on_sale"`
	} `json:"goods"`
	UserOrders struct {
		Completed [7]int64 `json:"completed"`
		Failed    [7]int64 `json:"failed"`
	} `json:"user_orders"`
	Income [7]float32 `json:"income"`

	Message string `json:"message"`
}
