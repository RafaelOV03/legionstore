import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Button, Modal, Form, Badge, Spinner, Alert, InputGroup, Tabs, Tab } from 'react-bootstrap';
import { Plus, Pencil, Trash, Search, ClipboardCheck, Person, Tools, Clock, CheckCircle, ExclamationTriangle } from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import * as api from '../services/resourceApi';

function OrdenesTrabajoPage() {
  const { hasPermission, user } = useAuth();
  const [ordenes, setOrdenes] = useState([]);
  const [tecnicos, setTecnicos] = useState([]);
  const [insumos, setInsumos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [statusFilter, setStatusFilter] = useState('');
  const [stats, setStats] = useState(null);
  
  const [showModal, setShowModal] = useState(false);
  const [showInsumoModal, setShowInsumoModal] = useState(false);
  const [editingOrden, setEditingOrden] = useState(null);
  const [selectedOrden, setSelectedOrden] = useState(null);
  const [formData, setFormData] = useState({
    cliente_nombre: '',
    cliente_telefono: '',
    equipo_tipo: '',
    equipo_marca: '',
    equipo_modelo: '',
    num_serie: '',
    problema_reportado: '',
    prioridad: 'media',
    fecha_estimada: ''
  });
  const [insumoData, setInsumoData] = useState({ insumo_id: '', cantidad: 1 });

  const canCreate = hasPermission('ordenes.create');
  const canEdit = hasPermission('ordenes.update');

  const estados = [
    { value: 'recibido', label: 'Recibido', color: 'secondary', icon: Clock },
    { value: 'en_diagnostico', label: 'En Diagnóstico', color: 'info', icon: Tools },
    { value: 'en_reparacion', label: 'En Reparación', color: 'warning', icon: Tools },
    { value: 'listo', label: 'Listo', color: 'success', icon: CheckCircle },
    { value: 'entregado', label: 'Entregado', color: 'primary', icon: CheckCircle }
  ];

  const prioridades = [
    { value: 'baja', label: 'Baja', color: 'secondary' },
    { value: 'media', label: 'Media', color: 'info' },
    { value: 'alta', label: 'Alta', color: 'warning' },
    { value: 'urgente', label: 'Urgente', color: 'danger' }
  ];

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [ordenesData, insumosData, statsData, tecnicosData] = await Promise.all([
        api.getOrdenesTrabajo(),
        api.getInsumos(),
        api.getOrdenesStats(),
        api.getTecnicos()
      ]);
      setOrdenes(ordenesData || []);
      setInsumos(insumosData || []);
      setStats(statsData);
      setTecnicos(tecnicosData || []);
    } catch (err) {
      setError('Error cargando datos');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        ...formData,
        usuario_recibe_id: user?.id,
        sede_id: user?.sede_id || 1
      };

      if (editingOrden) {
        await api.updateOrdenTrabajo(editingOrden.id, payload);
      } else {
        await api.createOrdenTrabajo(payload);
      }
      
      setShowModal(false);
      resetForm();
      loadData();
      setSuccess('Orden guardada exitosamente');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error guardando orden');
      console.error(err);
    }
  };

  const handleAsignarTecnico = async (ordenId, tecnicoId) => {
    try {
      await api.asignarTecnico(ordenId, { tecnico_id: parseInt(tecnicoId) });
      loadData();
      setSuccess('Técnico asignado');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error asignando técnico');
      console.error(err);
    }
  };

  const handleCambiarEstado = async (ordenId, nuevoEstado) => {
    try {
      await api.updateOrdenTrabajo(ordenId, { estado: nuevoEstado });
      loadData();
    } catch (err) {
      setError('Error cambiando estado');
      console.error(err);
    }
  };

  const handleAgregarInsumo = async (e) => {
    e.preventDefault();
    try {
      await api.agregarInsumoOrden(selectedOrden.id, {
        insumo_id: parseInt(insumoData.insumo_id),
        cantidad: parseInt(insumoData.cantidad)
      });
      setShowInsumoModal(false);
      setInsumoData({ insumo_id: '', cantidad: 1 });
      loadData();
      setSuccess('Insumo agregado');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error agregando insumo');
      console.error(err);
    }
  };

  const handleEdit = (orden) => {
    setEditingOrden(orden);
    setFormData({
      cliente_nombre: orden.cliente_nombre || '',
      cliente_telefono: orden.cliente_telefono || '',
      equipo_tipo: orden.equipo_tipo || '',
      equipo_marca: orden.equipo_marca || '',
      equipo_modelo: orden.equipo_modelo || '',
      num_serie: orden.num_serie || '',
      problema_reportado: orden.problema_reportado || '',
      prioridad: orden.prioridad || 'media',
      fecha_estimada: orden.fecha_estimada?.split('T')[0] || ''
    });
    setShowModal(true);
  };

  const resetForm = () => {
    setEditingOrden(null);
    setFormData({
      cliente_nombre: '',
      cliente_telefono: '',
      equipo_tipo: '',
      equipo_marca: '',
      equipo_modelo: '',
      num_serie: '',
      problema_reportado: '',
      prioridad: 'media',
      fecha_estimada: ''
    });
  };

  const getEstadoBadge = (estado) => {
    const est = estados.find(e => e.value === estado);
    if (!est) return <Badge bg="secondary">{estado}</Badge>;
    const Icon = est.icon;
    return (
      <Badge bg={est.color}>
        <Icon className="me-1" size={12} />
        {est.label}
      </Badge>
    );
  };

  const getPrioridadBadge = (prioridad) => {
    const prio = prioridades.find(p => p.value === prioridad);
    return <Badge bg={prio?.color || 'secondary'}>{prio?.label || prioridad}</Badge>;
  };

  const filteredOrdenes = ordenes.filter(o => 
    !statusFilter || o.estado === statusFilter
  );

  if (loading) {
    return (
      <Container className="py-5 text-center">
        <Spinner animation="border" variant="primary" />
      </Container>
    );
  }

  return (
    <Container fluid className="py-4">
      <Row className="mb-4 align-items-center">
        <Col>
          <h2><ClipboardCheck className="me-2" />Órdenes de Trabajo</h2>
          <p className="text-muted">Gestión de servicio técnico</p>
        </Col>
        <Col xs="auto">
          {canCreate && (
            <Button variant="primary" onClick={() => { resetForm(); setShowModal(true); }}>
              <Plus className="me-1" /> Nueva Orden
            </Button>
          )}
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}
      {success && <Alert variant="success" onClose={() => setSuccess(null)} dismissible>{success}</Alert>}

      {/* Estadísticas */}
      {stats && (
        <Row className="mb-4 g-3">
          {estados.slice(0, 4).map(est => (
            <Col xs={6} md={3} key={est.value}>
              <Card 
                style={{ 
                  background: 'var(--card-bg)', 
                  border: `1px solid ${statusFilter === est.value ? `var(--bs-${est.color})` : 'var(--border-color)'}`,
                  cursor: 'pointer'
                }}
                onClick={() => setStatusFilter(statusFilter === est.value ? '' : est.value)}
              >
                <Card.Body className="text-center py-3">
                  <h3 className="mb-0" style={{ color: `var(--bs-${est.color})` }}>
                    {stats[est.value] || 0}
                  </h3>
                  <small className="text-muted">{est.label}</small>
                </Card.Body>
              </Card>
            </Col>
          ))}
        </Row>
      )}

      {stats?.urgentes > 0 && (
        <Alert variant="warning" className="mb-4">
          <ExclamationTriangle className="me-2" />
          <strong>{stats.urgentes}</strong> orden(es) marcada(s) como urgente(s)
        </Alert>
      )}

      <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
        <Card.Body className="p-0">
          <Table responsive hover className="mb-0">
            <thead>
              <tr>
                <th>Orden</th>
                <th>Cliente</th>
                <th>Equipo</th>
                <th>Prioridad</th>
                <th>Estado</th>
                <th>Técnico</th>
                <th>Fecha Promesa</th>
                <th className="text-center">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {filteredOrdenes.length === 0 ? (
                <tr>
                  <td colSpan="8" className="text-center py-4 text-muted">
                    No hay órdenes para mostrar
                  </td>
                </tr>
              ) : (
                filteredOrdenes.map(orden => (
                  <tr key={orden.id}>
                    <td>
                      <strong>#{orden.numero || orden.id}</strong>
                      <div className="small text-muted">
                        {new Date(orden.created_at).toLocaleDateString()}
                      </div>
                    </td>
                    <td>
                      <div>{orden.cliente_nombre}</div>
                      {orden.cliente_telefono && (
                        <small className="text-muted">{orden.cliente_telefono}</small>
                      )}
                    </td>
                    <td>
                      <div>{orden.equipo_tipo} {orden.equipo_marca}</div>
                      <small className="text-muted">{orden.equipo_modelo}</small>
                    </td>
                    <td>{getPrioridadBadge(orden.prioridad)}</td>
                    <td>
                      {canEdit ? (
                        <Form.Select
                          size="sm"
                          value={orden.estado}
                          onChange={(e) => handleCambiarEstado(orden.id, e.target.value)}
                          style={{ width: 140 }}
                        >
                          {estados.map(est => (
                            <option key={est.value} value={est.value}>{est.label}</option>
                          ))}
                        </Form.Select>
                      ) : (
                        getEstadoBadge(orden.estado)
                      )}
                    </td>
                    <td>
                      {canEdit ? (
                        <Form.Select
                          size="sm"
                          value={orden.tecnico_id || ''}
                          onChange={(e) => handleAsignarTecnico(orden.id, e.target.value)}
                          style={{ width: 140 }}
                        >
                          <option value="">Sin asignar</option>
                          {tecnicos.map(t => (
                            <option key={t.id} value={t.id}>{t.name}</option>
                          ))}
                        </Form.Select>
                      ) : (
                        <span>{orden.tecnico?.name || 'Sin asignar'}</span>
                      )}
                    </td>
                    <td>
                      {orden.fecha_estimada 
                        ? new Date(orden.fecha_estimada).toLocaleDateString()
                        : '-'
                      }
                    </td>
                    <td className="text-center">
                      {canEdit && (
                        <>
                          <Button 
                            variant="outline-success" 
                            size="sm" 
                            className="me-1"
                            onClick={() => { setSelectedOrden(orden); setShowInsumoModal(true); }}
                            title="Agregar insumo"
                          >
                            <Tools />
                          </Button>
                          <Button 
                            variant="outline-primary" 
                            size="sm"
                            onClick={() => handleEdit(orden)}
                            title="Editar"
                          >
                            <Pencil />
                          </Button>
                        </>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </Table>
        </Card.Body>
        <Card.Footer className="bg-transparent text-muted">
          Total: {filteredOrdenes.length} orden(es)
        </Card.Footer>
      </Card>

      {/* Modal de creación/edición */}
      <Modal show={showModal} onHide={() => setShowModal(false)} size="lg">
        <Form onSubmit={handleSubmit}>
          <Modal.Header closeButton>
            <Modal.Title>
              {editingOrden ? 'Editar Orden' : 'Nueva Orden de Trabajo'}
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <h6 className="text-muted mb-3">Datos del Cliente</h6>
            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Nombre del Cliente *</Form.Label>
                  <Form.Control
                    required
                    value={formData.cliente_nombre}
                    onChange={(e) => setFormData({...formData, cliente_nombre: e.target.value})}
                    placeholder="Nombre completo"
                  />
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Teléfono</Form.Label>
                  <Form.Control
                    value={formData.cliente_telefono}
                    onChange={(e) => setFormData({...formData, cliente_telefono: e.target.value})}
                    placeholder="Teléfono de contacto"
                  />
                </Form.Group>
              </Col>
            </Row>

            <h6 className="text-muted mb-3 mt-3">Datos del Equipo</h6>
            <Row>
              <Col md={3}>
                <Form.Group className="mb-3">
                  <Form.Label>Tipo *</Form.Label>
                  <Form.Select
                    required
                    value={formData.equipo_tipo}
                    onChange={(e) => setFormData({...formData, equipo_tipo: e.target.value})}
                  >
                    <option value="">Seleccionar...</option>
                    <option value="laptop">Laptop</option>
                    <option value="pc">PC de Escritorio</option>
                    <option value="celular">Celular</option>
                    <option value="tablet">Tablet</option>
                    <option value="impresora">Impresora</option>
                    <option value="otro">Otro</option>
                  </Form.Select>
                </Form.Group>
              </Col>
              <Col md={3}>
                <Form.Group className="mb-3">
                  <Form.Label>Marca</Form.Label>
                  <Form.Control
                    value={formData.equipo_marca}
                    onChange={(e) => setFormData({...formData, equipo_marca: e.target.value})}
                    placeholder="Marca"
                  />
                </Form.Group>
              </Col>
              <Col md={3}>
                <Form.Group className="mb-3">
                  <Form.Label>Modelo</Form.Label>
                  <Form.Control
                    value={formData.equipo_modelo}
                    onChange={(e) => setFormData({...formData, equipo_modelo: e.target.value})}
                    placeholder="Modelo"
                  />
                </Form.Group>
              </Col>
              <Col md={3}>
                <Form.Group className="mb-3">
                  <Form.Label>N° Serie</Form.Label>
                  <Form.Control
                    value={formData.num_serie}
                    onChange={(e) => setFormData({...formData, num_serie: e.target.value})}
                    placeholder="Número de serie"
                  />
                </Form.Group>
              </Col>
            </Row>

            <Form.Group className="mb-3">
              <Form.Label>Problema Reportado *</Form.Label>
              <Form.Control
                as="textarea"
                rows={3}
                required
                value={formData.problema_reportado}
                onChange={(e) => setFormData({...formData, problema_reportado: e.target.value})}
                placeholder="Describir el problema que reporta el cliente"
              />
            </Form.Group>

            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Prioridad</Form.Label>
                  <Form.Select
                    value={formData.prioridad}
                    onChange={(e) => setFormData({...formData, prioridad: e.target.value})}
                  >
                    {prioridades.map(p => (
                      <option key={p.value} value={p.value}>{p.label}</option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Fecha Estimada Entrega</Form.Label>
                  <Form.Control
                    type="date"
                    value={formData.fecha_estimada}
                    onChange={(e) => setFormData({...formData, fecha_estimada: e.target.value})}
                  />
                </Form.Group>
              </Col>
            </Row>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowModal(false)}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              {editingOrden ? 'Guardar Cambios' : 'Crear Orden'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>

      {/* Modal de agregar insumo */}
      <Modal show={showInsumoModal} onHide={() => setShowInsumoModal(false)}>
        <Form onSubmit={handleAgregarInsumo}>
          <Modal.Header closeButton>
            <Modal.Title>Agregar Insumo a Orden #{selectedOrden?.id}</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Form.Group className="mb-3">
              <Form.Label>Insumo *</Form.Label>
              <Form.Select
                required
                value={insumoData.insumo_id}
                onChange={(e) => setInsumoData({...insumoData, insumo_id: e.target.value})}
              >
                <option value="">Seleccionar insumo</option>
                {insumos.map(i => (
                  <option key={i.id} value={i.id}>
                    {i.nombre} (Stock: {i.stock})
                  </option>
                ))}
              </Form.Select>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Cantidad</Form.Label>
              <Form.Control
                type="number"
                min="1"
                value={insumoData.cantidad}
                onChange={(e) => setInsumoData({...insumoData, cantidad: e.target.value})}
              />
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowInsumoModal(false)}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              Agregar Insumo
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

export default OrdenesTrabajoPage;
