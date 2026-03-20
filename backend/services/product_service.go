package services

import (
	"smartech/backend/models"
	"smartech/backend/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) ListProducts() ([]repositories.ProductWithStock, error) {
	return s.repo.ListActiveWithStock()
}

func (s *ProductService) GetProduct(id int64) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) CreateProduct(input repositories.ProductCreateInput) (*models.Product, error) {
	return s.repo.Create(input)
}

func (s *ProductService) UpdateProduct(id int64, input repositories.ProductUpdateInput) (*models.Product, error) {
	exists, err := s.repo.ExistsByID(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, repositories.ErrProductNotFound
	}

	if err := s.repo.UpdatePartial(id, input); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

func (s *ProductService) DeleteProduct(id int64) error {
	exists, err := s.repo.ExistsByID(id)
	if err != nil {
		return err
	}
	if !exists {
		return repositories.ErrProductNotFound
	}

	return s.repo.SoftDelete(id)
}

func (s *ProductService) GetRandomProducts(limit int) ([]models.Product, error) {
	return s.repo.ListRandom(limit)
}

func (s *ProductService) GetProductsByCategory(category string) ([]models.Product, error) {
	return s.repo.ListByCategory(category)
}
