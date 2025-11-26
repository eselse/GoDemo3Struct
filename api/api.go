package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"3-struct/bins"
	"3-struct/config"
	"3-struct/file"
)

type binCreationResonse struct {
	Record   bins.BinList `json:"record"`
	Metadata `json:"metadata"`
}
type Metadata struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	Private   bool   `json:"private"`
}
type BinIDs struct {
	BinName   string `json:"binName"`
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	Private   bool   `json:"private"`
}

func InitAPI() *config.Config {
	newConfig := config.NewConfig("KEY")
	return newConfig
}

func CreateBin(fileName string, binName string, key string) (string, error) {
	// Tead file from fileName
	fileDB := file.NewFileDB()
	fileBody, err := fileDB.Read(fileName)
	if err != nil {
		return "", err
	}

	// Create an url with headers and body for request
	url := "https://api.jsonbin.io/v3/b"
	payload := strings.NewReader(string(fileBody))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Master-Key", key)
	req.Header.Add("X-Bin-Name", binName)

	// Get result
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	// Get body of result
	body, _ := io.ReadAll(res.Body)

	// Marshall it to struct
	var bodyJSON binCreationResonse
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		return "", err
	}

	// Get an ID for return
	createdBinID := bodyJSON.ID
	fmt.Printf("New bin was created, title is %s ID is %s", binName, createdBinID)
	return createdBinID, nil
}
