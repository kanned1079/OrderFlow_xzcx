package dto

type FetchMerchantsRequestDto struct {
	Search string `form:"search" json:"search"`
	Page   int    `form:"page" json:"page"`
	Size   int    `form:"size" json:"size"`
	SortAs string `form:"sort_as" json:"sort_as"`
	Sort   string `form:"sort" json:"sort"`
}
