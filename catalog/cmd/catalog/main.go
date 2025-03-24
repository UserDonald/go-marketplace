package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/donaldnash/go-marketplace/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	ElasticsearchURL string `envconfig:"ELASTICSEARCH_URL" required:"true"`
	Port             int    `envconfig:"PORT" default:"8082"`
}

func main() {
	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Catalog service...")

	// Load configuration
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded:\n  Elasticsearch URL: %s\n  Port: %d",
		cfg.ElasticsearchURL, cfg.Port)

	// Initialize repository with retry
	var repository catalog.Repository
	var err error

	log.Println("Connecting to Elasticsearch...")
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		repository, err = catalog.NewElasticRepository(cfg.ElasticsearchURL)
		if err != nil {
			log.Printf("Failed to connect to Elasticsearch: %v", err)
			return err
		}
		return nil
	})
	defer repository.Close()

	log.Println("Successfully connected to Elasticsearch")

	// Create service
	service := catalog.NewService(repository)

	// Handle graceful shutdown
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Service is shutting down...")
		repository.Close()
		close(done)
	}()

	// Start gRPC server
	log.Printf("Starting gRPC server on port %d...", cfg.Port)
	if err := catalog.ListenGRPC(service, cfg.Port); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	<-done
	log.Println("Service stopped")
}
