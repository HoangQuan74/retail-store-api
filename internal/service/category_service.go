package service

import (
	"context"

	db "github.com/kainguyen/retail-store-api/db/sqlc"
	"github.com/kainguyen/retail-store-api/internal/model/request"
	"github.com/kainguyen/retail-store-api/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(queries *db.Queries) *CategoryService {
	return &CategoryService{repo: repository.NewCategoryRepository(queries)}
}

func (s *CategoryService) Create(ctx context.Context, req request.CreateCategoryRequest) (db.Category, error) {
	return s.repo.Create(ctx, req.Name)
}

func (s *CategoryService) GetByID(ctx context.Context, id int64) (db.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) List(ctx context.Context) ([]db.Category, error) {
	return s.repo.List(ctx)
}

func (s *CategoryService) Update(ctx context.Context, id int64, req request.UpdateCategoryRequest) (db.Category, error) {
	return s.repo.Update(ctx, db.UpdateCategoryParams{
		ID:   id,
		Name: req.Name,
	})
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
