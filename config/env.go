package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                     int
	GoogleProject            string
	GoogleLocation           string
	GoogleBucketName         string
	GoogleStorageClientEmail string
	GoogleStoragePrivateKey  string
	CORSAllowedOrigins       []string
	CORSAllowedMethods       []string
	CORSAllowedHeaders       []string
	CORSAllowCredentials     bool
}

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %v", err)
	}

	corsAllowCredentials, err := strconv.ParseBool(getEnv("CORS_ALLOW_CREDENTIALS", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid CORS_ALLOW_CREDENTIALS: %v", err)
	}

	return &Config{
		Port:                     port,
		GoogleProject:            getEnv("GOOGLE_PROJECT", ""),
		GoogleLocation:           getEnv("GOOGLE_LOCATION", ""),
		GoogleBucketName:         getEnv("GOOGLE_BUCKET_NAME", ""),
		GoogleStorageClientEmail: getEnv("GOOGLE_STORAGE_CLIENT_EMAIL", ""),
		GoogleStoragePrivateKey:  getEnv("GOOGLE_STORAGE_PRIVATE_KEY", ""),
		CORSAllowedOrigins:       strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "*"), ","),
		CORSAllowedMethods:       strings.Split(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"), ","),
		CORSAllowedHeaders:       strings.Split(getEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization"), ","),
		CORSAllowCredentials:     corsAllowCredentials,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
