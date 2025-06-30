package dto

type FetchGoodsListAsCategoryRequestDto struct {
	MerchantId int64  `form:"merchant_id"`
	GoodsName  string `form:"goods_name"`
}
