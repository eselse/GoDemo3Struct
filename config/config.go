package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Key string
}

func NewConfig(key string) *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
		return &Config{
			Key: "",
		}
	}
	newKey := os.Getenv(key)
	return &Config{
		Key: newKey,
	}
}
