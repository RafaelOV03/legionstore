import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Button, Modal, Form, Badge, Spinner, Alert, InputGroup } from 'react-bootstrap';
import { Plus, Pencil, Trash, Search, People, Building, Telephone, Envelope } from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import * as api from '../services/resourceApi';

function Proveedores() {
  const { hasPermission } = useAuth();
  const [proveedores, setProveedores] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [search, setSearch] = useState('');
  
  const [showModal, setShowModal] = useState(false);
  const [editingProveedor, setEditingProveedor] = useState(null);
  const [formData, setFormData] = useState({
    nombre: '',
    ruc_nit: '',
    direccion: '',
    telefono: '',
    email: '',
    contacto_nombre: '',
    notas: '',
    activo: true
  });

  const canCreate = hasPermission('proveedores.create');
  const canEdit = hasPermission('proveedores.update');
  const canDelete = hasPermission('proveedores.delete');

  useEffect(() => {
    loadProveedores();
  }, []);

  const loadProveedores = async () => {
    try {
      setLoading(true);
      const data = await api.getProveedores();
      setProveedores(data || []);
    } catch (err) {
      setError('Error cargando proveedores');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (editingProveedor) {
        await api.updateProveedor(editingProveedor.id, formData);
      } else {
        await api.createProveedor(formData);
      }
      setShowModal(false);
      resetForm();
      loadProveedores();
    } catch (err) {
      setError('Error guardando proveedor');
      console.error(err);
    }
  };

  const handleEdit = (proveedor) => {
    setEditingProveedor(proveedor);
    setFormData({
      nombre: proveedor.nombre || '',
      ruc_nit: proveedor.ruc_nit || '',
      direccion: proveedor.direccion || '',
      telefono: proveedor.telefono || '',
      email: proveedor.email || '',
      contacto_nombre: proveedor.contacto_nombre || '',
      notas: proveedor.notas || '',
      activo: proveedor.activo !== false
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('¿Está seguro de eliminar este proveedor?')) return;
    try {
      await api.deleteProveedor(id);
      loadProveedores();
    } catch (err) {
      setError('Error eliminando proveedor');
      console.error(err);
    }
  };

  const resetForm = () => {
    setEditingProveedor(null);
    setFormData({
      nombre: '',
      ruc_nit: '',
      direccion: '',
      telefono: '',
      email: '',
      contacto_nombre: '',
      notas: '',
      activo: true
    });
  };

  const filteredProveedores = proveedores.filter(p =>
    p.nombre?.toLowerCase().includes(search.toLowerCase()) ||
    p.ruc_nit?.toLowerCase().includes(search.toLowerCase())
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
          <h2><People className="me-2" />Proveedores</h2>
          <p className="text-muted">Gestión de proveedores</p>
        </Col>
        <Col xs="auto">
          {canCreate && (
            <Button variant="primary" onClick={() => { resetForm(); setShowModal(true); }}>
              <Plus className="me-1" /> Nuevo Proveedor
            </Button>
          )}
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}

      <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
        <Card.Header className="bg-transparent">
          <InputGroup>
            <InputGroup.Text><Search /></InputGroup.Text>
            <Form.Control
              placeholder="Buscar por nombre o RUC/NIT..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </InputGroup>
        </Card.Header>
        <Card.Body className="p-0">
          <Table responsive hover className="mb-0">
            <thead>
              <tr>
                <th>Proveedor</th>
                <th>RUC/NIT</th>
                <th>Contacto</th>
                <th>Teléfono</th>
                <th>Email</th>
                <th>Estado</th>
                <th className="text-center">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {filteredProveedores.length === 0 ? (
                <tr>
                  <td colSpan="7" className="text-center py-4 text-muted">
                    No hay proveedores para mostrar
                  </td>
                </tr>
              ) : (
                filteredProveedores.map(proveedor => (
                  <tr key={proveedor.id}>
                    <td>
                      <div className="d-flex align-items-center">
                        <div 
                          className="rounded-circle p-2 me-2" 
                          style={{ background: proveedor.activo !== false ? '#0d6efd20' : '#6c757d20' }}
                        >
                          <Building size={16} style={{ color: proveedor.activo !== false ? '#0d6efd' : '#6c757d' }} />
                        </div>
                        <div>
                          <strong>{proveedor.nombre}</strong>
                          {proveedor.direccion && (
                            <div className="small text-muted">{proveedor.direccion}</div>
                          )}
                        </div>
                      </div>
                    </td>
                    <td>{proveedor.ruc_nit || '-'}</td>
                    <td>{proveedor.contacto_nombre || '-'}</td>
                    <td>
                      {proveedor.telefono && (
                        <span><Telephone className="me-1" size={12} />{proveedor.telefono}</span>
                      )}
                    </td>
                    <td>
                      {proveedor.email && (
                        <span><Envelope className="me-1" size={12} />{proveedor.email}</span>
                      )}
                    </td>
                    <td>
                      <Badge bg={proveedor.activo !== false ? 'success' : 'secondary'}>
                        {proveedor.activo !== false ? 'Activo' : 'Inactivo'}
                      </Badge>
                    </td>
                    <td className="text-center">
                      {canEdit && (
                        <Button 
                          variant="outline-primary" 
                          size="sm" 
                          className="me-1"
                          onClick={() => handleEdit(proveedor)}
                        >
                          <Pencil />
                        </Button>
                      )}
                      {canDelete && (
                        <Button 
                          variant="outline-danger" 
                          size="sm"
                          onClick={() => handleDelete(proveedor.id)}
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
          Total: {filteredProveedores.length} proveedor(es)
        </Card.Footer>
      </Card>

      {/* Modal de creación/edición */}
      <Modal show={showModal} onHide={() => setShowModal(false)} size="lg">
        <Form onSubmit={handleSubmit}>
          <Modal.Header closeButton>
            <Modal.Title>
              {editingProveedor ? 'Editar Proveedor' : 'Nuevo Proveedor'}
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Row>
              <Col md={8}>
                <Form.Group className="mb-3">
                  <Form.Label>Nombre del Proveedor *</Form.Label>
                  <Form.Control
                    required
                    value={formData.nombre}
                    onChange={(e) => setFormData({...formData, nombre: e.target.value})}
                    placeholder="Razón social o nombre comercial"
                  />
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>RUC/NIT</Form.Label>
                  <Form.Control
                    value={formData.ruc_nit}
                    onChange={(e) => setFormData({...formData, ruc_nit: e.target.value})}
                    placeholder="Número de RUC o NIT"
                  />
                </Form.Group>
              </Col>
            </Row>
            <Form.Group className="mb-3">
              <Form.Label>Dirección</Form.Label>
              <Form.Control
                value={formData.direccion}
                onChange={(e) => setFormData({...formData, direccion: e.target.value})}
                placeholder="Dirección del proveedor"
              />
            </Form.Group>
            <Row>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Teléfono</Form.Label>
                  <Form.Control
                    value={formData.telefono}
                    onChange={(e) => setFormData({...formData, telefono: e.target.value})}
                    placeholder="Teléfono"
                  />
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Email</Form.Label>
                  <Form.Control
                    type="email"
                    value={formData.email}
                    onChange={(e) => setFormData({...formData, email: e.target.value})}
                    placeholder="Email"
                  />
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Nombre de Contacto</Form.Label>
                  <Form.Control
                    value={formData.contacto_nombre}
                    onChange={(e) => setFormData({...formData, contacto_nombre: e.target.value})}
                    placeholder="Persona de contacto"
                  />
                </Form.Group>
              </Col>
            </Row>
            <Form.Group className="mb-3">
              <Form.Label>Notas</Form.Label>
              <Form.Control
                as="textarea"
                rows={2}
                value={formData.notas}
                onChange={(e) => setFormData({...formData, notas: e.target.value})}
                placeholder="Notas adicionales"
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Check
                type="switch"
                label="Proveedor activo"
                checked={formData.activo}
                onChange={(e) => setFormData({...formData, activo: e.target.checked})}
              />
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowModal(false)}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              {editingProveedor ? 'Guardar Cambios' : 'Crear Proveedor'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

export default Proveedores;
