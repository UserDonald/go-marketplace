package catalog

import (
	"context"
	"fmt"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductByID(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type catalogService struct {
	repository Repository
}

func NewService(r Repository) Service {
	if r == nil {
		panic("repository cannot be nil")
	}
	return &catalogService{r}
}

func (s *catalogService) PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	if name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if description == "" {
		return nil, fmt.Errorf("product description is required")
	}
	if price < 0 {
		return nil, fmt.Errorf("product price cannot be negative")
	}

	p := &Product{
		ID:          ksuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
	}
	if err := s.repository.PutProduct(ctx, *p); err != nil {
		return nil, fmt.Errorf("failed to create product: %v", err)
	}
	return p, nil
}

func (s *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	if id == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	product, err := s.repository.GetProductByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get product: %v", err)
	}
	return product, nil
}

func (s *catalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	// Enforce pagination limits
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	products, err := s.repository.ListProducts(ctx, skip, take)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %v", err)
	}
	return products, nil
}

func (s *catalogService) GetProductByID(ctx context.Context, ids []string) ([]Product, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	if len(ids) == 0 {
		return []Product{}, nil
	}
	if len(ids) > 100 {
		return nil, fmt.Errorf("cannot request more than 100 products at once")
	}

	products, err := s.repository.ListProductsWithIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by IDs: %v", err)
	}
	return products, nil
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	// Enforce pagination limits
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	products, err := s.repository.SearchProducts(ctx, query, skip, take)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %v", err)
	}
	return products, nil
}
