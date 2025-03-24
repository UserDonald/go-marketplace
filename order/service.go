package order

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type Order struct {
	ID         string
	CreatedAt  time.Time
	AccountID  string
	TotalPrice float64
	Products   []OrderedProduct
}

type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}

type orderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	if r == nil {
		panic("repository cannot be nil")
	}
	return &orderService{r}
}

func (s *orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}
	if len(products) == 0 {
		return nil, fmt.Errorf("at least one product is required")
	}
	if len(products) > 100 {
		return nil, fmt.Errorf("cannot order more than 100 products at once")
	}

	// Validate products
	for i, p := range products {
		if p.ID == "" {
			return nil, fmt.Errorf("product ID is required for product at index %d", i)
		}
		if p.Quantity == 0 {
			return nil, fmt.Errorf("quantity must be greater than 0 for product %s", p.ID)
		}
		if p.Price < 0 {
			return nil, fmt.Errorf("price cannot be negative for product %s", p.ID)
		}
	}

	o := &Order{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountID: accountID,
		Products:  products,
	}

	// Calculate total price
	o.TotalPrice = 0.0
	for _, p := range products {
		o.TotalPrice += p.Price * float64(p.Quantity)
	}

	if err := s.repository.PutOrder(ctx, *o); err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}
	return o, nil
}

func (s *orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	orders, err := s.repository.GetOrdersForAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders for account: %v", err)
	}
	return orders, nil
}
