package main

import (
	"context"
	"fmt"
	"log"

	"gcs-proxy/config"
	"gcs-proxy/internal/server"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func initGCSClient(cfg *config.Config) (*storage.Client, error) {
	ctx := context.Background()
	opts := []option.ClientOption{
		option.WithCredentialsJSON([]byte(fmt.Sprintf(`{
			"type": "service_account",
			"project_id": "%s",
			"client_email": "%s",
			"private_key": "%s"
		}`, cfg.GoogleProject, cfg.GoogleStorageClientEmail, cfg.GoogleStoragePrivateKey))),
	}

	log.Println("Initializing GCS")

	return storage.NewClient(ctx, opts...)
}

func main() {
	log.Println("Starting...")

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded successfully. Using port: %d", cfg.Port)

	gcsClient, err := initGCSClient(cfg)

	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}

	defer gcsClient.Close()
	log.Println("GCS successfully created")

	log.Println("Starting server...")
	srv := server.InitServer(gcsClient, cfg)
	log.Println("Server initialized")

	log.Printf("Starting server on port %d", cfg.Port)
	err = srv.Start(fmt.Sprintf(":%d", cfg.Port))

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
