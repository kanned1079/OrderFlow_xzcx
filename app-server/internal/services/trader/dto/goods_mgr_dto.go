package dto

type GetGoodsListRequestDto struct {
	MerchantId int64  `form:"merchant_id"`
	GoodsName  string `form:"goods_name"`
	Sort       string `form:"sort"`
	Page       int    `form:"page"`
	Size       int    `form:"size"`
}

type AddNewGoodsRequestDto struct {
	MerchantId  int64   `json:"merchant_id"` // 对应商户id
	CategoryId  int64   `json:"categoryId"`  // 分类id
	GoodsName   string  `json:"goods_name"`  // 商品名
	Description string  `json:"description"` // 商品描述
	LogoUrl     string  `json:"logo_url"`    // 商品Logo图片
	Price       float32 `json:"price"`       // 价格
	Residue     int64   `json:"residue"`     // 剩余库存
}

type EditGoodsRequestDto struct {
	Id          int64   `json:"id"`          // 商品id
	MerchantId  int64   `json:"merchant_id"` // 对应商户id
	CategoryId  int64   `json:"categoryId"`  // 分类id
	GoodsName   string  `json:"goods_name"`  // 商品名
	Description string  `json:"description"` // 商品描述
	LogoUrl     string  `json:"logo_url"`    // 商品Logo图片
	Price       float32 `json:"price"`       // 价格
	Residue     int64   `json:"residue"`     // 剩余库存
}
