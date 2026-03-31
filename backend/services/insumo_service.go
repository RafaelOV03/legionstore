package services

import (
	"database/sql"
	"errors"
	"smartech/backend/models"
	"smartech/backend/repositories"
)

var (
	ErrInsumoNotFound     = errors.New("insumo not found")
	ErrStockInsuficiente  = errors.New("stock insuficiente")
	ErrTipoAjusteInvalido = errors.New("tipo invalido")
	ErrProductoNotFound   = errors.New("producto not found")
	ErrCompatibilidadDup  = errors.New("compatibilidad ya existe")
	ErrCompatibilidadSelf = errors.New("producto mismo")
)

type AjusteStockInput struct {
	Cantidad int
	Tipo     string
	Motivo   string
}

type CreateCompatibilidadInput struct {
	ProductoID    int64
	CompatibleCon int64
	Notas         string
}

type CompatiblesResult struct {
	Producto    repositories.ProductoBase             `json:"producto"`
	Compatibles []repositories.CompatibleProductoView `json:"compatibles"`
}

type InsumoService struct {
	repo *repositories.InsumoRepository
}

func NewInsumoService(repo *repositories.InsumoRepository) *InsumoService {
	return &InsumoService{repo: repo}
}

func (s *InsumoService) ListInsumos(categoria string, bajoStock bool) ([]models.Insumo, error) {
	return s.repo.ListInsumos(categoria, bajoStock)
}

func (s *InsumoService) GetInsumo(id int64) (models.Insumo, error) {
	i, err := s.repo.GetInsumoByID(id)
	if err == sql.ErrNoRows {
		return models.Insumo{}, ErrInsumoNotFound
	}
	return i, err
}

func (s *InsumoService) CreateInsumo(i models.Insumo) (int64, error) {
	return s.repo.InsertInsumo(i)
}

func (s *InsumoService) UpdateInsumo(id int64, i models.Insumo) error {
	return s.repo.UpdateInsumo(id, i)
}

func (s *InsumoService) AjustarStock(id int64, input AjusteStockInput) (int, int, error) {
	stockActual, err := s.repo.GetStock(id)
	if err == sql.ErrNoRows {
		return 0, 0, ErrInsumoNotFound
	}
	if err != nil {
		return 0, 0, err
	}

	nuevoStock := stockActual
	if input.Tipo == "entrada" {
		nuevoStock = stockActual + input.Cantidad
	} else if input.Tipo == "salida" {
		if stockActual < input.Cantidad {
			return stockActual, stockActual, ErrStockInsuficiente
		}
		nuevoStock = stockActual - input.Cantidad
	} else {
		return stockActual, stockActual, ErrTipoAjusteInvalido
	}

	if err := s.repo.UpdateStock(id, nuevoStock); err != nil {
		return stockActual, stockActual, err
	}
	return stockActual, nuevoStock, nil
}

func (s *InsumoService) DeleteInsumo(id int64) error {
	return s.repo.DeleteInsumo(id)
}

func (s *InsumoService) ListCompatibilidades(productoID string) ([]repositories.CompatibilidadView, error) {
	return s.repo.ListCompatibilidades(productoID)
}

func (s *InsumoService) BuscarCompatibles(productoID string) (CompatiblesResult, error) {
	producto, err := s.repo.GetProductoBase(productoID)
	if err == sql.ErrNoRows {
		return CompatiblesResult{}, ErrProductoNotFound
	}
	if err != nil {
		return CompatiblesResult{}, err
	}

	directos, _ := s.repo.ListCompatiblesDirectos(productoID)
	seen := make(map[int64]bool)
	for _, d := range directos {
		seen[d.ID] = true
	}

	adicionales, _ := s.repo.ListCompatiblesByCategoriaMarca(productoID, producto.Category, producto.Brand)
	for _, c := range adicionales {
		if seen[c.ID] {
			continue
		}
		if c.Category == producto.Category {
			c.TipoMatch = "misma_categoria"
		} else {
			c.TipoMatch = "mismo_fabricante"
		}
		directos = append(directos, c)
	}

	return CompatiblesResult{Producto: producto, Compatibles: directos}, nil
}

func (s *InsumoService) CreateCompatibilidad(input CreateCompatibilidadInput) (int64, error) {
	if input.ProductoID == input.CompatibleCon {
		return 0, ErrCompatibilidadSelf
	}
	count, _ := s.repo.CountCompatibilidadPair(input.ProductoID, input.CompatibleCon)
	if count > 0 {
		return 0, ErrCompatibilidadDup
	}
	return s.repo.InsertCompatibilidad(input.ProductoID, input.CompatibleCon, input.Notas)
}

func (s *InsumoService) DeleteCompatibilidad(id int64) error {
	return s.repo.DeleteCompatibilidad(id)
}

func (s *InsumoService) Stats() (repositories.InsumoStats, error) {
	return s.repo.Stats()
}
