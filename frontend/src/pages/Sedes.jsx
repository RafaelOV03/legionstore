import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Button, Modal, Form, Badge, Spinner, Alert } from 'react-bootstrap';
import { Plus, Pencil, Trash, Building, GeoAlt, Telephone } from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import * as api from '../services/inventarioApi';

function Sedes() {
  const { hasPermission } = useAuth();
  const [sedes, setSedes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  
  const [showModal, setShowModal] = useState(false);
  const [editingSede, setEditingSede] = useState(null);
  const [formData, setFormData] = useState({
    nombre: '',
    direccion: '',
    telefono: '',
    activa: true
  });

  const canCreate = hasPermission('sedes.create');
  const canEdit = hasPermission('sedes.update');
  const canDelete = hasPermission('sedes.delete');

  useEffect(() => {
    loadSedes();
  }, []);

  const loadSedes = async () => {
    try {
      setLoading(true);
      const data = await api.getSedes();
      setSedes(data || []);
    } catch (err) {
      setError('Error cargando sedes');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (editingSede) {
        await api.updateSede(editingSede.id, formData);
      } else {
        await api.createSede(formData);
      }
      setShowModal(false);
      resetForm();
      loadSedes();
    } catch (err) {
      setError('Error guardando sede');
      console.error(err);
    }
  };

  const handleEdit = (sede) => {
    setEditingSede(sede);
    setFormData({
      nombre: sede.nombre || '',
      direccion: sede.direccion || '',
      telefono: sede.telefono || '',
      activa: sede.activa !== false
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('¿Está seguro de eliminar esta sede?')) return;
    try {
      await api.deleteSede(id);
      loadSedes();
    } catch (err) {
      setError('Error eliminando sede');
      console.error(err);
    }
  };

  const resetForm = () => {
    setEditingSede(null);
    setFormData({
      nombre: '',
      direccion: '',
      telefono: '',
      activa: true
    });
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
          <h2><Building className="me-2" />Sedes</h2>
          <p className="text-muted">Gestión de sucursales</p>
        </Col>
        <Col xs="auto">
          {canCreate && (
            <Button variant="primary" onClick={() => { resetForm(); setShowModal(true); }}>
              <Plus className="me-1" /> Nueva Sede
            </Button>
          )}
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}

      <Row className="g-4">
        {sedes.length === 0 ? (
          <Col>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-5 text-muted">
                No hay sedes registradas
              </Card.Body>
            </Card>
          </Col>
        ) : (
          sedes.map(sede => (
            <Col md={6} lg={4} key={sede.id}>
              <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
                <Card.Body>
                  <div className="d-flex justify-content-between align-items-start mb-3">
                    <div className="d-flex align-items-center">
                      <div 
                        className="rounded-circle p-3 me-3" 
                        style={{ background: sede.activa !== false ? '#0d6efd20' : '#6c757d20' }}
                      >
                        <Building size={24} style={{ color: sede.activa !== false ? '#0d6efd' : '#6c757d' }} />
                      </div>
                      <div>
                        <h5 className="mb-0">{sede.nombre}</h5>
                        <Badge bg={sede.activa !== false ? 'success' : 'secondary'}>
                          {sede.activa !== false ? 'Activa' : 'Inactiva'}
                        </Badge>
                      </div>
                    </div>
                    <div>
                      {canEdit && (
                        <Button 
                          variant="outline-primary" 
                          size="sm" 
                          className="me-1"
                          onClick={() => handleEdit(sede)}
                        >
                          <Pencil />
                        </Button>
                      )}
                      {canDelete && (
                        <Button 
                          variant="outline-danger" 
                          size="sm"
                          onClick={() => handleDelete(sede.id)}
                        >
                          <Trash />
                        </Button>
                      )}
                    </div>
                  </div>
                  
                  {sede.direccion && (
                    <div className="d-flex align-items-center text-muted mb-2">
                      <GeoAlt className="me-2" />
                      <small>{sede.direccion}</small>
                    </div>
                  )}
                  
                  {sede.telefono && (
                    <div className="d-flex align-items-center text-muted">
                      <Telephone className="me-2" />
                      <small>{sede.telefono}</small>
                    </div>
                  )}
                </Card.Body>
              </Card>
            </Col>
          ))
        )}
      </Row>

      {/* Modal de creación/edición */}
      <Modal show={showModal} onHide={() => setShowModal(false)}>
        <Form onSubmit={handleSubmit}>
          <Modal.Header closeButton>
            <Modal.Title>
              {editingSede ? 'Editar Sede' : 'Nueva Sede'}
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Form.Group className="mb-3">
              <Form.Label>Nombre *</Form.Label>
              <Form.Control
                required
                value={formData.nombre}
                onChange={(e) => setFormData({...formData, nombre: e.target.value})}
                placeholder="Nombre de la sede"
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Dirección</Form.Label>
              <Form.Control
                value={formData.direccion}
                onChange={(e) => setFormData({...formData, direccion: e.target.value})}
                placeholder="Dirección"
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Teléfono</Form.Label>
              <Form.Control
                value={formData.telefono}
                onChange={(e) => setFormData({...formData, telefono: e.target.value})}
                placeholder="Teléfono"
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Check
                type="switch"
                label="Sede activa"
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
              {editingSede ? 'Guardar Cambios' : 'Crear Sede'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

export default Sedes;
