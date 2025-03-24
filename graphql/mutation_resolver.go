package main

import (
	"context"
	"errors"
	"log"
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

	acc, err := r.server.accountClient.PostAccount(ctx, in.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Account{
		ID:   acc.ID,
		Name: acc.Name,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	p, err := r.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.Price)
	if err != nil {
		log.Println(err)
		return nil, err
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

	products := make([]order.OrderedProduct, len(in.Products))
	for i, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		products[i] = order.OrderedProduct{
			ID:       p.ID,
			Quantity: uint32(p.Quantity),
		}
	}

	o, err := r.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil {
		log.Println(err)
		return nil, err
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
