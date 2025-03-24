package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/donaldnash/go-marketplace/order"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, in AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if in.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	acc, err := r.server.accountClient.PostAccount(ctx, in.Name)
	if err != nil {
		log.Printf("Error creating account: %v", err)
		if strings.Contains(err.Error(), "already exists") {
			return nil, fmt.Errorf("an account with this name already exists")
		}
		return nil, fmt.Errorf("failed to create account: %v", err)
	}

	return &Account{
		ID:   acc.ID,
		Name: acc.Name,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if in.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if in.Description == "" {
		return nil, fmt.Errorf("description is required")
	}
	if in.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than 0")
	}

	p, err := r.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.Price)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		if strings.Contains(err.Error(), "already exists") {
			return nil, fmt.Errorf("a product with this name already exists")
		}
		return nil, fmt.Errorf("failed to create product: %v", err)
	}

	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if in.AccountID == "" {
		return nil, fmt.Errorf("accountId is required")
	}

	if len(in.Products) == 0 {
		return nil, fmt.Errorf("order must contain at least one product")
	}

	products := make([]order.OrderedProduct, len(in.Products))
	for i, p := range in.Products {
		if p.ID == "" {
			return nil, fmt.Errorf("product ID is required")
		}
		if p.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %s: quantity must be greater than 0", p.ID)
		}
		products[i] = order.OrderedProduct{
			ID:       p.ID,
			Quantity: uint32(p.Quantity),
		}
	}

	o, err := r.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		switch {
		case strings.Contains(err.Error(), "account not found"):
			return nil, fmt.Errorf("account with ID %s does not exist", in.AccountID)
		case strings.Contains(err.Error(), "product not found"):
			return nil, fmt.Errorf("one or more products in your order could not be found")
		case strings.Contains(err.Error(), "insufficient stock"):
			return nil, fmt.Errorf("one or more products in your order are out of stock")
		case strings.Contains(err.Error(), "could not post order"):
			return nil, fmt.Errorf("failed to create order, please try again")
		default:
			return nil, fmt.Errorf("an unexpected error occurred while creating your order")
		}
	}

	op := make([]*OrderedProduct, len(o.Products))
	for i, p := range o.Products {
		op[i] = &OrderedProduct{
			ID:       p.ID,
			Quantity: int(p.Quantity),
		}
	}

	return &Order{
		ID:         o.ID,
		CreatedAt:  o.CreatedAt,
		TotalPrice: o.TotalPrice,
		Products:   op,
	}, nil
}
