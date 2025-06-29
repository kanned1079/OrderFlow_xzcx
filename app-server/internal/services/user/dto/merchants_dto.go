package dto

type FetchMerchantsRequestDto struct {
	Search string `form:"search" json:"search"`
}
