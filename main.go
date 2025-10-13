package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Bin struct {
	id        string
	isPrivate bool
	createdAt time.Time
	name      string
}

type BinList struct {
	bins []Bin
}

func generateID() string {
	id := uuid.New()
	return id.String()
}

func newBin(name string, isPrivate bool) (*Bin, error) {
	if name == "" {
		return &Bin{}, errors.New("name can't be empty")
	}
	id := generateID()
	createdAt := time.Now()
	result := Bin{
		id:        id,
		isPrivate: isPrivate,
		createdAt: createdAt,
		name:      name,
	}
	return &result, nil
}

func main() {
	bin, err := newBin("example", false)
	fmt.Println(bin, err)
}
