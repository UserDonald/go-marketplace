package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/donaldnash/go-marketplace/order"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL" required:"true"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL" required:"true"`
	Port        int    `envconfig:"PORT" default:"8083"`
}

func main() {
	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Order service...")

	// Load configuration
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded:\n  Database URL: %s\n  Account URL: %s\n  Catalog URL: %s\n  Port: %d",
		cfg.DatabaseURL, cfg.AccountURL, cfg.CatalogURL, cfg.Port)

	// Initialize repository with retry
	var repository order.Repository
	var err error

	log.Println("Connecting to database...")
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		repository, err = order.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return err
		}
		return nil
	})
	defer repository.Close()

	log.Println("Successfully connected to database")

	// Create service
	service := order.NewService(repository)

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
	if err := order.ListenGRPC(service, cfg.AccountURL, cfg.CatalogURL, cfg.Port); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	<-done
	log.Println("Service stopped")
}
