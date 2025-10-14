package main

import (
	"fmt"

	"3-struct/bins"
)

func main() {
	bin, err := bins.NewBin("example", false)
	fmt.Println(bin, err)
}
