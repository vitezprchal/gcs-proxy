package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port                     int
	GoogleProject            string
	GoogleLocation           string
	GoogleBucketName         string
	GoogleStorageClientEmail string
	GoogleStoragePrivateKey  string
}

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %v", err)
	}

	return &Config{
		Port:                     port,
		GoogleProject:            getEnv("GOOGLE_PROJECT", ""),
		GoogleLocation:           getEnv("GOOGLE_LOCATION", ""),
		GoogleBucketName:         getEnv("GOOGLE_BUCKET_NAME", ""),
		GoogleStorageClientEmail: getEnv("GOOGLE_STORAGE_CLIENT_EMAIL", ""),
		GoogleStoragePrivateKey:  getEnv("GOOGLE_STORAGE_PRIVATE_KEY", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
