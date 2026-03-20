import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Form, Button, Table, Spinner, Alert } from 'react-bootstrap';
import { GraphUp, Calendar, CashStack, GraphUpArrow, GraphDownArrow } from 'react-bootstrap-icons';
import * as api from '../services/resourceApi';

function Reportes() {
  const [reporte, setReporte] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  
  const [filters, setFilters] = useState({
    fecha_inicio: new Date(new Date().getFullYear(), new Date().getMonth(), 1).toISOString().split('T')[0],
    fecha_fin: new Date().toISOString().split('T')[0],
    sede_id: ''
  });

  const [sedes, setSedes] = useState([]);

  useEffect(() => {
    loadSedes();
    handleGenerarReporte();
  }, []);

  const loadSedes = async () => {
    try {
      const data = await api.getSedes();
      setSedes(data || []);
    } catch (err) {
      console.error(err);
    }
  };

  const handleGenerarReporte = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await api.getReporteGanancias(filters);
      setReporte(data);
    } catch (err) {
      setError('Error generando reporte');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (value) => {
    return `Bs. ${(value || 0).toFixed(2)}`;
  };

  return (
    <Container fluid className="py-4">
      <Row className="mb-4">
        <Col>
          <h2><GraphUp className="me-2" />Reporte de Ganancias</h2>
          <p className="text-muted">Análisis de ingresos, costos y ganancias</p>
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}

      {/* Filtros */}
      <Card className="mb-4" style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
        <Card.Body>
          <Row className="g-3 align-items-end">
            <Col md={3}>
              <Form.Group>
                <Form.Label><Calendar className="me-1" />Fecha Inicio</Form.Label>
                <Form.Control
                  type="date"
                  value={filters.fecha_inicio}
                  onChange={(e) => setFilters({...filters, fecha_inicio: e.target.value})}
                />
              </Form.Group>
            </Col>
            <Col md={3}>
              <Form.Group>
                <Form.Label><Calendar className="me-1" />Fecha Fin</Form.Label>
                <Form.Control
                  type="date"
                  value={filters.fecha_fin}
                  onChange={(e) => setFilters({...filters, fecha_fin: e.target.value})}
                />
              </Form.Group>
            </Col>
            <Col md={3}>
              <Form.Group>
                <Form.Label>Sede</Form.Label>
                <Form.Select
                  value={filters.sede_id}
                  onChange={(e) => setFilters({...filters, sede_id: e.target.value})}
                >
                  <option value="">Todas las sedes</option>
                  {sedes.map(s => (
                    <option key={s.id} value={s.id}>{s.nombre}</option>
                  ))}
                </Form.Select>
              </Form.Group>
            </Col>
            <Col md={3}>
              <Button 
                variant="primary" 
                onClick={handleGenerarReporte}
                disabled={loading}
                className="w-100"
              >
                {loading ? (
                  <Spinner animation="border" size="sm" />
                ) : (
                  <>Generar Reporte</>
                )}
              </Button>
            </Col>
          </Row>
        </Card.Body>
      </Card>

      {reporte && (
        <>
          {/* Resumen Principal */}
          <Row className="mb-4 g-3">
            <Col md={3}>
              <Card style={{ background: 'linear-gradient(135deg, #198754, #20c997)', border: 'none' }}>
                <Card.Body className="text-white">
                  <div className="d-flex justify-content-between align-items-center">
                    <div>
                      <h6 className="opacity-75 mb-1">Ingresos Totales</h6>
                      <h3 className="mb-0">{formatCurrency(reporte.ingresos_totales)}</h3>
                    </div>
                    <GraphUpArrow size={40} className="opacity-50" />
                  </div>
                </Card.Body>
              </Card>
            </Col>
            <Col md={3}>
              <Card style={{ background: 'linear-gradient(135deg, #dc3545, #fd7e14)', border: 'none' }}>
                <Card.Body className="text-white">
                  <div className="d-flex justify-content-between align-items-center">
                    <div>
                      <h6 className="opacity-75 mb-1">Costos Totales</h6>
                      <h3 className="mb-0">{formatCurrency(reporte.costos_totales)}</h3>
                    </div>
                    <GraphDownArrow size={40} className="opacity-50" />
                  </div>
                </Card.Body>
              </Card>
            </Col>
            <Col md={3}>
              <Card style={{ background: 'linear-gradient(135deg, #0d6efd, #6f42c1)', border: 'none' }}>
                <Card.Body className="text-white">
                  <div className="d-flex justify-content-between align-items-center">
                    <div>
                      <h6 className="opacity-75 mb-1">Ganancia Neta</h6>
                      <h3 className="mb-0">{formatCurrency(reporte.ganancia_neta)}</h3>
                    </div>
                    <CashStack size={40} className="opacity-50" />
                  </div>
                </Card.Body>
              </Card>
            </Col>
            <Col md={3}>
              <Card style={{ background: 'linear-gradient(135deg, #ffc107, #fd7e14)', border: 'none' }}>
                <Card.Body className="text-white">
                  <div className="d-flex justify-content-between align-items-center">
                    <div>
                      <h6 className="opacity-75 mb-1">Margen de Ganancia</h6>
                      <h3 className="mb-0">{(reporte.margen_ganancia || 0).toFixed(1)}%</h3>
                    </div>
                    <GraphUp size={40} className="opacity-50" />
                  </div>
                </Card.Body>
              </Card>
            </Col>
          </Row>

          {/* Desglose */}
          <Row className="g-4">
            <Col lg={6}>
              <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
                <Card.Header className="bg-transparent">
                  <h6 className="mb-0">Ingresos por Categoría</h6>
                </Card.Header>
                <Card.Body className="p-0">
                  <Table hover className="mb-0">
                    <thead>
                      <tr>
                        <th>Categoría</th>
                        <th className="text-end">Ventas</th>
                        <th className="text-end">Ingresos</th>
                      </tr>
                    </thead>
                    <tbody>
                      {(reporte.por_categoria || []).length === 0 ? (
                        <tr>
                          <td colSpan="3" className="text-center py-4 text-muted">
                            Sin datos en este período
                          </td>
                        </tr>
                      ) : (
                        (reporte.por_categoria || []).map((cat, idx) => (
                          <tr key={idx}>
                            <td>{cat.categoria || 'Sin categoría'}</td>
                            <td className="text-end">{cat.cantidad || 0}</td>
                            <td className="text-end">{formatCurrency(cat.ingresos)}</td>
                          </tr>
                        ))
                      )}
                    </tbody>
                  </Table>
                </Card.Body>
              </Card>
            </Col>

            <Col lg={6}>
              <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
                <Card.Header className="bg-transparent">
                  <h6 className="mb-0">Productos Más Vendidos</h6>
                </Card.Header>
                <Card.Body className="p-0">
                  <Table hover className="mb-0">
                    <thead>
                      <tr>
                        <th>Producto</th>
                        <th className="text-center">Cantidad</th>
                        <th className="text-end">Total</th>
                      </tr>
                    </thead>
                    <tbody>
                      {(reporte.top_productos || []).length === 0 ? (
                        <tr>
                          <td colSpan="3" className="text-center py-4 text-muted">
                            Sin datos en este período
                          </td>
                        </tr>
                      ) : (
                        (reporte.top_productos || []).slice(0, 10).map((prod, idx) => (
                          <tr key={idx}>
                            <td>
                              <div className="d-flex align-items-center">
                                <span className="badge bg-primary me-2">{idx + 1}</span>
                                {prod.nombre}
                              </div>
                            </td>
                            <td className="text-center">{prod.cantidad}</td>
                            <td className="text-end">{formatCurrency(prod.total)}</td>
                          </tr>
                        ))
                      )}
                    </tbody>
                  </Table>
                </Card.Body>
              </Card>
            </Col>
          </Row>

          {/* Resumen por Sede */}
          {(reporte.por_sede || []).length > 0 && (
            <Card className="mt-4" style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Header className="bg-transparent">
                <h6 className="mb-0">Resumen por Sede</h6>
              </Card.Header>
              <Card.Body className="p-0">
                <Table hover className="mb-0">
                  <thead>
                    <tr>
                      <th>Sede</th>
                      <th className="text-end">Ingresos</th>
                      <th className="text-end">Costos</th>
                      <th className="text-end">Ganancia</th>
                      <th className="text-end">Margen</th>
                    </tr>
                  </thead>
                  <tbody>
                    {(reporte.por_sede || []).map((sede, idx) => (
                      <tr key={idx}>
                        <td><strong>{sede.nombre}</strong></td>
                        <td className="text-end text-success">{formatCurrency(sede.ingresos)}</td>
                        <td className="text-end text-danger">{formatCurrency(sede.costos)}</td>
                        <td className="text-end text-primary">{formatCurrency(sede.ganancia)}</td>
                        <td className="text-end">{(sede.margen || 0).toFixed(1)}%</td>
                      </tr>
                    ))}
                  </tbody>
                </Table>
              </Card.Body>
            </Card>
          )}
        </>
      )}
    </Container>
  );
}

export default Reportes;
