package dto

type FetchAllMerchantsRequestDto struct {
	SearchAs string `form:"search_as"`
	Search   string `form:"search"`
	Sort     string `form:"sort"`
	Page     int64  `form:"page"`
	Size     int64  `form:"size"`
}

type CreateNewMerchantRequestDto struct {
	UserId       int64  `json:"user_id"`
	MerchantName string `json:"merchant_name"`
	Description  string `json:"description"`
	LogoUrl      string `json:"logo_url"`
	Address      string `json:"address"`
}

type DeleteMerchantRequestDto struct {
	Id int64 `form:"id"`
}
