package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/cors"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting GraphQL server with configuration:\nAccount URL: %s\nCatalog URL: %s\nOrder URL: %s",
		cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)

	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		log.Fatal("Failed to create GraphQL server:", err)
	}

	log.Println("GraphQL server created successfully")

	// Create a new CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Create the GraphQL handler
	h := handler.NewDefaultServer(s.ToExecutableSchema())

	// Apply CORS middleware
	http.Handle("/graphql", c.Handler(h))
	http.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting HTTP server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
