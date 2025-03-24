package account

import (
	"context"
	"fmt"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	repository Repository
}

func NewService(r Repository) Service {
	if r == nil {
		panic("repository cannot be nil")
	}
	return &accountService{r}
}

func (s *accountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	if name == "" {
		return nil, fmt.Errorf("account name is required")
	}

	a := &Account{
		Name: name,
		ID:   ksuid.New().String(),
	}
	if err := s.repository.PutAccount(ctx, *a); err != nil {
		return nil, fmt.Errorf("failed to create account: %v", err)
	}
	return a, nil
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	if id == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	account, err := s.repository.GetAccountByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %v", err)
	}
	return account, nil
}

func (s *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	// Enforce pagination limits
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	accounts, err := s.repository.ListAccounts(ctx, skip, take)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %v", err)
	}
	return accounts, nil
}
