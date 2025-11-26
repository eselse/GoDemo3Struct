package config

import (
	"fmt"
	"os"

	"3-struct/env"
)

type Config struct {
	Key string
}

func NewConfig(key string) *Config {
	err := env.LoadEnv(".env")
	if err != nil {
		fmt.Println(err.Error())
	}
	newKey := os.Getenv(key)
	return &Config{
		Key: newKey,
	}
}
