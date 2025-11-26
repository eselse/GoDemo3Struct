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

func Delete(id, key string) bool {
	url := "https://api.jsonbin.io/v3/b/" + id

	req, _ := http.NewRequest("DELETE", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Master-Key", key)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return false
	}

	return true
}

func Update(fileName, id, key string) bool {
	// Tead file from fileName
	fileDB := file.NewFileDB()
	fileBody, err := fileDB.ReadJSON(fileName)
	if err != nil {
		return false
	}

	// Create an url with headers and body for request
	url := "https://api.jsonbin.io/v3/b/" + id

	payload := strings.NewReader(string(fileBody))
	req, _ := http.NewRequest("PUT", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Master-Key", key)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return false
	}

	return true
}

func Get(id, key string) string {
	// Create an url with headers and body for request
	url := "https://api.jsonbin.io/v3/b/" + id
	fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Master-Key", key)

	// Get result
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	// Get body of result
	body, _ := io.ReadAll(res.Body)

	return string(body)
}

func List() {
	fileDB := file.NewFileDB()
	data, err := fileDB.ReadPlain("my-bin")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(data))
}

func CreateBin(fileName, binName, name, key string) error {
	// Tead file from fileName
	fileDB := file.NewFileDB()
	fileBody, err := fileDB.ReadJSON(fileName)
	if err != nil {
		return err
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
		return err
	}

	// Get an ID for return
	createdBinID := bodyJSON.ID
	fmt.Printf("New bin was created, title is %s ID is %s\n", binName, createdBinID)

	// Save ID to file
	resultString := createdBinID + "\n"
	err = fileDB.Append([]byte(resultString), name)
	return nil
}
