package main

import (
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
	log.Printf("Connecting to account service at %s...", accountUrl)
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}
	log.Println("Successfully connected to account service")

	log.Printf("Connecting to catalog service at %s...", catalogUrl)
	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		accountClient.Close()
		return nil, err
	}
	log.Println("Successfully connected to catalog service")

	log.Printf("Connecting to order service at %s...", orderUrl)
	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, err
	}
	log.Println("Successfully connected to order service")

	return &Server{
		accountClient,
		catalogClient,
		orderClient,
	}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
