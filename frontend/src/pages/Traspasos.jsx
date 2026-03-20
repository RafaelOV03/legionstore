import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Button, Modal, Form, Badge, Spinner, Alert, InputGroup } from 'react-bootstrap';
import { Plus, Pencil, Trash, Search, Truck, ArrowRight, Check, X, Building } from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import * as api from '../services/resourceApi';

function Traspasos() {
  const { hasPermission, user } = useAuth();
  const [traspasos, setTraspasos] = useState([]);
  const [productos, setProductos] = useState([]);
  const [sedes, setSedes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [statusFilter, setStatusFilter] = useState('');
  
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    product_id: '',
    sede_origen_id: '',
    sede_destino_id: '',
    cantidad: 1,
    motivo: ''
  });

  const canCreate = hasPermission('traspasos.create');
  const canManage = hasPermission('traspasos.update');

  const estados = [
    { value: 'pendiente', label: 'Pendiente', color: 'warning' },
    { value: 'enviado', label: 'Enviado', color: 'info' },
    { value: 'recibido', label: 'Recibido', color: 'success' },
    { value: 'cancelado', label: 'Cancelado', color: 'danger' }
  ];

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [traspasosData, productosData, sedesData] = await Promise.all([
        api.getTraspasos(),
        api.getProducts(),
        api.getSedes()
      ]);
      setTraspasos(traspasosData || []);
      setProductos(productosData || []);
      setSedes(sedesData || []);
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
      // Backend expects items array with producto_id and cantidad
      const payload = {
        sede_origen_id: parseInt(formData.sede_origen_id),
        sede_destino_id: parseInt(formData.sede_destino_id),
        notas: formData.motivo || '',
        items: [{
          producto_id: parseInt(formData.product_id),
          cantidad: parseInt(formData.cantidad)
        }]
      };

      await api.createTraspaso(payload);
      
      setShowModal(false);
      resetForm();
      loadData();
      setSuccess('Traspaso creado exitosamente');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error creando traspaso');
      console.error(err);
    }
  };

  const handleEnviar = async (id) => {
    try {
      await api.enviarTraspaso(id);
      loadData();
      setSuccess('Traspaso enviado');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error enviando traspaso');
      console.error(err);
    }
  };

  const handleRecibir = async (id) => {
    try {
      await api.recibirTraspaso(id, { usuario_destino_id: user?.id });
      loadData();
      setSuccess('Traspaso recibido exitosamente');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error recibiendo traspaso');
      console.error(err);
    }
  };

  const handleCancelar = async (id) => {
    if (!window.confirm('¿Cancelar este traspaso?')) return;
    try {
      await api.cancelarTraspaso(id);
      loadData();
      setSuccess('Traspaso cancelado');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error cancelando traspaso');
      console.error(err);
    }
  };

  const resetForm = () => {
    setFormData({
      product_id: '',
      sede_origen_id: user?.sede_id?.toString() || '',
      sede_destino_id: '',
      cantidad: 1,
      motivo: ''
    });
  };

  const getEstadoBadge = (estado) => {
    const est = estados.find(e => e.value === estado);
    return <Badge bg={est?.color || 'secondary'}>{est?.label || estado}</Badge>;
  };

  const filteredTraspasos = traspasos.filter(t => 
    !statusFilter || t.estado === statusFilter
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
          <h2><Truck className="me-2" />Traspasos</h2>
          <p className="text-muted">Movimiento de inventario entre sedes</p>
        </Col>
        <Col xs="auto">
          {canCreate && (
            <Button variant="primary" onClick={() => { resetForm(); setShowModal(true); }}>
              <Plus className="me-1" /> Nuevo Traspaso
            </Button>
          )}
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}
      {success && <Alert variant="success" onClose={() => setSuccess(null)} dismissible>{success}</Alert>}

      {/* Resumen por estado */}
      <Row className="mb-4 g-3">
        {estados.map(est => (
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
                  {traspasos.filter(t => t.estado === est.value).length}
                </h3>
                <small className="text-muted">{est.label}</small>
              </Card.Body>
            </Card>
          </Col>
        ))}
      </Row>

      <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
        <Card.Body className="p-0">
          <Table responsive hover className="mb-0">
            <thead>
              <tr>
                <th>#</th>
                <th>Producto</th>
                <th className="text-center">Cantidad</th>
                <th>Origen → Destino</th>
                <th>Estado</th>
                <th>Fecha</th>
                <th className="text-center">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {filteredTraspasos.length === 0 ? (
                <tr>
                  <td colSpan="7" className="text-center py-4 text-muted">
                    No hay traspasos para mostrar
                  </td>
                </tr>
              ) : (
                filteredTraspasos.map(traspaso => (
                  <tr key={traspaso.id}>
                    <td><strong>#{traspaso.id}</strong></td>
                    <td>{traspaso.producto?.name || 'N/A'}</td>
                    <td className="text-center">
                      <Badge bg="primary">{traspaso.cantidad}</Badge>
                    </td>
                    <td>
                      <div className="d-flex align-items-center">
                        <Building className="text-muted me-1" size={14} />
                        <span>{traspaso.sede_origen?.nombre || 'N/A'}</span>
                        <ArrowRight className="mx-2 text-primary" />
                        <Building className="text-muted me-1" size={14} />
                        <span>{traspaso.sede_destino?.nombre || 'N/A'}</span>
                      </div>
                    </td>
                    <td>{getEstadoBadge(traspaso.estado)}</td>
                    <td>{new Date(traspaso.created_at).toLocaleDateString()}</td>
                    <td className="text-center">
                      {canManage && traspaso.estado === 'pendiente' && (
                        <>
                          <Button 
                            variant="outline-info" 
                            size="sm" 
                            className="me-1"
                            onClick={() => handleEnviar(traspaso.id)}
                            title="Enviar"
                          >
                            <Truck />
                          </Button>
                          <Button 
                            variant="outline-danger" 
                            size="sm"
                            onClick={() => handleCancelar(traspaso.id)}
                            title="Cancelar"
                          >
                            <X />
                          </Button>
                        </>
                      )}
                      {canManage && traspaso.estado === 'enviado' && (
                        <Button 
                          variant="outline-success" 
                          size="sm"
                          onClick={() => handleRecibir(traspaso.id)}
                          title="Recibir"
                        >
                          <Check />
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
          Total: {filteredTraspasos.length} traspaso(s)
        </Card.Footer>
      </Card>

      {/* Modal de creación */}
      <Modal show={showModal} onHide={() => setShowModal(false)} size="lg">
        <Form onSubmit={handleSubmit}>
          <Modal.Header closeButton>
            <Modal.Title>Nuevo Traspaso</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Form.Group className="mb-3">
              <Form.Label>Producto *</Form.Label>
              <Form.Select
                required
                value={formData.product_id}
                onChange={(e) => setFormData({...formData, product_id: e.target.value})}
              >
                <option value="">Seleccionar producto</option>
                {productos.map(p => (
                  <option key={p.id} value={p.id}>{p.name} (Stock: {p.stock})</option>
                ))}
              </Form.Select>
            </Form.Group>
            <Row>
              <Col md={5}>
                <Form.Group className="mb-3">
                  <Form.Label>Sede Origen *</Form.Label>
                  <Form.Select
                    required
                    value={formData.sede_origen_id}
                    onChange={(e) => setFormData({...formData, sede_origen_id: e.target.value})}
                  >
                    <option value="">Seleccionar sede</option>
                    {sedes.map(s => (
                      <option key={s.id} value={s.id}>{s.nombre}</option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
              <Col md={2} className="d-flex align-items-center justify-content-center">
                <ArrowRight size={24} className="text-primary mt-3" />
              </Col>
              <Col md={5}>
                <Form.Group className="mb-3">
                  <Form.Label>Sede Destino *</Form.Label>
                  <Form.Select
                    required
                    value={formData.sede_destino_id}
                    onChange={(e) => setFormData({...formData, sede_destino_id: e.target.value})}
                  >
                    <option value="">Seleccionar sede</option>
                    {sedes.filter(s => s.id.toString() !== formData.sede_origen_id).map(s => (
                      <option key={s.id} value={s.id}>{s.nombre}</option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
            </Row>
            <Form.Group className="mb-3">
              <Form.Label>Cantidad *</Form.Label>
              <Form.Control
                type="number"
                min="1"
                required
                value={formData.cantidad}
                onChange={(e) => setFormData({...formData, cantidad: e.target.value})}
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Motivo</Form.Label>
              <Form.Control
                as="textarea"
                rows={2}
                value={formData.motivo}
                onChange={(e) => setFormData({...formData, motivo: e.target.value})}
                placeholder="Motivo del traspaso"
              />
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowModal(false)}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              Crear Traspaso
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

export default Traspasos;
