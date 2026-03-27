package service

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/hoangquan/retail-store-api/db/sqlc"
	"github.com/hoangquan/retail-store-api/internal/model/request"
	"github.com/hoangquan/retail-store-api/internal/repository"
	pkgNats "github.com/hoangquan/retail-store-api/pkg/nats"
)

type ProductService struct {
	repo      *repository.ProductRepository
	publisher *pkgNats.Publisher
}

func NewProductService(queries *db.Queries, publisher *pkgNats.Publisher) *ProductService {
	return &ProductService{
		repo:      repository.NewProductRepository(queries),
		publisher: publisher,
	}
}

func (s *ProductService) Create(ctx context.Context, req request.CreateProductRequest) (db.Product, error) {
	price := pgtype.Numeric{}
	_ = price.Scan(strconv.FormatFloat(req.Price, 'f', 2, 64))

	var categoryID pgtype.Int8
	if req.CategoryID > 0 {
		categoryID = pgtype.Int8{Int64: req.CategoryID, Valid: true}
	}

	product, err := s.repo.Create(ctx, db.CreateProductParams{
		Name:        req.Name,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		Price:       price,
		Quantity:    req.Quantity,
		CategoryID:  categoryID,
	})
	if err != nil {
		return product, err
	}

	s.publishEvent(ctx, pkgNats.SubjectProductCreated, map[string]interface{}{
		"id":          product.ID,
		"name":        product.Name,
		"description": req.Description,
		"price":       req.Price,
		"quantity":    product.Quantity,
		"category_id": req.CategoryID,
	})

	return product, nil
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

func (s *ProductService) Update(ctx context.Context, id int64, req request.UpdateProductRequest) (db.Product, error) {
	price := pgtype.Numeric{}
	_ = price.Scan(strconv.FormatFloat(req.Price, 'f', 2, 64))

	var categoryID pgtype.Int8
	if req.CategoryID > 0 {
		categoryID = pgtype.Int8{Int64: req.CategoryID, Valid: true}
	}

	product, err := s.repo.Update(ctx, db.UpdateProductParams{
		ID:          id,
		Name:        req.Name,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		Price:       price,
		Quantity:    req.Quantity,
		CategoryID:  categoryID,
	})
	if err != nil {
		return product, err
	}

	s.publishEvent(ctx, pkgNats.SubjectProductUpdated, map[string]interface{}{
		"id":          product.ID,
		"name":        product.Name,
		"description": req.Description,
		"price":       req.Price,
		"quantity":    product.Quantity,
		"category_id": req.CategoryID,
	})

	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.publishEvent(ctx, pkgNats.SubjectProductDeleted, map[string]interface{}{
		"id": id,
	})

	return nil
}

func (s *ProductService) publishEvent(ctx context.Context, subject string, data interface{}) {
	if s.publisher == nil {
		return
	}
	if err := s.publisher.Publish(ctx, subject, data); err != nil {
		slog.Error("Failed to publish NATS event", "subject", subject, "error", err)
	}
}
