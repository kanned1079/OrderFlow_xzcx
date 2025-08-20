package dto

import "stay-server/internal/models"

type GetOrderInfoByAddressIdRequestSto struct {
}

type GetOrderInfoByAddressIdResponseDto struct {
	Address models.Address `json:"address"`
	Message string         `json:"message"`
}
