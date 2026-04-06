import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Button, Modal, Form, Badge, Spinner, Alert, InputGroup } from 'react-bootstrap';
import { Plus, Pencil, Trash, Search, Receipt, FileEarmarkPdf, Cart } from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import * as api from '../services/resourceApi';

function Cotizaciones() {
  const { hasPermission, user } = useAuth();
  const [cotizaciones, setCotizaciones] = useState([]);
  const [productos, setProductos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [search, setSearch] = useState('');
  
  const [showModal, setShowModal] = useState(false);
  const [editingCotizacion, setEditingCotizacion] = useState(null);
  const [formData, setFormData] = useState({
    cliente_nombre: '',
    cliente_telefono: '',
    cliente_email: '',
    validez_dias: 15,
    notas: '',
    items: []
  });

  const canCreate = hasPermission('cotizaciones.create');
  const canEdit = hasPermission('cotizaciones.update');
  const canDelete = hasPermission('cotizaciones.delete');

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [cotizacionesData, productosData] = await Promise.all([
        api.getCotizaciones(),
        api.getProducts()
      ]);
      setCotizaciones(cotizacionesData || []);
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
      // Transform items to use producto_id instead of product_id
      const transformedItems = formData.items.map(item => ({
        producto_id: parseInt(item.product_id),
        cantidad: parseInt(item.cantidad),
        precio_unitario: parseFloat(item.precio_unitario)
      }));

      const payload = {
        cliente_nombre: formData.cliente_nombre,
        cliente_telefono: formData.cliente_telefono,
        cliente_email: formData.cliente_email,
        validez_dias: parseInt(formData.validez_dias),
        notas: formData.notas,
        items: transformedItems,
        sede_id: user?.sede_id || 1
      };

      if (editingCotizacion) {
        await api.updateCotizacion(editingCotizacion.id, payload);
      } else {
        await api.createCotizacion(payload);
      }
      
      setShowModal(false);
      resetForm();
      loadData();
      setSuccess('Cotización guardada exitosamente');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error guardando cotización');
      console.error(err);
    }
  };

  const handleEdit = (cotizacion) => {
    setEditingCotizacion(cotizacion);
    setFormData({
      cliente_nombre: cotizacion.cliente_nombre || '',
      cliente_telefono: cotizacion.cliente_telefono || '',
      cliente_email: cotizacion.cliente_email || '',
      validez_dias: cotizacion.validez_dias || 15,
      notas: cotizacion.notas || '',
      items: cotizacion.items || []
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('¿Está seguro de eliminar esta cotización?')) return;
    try {
      await api.deleteCotizacion(id);
      loadData();
    } catch (err) {
      setError('Error eliminando cotización');
      console.error(err);
    }
  };

  const handleConvertirAVenta = async (id) => {
    if (!window.confirm('¿Convertir esta cotización en venta? Se descontará el stock.')) return;
    try {
      await api.convertirCotizacionAVenta(id);
      loadData();
      setSuccess('Cotización convertida a venta exitosamente');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error convirtiendo cotización a venta');
      console.error(err);
    }
  };

  const resetForm = () => {
    setEditingCotizacion(null);
    setFormData({
      cliente_nombre: '',
      cliente_telefono: '',
      cliente_email: '',
      validez_dias: 15,
      notas: '',
      items: []
    });
  };

  const addItem = () => {
    setFormData({
      ...formData,
      items: [...formData.items, { product_id: '', cantidad: 1, precio_unitario: 0 }]
    });
  };

  const removeItem = (index) => {
    const newItems = [...formData.items];
    newItems.splice(index, 1);
    setFormData({ ...formData, items: newItems });
  };

  const updateItem = (index, field, value) => {
    const newItems = [...formData.items];
    newItems[index][field] = value;
    
    // Autocompletar precio si se selecciona producto
    if (field === 'product_id') {
      const producto = productos.find(p => p.id.toString() === value);
      if (producto) {
        newItems[index].precio_unitario = producto.precio_venta || 0;
      }
    }
    
    setFormData({ ...formData, items: newItems });
  };

  const calcularTotal = (items) => {
    return items.reduce((sum, item) => sum + (item.cantidad * item.precio_unitario), 0);
  };

  const getEstadoBadge = (estado) => {
    const colors = {
      pendiente: 'warning',
      aceptada: 'success',
      rechazada: 'danger',
      vencida: 'secondary',
      convertida: 'info'
    };
    return <Badge bg={colors[estado] || 'secondary'}>{estado}</Badge>;
  };

  const filteredCotizaciones = cotizaciones.filter(c =>
    c.cliente_nombre?.toLowerCase().includes(search.toLowerCase()) ||
    c.numero?.toLowerCase().includes(search.toLowerCase())
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
          <h2><Receipt className="me-2" />Cotizaciones</h2>
          <p className="text-muted">Gestión de cotizaciones a clientes</p>
        </Col>
        <Col xs="auto">
          {canCreate && (
            <Button variant="primary" onClick={() => { resetForm(); setShowModal(true); }}>
              <Plus className="me-1" /> Nueva Cotización
            </Button>
          )}
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}
      {success && <Alert variant="success" onClose={() => setSuccess(null)} dismissible>{success}</Alert>}

      <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
        <Card.Header className="bg-transparent">
          <InputGroup>
            <InputGroup.Text><Search /></InputGroup.Text>
            <Form.Control
              placeholder="Buscar por cliente o número..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </InputGroup>
        </Card.Header>
        <Card.Body className="p-0">
          <Table responsive hover className="mb-0">
            <thead>
              <tr>
                <th>Número</th>
                <th>Cliente</th>
                <th className="text-end">Total</th>
                <th>Estado</th>
                <th>Válida hasta</th>
                <th>Vendedor</th>
                <th className="text-center">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {filteredCotizaciones.length === 0 ? (
                <tr>
                  <td colSpan="7" className="text-center py-4 text-muted">
                    No hay cotizaciones para mostrar
                  </td>
                </tr>
              ) : (
                filteredCotizaciones.map(cotizacion => {
                  const fechaVencimiento = new Date(cotizacion.created_at);
                  fechaVencimiento.setDate(fechaVencimiento.getDate() + (cotizacion.validez_dias || 15));
                  
                  return (
                    <tr key={cotizacion.id}>
                      <td><strong>{cotizacion.numero || `COT-${cotizacion.id}`}</strong></td>
                      <td>
                        <div>{cotizacion.cliente_nombre}</div>
                        {cotizacion.cliente_telefono && (
                          <small className="text-muted">{cotizacion.cliente_telefono}</small>
                        )}
                      </td>
                      <td className="text-end">
                        <strong>Bs. {cotizacion.total?.toFixed(2) || '0.00'}</strong>
                      </td>
                      <td>{getEstadoBadge(cotizacion.estado)}</td>
                      <td>{fechaVencimiento.toLocaleDateString()}</td>
                      <td>{cotizacion.usuario?.name || 'N/A'}</td>
                      <td className="text-center">
                        {cotizacion.estado === 'pendiente' && (
                          <>
                            <Button 
                              variant="outline-success" 
                              size="sm" 
                              className="me-1"
                              onClick={() => handleConvertirAVenta(cotizacion.id)}
                              title="Convertir a venta"
                            >
                              <Cart />
                            </Button>
                          </>
                        )}
                        {canEdit && (
                          <Button 
                            variant="outline-primary" 
                            size="sm" 
                            className="me-1"
                            onClick={() => handleEdit(cotizacion)}
                          >
                            <Pencil />
                          </Button>
                        )}
                        {canDelete && (
                          <Button 
                            variant="outline-danger" 
                            size="sm"
                            onClick={() => handleDelete(cotizacion.id)}
                          >
                            <Trash />
                          </Button>
                        )}
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </Table>
        </Card.Body>
        <Card.Footer className="bg-transparent text-muted">
          Total: {filteredCotizaciones.length} cotización(es)
        </Card.Footer>
      </Card>

      {/* Modal de creación/edición */}
      <Modal show={showModal} onHide={() => setShowModal(false)} size="xl">
        <Form onSubmit={handleSubmit}>
          <Modal.Header closeButton>
            <Modal.Title>
              {editingCotizacion ? 'Editar Cotización' : 'Nueva Cotización'}
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Row>
              <Col md={4}>
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
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Teléfono</Form.Label>
                  <Form.Control
                    value={formData.cliente_telefono}
                    onChange={(e) => setFormData({...formData, cliente_telefono: e.target.value})}
                    placeholder="Teléfono"
                  />
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Email</Form.Label>
                  <Form.Control
                    type="email"
                    value={formData.cliente_email}
                    onChange={(e) => setFormData({...formData, cliente_email: e.target.value})}
                    placeholder="Email"
                  />
                </Form.Group>
              </Col>
            </Row>
            <Row>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Días de validez</Form.Label>
                  <Form.Control
                    type="number"
                    min="1"
                    value={formData.validez_dias}
                    onChange={(e) => setFormData({...formData, validez_dias: e.target.value})}
                  />
                </Form.Group>
              </Col>
              <Col md={8}>
                <Form.Group className="mb-3">
                  <Form.Label>Notas</Form.Label>
                  <Form.Control
                    value={formData.notas}
                    onChange={(e) => setFormData({...formData, notas: e.target.value})}
                    placeholder="Notas adicionales"
                  />
                </Form.Group>
              </Col>
            </Row>

            <hr />
            <div className="d-flex justify-content-between align-items-center mb-3">
              <h6 className="mb-0">Items de la Cotización</h6>
              <Button variant="outline-primary" size="sm" onClick={addItem}>
                <Plus className="me-1" /> Agregar Item
              </Button>
            </div>

            {formData.items.length === 0 ? (
              <Alert variant="info">No hay items. Click en "Agregar Item" para comenzar.</Alert>
            ) : (
              <Table responsive size="sm">
                <thead>
                  <tr>
                    <th style={{ width: '40%' }}>Producto</th>
                    <th style={{ width: '15%' }}>Cantidad</th>
                    <th style={{ width: '20%' }}>Precio Unit.</th>
                    <th style={{ width: '15%' }}>Subtotal</th>
                    <th style={{ width: '10%' }}></th>
                  </tr>
                </thead>
                <tbody>
                  {formData.items.map((item, index) => (
                    <tr key={index}>
                      <td>
                        <Form.Select
                          size="sm"
                          value={item.product_id}
                          onChange={(e) => updateItem(index, 'product_id', e.target.value)}
                        >
                          <option value="">Seleccionar...</option>
                          {productos.map(p => (
                            <option key={p.id} value={p.id}>{p.name} - Bs. {(p.precio_venta || 0).toFixed(2)}</option>
                          ))}
                        </Form.Select>
                      </td>
                      <td>
                        <Form.Control
                          type="number"
                          size="sm"
                          min="1"
                          value={item.cantidad}
                          onChange={(e) => updateItem(index, 'cantidad', parseInt(e.target.value) || 1)}
                        />
                      </td>
                      <td>
                        <Form.Control
                          type="number"
                          size="sm"
                          step="0.01"
                          value={item.precio_unitario}
                          onChange={(e) => updateItem(index, 'precio_unitario', parseFloat(e.target.value) || 0)}
                        />
                      </td>
                      <td className="text-end">
                        Bs. {(item.cantidad * item.precio_unitario).toFixed(2)}
                      </td>
                      <td>
                        <Button 
                          variant="outline-danger" 
                          size="sm"
                          onClick={() => removeItem(index)}
                        >
                          <Trash />
                        </Button>
                      </td>
                    </tr>
                  ))}
                  <tr>
                    <td colSpan="3" className="text-end"><strong>TOTAL:</strong></td>
                    <td className="text-end"><strong>Bs. {calcularTotal(formData.items).toFixed(2)}</strong></td>
                    <td></td>
                  </tr>
                </tbody>
              </Table>
            )}
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowModal(false)}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              {editingCotizacion ? 'Guardar Cambios' : 'Crear Cotización'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

export default Cotizaciones;
