package bins

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type bin struct {
	id        string
	isPrivate bool
	createdAt time.Time
	name      string
}

type binList struct {
	bins []bin
}

func generateID() string {
	id := uuid.New()
	return id.String()
}

func NewBin(name string, isPrivate bool) (*bin, error) {
	if name == "" {
		return &bin{}, errors.New("name can't be empty")
	}
	id := generateID()
	createdAt := time.Now()
	result := bin{
		id:        id,
		isPrivate: isPrivate,
		createdAt: createdAt,
		name:      name,
	}
	return &result, nil
}
