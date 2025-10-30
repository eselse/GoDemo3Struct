package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"3-struct/bins"
	"3-struct/file"
)

func SaveToFile(data []byte, name string) {
	file, err := os.Create(name)
	if err != nil {
		fmt.Println("Error creating file:", err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("File written successfully")
}

func GetFromFile(db file.DB) *bins.BinList {
	file, err := db.Read("bins.json")
	if err != nil {
		return &bins.BinList{
			Bins: []bins.Bin{},
		}
	}
	var bins bins.BinList
	err = json.Unmarshal(file, &bins)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &bins
}
