package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%w: context is required", ErrInvalidContext)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		if *id == "" {
			return nil, fmt.Errorf("%w: id cannot be empty when provided", ErrInvalidParameter)
		}

		account, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Printf("Error fetching account with ID %s: %v", *id, err)
			if strings.Contains(err.Error(), "not found") {
				return []*Account{}, nil
			}
			return nil, fmt.Errorf("failed to fetch account: %v", err)
		}

		if account == nil {
			return []*Account{}, nil
		}

		return []*Account{{
			ID:   account.ID,
			Name: account.Name,
		}}, nil
	}

	skip, take := uint64(0), uint64(100)
	if pagination != nil {
		if pagination.Skip != nil && *pagination.Skip < 0 {
			return nil, fmt.Errorf("%w: skip cannot be negative", ErrInvalidParameter)
		}
		if pagination.Take != nil && *pagination.Take < 0 {
			return nil, fmt.Errorf("%w: take cannot be negative", ErrInvalidParameter)
		}
		if pagination.Take != nil && *pagination.Take > 100 {
			return nil, fmt.Errorf("%w: take cannot exceed 100", ErrInvalidParameter)
		}

		if pagination.Skip != nil {
			skip = uint64(*pagination.Skip)
		}
		if pagination.Take != nil {
			take = uint64(*pagination.Take)
		}
	}

	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		log.Printf("Error fetching accounts: %v", err)
		if strings.Contains(err.Error(), "not found") {
			return []*Account{}, nil
		}
		return nil, fmt.Errorf("failed to fetch accounts: %v", err)
	}

	if accountList == nil {
		return []*Account{}, nil
	}

	accounts := make([]*Account, len(accountList))
	for i, acc := range accountList {
		accounts[i] = &Account{
			ID:   acc.ID,
			Name: acc.Name,
		}
	}
	return accounts, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string, ids []string) ([]*Product, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%w: context is required", ErrInvalidContext)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Handle single product lookup by ID
	if id != nil {
		if *id == "" {
			return nil, fmt.Errorf("%w: id cannot be empty when provided", ErrInvalidParameter)
		}

		product, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			log.Printf("Error fetching product with ID %s: %v", *id, err)
			if strings.Contains(err.Error(), "not found") {
				return []*Product{}, nil
			}
			return nil, fmt.Errorf("failed to fetch product: %v", err)
		}

		if product == nil {
			return []*Product{}, nil
		}

		return []*Product{{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		}}, nil
	}

	// Handle multiple products lookup by IDs
	if len(ids) > 0 {
		if len(ids) > 100 {
			return nil, fmt.Errorf("%w: cannot request more than 100 products at once", ErrInvalidParameter)
		}

		for i, productID := range ids {
			if productID == "" {
				return nil, fmt.Errorf("%w: product ID at index %d cannot be empty", ErrInvalidParameter, i)
			}
		}

		products, err := r.server.catalogClient.GetProducts(ctx, 0, 100, ids, "")
		if err != nil {
			log.Printf("Error fetching products by IDs: %v", err)
			if strings.Contains(err.Error(), "not found") {
				return []*Product{}, nil
			}
			return nil, fmt.Errorf("failed to fetch products by IDs: %v", err)
		}

		if len(products) == 0 {
			return []*Product{}, nil
		}

		result := make([]*Product, len(products))
		for i, p := range products {
			result[i] = &Product{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			}
		}
		return result, nil
	}

	// Handle product search or listing
	skip, take := uint64(0), uint64(100)
	if pagination != nil {
		if pagination.Skip != nil && *pagination.Skip < 0 {
			return nil, fmt.Errorf("%w: skip cannot be negative", ErrInvalidParameter)
		}
		if pagination.Take != nil && *pagination.Take < 0 {
			return nil, fmt.Errorf("%w: take cannot be negative", ErrInvalidParameter)
		}
		if pagination.Take != nil && *pagination.Take > 100 {
			return nil, fmt.Errorf("%w: take cannot exceed 100", ErrInvalidParameter)
		}

		if pagination.Skip != nil {
			skip = uint64(*pagination.Skip)
		}
		if pagination.Take != nil {
			take = uint64(*pagination.Take)
		}
	}

	queryStr := ""
	if query != nil {
		queryStr = *query
	}

	products, err := r.server.catalogClient.GetProducts(ctx, skip, take, nil, queryStr)
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		if strings.Contains(err.Error(), "not found") {
			return []*Product{}, nil
		}
		return nil, fmt.Errorf("failed to fetch products: %v", err)
	}

	if len(products) == 0 {
		return []*Product{}, nil
	}

	result := make([]*Product, len(products))
	for i, p := range products {
		result[i] = &Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}
	}
	return result, nil
}

func (p *PaginationInput) bounds() (uint64, uint64) {
	skip, take := uint64(0), uint64(0)

	if p.Skip != nil {
		skip = uint64(*p.Skip)
	}

	if p.Take != nil {
		take = uint64(*p.Take)
	}

	return skip, take
}
