package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/olivere/elastic/v7"
)

var (
	ErrNotFound = fmt.Errorf("product not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %v", err)
	}

	// Check if the client can connect
	_, err = client.ClusterHealth().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Elasticsearch: %v", err)
	}

	return &elasticRepository{client}, nil
}

func (r *elasticRepository) Close() {
	// No explicit close needed for Elasticsearch client
}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	if p.ID == "" {
		return fmt.Errorf("product ID is required")
	}
	if p.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if p.Price < 0 {
		return fmt.Errorf("product price cannot be negative")
	}

	_, err := r.client.Index().
		Index("catalog").
		Id(p.ID).
		BodyJson(productDocument{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to index product: %v", err)
	}
	return nil
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	if id == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	res, err := r.client.Get().
		Index("catalog").
		Id(id).
		Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get product: %v", err)
	}
	if !res.Found {
		return nil, ErrNotFound
	}

	p := productDocument{}
	if err = json.Unmarshal(res.Source, &p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product data: %v", err)
	}

	return &Product{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	// Validate pagination parameters
	if take > 100 {
		take = 100 // Enforce maximum limit
	}
	if take == 0 {
		take = 10 // Default limit
	}

	res, err := r.client.Search().
		Index("catalog").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %v", err)
	}

	return r.extractProducts(res)
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	if len(ids) == 0 {
		return []Product{}, nil
	}
	if len(ids) > 100 {
		return nil, fmt.Errorf("cannot request more than 100 products at once")
	}

	items := make([]*elastic.MultiGetItem, len(ids))
	for i, id := range ids {
		items[i] = elastic.NewMultiGetItem().
			Index("catalog").
			Id(id)
	}

	res, err := r.client.MultiGet().
		Add(items...).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple products: %v", err)
	}

	products := make([]Product, 0, len(res.Docs))
	for _, doc := range res.Docs {
		if !doc.Found {
			continue // Skip not found products
		}

		p := productDocument{}
		if err := json.Unmarshal(doc.Source, &p); err != nil {
			log.Printf("Warning: failed to unmarshal product %s: %v", doc.Id, err)
			continue
		}

		products = append(products, Product{
			ID:          doc.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	// Validate pagination parameters
	if take > 100 {
		take = 100 // Enforce maximum limit
	}
	if take == 0 {
		take = 10 // Default limit
	}

	searchQuery := elastic.NewMultiMatchQuery(query, "name", "description")
	res, err := r.client.Search().
		Index("catalog").
		Query(searchQuery).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %v", err)
	}

	return r.extractProducts(res)
}

// Helper function to extract products from search results
func (r *elasticRepository) extractProducts(res *elastic.SearchResult) ([]Product, error) {
	products := make([]Product, 0, len(res.Hits.Hits))
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err := json.Unmarshal(hit.Source, &p); err != nil {
			log.Printf("Warning: failed to unmarshal product %s: %v", hit.Id, err)
			continue
		}

		products = append(products, Product{
			ID:          hit.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}
