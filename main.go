package main

import (
	"fmt"

	"3-struct/bins"
)

func main() {
	bin, err := bins.NewBin("example", false)
	if err != nil {
		fmt.Println(err.Error())
	}
	print(bin.Name, "\n")
}
