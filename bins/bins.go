package bins

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"3-struct/file"

	"github.com/google/uuid"
)

type Bin struct {
	ID        string    `json:"id"`
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
}

type BinList struct {
	Bins []Bin `json:"bins"`
}

func (bins *BinList) ToBytes() ([]byte, error) {
	file, err := json.Marshal(bins)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err.Error())
	}
	return file, err
}

func generateID() string {
	id := uuid.New()
	return id.String()
}

func NewBin(name string, isPrivate bool) (*Bin, error) {
	if name == "" {
		return &Bin{}, errors.New("name can't be empty")
	}
	id := generateID()
	createdAt := time.Now()
	result := Bin{
		ID:        id,
		IsPrivate: isPrivate,
		CreatedAt: createdAt,
		Name:      name,
	}
	return &result, nil
}

func NewBins(db file.DB) *BinList {
	file, err := db.Read("bins.json")
	if err != nil {
		return &BinList{
			Bins: []Bin{},
		}
	}
	var bins BinList
	err = json.Unmarshal(file, &bins)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &bins
}
