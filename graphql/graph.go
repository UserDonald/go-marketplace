package main

import (
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/donaldnash/go-marketplace/account"
	"github.com/donaldnash/go-marketplace/catalog"
	"github.com/donaldnash/go-marketplace/order"
)

type Server struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}

func NewGraphQLServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {
	if accountUrl == "" {
		return nil, fmt.Errorf("%w: account service URL is required", ErrInvalidParameter)
	}
	if catalogUrl == "" {
		return nil, fmt.Errorf("%w: catalog service URL is required", ErrInvalidParameter)
	}
	if orderUrl == "" {
		return nil, fmt.Errorf("%w: order service URL is required", ErrInvalidParameter)
	}

	log.Printf("Connecting to account service at %s...", accountUrl)
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to account service: %v", err)
	}
	log.Println("Successfully connected to account service")

	log.Printf("Connecting to catalog service at %s...", catalogUrl)
	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		accountClient.Close()
		return nil, fmt.Errorf("failed to connect to catalog service: %v", err)
	}
	log.Println("Successfully connected to catalog service")

	log.Printf("Connecting to order service at %s...", orderUrl)
	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, fmt.Errorf("failed to connect to order service: %v", err)
	}
	log.Println("Successfully connected to order service")

	return &Server{
		accountClient: accountClient,
		catalogClient: catalogClient,
		orderClient:   orderClient,
	}, nil
}

func (s *Server) Close() {
	if s.accountClient != nil {
		s.accountClient.Close()
	}
	if s.catalogClient != nil {
		s.catalogClient.Close()
	}
	if s.orderClient != nil {
		s.orderClient.Close()
	}
}

func (s *Server) Mutation() MutationResolver {
	if s == nil {
		panic("server cannot be nil")
	}
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	if s == nil {
		panic("server cannot be nil")
	}
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	if s == nil {
		panic("server cannot be nil")
	}
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	if s == nil {
		panic("server cannot be nil")
	}
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
