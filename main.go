package main

import (
	"3-struct/bins"
	"3-struct/file"
)

func main() {
	db := file.NewFileDB()
	binList := bins.NewBins(db)
	print(binList, "\n")
}
