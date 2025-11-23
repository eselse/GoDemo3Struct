package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"3-struct/bins"
	"3-struct/file"
)

func SaveToFile(data []byte, name string) (bool, error) {
	file, err := os.Create(name)
	if err != nil {
		return false, err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ReadFile(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func GetBinsFromFile(db file.DB) *bins.BinList {
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
