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
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Printf("Error fetching account with ID %s: %v", *id, err)
			if strings.Contains(err.Error(), "not found") {
				return []*Account{}, nil
			}
			return nil, fmt.Errorf("failed to fetch account: %v", err)
		}
		return []*Account{{
			ID:   r.ID,
			Name: r.Name,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		log.Printf("Error fetching accounts: %v", err)
		if strings.Contains(err.Error(), "not found") {
			return []*Account{}, nil
		}
		return nil, fmt.Errorf("failed to fetch accounts: %v", err)
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
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Handle single ID lookup
	if id != nil {
		r, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			log.Printf("Error fetching product with ID %s: %v", *id, err)
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				return []*Product{}, nil
			}
			return nil, fmt.Errorf("failed to fetch product: %v", err)
		}
		return []*Product{{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
		}}, nil
	}

	// Handle multiple IDs lookup
	if ids != nil && len(ids) > 0 {
		var products []*Product
		for _, productID := range ids {
			r, err := r.server.catalogClient.GetProduct(ctx, productID)
			if err != nil {
				log.Printf("Error fetching product with ID %s: %v", productID, err)
				// Skip products that are not found
				if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
					continue
				}
				return nil, fmt.Errorf("failed to fetch products: %v", err)
			}
			products = append(products, &Product{
				ID:          r.ID,
				Name:        r.Name,
				Description: r.Description,
				Price:       r.Price,
			})
		}
		return products, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	queryStr := ""
	if query != nil {
		queryStr = *query
	}

	productList, err := r.server.catalogClient.GetProducts(ctx, skip, take, nil, queryStr)
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			return []*Product{}, nil
		}
		return nil, fmt.Errorf("failed to fetch products: %v", err)
	}

	var products []*Product
	for _, a := range productList {
		products = append(products, &Product{
			ID:          a.ID,
			Name:        a.Name,
			Description: a.Description,
			Price:       a.Price,
		})
	}

	return products, nil
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
