package service

import (
	"context"

	db "github.com/hoangquan/retail-store-api/db/sqlc"
	"github.com/hoangquan/retail-store-api/pkg/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(queries *db.Queries) *ProductService {
	return &ProductService{
		repo: repository.NewProductRepository(queries),
	}
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (db.GetProductByIDRow, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context, limit, offset int32) ([]db.ListProductsRow, error) {
	return s.repo.List(ctx, db.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	})
}
