package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"smartech/backend/models"
	"smartech/backend/repositories"
)

var ErrUploadProductNotFound = errors.New("product not found")

type UploadService struct {
	repo *repositories.UploadRepository
}

func NewUploadService(repo *repositories.UploadRepository) *UploadService {
	return &UploadService{repo: repo}
}

func (s *UploadService) GetProductImages(productID int64) ([]string, error) {
	images, err := s.repo.GetProductImages(productID)
	if err == sql.ErrNoRows {
		return nil, ErrUploadProductNotFound
	}
	if err != nil {
		return nil, err
	}

	list := make([]string, 0)
	if images != "" {
		_ = json.Unmarshal([]byte(images), &list)
	}
	return list, nil
}

func (s *UploadService) UpdateProductImages(productID int64, images []string) (models.Product, error) {
	count, err := s.repo.CountProductByID(productID)
	if err != nil || count == 0 {
		return models.Product{}, ErrUploadProductNotFound
	}

	payload, err := json.Marshal(images)
	if err != nil {
		return models.Product{}, err
	}

	if err := s.repo.UpdateProductImages(productID, string(payload)); err != nil {
		return models.Product{}, err
	}

	return s.repo.GetProductByID(productID)
}
