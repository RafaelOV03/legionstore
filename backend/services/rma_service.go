package services

import (
	"database/sql"
	"errors"
	"smartech/backend/models"
	"smartech/backend/repositories"
	"time"
)

var ErrRMANotFound = errors.New("rma not found")

type CreateRMAInput struct {
	ProductoID       int64
	ClienteNombre    string
	ClienteTelefono  string
	ClienteEmail     string
	NumSerie         string
	FechaCompra      string
	MotivoDevolucion string
	SedeID           int64
	Notas            string
}

type UpdateRMAInput struct {
	Diagnostico string
	Estado      string
	Solucion    string
	Notas       string
	Comentario  string
}

type RMAService struct {
	repo *repositories.RMARepository
}

func NewRMAService(repo *repositories.RMARepository) *RMAService {
	return &RMAService{repo: repo}
}

func (s *RMAService) ListRMAs(estado, sedeID string) ([]repositories.RMAView, error) {
	return s.repo.ListRMAs(estado, sedeID)
}

func (s *RMAService) GetRMA(id int64) (models.RMA, []repositories.HistorialItem, error) {
	rma, err := s.repo.GetRMAByID(id)
	if err == sql.ErrNoRows {
		return models.RMA{}, nil, ErrRMANotFound
	}
	if err != nil {
		return models.RMA{}, nil, err
	}

	historial, err := s.repo.ListHistorial(id)
	if err != nil {
		historial = []repositories.HistorialItem{}
	}
	return rma, historial, nil
}

func (s *RMAService) CreateRMA(input CreateRMAInput, userID interface{}) (int64, string, error) {
	count, _ := s.repo.CountRMAs()
	numero := s.repo.BuildRMANumber(count)

	var fechaCompra interface{}
	if input.FechaCompra != "" {
		fechaCompra = input.FechaCompra
	} else {
		fechaCompra = nil
	}

	rmaID, err := s.repo.InsertRMA(numero, models.RMA{
		ProductoID:       input.ProductoID,
		ClienteNombre:    input.ClienteNombre,
		ClienteTelefono:  input.ClienteTelefono,
		ClienteEmail:     input.ClienteEmail,
		NumSerie:         input.NumSerie,
		MotivoDevolucion: input.MotivoDevolucion,
		SedeID:           input.SedeID,
		Notas:            input.Notas,
	}, fechaCompra, userID)
	if err != nil {
		return 0, "", err
	}

	_ = s.repo.InsertHistorial(rmaID, "", "recibido", "RMA creada", userID)
	return rmaID, numero, nil
}

func (s *RMAService) UpdateRMA(id int64, input UpdateRMAInput, userID interface{}) (string, error) {
	estadoAnterior, err := s.repo.GetEstadoByID(id)
	if err == sql.ErrNoRows {
		return "", ErrRMANotFound
	}
	if err != nil {
		return "", err
	}

	var fechaResolucion interface{}
	if input.Estado == "resuelto" && estadoAnterior != "resuelto" {
		fechaResolucion = time.Now()
	} else {
		fechaResolucion = nil
	}

	if err := s.repo.UpdateRMA(id, input.Diagnostico, input.Estado, input.Solucion, input.Notas, fechaResolucion); err != nil {
		return "", err
	}

	if input.Estado != estadoAnterior {
		_ = s.repo.InsertHistorial(id, estadoAnterior, input.Estado, input.Comentario, userID)
	}

	return estadoAnterior, nil
}

func (s *RMAService) DeleteRMA(id int64) error {
	return s.repo.DeleteRMA(id)
}

func (s *RMAService) Stats() (repositories.RMAStats, error) {
	return s.repo.Stats()
}
