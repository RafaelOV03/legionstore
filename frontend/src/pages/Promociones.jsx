import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Button, Modal, Form, Badge, Spinner, Alert, InputGroup } from 'react-bootstrap';
import { Plus, Pencil, Trash, Search, Tag, Calendar, Percent } from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import * as api from '../services/resourceApi';

function Promociones() {
  const { hasPermission } = useAuth();
  const [promociones, setPromociones] = useState([]);
  const [productos, setProductos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  
  const [showModal, setShowModal] = useState(false);
  const [editingPromocion, setEditingPromocion] = useState(null);
  const [formData, setFormData] = useState({
    nombre: '',
    descripcion: '',
    tipo: 'descuento_porcentaje',
    valor: '',
    fecha_inicio: '',
    fecha_fin: '',
    producto_ids: [],
    categoria: '',
    codigo: '',
    activa: true
  });

  const canCreate = hasPermission('promociones.create');
  const canEdit = hasPermission('promociones.update');
  const canDelete = hasPermission('promociones.delete');

  const tipos = [
    { value: 'descuento_porcentaje', label: 'Descuento %' },
    { value: 'descuento_fijo', label: 'Descuento Fijo Bs.' },
    { value: '2x1', label: '2x1' },
    { value: '3x2', label: '3x2' }
  ];

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [promocionesData, productosData] = await Promise.all([
        api.getPromociones(),
        api.getProducts()
      ]);
      setPromociones(promocionesData || []);
      setProductos(productosData || []);
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
        valor: parseFloat(formData.valor) || 0
      };

      if (editingPromocion) {
        await api.updatePromocion(editingPromocion.id, payload);
      } else {
        await api.createPromocion(payload);
      }
      
      setShowModal(false);
      resetForm();
      loadData();
      setSuccess('Promoción guardada exitosamente');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error guardando promoción');
      console.error(err);
    }
  };

  const handleEdit = (promocion) => {
    setEditingPromocion(promocion);
    setFormData({
      nombre: promocion.nombre || '',
      descripcion: promocion.descripcion || '',
      tipo: promocion.tipo || 'descuento_porcentaje',
      valor: promocion.valor?.toString() || '',
      fecha_inicio: promocion.fecha_inicio?.split('T')[0] || '',
      fecha_fin: promocion.fecha_fin?.split('T')[0] || '',
      producto_ids: promocion.producto_ids || [],
      categoria: promocion.categoria || '',
      codigo: promocion.codigo || '',
      activa: promocion.activa !== false
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('¿Está seguro de eliminar esta promoción?')) return;
    try {
      await api.deletePromocion(id);
      loadData();
    } catch (err) {
      setError('Error eliminando promoción');
      console.error(err);
    }
  };

  const resetForm = () => {
    setEditingPromocion(null);
    setFormData({
      nombre: '',
      descripcion: '',
      tipo: 'descuento_porcentaje',
      valor: '',
      fecha_inicio: '',
      fecha_fin: '',
      producto_ids: [],
      categoria: '',
      codigo: '',
      activa: true
    });
  };

  const isVigente = (promocion) => {
    const now = new Date();
    const inicio = promocion.fecha_inicio ? new Date(promocion.fecha_inicio) : null;
    const fin = promocion.fecha_fin ? new Date(promocion.fecha_fin) : null;
    
    if (!promocion.activa) return false;
    if (inicio && inicio > now) return false;
    if (fin && fin < now) return false;
    return true;
  };

  const getValorDisplay = (promocion) => {
    switch (promocion.tipo) {
      case 'descuento_porcentaje':
        return `${promocion.valor}%`;
      case 'descuento_fijo':
        return `Bs. ${promocion.valor}`;
      case '2x1':
        return '2x1';
      case '3x2':
        return '3x2';
      default:
        return promocion.valor;
    }
  };

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
          <h2><Tag className="me-2" />Promociones</h2>
          <p className="text-muted">Gestión de promociones y descuentos</p>
        </Col>
        <Col xs="auto">
          {canCreate && (
            <Button variant="primary" onClick={() => { resetForm(); setShowModal(true); }}>
              <Plus className="me-1" /> Nueva Promoción
            </Button>
          )}
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}
      {success && <Alert variant="success" onClose={() => setSuccess(null)} dismissible>{success}</Alert>}

      {/* Resumen */}
      <Row className="mb-4 g-3">
        <Col xs={6} md={3}>
          <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
            <Card.Body className="text-center py-3">
              <h3 className="mb-0 text-primary">{promociones.length}</h3>
              <small className="text-muted">Total Promociones</small>
            </Card.Body>
          </Card>
        </Col>
        <Col xs={6} md={3}>
          <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
            <Card.Body className="text-center py-3">
              <h3 className="mb-0 text-success">{promociones.filter(isVigente).length}</h3>
              <small className="text-muted">Vigentes</small>
            </Card.Body>
          </Card>
        </Col>
        <Col xs={6} md={3}>
          <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
            <Card.Body className="text-center py-3">
              <h3 className="mb-0 text-warning">
                {promociones.filter(p => p.activa && !isVigente(p)).length}
              </h3>
              <small className="text-muted">Programadas</small>
            </Card.Body>
          </Card>
        </Col>
        <Col xs={6} md={3}>
          <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
            <Card.Body className="text-center py-3">
              <h3 className="mb-0 text-secondary">{promociones.filter(p => !p.activa).length}</h3>
              <small className="text-muted">Inactivas</small>
            </Card.Body>
          </Card>
        </Col>
      </Row>

      <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
        <Card.Body className="p-0">
          <Table responsive hover className="mb-0">
            <thead>
              <tr>
                <th>Promoción</th>
                <th>Tipo</th>
                <th className="text-center">Descuento</th>
                <th>Vigencia</th>
                <th>Código</th>
                <th>Estado</th>
                <th className="text-center">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {promociones.length === 0 ? (
                <tr>
                  <td colSpan="7" className="text-center py-4 text-muted">
                    No hay promociones para mostrar
                  </td>
                </tr>
              ) : (
                promociones.map(promocion => (
                  <tr key={promocion.id}>
                    <td>
                      <div>
                        <strong>{promocion.nombre}</strong>
                      </div>
                      {promocion.descripcion && (
                        <small className="text-muted">{promocion.descripcion}</small>
                      )}
                    </td>
                    <td>
                      <Badge bg="secondary">
                        {tipos.find(t => t.value === promocion.tipo)?.label || promocion.tipo}
                      </Badge>
                    </td>
                    <td className="text-center">
                      <Badge bg="primary" style={{ fontSize: '1em' }}>
                        {getValorDisplay(promocion)}
                      </Badge>
                    </td>
                    <td>
                      <div className="small">
                        <Calendar className="me-1" size={12} />
                        {promocion.fecha_inicio 
                          ? new Date(promocion.fecha_inicio).toLocaleDateString()
                          : 'Sin fecha'
                        }
                        {' - '}
                        {promocion.fecha_fin 
                          ? new Date(promocion.fecha_fin).toLocaleDateString()
                          : 'Sin límite'
                        }
                      </div>
                    </td>
                    <td>
                      {promocion.codigo ? (
                        <code className="bg-light px-2 py-1 rounded">{promocion.codigo}</code>
                      ) : '-'}
                    </td>
                    <td>
                      {isVigente(promocion) ? (
                        <Badge bg="success">Vigente</Badge>
                      ) : promocion.activa ? (
                        <Badge bg="warning">Programada</Badge>
                      ) : (
                        <Badge bg="secondary">Inactiva</Badge>
                      )}
                    </td>
                    <td className="text-center">
                      {canEdit && (
                        <Button 
                          variant="outline-primary" 
                          size="sm" 
                          className="me-1"
                          onClick={() => handleEdit(promocion)}
                        >
                          <Pencil />
                        </Button>
                      )}
                      {canDelete && (
                        <Button 
                          variant="outline-danger" 
                          size="sm"
                          onClick={() => handleDelete(promocion.id)}
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
      </Card>

      {/* Modal de creación/edición */}
      <Modal show={showModal} onHide={() => setShowModal(false)} size="lg">
        <Form onSubmit={handleSubmit}>
          <Modal.Header closeButton>
            <Modal.Title>
              {editingPromocion ? 'Editar Promoción' : 'Nueva Promoción'}
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Row>
              <Col md={8}>
                <Form.Group className="mb-3">
                  <Form.Label>Nombre *</Form.Label>
                  <Form.Control
                    required
                    value={formData.nombre}
                    onChange={(e) => setFormData({...formData, nombre: e.target.value})}
                    placeholder="Nombre de la promoción"
                  />
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Código</Form.Label>
                  <Form.Control
                    value={formData.codigo}
                    onChange={(e) => setFormData({...formData, codigo: e.target.value.toUpperCase()})}
                    placeholder="CODIGO10"
                  />
                </Form.Group>
              </Col>
            </Row>
            <Form.Group className="mb-3">
              <Form.Label>Descripción</Form.Label>
              <Form.Control
                as="textarea"
                rows={2}
                value={formData.descripcion}
                onChange={(e) => setFormData({...formData, descripcion: e.target.value})}
                placeholder="Descripción de la promoción"
              />
            </Form.Group>
            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Tipo de Descuento</Form.Label>
                  <Form.Select
                    value={formData.tipo}
                    onChange={(e) => setFormData({...formData, tipo: e.target.value})}
                  >
                    {tipos.map(t => (
                      <option key={t.value} value={t.value}>{t.label}</option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>
                    Valor {formData.tipo === 'descuento_porcentaje' ? '(%)' : '(Bs.)'}
                  </Form.Label>
                  <Form.Control
                    type="number"
                    step="0.01"
                    value={formData.valor}
                    onChange={(e) => setFormData({...formData, valor: e.target.value})}
                    placeholder="10"
                    disabled={['2x1', '3x2'].includes(formData.tipo)}
                  />
                </Form.Group>
              </Col>
            </Row>
            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Fecha Inicio</Form.Label>
                  <Form.Control
                    type="date"
                    value={formData.fecha_inicio}
                    onChange={(e) => setFormData({...formData, fecha_inicio: e.target.value})}
                  />
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Fecha Fin</Form.Label>
                  <Form.Control
                    type="date"
                    value={formData.fecha_fin}
                    onChange={(e) => setFormData({...formData, fecha_fin: e.target.value})}
                  />
                </Form.Group>
              </Col>
            </Row>
            <Form.Group className="mb-3">
              <Form.Label>Categoría (aplica a todos los productos de esta categoría)</Form.Label>
              <Form.Control
                value={formData.categoria}
                onChange={(e) => setFormData({...formData, categoria: e.target.value})}
                placeholder="Ej: Laptops, Celulares..."
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Check
                type="switch"
                label="Promoción activa"
                checked={formData.activa}
                onChange={(e) => setFormData({...formData, activa: e.target.checked})}
              />
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowModal(false)}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              {editingPromocion ? 'Guardar Cambios' : 'Crear Promoción'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

export default Promociones;
