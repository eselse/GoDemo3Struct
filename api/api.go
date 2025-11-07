package api

import (
	"fmt"

	"3-struct/config"
)

func InitAPI() {
	newConfig := config.NewConfig("KEY")
	fmt.Println(newConfig.Key)
}
