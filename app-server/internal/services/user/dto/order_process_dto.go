package dto

type OrderGoodsItemInput struct {
	GoodsId int64 `json:"goods_id"`
	Count   int64 `json:"count"`
}

type UserCommitOrderRequestDto struct {
	UserId     int64                 `json:"user_id"`
	MerchantId int64                 `json:"merchant_id"`
	AddressId  int64                 `json:"address_id"`
	GoodsList  []OrderGoodsItemInput `json:"goods_list"`
}

// PATCH
type CancelOrderByUserRequestDto struct {
	UserId     int64  `json:"user_id"`
	OrderId    string `json:"order_id"`
	MerchantId int64  `json:"merchant_id"`
}
