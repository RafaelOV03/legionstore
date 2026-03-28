package services

import "smartech/backend/repositories"

type StockService struct {
	repo *repositories.StockRepository
}

func NewStockService(repo *repositories.StockRepository) *StockService {
	return &StockService{repo: repo}
}

func (s *StockService) GetStockMultisede() ([]repositories.StockItem, error) {
	return s.repo.ListMultisede()
}

func (s *StockService) GetStockBySede(sedeID int64) ([]repositories.StockProducto, error) {
	return s.repo.ListBySede(sedeID)
}

func (s *StockService) UpdateStock(input repositories.StockUpdateInput) error {
	return s.repo.Upsert(input)
}
