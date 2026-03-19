package repository

import (
	"context"

	db "github.com/kainguyen/retail-store-api/db/sqlc"
)

type CategoryRepository struct {
	queries *db.Queries
}

func NewCategoryRepository(queries *db.Queries) *CategoryRepository {
	return &CategoryRepository{queries: queries}
}

func (r *CategoryRepository) Create(ctx context.Context, name string) (db.Category, error) {
	return r.queries.CreateCategory(ctx, name)
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int64) (db.Category, error) {
	return r.queries.GetCategoryByID(ctx, id)
}

func (r *CategoryRepository) List(ctx context.Context) ([]db.Category, error) {
	return r.queries.ListCategories(ctx)
}

func (r *CategoryRepository) Update(ctx context.Context, arg db.UpdateCategoryParams) (db.Category, error) {
	return r.queries.UpdateCategory(ctx, arg)
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteCategory(ctx, id)
}
