package utils

import "stay-server/internal/config"

type Logger struct {
}

type Utils struct {
	Logger               Logger
	jwtSecret            []byte
	accessTokenExpiredIn int
}

func NewUtils() *Utils {
	return &Utils{
		Logger:               Logger{},
		jwtSecret:            []byte(config.AppCfg.Runtime.JwtSecret),
		accessTokenExpiredIn: config.AppCfg.Runtime.AccessTokenExpiredIn,
	}
}
