package service

import (
	"context"

	db "github.com/hoangquan/retail-store-api/db/sqlc"
	"github.com/hoangquan/retail-store-api/pkg/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(queries *db.Queries) *CategoryService {
	return &CategoryService{repo: repository.NewCategoryRepository(queries)}
}

func (s *CategoryService) GetByID(ctx context.Context, id int64) (db.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) List(ctx context.Context) ([]db.Category, error) {
	return s.repo.List(ctx)
}
