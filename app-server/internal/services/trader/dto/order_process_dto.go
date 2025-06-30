package dto

type CancelOrderByTraderRequestDto struct {
	UserId       int64  `json:"user_id"`
	OrderId      string `json:"order_id"`
	MerchantId   int64  `json:"merchant_id"`
	CancelReason string `json:"cancel_reason"`
}

type AcceptOrderByTraderRequestDto struct {
	OrderId    string `json:"order_id"`
	MerchantId int64  `json:"merchant_id"`
}

type CompleteOrderByTraderRequestDto struct {
	OrderId    string `json:"order_id"`
	MerchantId int64  `json:"merchant_id"`
}

type GetOrderListRequestDto struct {
	MerchantId int64  `form:"merchant_id" json:"merchant_id"`
	IsActive   bool   `form:"is_active" json:"is_active"`
	Page       int    `form:"page" json:"page"`
	Size       int    `form:"size" json:"size"`
	Sort       string `form:"sort" json:"sort"`
}

type GetOrderByIdRequestDto struct {
	MerchantId int64  `form:"merchant_id" json:"merchant_id"`
	OrderId    string `form:"order_id" json:"order_id"`
}
