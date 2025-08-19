package dto

type FetchAllMerchantsRequestDto struct {
	SearchAs    string `form:"search_as"`
	ShowDeleted bool   `form:"show_deleted"`
	Search      string `form:"search"`
	Sort        string `form:"sort"`
	Page        int64  `form:"page"`
	Size        int64  `form:"size"`
}

type CreateNewMerchantRequestDto struct {
	//UserId       int64  `json:"user_id"`
	PhoneNumber  string `json:"phone_number"`
	MerchantName string `json:"merchant_name"`
	Description  string `json:"description"`
	LogoUrl      string `json:"logo_url"`
	Address      string `json:"address"`
}

type UpdateMerchantByIdRequestDto struct {
	MerchantName string `json:"merchant_name"`
	Description  string `json:"description"`
	LogoUrl      string `json:"logo_url"`
	Address      string `json:"address"`
}

type DeleteMerchantRequestDto struct {
	Id int64 `form:"id"`
}
