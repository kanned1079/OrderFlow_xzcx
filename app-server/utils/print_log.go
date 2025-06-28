package utils

import (
	"fmt"
	"log"
)

const (
	RedColor    = "\033[31m"
	GreenColor  = "\033[32m"
	OrangeColor = "\033[33m"
	BlueColor   = "\033[34m"
	ResetColor  = "\033[0m"

	ErrorPrefix   = "[ERROR]"
	InfoPrefix    = "[INFO]"
	WarnPrefix    = "[WARN]"
	SuccessPrefix = "[SUCCESS]"
)

func (Logger) PrintError(v ...interface{}) {
	log.Printf("%s%s%s %s", RedColor, ErrorPrefix, ResetColor, fmt.Sprint(v...))
}

func (Logger) PrintInfo(v ...interface{}) {
	log.Printf("%s%s%s %s", BlueColor, InfoPrefix, ResetColor, fmt.Sprint(v...))
}

func (Logger) PrintWarn(v ...interface{}) {
	log.Printf("%s%s%s %s", OrangeColor, WarnPrefix, ResetColor, fmt.Sprint(v...))
}

func (Logger) PrintSuccess(v ...interface{}) {
	log.Printf("%s%s%s %s", GreenColor, SuccessPrefix, ResetColor, fmt.Sprint(v...))
}
