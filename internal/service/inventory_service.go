package service

import (
	"context"
	"server/internal/models"
	"server/internal/repository"
)

type InventoryService struct{
	repo *repository.InventoryRepository
}

func NewInventoryService(r *repository.InventoryRepository) *InventoryService{
	return &InventoryService{repo: r}
}

func (s *InventoryService) CreateLog(ctx context.Context, log models.InventoryLog) error{
	return s.repo.CreateLog(ctx,log)
}

func (s *InventoryService) GetAllLogs(ctx context.Context) ([]models.InventoryLog, error){
	return s.repo.GetAllLogs(ctx)
}