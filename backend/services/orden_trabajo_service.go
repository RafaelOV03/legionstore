package services

import (
	"database/sql"
	"errors"
	"fmt"
	"smartech/backend/models"
	"smartech/backend/repositories"
	"time"
)

var (
	ErrOrdenTrabajoNotFound      = errors.New("orden trabajo not found")
	ErrOrdenTrabajoDeleteBlocked = errors.New("orden trabajo delete blocked")
	ErrOTInsumoNotFound          = errors.New("insumo not found")
	ErrInsumoSinStock            = errors.New("insumo sin stock")
)

type CreateOrdenTrabajoInput struct {
	ClienteNombre     string
	ClienteTelefono   string
	Equipo            string
	Marca             string
	Modelo            string
	NumSerie          string
	ProblemaReportado string
	Prioridad         string
	SedeID            int64
	TecnicoID         *int64
	FechaPromesa      string
	Notas             string
	UsuarioID         interface{}
}

type UpdateOrdenTrabajoInput struct {
	DiagnosticoTecnico string
	SolucionAplicada   string
	Estado             string
	Prioridad          string
	TecnicoID          *int64
	CostoManoObra      float64
	CostoRepuestos     float64
	Notas              string
	FechaPromesa       string
	UsuarioID          interface{}
}

type OrdenTrabajoService struct {
	repo *repositories.OrdenTrabajoRepository
}

func NewOrdenTrabajoService(repo *repositories.OrdenTrabajoRepository) *OrdenTrabajoService {
	return &OrdenTrabajoService{repo: repo}
}

func (s *OrdenTrabajoService) ListOrdenesTrabajo(estado, prioridad, tecnicoID, sedeID string) ([]repositories.OrdenTrabajoView, error) {
	return s.repo.ListOrdenesTrabajo(estado, prioridad, tecnicoID, sedeID)
}

func (s *OrdenTrabajoService) GetOrdenTrabajo(id int64) (models.OrdenTrabajo, []repositories.InsumoOrdenView, []repositories.TrazabilidadItem, error) {
	orden, err := s.repo.GetOrdenTrabajoByID(id)
	if err == sql.ErrNoRows {
		return models.OrdenTrabajo{}, nil, nil, ErrOrdenTrabajoNotFound
	}
	if err != nil {
		return models.OrdenTrabajo{}, nil, nil, err
	}

	insumos, err := s.repo.ListInsumosByOrdenID(id)
	if err != nil {
		insumos = []repositories.InsumoOrdenView{}
	}

	trazabilidad, err := s.repo.ListTrazabilidadByOrdenID(id)
	if err != nil {
		trazabilidad = []repositories.TrazabilidadItem{}
	}

	return orden, insumos, trazabilidad, nil
}

func (s *OrdenTrabajoService) CreateOrdenTrabajo(input CreateOrdenTrabajoInput) (int64, string, error) {
	prioridad := input.Prioridad
	if prioridad == "" {
		prioridad = "media"
	}

	var fechaPromesa interface{}
	if input.FechaPromesa != "" {
		fechaPromesa = input.FechaPromesa
	} else {
		fechaPromesa = nil
	}

	count, err := s.repo.CountOrdenesTrabajo()
	if err != nil {
		return 0, "", err
	}
	numeroOrden := s.repo.BuildNumeroOrden(count)

	ordenID, err := s.repo.CreateOrdenTrabajo(repositories.CreateOrdenTrabajoParams{
		NumeroOrden:       numeroOrden,
		ClienteNombre:     input.ClienteNombre,
		ClienteTelefono:   input.ClienteTelefono,
		Equipo:            input.Equipo,
		Marca:             input.Marca,
		Modelo:            input.Modelo,
		NumSerie:          input.NumSerie,
		ProblemaReportado: input.ProblemaReportado,
		Prioridad:         prioridad,
		FechaPromesa:      fechaPromesa,
		TecnicoID:         input.TecnicoID,
		SedeID:            input.SedeID,
		Notas:             input.Notas,
	})
	if err != nil {
		return 0, "", err
	}

	_ = s.repo.InsertTrazabilidad(ordenID, "ingreso", "Equipo ingresado al servicio técnico", input.UsuarioID)
	return ordenID, numeroOrden, nil
}

func (s *OrdenTrabajoService) UpdateOrdenTrabajo(id int64, input UpdateOrdenTrabajoInput) (string, string, error) {
	estadoAnterior, err := s.repo.GetEstadoByOrdenID(id)
	if err == sql.ErrNoRows {
		return "", "", ErrOrdenTrabajoNotFound
	}
	if err != nil {
		return "", "", err
	}

	var fechaPromesa interface{}
	if input.FechaPromesa != "" {
		fechaPromesa = input.FechaPromesa
	} else {
		fechaPromesa = nil
	}

	var fechaEntrega interface{}
	if input.Estado == "entregado" && estadoAnterior != "entregado" {
		fechaEntrega = time.Now()
	} else {
		fechaEntrega = nil
	}

	err = s.repo.UpdateOrdenTrabajo(id, repositories.UpdateOrdenTrabajoParams{
		DiagnosticoTecnico: input.DiagnosticoTecnico,
		SolucionAplicada:   input.SolucionAplicada,
		Estado:             input.Estado,
		Prioridad:          input.Prioridad,
		TecnicoID:          input.TecnicoID,
		CostoManoObra:      input.CostoManoObra,
		CostoRepuestos:     input.CostoRepuestos,
		Notas:              input.Notas,
		FechaPromesa:       fechaPromesa,
		FechaEntrega:       fechaEntrega,
	})
	if err != nil {
		return "", "", err
	}

	if input.Estado != estadoAnterior {
		detalle := fmt.Sprintf("Estado cambiado de %s a %s", estadoAnterior, input.Estado)
		_ = s.repo.InsertTrazabilidad(id, "cambio_estado", detalle, input.UsuarioID)
	}

	return estadoAnterior, input.Estado, nil
}

func (s *OrdenTrabajoService) AsignarTecnico(id, tecnicoID int64, usuarioID interface{}) (string, error) {
	_, err := s.repo.GetEstadoByOrdenID(id)
	if err == sql.ErrNoRows {
		return "", ErrOrdenTrabajoNotFound
	}
	if err != nil {
		return "", err
	}

	tecnicoNombre, _ := s.repo.GetTecnicoNombre(tecnicoID)
	if err := s.repo.AssignTecnico(id, tecnicoID); err != nil {
		return "", err
	}

	detalle := fmt.Sprintf("Técnico asignado: %s", tecnicoNombre)
	_ = s.repo.InsertTrazabilidad(id, "asignacion", detalle, usuarioID)
	return tecnicoNombre, nil
}

func (s *OrdenTrabajoService) AgregarInsumo(ordenID, insumoID int64, cantidad int, usuarioID interface{}) (string, int, error) {
	insumoNombre, stockActual, err := s.repo.GetInsumoStockAndNombre(insumoID)
	if err == sql.ErrNoRows {
		return "", 0, ErrOTInsumoNotFound
	}
	if err != nil {
		return "", 0, err
	}

	if stockActual < cantidad {
		return insumoNombre, stockActual, ErrInsumoSinStock
	}

	if err := s.repo.AddInsumoToOrden(ordenID, insumoID, cantidad); err != nil {
		return "", 0, err
	}

	detalle := fmt.Sprintf("Insumo utilizado: %s x%d", insumoNombre, cantidad)
	_ = s.repo.InsertTrazabilidad(ordenID, "insumo_usado", detalle, usuarioID)
	return insumoNombre, stockActual, nil
}

func (s *OrdenTrabajoService) RegistrarTrazabilidad(ordenID int64, accion, detalle string, usuarioID interface{}) error {
	return s.repo.InsertTrazabilidad(ordenID, accion, detalle, usuarioID)
}

func (s *OrdenTrabajoService) DeleteOrdenTrabajo(id int64) error {
	estado, err := s.repo.GetEstadoByOrdenID(id)
	if err == sql.ErrNoRows {
		return ErrOrdenTrabajoNotFound
	}
	if err != nil {
		return err
	}

	if estado != "recibido" && estado != "cancelado" {
		return ErrOrdenTrabajoDeleteBlocked
	}

	return s.repo.DeleteOrdenTrabajo(id)
}

func (s *OrdenTrabajoService) Stats() (repositories.OrdenesStats, error) {
	return s.repo.Stats()
}

func (s *OrdenTrabajoService) ListTecnicos() ([]repositories.TecnicoView, error) {
	return s.repo.ListTecnicos()
}
