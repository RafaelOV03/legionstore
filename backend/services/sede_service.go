package services

import (
	"smartech/backend/models"
	"smartech/backend/repositories"
)

type SedeService struct {
	repo *repositories.SedeRepository
}

func NewSedeService(repo *repositories.SedeRepository) *SedeService {
	return &SedeService{repo: repo}
}

func (s *SedeService) ListSedes() ([]models.Sede, error) {
	return s.repo.List()
}

func (s *SedeService) GetSede(id int64) (*models.Sede, error) {
	return s.repo.GetByID(id)
}

func (s *SedeService) CreateSede(input models.Sede) (*models.Sede, error) {
	return s.repo.Create(input)
}

func (s *SedeService) UpdateSede(id int64, input models.Sede) (*models.Sede, error) {
	return s.repo.Update(id, input)
}

func (s *SedeService) DeleteSede(id int64) error {
	hasUsers, err := s.repo.HasAssociatedUsers(id)
	if err != nil {
		return err
	}
	if hasUsers {
		return repositories.ErrSedeHasAssociatedUsers
	}

	return s.repo.Delete(id)
}
