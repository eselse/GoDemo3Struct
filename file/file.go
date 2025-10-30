package file

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type DB interface {
	Read(fileName string) ([]byte, error)
	Write(content []byte, name string)
}

type FileDB struct{}

func NewFileDB() DB {
	return &FileDB{}
}

func (fd FileDB) Read(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(fileName, ".json") {
		return nil, errors.New("file isn't valid json file")
	}
	return data, nil
}

func (fd FileDB) Write(content []byte, name string) {
	file, err := os.Create(name)
	if err != nil {
		fmt.Println("Error creating file:", err)
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("\nFile written successfully")
}
