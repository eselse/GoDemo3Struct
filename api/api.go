package api

import (
	"3-struct/config"
)

func InitAPI() *config.Config {
	newConfig := config.NewConfig("KEY")
	return newConfig
}
