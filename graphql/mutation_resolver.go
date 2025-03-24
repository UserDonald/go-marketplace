package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/donaldnash/go-marketplace/catalog"
	"github.com/donaldnash/go-marketplace/order"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrInvalidContext   = errors.New("invalid context")
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, in AccountInput) (*Account, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%w: context is required", ErrInvalidContext)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if in.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidParameter)
	}

	acc, err := r.server.accountClient.PostAccount(ctx, in.Name)
	if err != nil {
		log.Printf("Error creating account: %v", err)
		if strings.Contains(err.Error(), "already exists") {
			return nil, fmt.Errorf("%w: an account with this name already exists", ErrAlreadyExists)
		}
		return nil, fmt.Errorf("failed to create account: %v", err)
	}

	if acc == nil {
		return nil, fmt.Errorf("unexpected error: account creation succeeded but returned nil")
	}

	return &Account{
		ID:   acc.ID,
		Name: acc.Name,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%w: context is required", ErrInvalidContext)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if in.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidParameter)
	}
	if in.Description == "" {
		return nil, fmt.Errorf("%w: description is required", ErrInvalidParameter)
	}
	if in.Price < 0 {
		return nil, fmt.Errorf("%w: price cannot be negative", ErrInvalidParameter)
	}

	p, err := r.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.Price)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		if strings.Contains(err.Error(), "already exists") {
			return nil, fmt.Errorf("%w: a product with this name already exists", ErrAlreadyExists)
		}
		return nil, fmt.Errorf("failed to create product: %v", err)
	}

	if p == nil {
		return nil, fmt.Errorf("unexpected error: product creation succeeded but returned nil")
	}

	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%w: context is required", ErrInvalidContext)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if in.AccountID == "" {
		return nil, fmt.Errorf("%w: accountId is required", ErrInvalidParameter)
	}

	if len(in.Products) == 0 {
		return nil, fmt.Errorf("%w: order must contain at least one product", ErrInvalidParameter)
	}

	if len(in.Products) > 100 {
		return nil, fmt.Errorf("%w: order cannot contain more than 100 products", ErrInvalidParameter)
	}

	// Filter out products with quantity <= 0
	var validProducts []order.OrderedProduct
	for i, p := range in.Products {
		if p == nil {
			return nil, fmt.Errorf("%w: product at index %d is nil", ErrInvalidParameter, i)
		}
		if p.ID == "" {
			return nil, fmt.Errorf("%w: product ID is required at index %d", ErrInvalidParameter, i)
		}
		if p.Quantity > 0 {
			validProducts = append(validProducts, order.OrderedProduct{
				ID:       p.ID,
				Quantity: uint32(p.Quantity),
			})
		}
	}

	// Check if there are any valid products after filtering
	if len(validProducts) == 0 {
		return nil, fmt.Errorf("%w: order must contain at least one product with quantity greater than 0", ErrInvalidParameter)
	}

	// Get product details from catalog service
	productIDs := make([]string, len(validProducts))
	for i, p := range validProducts {
		productIDs[i] = p.ID
	}

	catalogProducts, err := r.server.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Printf("Error getting product details: %v", err)
		return nil, fmt.Errorf("failed to get product details: %v", err)
	}

	// Update product details
	productMap := make(map[string]*catalog.Product)
	for _, p := range catalogProducts {
		productMap[p.ID] = &p
	}

	for i := range validProducts {
		if p, ok := productMap[validProducts[i].ID]; ok {
			validProducts[i].Name = p.Name
			validProducts[i].Description = p.Description
			validProducts[i].Price = p.Price
		}
	}

	o, err := r.server.orderClient.PostOrder(ctx, in.AccountID, validProducts)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		switch {
		case strings.Contains(err.Error(), "account with ID"):
			return nil, fmt.Errorf("%w: account with ID %s does not exist", ErrNotFound, in.AccountID)
		case strings.Contains(err.Error(), "one or more products in your order could not be found"):
			return nil, fmt.Errorf("%w: one or more products in your order could not be found", ErrNotFound)
		case strings.Contains(err.Error(), "failed to create order"):
			return nil, fmt.Errorf("failed to create order: %v", err)
		default:
			return nil, fmt.Errorf("an unexpected error occurred while creating your order: %v", err)
		}
	}

	if o == nil {
		return nil, fmt.Errorf("unexpected error: order creation succeeded but returned nil")
	}

	op := make([]*OrderedProduct, len(o.Products))
	for i, p := range o.Products {
		op[i] = &OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    int(p.Quantity),
		}
	}

	return &Order{
		ID:         o.ID,
		CreatedAt:  o.CreatedAt,
		TotalPrice: o.TotalPrice,
		Products:   op,
	}, nil
}
