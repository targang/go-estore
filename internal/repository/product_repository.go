package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go_store/internal/model"
)

var _ ProductRepository = (*productRepositoryImpl)(nil)

type productRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (p *productRepositoryImpl) Create(ctx context.Context, product *model.Product) (string, error) {
	const query = `
INSERT INTO product (name, description, price)
VALUES ($1, $2, $3)
RETURNING id
`
	var result string
	err := p.db.QueryRow(ctx, query, product.Name, product.Description, product.Price).
		Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (p *productRepositoryImpl) GetByID(ctx context.Context, id string) (*model.Product, error) {
	const query = `
SELECT id, name, description, price FROM product WHERE id = $1
`
	var product model.Product
	err := p.db.QueryRow(ctx, query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.Price,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *productRepositoryImpl) Delete(ctx context.Context, id string) error {
	const query = `
DELETE FROM product WHERE id = $1
`
	_, err := p.db.Exec(ctx, query, id)
	return err
}

func (p *productRepositoryImpl) List(ctx context.Context, limit, offset int32) ([]model.Product, error) {
	const query = `
SELECT id, name, description, price 
FROM product 
ORDER BY name 
LIMIT $1 OFFSET $2
`
	rows, err := p.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}
