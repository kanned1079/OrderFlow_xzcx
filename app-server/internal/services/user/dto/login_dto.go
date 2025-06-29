package dto

//type CheckUserIsValidPreRegisterRequestDto struct {
//	PlatformId string `json:"platform_id"`
//	Name       string `json:"name"`
//}

type UserLoginRequestDto struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserRegisterRequestDto struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Username    string `json:"username"`
}

//type UserLoginResponseDto struct {
//	User         models.User `json:"user"`
//	AccessToken  string      `json:"access_token"`
//	RefreshToken string      `json:"refresh_token"`
//}
