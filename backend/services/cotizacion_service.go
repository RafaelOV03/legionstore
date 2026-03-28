package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"smartech/backend/models"
	"smartech/backend/repositories"
)

var (
	ErrCotizacionNotFound   = errors.New("cotizacion not found")
	ErrCotizacionSinItems   = errors.New("cotizacion sin items")
	ErrEstadoCotizacionBad  = errors.New("estado invalido")
	ErrCotizacionNoAprobada = errors.New("cotizacion no aprobada")
)

type CotizacionCreateItemInput struct {
	ProductoID     int64
	Cantidad       int
	PrecioUnitario float64
}

type CreateCotizacionInput struct {
	ClienteNombre   string
	ClienteTelefono string
	ClienteEmail    string
	ValidezDias     int
	Descuento       float64
	Notas           string
	SedeID          int64
	UsuarioID       interface{}
	Items           []CotizacionCreateItemInput
}

type CotizacionService struct {
	repo *repositories.CotizacionRepository
}

func NewCotizacionService(repo *repositories.CotizacionRepository) *CotizacionService {
	return &CotizacionService{repo: repo}
}

func (s *CotizacionService) ListCotizaciones(estado, sedeID string) ([]repositories.CotizacionView, error) {
	return s.repo.ListCotizaciones(estado, sedeID)
}

func (s *CotizacionService) GetCotizacion(id int64) (models.Cotizacion, []repositories.CotizacionItemView, error) {
	c, err := s.repo.GetCotizacionByID(id)
	if err == sql.ErrNoRows {
		return models.Cotizacion{}, nil, ErrCotizacionNotFound
	}
	if err != nil {
		return models.Cotizacion{}, nil, err
	}
	items, _ := s.repo.ListCotizacionItems(id)
	return c, items, nil
}

func (s *CotizacionService) CreateCotizacion(input CreateCotizacionInput) (int64, string, float64, error) {
	if len(input.Items) == 0 {
		return 0, "", 0, ErrCotizacionSinItems
	}

	count, _ := s.repo.CountCotizaciones()
	numero := s.repo.BuildNumeroCotizacion(count)

	var total float64
	items := make([]repositories.CotizacionCreateItem, 0, len(input.Items))
	for _, item := range input.Items {
		total += float64(item.Cantidad) * item.PrecioUnitario
		items = append(items, repositories.CotizacionCreateItem{ProductoID: item.ProductoID, Cantidad: item.Cantidad, PrecioUnitario: item.PrecioUnitario})
	}
	total = total - input.Descuento

	validez := input.ValidezDias
	if validez == 0 {
		validez = 30
	}

	id, err := s.repo.CreateCotizacionWithItems(numero, input.ClienteNombre, input.ClienteTelefono, input.ClienteEmail, validez, total, input.Descuento, input.Notas, input.UsuarioID, input.SedeID, items)
	if err != nil {
		return 0, "", 0, err
	}
	return id, numero, total, nil
}

func (s *CotizacionService) UpdateEstado(id int64, estado string) (string, error) {
	validEstados := map[string]bool{"pendiente": true, "aprobada": true, "rechazada": true, "vencida": true, "convertida": true}
	if !validEstados[estado] {
		return "", ErrEstadoCotizacionBad
	}

	estadoAnterior, err := s.repo.GetEstadoByID(id)
	if err != nil {
		return "", err
	}
	if err := s.repo.UpdateEstado(id, estado); err != nil {
		return "", err
	}
	return estadoAnterior, nil
}

func (s *CotizacionService) DeleteCotizacion(id int64) error {
	return s.repo.DeleteCotizacion(id)
}

func (s *CotizacionService) ConvertirAVenta(id int64, userID interface{}) (int64, string, string, error) {
	estado, sedeID, total, clienteNombre, err := s.repo.GetCotizacionConversionData(id)
	if err == sql.ErrNoRows {
		return 0, "", "", ErrCotizacionNotFound
	}
	if err != nil {
		return 0, "", "", err
	}
	if estado != "aprobada" {
		return 0, "", "", ErrCotizacionNoAprobada
	}

	ventaID, numeroVenta, items, err := s.repo.ConvertToVenta(id, userID, sedeID, total, clienteNombre)
	if err != nil {
		return 0, "", "", err
	}
	itemsJSON, _ := json.Marshal(items)
	return ventaID, numeroVenta, string(itemsJSON), nil
}

func (s *CotizacionService) GetPDFData(id int64) (repositories.CotizacionPDF, []repositories.CotizacionPDFItem, error) {
	cot, err := s.repo.GetCotizacionPDF(id)
	if err == sql.ErrNoRows {
		return repositories.CotizacionPDF{}, nil, ErrCotizacionNotFound
	}
	if err != nil {
		return repositories.CotizacionPDF{}, nil, err
	}
	items, _ := s.repo.ListCotizacionPDFItems(id)
	return cot, items, nil
}
