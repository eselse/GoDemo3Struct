package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err // file not found is acceptable in many cases
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle inline comments: KEY=value # comment
		if strings.Contains(line, "#") {
			line = strings.Split(line, "#")[0]
			line = strings.TrimSpace(line)
		}

		// Split on first = only
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // invalid line
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		value = strings.Trim(value, `"'`)

		os.Setenv(key, value)
	}

	return scanner.Err()
}

func InitAPI() *config.Config {
	newConfig := config.NewConfig("KEY")
	return newConfig
}

func CreateBin(fileName string, binName string, key string) (string, error) {
	fmt.Printf("CreateBin called with file name %s and bin name is %s\n", fileName, binName)
	// Tead file from fileName
	fileDB := file.NewFileDB()
	fileBody, err := fileDB.Read(fileName)
	if err != nil {
		return "", err
	}
	// fmt.Println(fileBody)
	// Transform it to JSON
	var fileBodyJSON bins.BinList
	err = json.Unmarshal(fileBody, &fileBodyJSON)
	if err != nil {
		return "", err
	}
	// fmt.Println(fileBodyJSON)
	//
	// Create a body for request
	url := "https://api.jsonbin.io/v3/b"
	config := InitAPI()
	fmt.Println(config.Key)
	payload := strings.NewReader(string(fileBody))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Master-Key", config.Key)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))

	return "", nil
}
