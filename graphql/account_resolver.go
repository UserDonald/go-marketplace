package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	if obj == nil {
		return nil, fmt.Errorf("account object is required")
	}
	if obj.ID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := r.server.orderClient.GetOrdersForAccount(ctx, obj.ID)
	if err != nil {
		log.Printf("Error fetching orders for account %s: %v", obj.ID, err)
		if strings.Contains(err.Error(), "not found") {
			return []*Order{}, nil
		}
		return nil, fmt.Errorf("failed to fetch orders for account: %v", err)
	}

	if orderList == nil {
		return []*Order{}, nil
	}

	orders := make([]*Order, len(orderList))
	for i, o := range orderList {
		products := make([]*OrderedProduct, len(o.Products))
		for j, p := range o.Products {
			products[j] = &OrderedProduct{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    int(p.Quantity),
			}
		}

		orders[i] = &Order{
			ID:         o.ID,
			CreatedAt:  o.CreatedAt,
			TotalPrice: o.TotalPrice,
			Products:   products,
		}
	}

	return orders, nil
}
