import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Button, Modal, Form, Badge, Spinner, Alert, InputGroup, Tabs, Tab } from 'react-bootstrap';
import { Plus, Pencil, Trash, Search, ClipboardCheck, Clock, CheckCircle, XCircle, ArrowRepeat } from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import * as api from '../services/inventarioApi';

function RMA() {
  const { hasPermission, user } = useAuth();
  const [rmas, setRmas] = useState([]);
  const [productos, setProductos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [stats, setStats] = useState(null);
  
  const [showModal, setShowModal] = useState(false);
  const [editingRMA, setEditingRMA] = useState(null);
  const [formData, setFormData] = useState({
    producto_id: '',
    cliente_nombre: '',
    cliente_telefono: '',
    cliente_email: '',
    num_serie: '',
    fecha_compra: '',
    motivo_devolucion: '',
    sede_id: '',
    notas: ''
  });

  const canCreate = hasPermission('rmas.create');
  const canEdit = hasPermission('rmas.update');
  const canDelete = hasPermission('rmas.delete');

  const estados = [
    { value: 'recibido', label: 'Recibido', color: 'secondary', icon: Clock },
    { value: 'en_revision', label: 'En Revisión', color: 'warning', icon: ArrowRepeat },
    { value: 'aprobado', label: 'Aprobado', color: 'success', icon: CheckCircle },
    { value: 'rechazado', label: 'Rechazado', color: 'danger', icon: XCircle },
    { value: 'reemplazado', label: 'Reemplazado', color: 'info', icon: ArrowRepeat },
    { value: 'reembolsado', label: 'Reembolsado', color: 'primary', icon: CheckCircle }
  ];

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [rmasData, productosData, statsData] = await Promise.all([
        api.getRMAs(),
        api.getProducts(),
        api.getRMAStats()
      ]);
      setRmas(rmasData || []);
      setProductos(productosData || []);
      setStats(statsData);
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
        producto_id: parseInt(formData.producto_id),
        cliente_nombre: formData.cliente_nombre,
        cliente_telefono: formData.cliente_telefono,
        cliente_email: formData.cliente_email,
        num_serie: formData.num_serie,
        fecha_compra: formData.fecha_compra,
        motivo_devolucion: formData.motivo_devolucion,
        sede_id: parseInt(formData.sede_id) || user?.sede_id || 1,
        notas: formData.notas
      };

      if (editingRMA) {
        await api.updateRMA(editingRMA.id, payload);
      } else {
        await api.createRMA(payload);
      }
      
      setShowModal(false);
      resetForm();
      loadData();
    } catch (err) {
      setError('Error guardando RMA');
      console.error(err);
    }
  };

  const handleEdit = (rma) => {
    setEditingRMA(rma);
    setFormData({
      producto_id: rma.producto_id?.toString() || '',
      cliente_nombre: rma.cliente_nombre || '',
      cliente_telefono: rma.cliente_telefono || '',
      cliente_email: rma.cliente_email || '',
      num_serie: rma.num_serie || '',
      fecha_compra: rma.fecha_compra ? rma.fecha_compra.split('T')[0] : '',
      motivo_devolucion: rma.motivo_devolucion || '',
      sede_id: rma.sede_id?.toString() || '',
      notas: rma.notas || ''
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('¿Está seguro de eliminar este RMA?')) return;
    try {
      await api.deleteRMA(id);
      loadData();
    } catch (err) {
      setError('Error eliminando RMA');
      console.error(err);
    }
  };

  const resetForm = () => {
    setEditingRMA(null);
    setFormData({
      producto_id: '',
      cliente_nombre: '',
      cliente_telefono: '',
      cliente_email: '',
      num_serie: '',
      fecha_compra: '',
      motivo_devolucion: '',
      sede_id: user?.sede_id?.toString() || '',
      notas: ''
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

  const filteredRMAs = rmas.filter(rma => {
    const matchSearch = 
      rma.producto?.name?.toLowerCase().includes(search.toLowerCase()) ||
      rma.motivo?.toLowerCase().includes(search.toLowerCase());
    const matchStatus = !statusFilter || rma.estado === statusFilter;
    return matchSearch && matchStatus;
  });

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
          <h2><ClipboardCheck className="me-2" />RMA / Garantías</h2>
          <p className="text-muted">Gestión de devoluciones y garantías</p>
        </Col>
        <Col xs="auto">
          {canCreate && (
            <Button variant="primary" onClick={() => { resetForm(); setShowModal(true); }}>
              <Plus className="me-1" /> Nuevo RMA
            </Button>
          )}
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}

      {/* Estadísticas */}
      {stats && (
        <Row className="mb-4 g-3">
          {estados.slice(0, 4).map(est => (
            <Col xs={6} md={3} key={est.value}>
              <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
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

      <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
        <Card.Header className="bg-transparent">
          <Row className="g-2">
            <Col md={6}>
              <InputGroup>
                <InputGroup.Text><Search /></InputGroup.Text>
                <Form.Control
                  placeholder="Buscar por producto o motivo..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
              </InputGroup>
            </Col>
            <Col md={6}>
              <Form.Select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
              >
                <option value="">Todos los estados</option>
                {estados.map(est => (
                  <option key={est.value} value={est.value}>{est.label}</option>
                ))}
              </Form.Select>
            </Col>
          </Row>
        </Card.Header>
        <Card.Body className="p-0">
          <Table responsive hover className="mb-0">
            <thead>
              <tr>
                <th>#</th>
                <th>Producto</th>
                <th className="text-center">Cantidad</th>
                <th>Motivo</th>
                <th>Estado</th>
                <th>Fecha</th>
                <th className="text-center">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {filteredRMAs.length === 0 ? (
                <tr>
                  <td colSpan="7" className="text-center py-4 text-muted">
                    No hay RMAs para mostrar
                  </td>
                </tr>
              ) : (
                filteredRMAs.map(rma => (
                  <tr key={rma.id}>
                    <td><strong>#{rma.id}</strong></td>
                    <td>{rma.producto?.name || 'N/A'}</td>
                    <td className="text-center">{rma.cantidad}</td>
                    <td>{rma.motivo}</td>
                    <td>{getEstadoBadge(rma.estado)}</td>
                    <td>{new Date(rma.created_at).toLocaleDateString()}</td>
                    <td className="text-center">
                      {canEdit && (
                        <Button 
                          variant="outline-primary" 
                          size="sm" 
                          className="me-1"
                          onClick={() => handleEdit(rma)}
                        >
                          <Pencil />
                        </Button>
                      )}
                      {canDelete && (
                        <Button 
                          variant="outline-danger" 
                          size="sm"
                          onClick={() => handleDelete(rma.id)}
                        >
                          <Trash />
                        </Button>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </Table>
        </Card.Body>
        <Card.Footer className="bg-transparent text-muted">
          Total: {filteredRMAs.length} RMA(s)
        </Card.Footer>
      </Card>

      {/* Modal de creación/edición */}
      <Modal show={showModal} onHide={() => setShowModal(false)} size="lg">
        <Form onSubmit={handleSubmit}>
          <Modal.Header closeButton>
            <Modal.Title>
              {editingRMA ? 'Editar RMA' : 'Nuevo RMA / Garantía'}
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Producto *</Form.Label>
                  <Form.Select
                    required
                    value={formData.producto_id}
                    onChange={(e) => setFormData({...formData, producto_id: e.target.value})}
                  >
                    <option value="">Seleccionar producto</option>
                    {productos.map(p => (
                      <option key={p.id} value={p.id}>{p.name}</option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Cliente Nombre *</Form.Label>
                  <Form.Control
                    required
                    value={formData.cliente_nombre}
                    onChange={(e) => setFormData({...formData, cliente_nombre: e.target.value})}
                    placeholder="Nombre del cliente"
                  />
                </Form.Group>
              </Col>
            </Row>
            <Row>
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
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Email</Form.Label>
                  <Form.Control
                    type="email"
                    value={formData.cliente_email}
                    onChange={(e) => setFormData({...formData, cliente_email: e.target.value})}
                    placeholder="Email del cliente"
                  />
                </Form.Group>
              </Col>
            </Row>
            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Número de Serie</Form.Label>
                  <Form.Control
                    value={formData.num_serie}
                    onChange={(e) => setFormData({...formData, num_serie: e.target.value})}
                    placeholder="Serie del equipo"
                  />
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Fecha de Compra</Form.Label>
                  <Form.Control
                    type="date"
                    value={formData.fecha_compra}
                    onChange={(e) => setFormData({...formData, fecha_compra: e.target.value})}
                  />
                </Form.Group>
              </Col>
            </Row>
            <Form.Group className="mb-3">
              <Form.Label>Motivo de Devolución *</Form.Label>
              <Form.Control
                as="textarea"
                rows={2}
                required
                value={formData.motivo_devolucion}
                onChange={(e) => setFormData({...formData, motivo_devolucion: e.target.value})}
                placeholder="Describir el motivo de la devolución o problema"
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Notas Adicionales</Form.Label>
              <Form.Control
                as="textarea"
                rows={2}
                value={formData.notas}
                onChange={(e) => setFormData({...formData, notas: e.target.value})}
                placeholder="Notas adicionales"
              />
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowModal(false)}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              {editingRMA ? 'Guardar Cambios' : 'Crear RMA'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

export default RMA;
