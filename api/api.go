package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"3-struct/bins"
	"3-struct/config"
	"3-struct/file"
)

type binCreationResponse struct {
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

	defer func() {
		_, _ = io.Copy(io.Discard, res.Body) // Drain body for connection reuse
		if err := res.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	// Get body of result
	body, _ := io.ReadAll(res.Body)

	// Marshall it to struct
	var bodyJSON binCreationResponse
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
	if err != nil {
		return err
	}
	return nil
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

	defer func() {
		_, _ = io.Copy(io.Discard, res.Body) // Drain body for connection reuse
		if err := res.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	defer func() {
		_, _ = io.Copy(io.Discard, res.Body) // Drain body for connection reuse
		if err := res.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	return res.StatusCode == 200
}

func Delete(id, key string) bool {
	url := "https://api.jsonbin.io/v3/b/" + id

	req, _ := http.NewRequest("DELETE", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Master-Key", key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	defer func() {
		_, _ = io.Copy(io.Discard, res.Body) // Drain body for connection reuse
		if err := res.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	return res.StatusCode == 200
}
