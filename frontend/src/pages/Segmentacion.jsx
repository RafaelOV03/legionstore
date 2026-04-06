import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Button, Modal, Form, Badge, Spinner, Alert } from 'react-bootstrap';
import { Plus, Pencil, Trash, People, PieChart } from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import * as api from '../services/resourceApi';

function Segmentacion() {
  const { hasPermission } = useAuth();
  const [segmentaciones, setSegmentaciones] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  
  const [showModal, setShowModal] = useState(false);
  const [editingSegmentacion, setEditingSegmentacion] = useState(null);
  const [formData, setFormData] = useState({
    nombre: '',
    descripcion: '',
    criterios: '',
    monto_minimo: '',
    monto_maximo: '',
    frecuencia_compra: '',
    activa: true
  });

  const canCreate = hasPermission('segmentacion.create');
  const canEdit = hasPermission('segmentacion.update');
  const canDelete = hasPermission('segmentacion.delete');

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const data = await api.getSegmentaciones();
      setSegmentaciones(data || []);
    } catch (err) {
      setError('Error cargando segmentaciones');
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
        monto_minimo: formData.monto_minimo ? parseFloat(formData.monto_minimo) : null,
        monto_maximo: formData.monto_maximo ? parseFloat(formData.monto_maximo) : null,
        frecuencia_compra: formData.frecuencia_compra ? parseInt(formData.frecuencia_compra) : null
      };

      if (editingSegmentacion) {
        await api.updateSegmentacion(editingSegmentacion.id, payload);
      } else {
        await api.createSegmentacion(payload);
      }
      
      setShowModal(false);
      resetForm();
      loadData();
      setSuccess('Segmentación guardada exitosamente');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error guardando segmentación');
      console.error(err);
    }
  };

  const handleEdit = (segmentacion) => {
    setEditingSegmentacion(segmentacion);
    setFormData({
      nombre: segmentacion.nombre || '',
      descripcion: segmentacion.descripcion || '',
      criterios: segmentacion.criterios || '',
      monto_minimo: segmentacion.monto_minimo?.toString() || '',
      monto_maximo: segmentacion.monto_maximo?.toString() || '',
      frecuencia_compra: segmentacion.frecuencia_compra?.toString() || '',
      activa: segmentacion.activa !== false
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('¿Está seguro de eliminar esta segmentación?')) return;
    try {
      await api.deleteSegmentacion(id);
      loadData();
    } catch (err) {
      setError('Error eliminando segmentación');
      console.error(err);
    }
  };

  const resetForm = () => {
    setEditingSegmentacion(null);
    setFormData({
      nombre: '',
      descripcion: '',
      criterios: '',
      monto_minimo: '',
      monto_maximo: '',
      frecuencia_compra: '',
      activa: true
    });
  };

  // Colores predefinidos para segmentos
  const segmentColors = ['primary', 'success', 'warning', 'info', 'danger', 'secondary'];

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
          <h2><PieChart className="me-2" />Segmentación de Clientes</h2>
          <p className="text-muted">Definir segmentos para análisis y promociones</p>
        </Col>
        <Col xs="auto">
          {canCreate && (
            <Button variant="primary" onClick={() => { resetForm(); setShowModal(true); }}>
              <Plus className="me-1" /> Nueva Segmentación
            </Button>
          )}
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}
      {success && <Alert variant="success" onClose={() => setSuccess(null)} dismissible>{success}</Alert>}

      <Row className="g-4">
        {segmentaciones.length === 0 ? (
          <Col>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-5 text-muted">
                No hay segmentaciones definidas. Crea una para comenzar a clasificar clientes.
              </Card.Body>
            </Card>
          </Col>
        ) : (
          segmentaciones.map((seg, index) => (
            <Col md={6} lg={4} key={seg.id}>
              <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
                <Card.Header className="bg-transparent d-flex justify-content-between align-items-center">
                  <div className="d-flex align-items-center">
                    <Badge bg={segmentColors[index % segmentColors.length]} className="me-2">
                      <People className="me-1" />
                      {seg.clientes_count || 0}
                    </Badge>
                    <h6 className="mb-0">{seg.nombre}</h6>
                  </div>
                  <Badge bg={seg.activa !== false ? 'success' : 'secondary'}>
                    {seg.activa !== false ? 'Activa' : 'Inactiva'}
                  </Badge>
                </Card.Header>
                <Card.Body>
                  {seg.descripcion && (
                    <p className="text-muted small mb-3">{seg.descripcion}</p>
                  )}
                  
                  <div className="mb-3">
                    <strong className="small text-muted">Criterios:</strong>
                    <ul className="mb-0 small">
                      {seg.monto_minimo && (
                        <li>Monto mínimo: Bs. {seg.monto_minimo}</li>
                      )}
                      {seg.monto_maximo && (
                        <li>Monto máximo: Bs. {seg.monto_maximo}</li>
                      )}
                      {seg.frecuencia_compra && (
                        <li>Frecuencia: {seg.frecuencia_compra} compras/mes</li>
                      )}
                      {seg.criterios && (
                        <li>{seg.criterios}</li>
                      )}
                      {!seg.monto_minimo && !seg.monto_maximo && !seg.frecuencia_compra && !seg.criterios && (
                        <li className="text-muted">Sin criterios definidos</li>
                      )}
                    </ul>
                  </div>

                  <div className="d-flex gap-2">
                    {canEdit && (
                      <Button 
                        variant="outline-primary" 
                        size="sm"
                        onClick={() => handleEdit(seg)}
                      >
                        <Pencil className="me-1" /> Editar
                      </Button>
                    )}
                    {canDelete && (
                      <Button 
                        variant="outline-danger" 
                        size="sm"
                        onClick={() => handleDelete(seg.id)}
                      >
                        <Trash className="me-1" /> Eliminar
                      </Button>
                    )}
                  </div>
                </Card.Body>
              </Card>
            </Col>
          ))
        )}
      </Row>

      {/* Modal de creación/edición */}
      <Modal show={showModal} onHide={() => setShowModal(false)} size="lg">
        <Form onSubmit={handleSubmit}>
          <Modal.Header closeButton>
            <Modal.Title>
              {editingSegmentacion ? 'Editar Segmentación' : 'Nueva Segmentación'}
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Form.Group className="mb-3">
              <Form.Label>Nombre *</Form.Label>
              <Form.Control
                required
                value={formData.nombre}
                onChange={(e) => setFormData({...formData, nombre: e.target.value})}
                placeholder="Ej: Clientes VIP, Compradores Frecuentes"
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Descripción</Form.Label>
              <Form.Control
                as="textarea"
                rows={2}
                value={formData.descripcion}
                onChange={(e) => setFormData({...formData, descripcion: e.target.value})}
                placeholder="Descripción del segmento"
              />
            </Form.Group>
            
            <h6 className="mt-4 mb-3">Criterios de Segmentación</h6>
            <Row>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Monto Mínimo (Bs.)</Form.Label>
                  <Form.Control
                    type="number"
                    step="0.01"
                    value={formData.monto_minimo}
                    onChange={(e) => setFormData({...formData, monto_minimo: e.target.value})}
                    placeholder="0.00"
                  />
                  <Form.Text className="text-muted">
                    Compras totales mínimas
                  </Form.Text>
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Monto Máximo (Bs.)</Form.Label>
                  <Form.Control
                    type="number"
                    step="0.01"
                    value={formData.monto_maximo}
                    onChange={(e) => setFormData({...formData, monto_maximo: e.target.value})}
                    placeholder="0.00"
                  />
                  <Form.Text className="text-muted">
                    Compras totales máximas
                  </Form.Text>
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Frecuencia de Compra</Form.Label>
                  <Form.Control
                    type="number"
                    value={formData.frecuencia_compra}
                    onChange={(e) => setFormData({...formData, frecuencia_compra: e.target.value})}
                    placeholder="0"
                  />
                  <Form.Text className="text-muted">
                    Compras por mes
                  </Form.Text>
                </Form.Group>
              </Col>
            </Row>
            <Form.Group className="mb-3">
              <Form.Label>Otros Criterios</Form.Label>
              <Form.Control
                as="textarea"
                rows={2}
                value={formData.criterios}
                onChange={(e) => setFormData({...formData, criterios: e.target.value})}
                placeholder="Criterios adicionales de segmentación"
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Check
                type="switch"
                label="Segmentación activa"
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
              {editingSegmentacion ? 'Guardar Cambios' : 'Crear Segmentación'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

export default Segmentacion;
