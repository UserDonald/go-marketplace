package account

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutAccount(ctx context.Context, a Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	if r.db != nil {
		r.db.Close()
	}
}

func (r *postgresRepository) PutAccount(ctx context.Context, a Account) error {
	if a.ID == "" {
		return fmt.Errorf("account ID is required")
	}
	if a.Name == "" {
		return fmt.Errorf("account name is required")
	}

	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts(id, name) VALUES($1, $2)", a.ID, a.Name)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return fmt.Errorf("account with ID %s already exists", a.ID)
			}
		}
		return fmt.Errorf("failed to insert account: %v", err)
	}
	return nil
}

func (r *postgresRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	if id == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	row := r.db.QueryRowContext(ctx, "SELECT id, name FROM accounts WHERE id = $1", id)
	a := &Account{}
	if err := row.Scan(&a.ID, &a.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to scan account row: %v", err)
	}
	return a, nil
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	// Validate pagination parameters
	if take > 100 {
		take = 100 // Enforce maximum limit
	}
	if take == 0 {
		take = 10 // Default limit
	}

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2",
		skip,
		take,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %v", err)
	}
	defer rows.Close()

	accounts := []Account{}
	for rows.Next() {
		a := Account{}
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, fmt.Errorf("failed to scan account row: %v", err)
		}
		accounts = append(accounts, a)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over account rows: %v", err)
	}

	return accounts, nil
}
