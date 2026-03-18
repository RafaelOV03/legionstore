import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Form, Badge, Spinner, Alert, InputGroup, Button } from 'react-bootstrap';
import { Search, ClipboardData, Clock, Person, Activity } from 'react-bootstrap-icons';
import * as api from '../services/inventarioApi';

function Auditoria() {
  const [logs, setLogs] = useState([]);
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  
  const [filters, setFilters] = useState({
    accion: '',
    entidad: '',
    fecha_inicio: '',
    fecha_fin: ''
  });

  const acciones = ['crear', 'actualizar', 'eliminar', 'login', 'logout'];
  const entidades = ['producto', 'usuario', 'orden_trabajo', 'rma', 'cotizacion', 'traspaso', 'insumo', 'proveedor', 'deuda'];

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [logsData, statsData] = await Promise.all([
        api.getAuditoriaLogs(filters),
        api.getLogStats()
      ]);
      setLogs(logsData || []);
      setStats(statsData);
    } catch (err) {
      setError('Error cargando logs de auditoría');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleFilter = () => {
    loadData();
  };

  const clearFilters = () => {
    setFilters({
      accion: '',
      entidad: '',
      fecha_inicio: '',
      fecha_fin: ''
    });
    setTimeout(() => loadData(), 100);
  };

  const getAccionBadge = (accion) => {
    const colors = {
      crear: 'success',
      actualizar: 'info',
      eliminar: 'danger',
      login: 'primary',
      logout: 'secondary'
    };
    return <Badge bg={colors[accion] || 'secondary'}>{accion}</Badge>;
  };

  if (loading && !logs.length) {
    return (
      <Container className="py-5 text-center">
        <Spinner animation="border" variant="primary" />
      </Container>
    );
  }

  return (
    <Container fluid className="py-4">
      <Row className="mb-4">
        <Col>
          <h2><ClipboardData className="me-2" />Auditoría del Sistema</h2>
          <p className="text-muted">Registro de actividad del sistema</p>
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}

      {/* Estadísticas */}
      {stats && (
        <Row className="mb-4 g-3">
          <Col xs={6} md={3}>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-3">
                <Activity className="text-primary mb-2" size={24} />
                <h4 className="mb-0">{stats.total_logs || 0}</h4>
                <small className="text-muted">Total de Registros</small>
              </Card.Body>
            </Card>
          </Col>
          <Col xs={6} md={3}>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-3">
                <h4 className="mb-0 text-success">{stats.acciones?.crear || 0}</h4>
                <small className="text-muted">Creaciones</small>
              </Card.Body>
            </Card>
          </Col>
          <Col xs={6} md={3}>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-3">
                <h4 className="mb-0 text-info">{stats.acciones?.actualizar || 0}</h4>
                <small className="text-muted">Actualizaciones</small>
              </Card.Body>
            </Card>
          </Col>
          <Col xs={6} md={3}>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-3">
                <h4 className="mb-0 text-danger">{stats.acciones?.eliminar || 0}</h4>
                <small className="text-muted">Eliminaciones</small>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      )}

      {/* Filtros */}
      <Card className="mb-4" style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
        <Card.Body>
          <Row className="g-2 align-items-end">
            <Col md={2}>
              <Form.Group>
                <Form.Label className="small">Acción</Form.Label>
                <Form.Select
                  size="sm"
                  value={filters.accion}
                  onChange={(e) => setFilters({...filters, accion: e.target.value})}
                >
                  <option value="">Todas</option>
                  {acciones.map(a => (
                    <option key={a} value={a}>{a}</option>
                  ))}
                </Form.Select>
              </Form.Group>
            </Col>
            <Col md={2}>
              <Form.Group>
                <Form.Label className="small">Entidad</Form.Label>
                <Form.Select
                  size="sm"
                  value={filters.entidad}
                  onChange={(e) => setFilters({...filters, entidad: e.target.value})}
                >
                  <option value="">Todas</option>
                  {entidades.map(e => (
                    <option key={e} value={e}>{e}</option>
                  ))}
                </Form.Select>
              </Form.Group>
            </Col>
            <Col md={2}>
              <Form.Group>
                <Form.Label className="small">Desde</Form.Label>
                <Form.Control
                  type="date"
                  size="sm"
                  value={filters.fecha_inicio}
                  onChange={(e) => setFilters({...filters, fecha_inicio: e.target.value})}
                />
              </Form.Group>
            </Col>
            <Col md={2}>
              <Form.Group>
                <Form.Label className="small">Hasta</Form.Label>
                <Form.Control
                  type="date"
                  size="sm"
                  value={filters.fecha_fin}
                  onChange={(e) => setFilters({...filters, fecha_fin: e.target.value})}
                />
              </Form.Group>
            </Col>
            <Col md={4}>
              <Button variant="primary" size="sm" className="me-2" onClick={handleFilter}>
                <Search className="me-1" /> Filtrar
              </Button>
              <Button variant="outline-secondary" size="sm" onClick={clearFilters}>
                Limpiar
              </Button>
            </Col>
          </Row>
        </Card.Body>
      </Card>

      {/* Tabla de logs */}
      <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
        <Card.Body className="p-0">
          <Table responsive hover className="mb-0">
            <thead>
              <tr>
                <th style={{ width: 160 }}>Fecha/Hora</th>
                <th style={{ width: 120 }}>Usuario</th>
                <th style={{ width: 100 }}>Acción</th>
                <th style={{ width: 120 }}>Entidad</th>
                <th>Detalles</th>
                <th style={{ width: 120 }}>IP</th>
              </tr>
            </thead>
            <tbody>
              {logs.length === 0 ? (
                <tr>
                  <td colSpan="6" className="text-center py-4 text-muted">
                    No hay registros de auditoría
                  </td>
                </tr>
              ) : (
                logs.map(log => (
                  <tr key={log.id}>
                    <td>
                      <div className="d-flex align-items-center">
                        <Clock className="text-muted me-2" size={14} />
                        <div>
                          <div className="small">{new Date(log.created_at).toLocaleDateString()}</div>
                          <div className="small text-muted">{new Date(log.created_at).toLocaleTimeString()}</div>
                        </div>
                      </div>
                    </td>
                    <td>
                      <div className="d-flex align-items-center">
                        <Person className="text-muted me-1" size={14} />
                        <span>{log.usuario?.name || 'Sistema'}</span>
                      </div>
                    </td>
                    <td>{getAccionBadge(log.accion)}</td>
                    <td>
                      <Badge bg="secondary">{log.entidad}</Badge>
                      {log.entidad_id && (
                        <span className="ms-1 text-muted small">#{log.entidad_id}</span>
                      )}
                    </td>
                    <td>
                      <small className="text-muted">{log.descripcion || log.detalles || '-'}</small>
                    </td>
                    <td>
                      <small className="text-muted">{log.ip || '-'}</small>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </Table>
        </Card.Body>
        <Card.Footer className="bg-transparent text-muted">
          Mostrando {logs.length} registro(s)
        </Card.Footer>
      </Card>
    </Container>
  );
}

export default Auditoria;
