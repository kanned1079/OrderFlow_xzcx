package dto

type GetCategoryListRequestDto struct {
	CategoryTitle string `form:"category_title" json:"category_title"`
	Page          int    `form:"page" json:"page"`
	Size          int    `form:"size" json:"size"`
}

type AddNewCategoryRequestDto struct {
	MerchantId    int64  `json:"merchant_id"` // 对应商户id
	CategoryTitle string `form:"category_title" json:"category_title"`
}

type EditCategoryRequestDto struct {
	MerchantId    int64  `json:"merchant_id"` // 对应商户id
	CategoryTitle string `form:"category_title" json:"category_title"`
}
