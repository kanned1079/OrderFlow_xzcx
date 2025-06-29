package dto

type FetchAllTradersRequestDto struct {
	PhoneNumber string `form:"phone_number"`
	Page        int    `form:"page"`
	Size        int    `form:"size"`
}
