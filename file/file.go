package file

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type DB interface {
	ReadJSON(fileName string) ([]byte, error)
	ReadPlain(fileName string) ([]byte, error)
	Write(content []byte, name string)
	Append(content []byte, name string) error
}

type FileDB struct{}

func NewFileDB() DB {
	return &FileDB{}
}

func (fd FileDB) ReadJSON(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(fileName, ".json") {
		return nil, errors.New("file isn't valid json file")
	}
	return data, nil
}

func (fd FileDB) ReadPlain(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
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

func (fd FileDB) Append(content []byte, name string) error {
	// Open the file in append mode, create it if it doesn't exist
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write data to the file
	if _, err := f.WriteString(string(content)); err != nil {
		return err
	}
	return nil
}
