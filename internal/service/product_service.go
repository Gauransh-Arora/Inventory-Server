package service

import (
	"context"
	"server/internal/models"
	"server/internal/repository"
)

type ProductService struct{
	repo *repository.ProductRepository
}

func NewProductService(r *repository.ProductRepository) *ProductService{
	return &ProductService{repo: r}
}

func(s *ProductService) CreateProduct(ctx context.Context, p models.Product) error{
	return s.repo.CreateProduct(ctx,p)
}

func(s *ProductService) GetAllProducts(ctx context.Context)([]models.Product, error){
	return s.repo.GetAllProducts(ctx)
}

func(s *ProductService) GetProductByBarcode(ctx context.Context, barcode string)([]models.Product, error){
	return s.repo.GetProductByBarcode(ctx,barcode)
}