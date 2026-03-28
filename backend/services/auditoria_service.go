package services

import (
	"smartech/backend/models"
	"smartech/backend/repositories"
)

type GetLogsInput struct {
	Accion     string
	Entidad    string
	UsuarioID  string
	FechaDesde string
	FechaHasta string
	Limit      string
}

type GananciasInput struct {
	FechaDesde string
	FechaHasta string
	SedeID     string
}

type GananciasReporte struct {
	TotalVentas        float64                  `json:"total_ventas"`
	CostoProductos     float64                  `json:"costo_productos"`
	GananciaProductos  float64                  `json:"ganancia_productos"`
	TotalServicios     float64                  `json:"total_servicios"`
	GananciaTotal      float64                  `json:"ganancia_total"`
	VentasPorCategoria []repositories.CatStats  `json:"ventas_por_categoria"`
	VentasPorSede      []repositories.SedeStats `json:"ventas_por_sede"`
}

type LogStatsResponse struct {
	AccionesPorTipo     []repositories.AccionStats `json:"acciones_por_tipo"`
	ActividadPorUsuario []repositories.UserStats   `json:"actividad_por_usuario"`
	ActividadPorDia     []repositories.DayStats    `json:"actividad_por_dia"`
}

type CreateSegmentacionInput struct {
	Nombre      string
	Descripcion string
	Criterios   string
}

type CreatePromocionInput struct {
	Nombre         string
	Descripcion    string
	Tipo           string
	Valor          float64
	FechaInicio    string
	FechaFin       string
	ProductosIDs   string
	SegmentacionID *int64
}

type UpdatePromocionInput struct {
	Nombre         string
	Descripcion    string
	Tipo           string
	Valor          float64
	FechaInicio    string
	FechaFin       string
	ProductosIDs   string
	SegmentacionID *int64
	Activa         bool
}

type AuditoriaService struct {
	repo *repositories.AuditoriaRepository
}

func NewAuditoriaService(repo *repositories.AuditoriaRepository) *AuditoriaService {
	return &AuditoriaService{repo: repo}
}

func (s *AuditoriaService) ListLogs(input GetLogsInput) ([]repositories.LogView, error) {
	return s.repo.ListLogs(input.Accion, input.Entidad, input.UsuarioID, input.FechaDesde, input.FechaHasta, input.Limit)
}

func (s *AuditoriaService) LogStats() (LogStatsResponse, error) {
	accionesPorTipo, err := s.repo.ListAccionStats()
	if err != nil {
		return LogStatsResponse{}, err
	}

	actividadPorUsuario, err := s.repo.ListActividadPorUsuario()
	if err != nil {
		return LogStatsResponse{}, err
	}

	actividadPorDia, err := s.repo.ListActividadPorDia()
	if err != nil {
		return LogStatsResponse{}, err
	}

	return LogStatsResponse{
		AccionesPorTipo:     accionesPorTipo,
		ActividadPorUsuario: actividadPorUsuario,
		ActividadPorDia:     actividadPorDia,
	}, nil
}

func (s *AuditoriaService) ReporteGanancias(input GananciasInput) (GananciasReporte, error) {
	totalVentas, err := s.repo.TotalVentas(input.FechaDesde, input.FechaHasta, input.SedeID)
	if err != nil {
		return GananciasReporte{}, err
	}

	costoProductos, err := s.repo.CostoProductos(input.FechaDesde, input.FechaHasta, input.SedeID)
	if err != nil {
		return GananciasReporte{}, err
	}

	totalServicios, err := s.repo.TotalServicios(input.FechaDesde, input.FechaHasta, input.SedeID)
	if err != nil {
		return GananciasReporte{}, err
	}

	ventasPorCategoria, err := s.repo.ListVentasPorCategoria()
	if err != nil {
		return GananciasReporte{}, err
	}

	ventasPorSede, err := s.repo.ListVentasPorSede()
	if err != nil {
		return GananciasReporte{}, err
	}

	gananciaProductos := totalVentas - costoProductos
	return GananciasReporte{
		TotalVentas:        totalVentas,
		CostoProductos:     costoProductos,
		GananciaProductos:  gananciaProductos,
		TotalServicios:     totalServicios,
		GananciaTotal:      gananciaProductos + totalServicios,
		VentasPorCategoria: ventasPorCategoria,
		VentasPorSede:      ventasPorSede,
	}, nil
}

func (s *AuditoriaService) ListSegmentaciones() ([]models.Segmentacion, error) {
	return s.repo.ListSegmentaciones()
}

func (s *AuditoriaService) CreateSegmentacion(input CreateSegmentacionInput) (int64, error) {
	return s.repo.CreateSegmentacion(input.Nombre, input.Descripcion, input.Criterios)
}

func (s *AuditoriaService) ListPromociones(activas bool) ([]models.Promocion, error) {
	return s.repo.ListPromociones(activas)
}

func (s *AuditoriaService) CreatePromocion(input CreatePromocionInput) (int64, error) {
	return s.repo.CreatePromocion(repositories.PromocionCreateInput{
		Nombre:         input.Nombre,
		Descripcion:    input.Descripcion,
		Tipo:           input.Tipo,
		Valor:          input.Valor,
		FechaInicio:    input.FechaInicio,
		FechaFin:       input.FechaFin,
		ProductosIDs:   input.ProductosIDs,
		SegmentacionID: input.SegmentacionID,
	})
}

func (s *AuditoriaService) UpdatePromocion(id int64, input UpdatePromocionInput) error {
	return s.repo.UpdatePromocion(id, repositories.PromocionUpdateInput{
		Nombre:         input.Nombre,
		Descripcion:    input.Descripcion,
		Tipo:           input.Tipo,
		Valor:          input.Valor,
		FechaInicio:    input.FechaInicio,
		FechaFin:       input.FechaFin,
		ProductosIDs:   input.ProductosIDs,
		SegmentacionID: input.SegmentacionID,
		Activa:         input.Activa,
	})
}

func (s *AuditoriaService) DeletePromocion(id int64) error {
	return s.repo.DeletePromocion(id)
}
