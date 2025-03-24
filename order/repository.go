package order

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
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

func (r *postgresRepository) PutOrder(ctx context.Context, o Order) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Ensure transaction is handled properly
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("error rolling back transaction: %v, original error: %v", rbErr, err)
			}
		}
	}()

	// Insert order
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES ($1, $2, $3, $4)",
		o.ID,
		o.CreatedAt,
		o.AccountID,
		o.TotalPrice,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				return fmt.Errorf("account with ID %s does not exist", o.AccountID)
			case "unique_violation":
				return fmt.Errorf("order with ID %s already exists", o.ID)
			}
		}
		return fmt.Errorf("failed to insert order: %v", err)
	}

	// Prepare statement for order products
	stmt, err := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity", "price"))
	if err != nil {
		return fmt.Errorf("failed to prepare order products statement: %v", err)
	}
	defer stmt.Close()

	// Insert order products
	for _, p := range o.Products {
		_, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity, p.Price)
		if err != nil {
			return fmt.Errorf("failed to insert order product (ID: %s): %v", p.ID, err)
		}
	}

	// Execute the prepared statement
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "foreign_key_violation" {
			return fmt.Errorf("one or more products in your order could not be found")
		}
		return fmt.Errorf("failed to execute order products statement: %v", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	// Query orders and their products
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			o.id,
			o.created_at,
			o.account_id,
			o.total_price::money::numeric::float8,
			op.product_id,
			op.quantity,
			op.price
		FROM orders o 
		JOIN order_products op ON o.id = op.order_id
		WHERE o.account_id = $1
		ORDER BY o.id`,
		accountID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return []Order{}, nil
		}
		return nil, fmt.Errorf("failed to query orders: %v", err)
	}
	defer rows.Close()

	orderMap := make(map[string]*Order)
	var orders []Order

	// Scan rows into orders
	for rows.Next() {
		var order Order
		var product OrderedProduct
		if err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.AccountID,
			&order.TotalPrice,
			&product.ID,
			&product.Quantity,
			&product.Price,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order row: %v", err)
		}

		if existingOrder, ok := orderMap[order.ID]; !ok {
			order.Products = []OrderedProduct{product}
			orderMap[order.ID] = &order
		} else {
			existingOrder.Products = append(existingOrder.Products, product)
		}
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over order rows: %v", err)
	}

	// Convert map to slice
	for _, order := range orderMap {
		orders = append(orders, *order)
	}

	return orders, nil
}
