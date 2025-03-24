package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/cors"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL" required:"true"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL" required:"true"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL" required:"true"`
	Port       string `envconfig:"PORT" default:"8080"`
}

func main() {
	// Initialize logger with file and line numbers
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting GraphQL service...")

	// Load configuration
	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded:\n  Account URL: %s\n  Catalog URL: %s\n  Order URL: %s\n  Port: %s",
		cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL, cfg.Port)

	// Create GraphQL server
	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		log.Fatalf("Failed to create GraphQL server: %v", err)
	}
	defer func() {
		log.Println("Shutting down GraphQL server...")
		s.Close()
	}()

	// Create CORS middleware with secure defaults
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // In production, replace with specific origins
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browser
		Debug:            false,
	})

	// Create HTTP server with timeouts
	mux := http.NewServeMux()
	mux.Handle("/graphql", c.Handler(handler.NewDefaultServer(s.ToExecutableSchema())))
	mux.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Handle graceful shutdown
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		log.Printf("Received signal: %v\n", sig)
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
		close(done)
	}()

	// Start server
	log.Printf("Server is starting on port %s...", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	<-done
	log.Println("Server stopped gracefully")
}
