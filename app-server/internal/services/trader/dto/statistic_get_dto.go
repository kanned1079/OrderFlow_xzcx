package dto

import "stay-server/internal/models"

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

//type FetchMerchantStatisticResponse struct {
//	Goods struct {
//		Total   int64 `json:"total"`
//		Deleted int64 `json:"deleted"`
//		OnSale  int64 `json:"on_sale"`
//	} `json:"goods"`
//	UserOrders struct {
//		Completed [7]int64 `json:"completed"`
//		Failed    [7]int64 `json:"failed"`
//	} `json:"user_orders"`
//	Income [7]float32 `json:"income"`
//
//	Message string `json:"message"`
//}

type OverviewItem struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

//type FetchMerchantStatisticResponse struct {
//	Goods struct {
//		Total   int64 `json:"total"`
//		Deleted int64 `json:"deleted"`
//		OnSale  int64 `json:"on_sale"`
//	} `json:"goods"`
//
//	UserOrders struct {
//		Completed []OverviewItem `json:"completed"`
//		Failed    []OverviewItem `json:"failed"`
//	} `json:"user_orders"`
//
//	Income []OverviewItem `json:"income"`
//
//	Message string `json:"message"`
//}
//
//type OverviewItem struct {
//	Date  string  `json:"date"`
//	Value float64 `json:"value"`
//}

type FetchMerchantStatisticResponse struct {
	GoodsCount struct { // 商品计数
		Total   int64 `json:"total"`   // 总商品
		Deleted int64 `json:"deleted"` // 删除了的
		OnSale  int64 `json:"on_sale"` // 总商品-删除了的
	} `json:"goods_count"`
	UserOrders struct {
		TodayCount       int64          `json:"today_count"`        // 今天订单数
		TodayProcessing  int64          `json:"today_processing"`   // 今天正在处理中的订单
		WeekFinishedRate float32        `json:"week_finished_rate"` // 最近一周的订单完成率
		CompletedGraph   []OverviewItem `json:"completed_graph"`    // 最近一周完成订单数量的图表数据
		FailedGraph      []OverviewItem `json:"failed_graph"`       // 最近一周失败图表数据
	} `json:"user_orders"`
	IncomeGraph     []OverviewItem `json:"income_graph"`      // 最近一周收入图表
	ActiveOrderList []models.Order `json:"active_order_list"` // 现在还没接单和正在处理的订单列表

	Message string `json:"message"`
}
