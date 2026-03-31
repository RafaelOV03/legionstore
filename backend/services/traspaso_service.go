package services

import (
	"database/sql"
	"errors"
	"fmt"
	"smartech/backend/models"
	"smartech/backend/repositories"
)

var (
	ErrTraspasoNotFound              = errors.New("traspaso not found")
	ErrTraspasoSedesIguales          = errors.New("sedes iguales")
	ErrTraspasoSinItems              = errors.New("traspaso sin items")
	ErrTraspasoSoloPendienteEnviar   = errors.New("solo pendientes enviar")
	ErrTraspasoSoloEnviadoRecibir    = errors.New("solo enviados recibir")
	ErrTraspasoSoloPendienteCancelar = errors.New("solo pendientes cancelar")
	ErrTraspasoSoloCanceladoEliminar = errors.New("solo cancelados eliminar")
)

type TraspasoCreateItemInput struct {
	ProductoID int64
	Cantidad   int
}

type TraspasoRecibirItemInput struct {
	ItemID           int64
	CantidadRecibida int
}

type CreateTraspasoInput struct {
	SedeOrigenID  int64
	SedeDestinoID int64
	Notas         string
	Items         []TraspasoCreateItemInput
	UsuarioID     interface{}
}

type RecibirTraspasoInput struct {
	ID        int64
	UsuarioID interface{}
	Notas     string
	Items     []TraspasoRecibirItemInput
}

type TraspasoService struct {
	repo *repositories.TraspasoRepository
}

func NewTraspasoService(repo *repositories.TraspasoRepository) *TraspasoService {
	return &TraspasoService{repo: repo}
}

func (s *TraspasoService) ListTraspasos(estado, sedeOrigenID, sedeDestinoID string) ([]repositories.TraspasoView, error) {
	return s.repo.ListTraspasos(estado, sedeOrigenID, sedeDestinoID)
}

func (s *TraspasoService) GetTraspaso(id int64) (models.Traspaso, []repositories.TraspasoItemView, error) {
	t, err := s.repo.GetTraspasoByID(id)
	if err == sql.ErrNoRows {
		return models.Traspaso{}, nil, ErrTraspasoNotFound
	}
	if err != nil {
		return models.Traspaso{}, nil, err
	}

	items, err := s.repo.ListTraspasoItems(id)
	if err != nil {
		return models.Traspaso{}, nil, err
	}

	return t, items, nil
}

func (s *TraspasoService) CreateTraspaso(input CreateTraspasoInput) (int64, string, error) {
	if input.SedeOrigenID == input.SedeDestinoID {
		return 0, "", ErrTraspasoSedesIguales
	}
	if len(input.Items) == 0 {
		return 0, "", ErrTraspasoSinItems
	}

	sedeOrigenExists, err := s.repo.CountSedeByID(input.SedeOrigenID)
	if err != nil {
		return 0, "", err
	}
	if sedeOrigenExists == 0 {
		return 0, "", fmt.Errorf("Sede origen con ID %d no existe", input.SedeOrigenID)
	}

	sedeDestinoExists, err := s.repo.CountSedeByID(input.SedeDestinoID)
	if err != nil {
		return 0, "", err
	}
	if sedeDestinoExists == 0 {
		return 0, "", fmt.Errorf("Sede destino con ID %d no existe", input.SedeDestinoID)
	}

	for _, item := range input.Items {
		stockDisponible, stockErr := s.repo.GetStockDisponible(item.ProductoID, input.SedeOrigenID)
		nombreProducto, _ := s.repo.GetProductName(item.ProductoID)
		if stockErr == sql.ErrNoRows || stockDisponible == 0 {
			return 0, "", fmt.Errorf("El producto '%s' no tiene stock en la sede origen", nombreProducto)
		}
		if stockErr != nil {
			return 0, "", stockErr
		}
		if stockDisponible < item.Cantidad {
			return 0, "", fmt.Errorf("Stock insuficiente de %s en sede origen (disponible: %d, solicitado: %d)", nombreProducto, stockDisponible, item.Cantidad)
		}
	}

	count, err := s.repo.CountTraspasos()
	if err != nil {
		return 0, "", err
	}
	numero := s.repo.BuildNumeroTraspaso(count)

	repoItems := make([]repositories.TraspasoCreateItemInput, 0, len(input.Items))
	for _, item := range input.Items {
		repoItems = append(repoItems, repositories.TraspasoCreateItemInput{ProductoID: item.ProductoID, Cantidad: item.Cantidad})
	}

	id, err := s.repo.CreateTraspasoWithItems(numero, input.SedeOrigenID, input.SedeDestinoID, input.Notas, input.UsuarioID, repoItems)
	if err != nil {
		return 0, "", err
	}

	return id, numero, nil
}

func (s *TraspasoService) EnviarTraspaso(id int64) error {
	estado, sedeOrigenID, err := s.repo.GetEstadoAndSedeOrigen(id)
	if err == sql.ErrNoRows {
		return ErrTraspasoNotFound
	}
	if err != nil {
		return err
	}

	if estado != "pendiente" {
		return ErrTraspasoSoloPendienteEnviar
	}

	return s.repo.EnviarTraspaso(id, sedeOrigenID)
}

func (s *TraspasoService) RecibirTraspaso(input RecibirTraspasoInput) error {
	estado, sedeDestinoID, err := s.repo.GetEstadoAndSedeDestino(input.ID)
	if err == sql.ErrNoRows {
		return ErrTraspasoNotFound
	}
	if err != nil {
		return err
	}

	if estado != "enviado" {
		return ErrTraspasoSoloEnviadoRecibir
	}

	if len(input.Items) == 0 {
		return s.repo.RecibirTraspasoAuto(input.ID, sedeDestinoID, input.UsuarioID, input.Notas)
	}

	repoItems := make([]repositories.TraspasoRecibirItemInput, 0, len(input.Items))
	for _, item := range input.Items {
		repoItems = append(repoItems, repositories.TraspasoRecibirItemInput{ItemID: item.ItemID, CantidadRecibida: item.CantidadRecibida})
	}
	return s.repo.RecibirTraspasoConItems(input.ID, sedeDestinoID, input.UsuarioID, input.Notas, repoItems)
}

func (s *TraspasoService) CancelarTraspaso(id int64) error {
	estado, err := s.repo.GetEstado(id)
	if err == sql.ErrNoRows {
		return ErrTraspasoNotFound
	}
	if err != nil {
		return err
	}

	if estado != "pendiente" {
		return ErrTraspasoSoloPendienteCancelar
	}

	return s.repo.CancelarTraspaso(id)
}

func (s *TraspasoService) DeleteTraspaso(id int64) error {
	estado, err := s.repo.GetEstado(id)
	if err == sql.ErrNoRows {
		return ErrTraspasoNotFound
	}
	if err != nil {
		return err
	}

	if estado != "cancelado" {
		return ErrTraspasoSoloCanceladoEliminar
	}

	return s.repo.DeleteTraspaso(id)
}
