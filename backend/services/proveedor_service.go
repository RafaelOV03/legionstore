package services

import (
	"database/sql"
	"errors"
	"smartech/backend/models"
	"smartech/backend/repositories"
)

var (
	ErrProveedorNotFound         = errors.New("proveedor not found")
	ErrProveedorHasPendingDeudas = errors.New("proveedor has pending deudas")
	ErrDeudaNotFound             = errors.New("deuda not found")
	ErrPagoExcedeSaldo           = errors.New("pago excede saldo")
)

type ProveedorService struct {
	repo *repositories.ProveedorRepository
}

func NewProveedorService(repo *repositories.ProveedorRepository) *ProveedorService {
	return &ProveedorService{repo: repo}
}

func (s *ProveedorService) ListProveedores() ([]models.Proveedor, error) {
	return s.repo.ListProveedores()
}

func (s *ProveedorService) GetProveedor(id int64) (models.Proveedor, []repositories.DeudaSimple, float64, error) {
	p, err := s.repo.GetProveedorByID(id)
	if err == sql.ErrNoRows {
		return models.Proveedor{}, nil, 0, ErrProveedorNotFound
	}
	if err != nil {
		return models.Proveedor{}, nil, 0, err
	}

	deudas, err := s.repo.ListPendingDeudasByProveedor(id)
	if err != nil {
		return models.Proveedor{}, nil, 0, err
	}

	var total float64
	for _, d := range deudas {
		total += d.MontoTotal - d.MontoPagado
	}

	return p, deudas, total, nil
}

func (s *ProveedorService) CreateProveedor(p models.Proveedor) (int64, error) {
	return s.repo.InsertProveedor(p)
}

func (s *ProveedorService) UpdateProveedor(id int64, p models.Proveedor) error {
	return s.repo.UpdateProveedor(id, p)
}

func (s *ProveedorService) DeleteProveedor(id int64) error {
	count, err := s.repo.CountPendingDeudas(id)
	if err == nil && count > 0 {
		return ErrProveedorHasPendingDeudas
	}
	return s.repo.DeleteProveedorCascade(id)
}

func (s *ProveedorService) ListDeudas(estado, proveedorID string) ([]repositories.DeudaView, error) {
	return s.repo.ListDeudas(estado, proveedorID)
}

func (s *ProveedorService) CreateDeuda(proveedorID int64, numeroFactura string, montoTotal float64, fechaVenc interface{}, notas string) (int64, error) {
	return s.repo.InsertDeuda(proveedorID, numeroFactura, montoTotal, fechaVenc, notas)
}

func (s *ProveedorService) RegistrarPago(deudaID int64, monto float64, metodo, referencia string, usuarioID interface{}) (float64, string, error) {
	montoTotal, montoPagado, err := s.repo.GetDeudaMontos(deudaID)
	if err == sql.ErrNoRows {
		return 0, "", ErrDeudaNotFound
	}
	if err != nil {
		return 0, "", err
	}

	saldoPendiente := montoTotal - montoPagado
	if monto > saldoPendiente {
		return 0, "", ErrPagoExcedeSaldo
	}

	if err := s.repo.RegisterPago(deudaID, monto, metodo, referencia, usuarioID); err != nil {
		return 0, "", err
	}

	nuevoMontoPagado := montoPagado + monto
	nuevoEstado := "pendiente"
	if nuevoMontoPagado >= montoTotal {
		nuevoEstado = "pagada"
	} else if nuevoMontoPagado > 0 {
		nuevoEstado = "parcial"
	}

	return montoTotal - nuevoMontoPagado, nuevoEstado, nil
}

func (s *ProveedorService) ListPagosDeuda(deudaID int64) ([]repositories.PagoView, error) {
	return s.repo.ListPagosByDeuda(deudaID)
}

func (s *ProveedorService) ResumenDeudas() ([]repositories.ResumenProveedor, float64, int, int, int, error) {
	total, pendientes, vencidas, pagadas, err := s.repo.ResumenTotales()
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	resumen, err := s.repo.ResumenPorProveedor()
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return resumen, total, pendientes, vencidas, pagadas, nil
}
