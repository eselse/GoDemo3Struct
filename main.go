package main

import (
	"fmt"

	"3-struct/config"
)

func main() {
	// db := file.NewFileDB()
	// binList := bins.NewBins(db)
	someConfig := config.NewConfig("KEY")
	someKey := someConfig.Key
	fmt.Println(someKey)
}
