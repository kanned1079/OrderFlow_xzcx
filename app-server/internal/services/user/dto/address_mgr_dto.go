package dto

type AddNewAddressRequestDto struct {
	UserId      int64  `json:"user_id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	FullAddress string `json:"full_address"`
}

type EditAddressRequestDto struct {
	Id          int64  `json:"id"`
	UserId      int64  `json:"user_id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	FullAddress string `json:"full_address"`
}
