package repository

import (
	"context"

	db "github.com/kainguyen/retail-store-api/db/sqlc"
)

type ProductRepository struct {
	queries *db.Queries
}

func NewProductRepository(queries *db.Queries) *ProductRepository {
	return &ProductRepository{queries: queries}
}

func (r *ProductRepository) Create(ctx context.Context, arg db.CreateProductParams) (db.Product, error) {
	return r.queries.CreateProduct(ctx, arg)
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (db.GetProductByIDRow, error) {
	return r.queries.GetProductByID(ctx, id)
}

func (r *ProductRepository) List(ctx context.Context, arg db.ListProductsParams) ([]db.ListProductsRow, error) {
	return r.queries.ListProducts(ctx, arg)
}

func (r *ProductRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountProducts(ctx)
}

func (r *ProductRepository) Update(ctx context.Context, arg db.UpdateProductParams) (db.Product, error) {
	return r.queries.UpdateProduct(ctx, arg)
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteProduct(ctx, id)
}
