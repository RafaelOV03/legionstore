import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Button, Modal, Form, Badge, Spinner, Alert, InputGroup } from 'react-bootstrap';
import { Plus, Pencil, Trash, Search, BoxSeam } from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import * as api from '../services/inventarioApi';

function Productos() {
  const { hasPermission } = useAuth();
  const [productos, setProductos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [search, setSearch] = useState('');
  
  const [showModal, setShowModal] = useState(false);
  const [editingProduct, setEditingProduct] = useState(null);
  const [formData, setFormData] = useState({
    codigo: '',
    name: '',
    description: '',
    precio_venta: '',
    precio_compra: '',
    category: '',
    brand: '',
    image_url: ''
  });

  const canCreate = hasPermission('products.create');
  const canEdit = hasPermission('products.update');
  const canDelete = hasPermission('products.delete');

  useEffect(() => {
    loadProductos();
  }, []);

  const loadProductos = async () => {
    try {
      setLoading(true);
      const data = await api.getProducts();
      setProductos(data || []);
    } catch (err) {
      setError('Error cargando productos');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        codigo: formData.codigo,
        name: formData.name,
        description: formData.description,
        precio_venta: parseFloat(formData.precio_venta) || 0,
        precio_compra: parseFloat(formData.precio_compra) || 0,
        category: formData.category,
        brand: formData.brand,
        image_url: formData.image_url
      };

      if (editingProduct) {
        await api.updateProduct(editingProduct.id, payload);
      } else {
        await api.createProduct(payload);
      }
      
      setShowModal(false);
      resetForm();
      loadProductos();
    } catch (err) {
      setError('Error guardando producto');
      console.error(err);
    }
  };

  const handleEdit = (producto) => {
    setEditingProduct(producto);
    setFormData({
      codigo: producto.codigo || '',
      name: producto.name || '',
      description: producto.description || '',
      precio_venta: producto.precio_venta?.toString() || '',
      precio_compra: producto.precio_compra?.toString() || '',
      category: producto.category || '',
      brand: producto.brand || '',
      image_url: producto.image_url || ''
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('¿Está seguro de eliminar este producto?')) return;
    try {
      await api.deleteProduct(id);
      loadProductos();
    } catch (err) {
      setError('Error eliminando producto');
      console.error(err);
    }
  };

  const resetForm = () => {
    setEditingProduct(null);
    setFormData({
      codigo: '',
      name: '',
      description: '',
      precio_venta: '',
      precio_compra: '',
      category: '',
      brand: '',
      image_url: ''
    });
  };

  const filteredProducts = productos.filter(p => 
    p.name?.toLowerCase().includes(search.toLowerCase()) ||
    p.category?.toLowerCase().includes(search.toLowerCase())
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
          <h2><BoxSeam className="me-2" />Productos</h2>
          <p className="text-muted">Gestión del catálogo de productos</p>
        </Col>
        <Col xs="auto">
          {canCreate && (
            <Button variant="primary" onClick={() => { resetForm(); setShowModal(true); }}>
              <Plus className="me-1" /> Nuevo Producto
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
              placeholder="Buscar productos..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </InputGroup>
        </Card.Header>
        <Card.Body className="p-0">
          <Table responsive hover className="mb-0">
            <thead>
              <tr>
                <th>Código</th>
                <th>Nombre</th>
                <th>Categoría</th>
                <th className="text-center">Stock</th>
                <th className="text-end">Costo</th>
                <th className="text-end">Precio</th>
                <th className="text-end">Margen</th>
                <th className="text-center">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {filteredProducts.length === 0 ? (
                <tr>
                  <td colSpan="8" className="text-center py-4 text-muted">
                    No hay productos para mostrar
                  </td>
                </tr>
              ) : (
                filteredProducts.map(producto => {
                  const margen = producto.precio_venta && producto.precio_compra 
                    ? (((producto.precio_venta - producto.precio_compra) / producto.precio_compra) * 100).toFixed(1) 
                    : 0;
                  return (
                    <tr key={producto.id}>
                      <td><code>{producto.codigo}</code></td>
                      <td>
                        <div className="d-flex align-items-center">
                          {producto.image_url && (
                            <img 
                              src={producto.image_url.startsWith('http') ? producto.image_url : `http://localhost:8080${producto.image_url}`}
                              alt={producto.name}
                              style={{ width: 40, height: 40, objectFit: 'cover', borderRadius: 4 }}
                              className="me-2"
                            />
                          )}
                          <div>
                            <strong>{producto.name}</strong>
                            {producto.description && (
                              <div className="small text-muted text-truncate" style={{ maxWidth: 200 }}>
                                {producto.description}
                              </div>
                            )}
                          </div>
                        </div>
                      </td>
                      <td>
                        <Badge bg="secondary">{producto.category || 'Sin categoría'}</Badge>
                      </td>
                      <td className="text-center">
                        <Badge bg={producto.stock_total > 10 ? 'success' : producto.stock_total > 0 ? 'warning' : 'danger'}>
                          {producto.stock_total || 0}
                        </Badge>
                      </td>
                      <td className="text-end">Bs. {producto.precio_compra?.toFixed(2) || '0.00'}</td>
                      <td className="text-end">Bs. {producto.precio_venta?.toFixed(2) || '0.00'}</td>
                      <td className="text-end">
                        <Badge bg={margen > 20 ? 'success' : margen > 0 ? 'warning' : 'danger'}>
                          {margen}%
                        </Badge>
                      </td>
                      <td className="text-center">
                        {canEdit && (
                          <Button 
                            variant="outline-primary" 
                            size="sm" 
                            className="me-1"
                            onClick={() => handleEdit(producto)}
                          >
                            <Pencil />
                          </Button>
                        )}
                        {canDelete && (
                          <Button 
                            variant="outline-danger" 
                            size="sm"
                            onClick={() => handleDelete(producto.id)}
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
          Total: {filteredProducts.length} producto(s)
        </Card.Footer>
      </Card>

      {/* Modal de creación/edición */}
      <Modal show={showModal} onHide={() => setShowModal(false)} size="lg">
        <Form onSubmit={handleSubmit}>
          <Modal.Header closeButton>
            <Modal.Title>
              {editingProduct ? 'Editar Producto' : 'Nuevo Producto'}
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Row>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Código *</Form.Label>
                  <Form.Control
                    required
                    value={formData.codigo}
                    onChange={(e) => setFormData({...formData, codigo: e.target.value})}
                    placeholder="SKU-001"
                  />
                </Form.Group>
              </Col>
              <Col md={8}>
                <Form.Group className="mb-3">
                  <Form.Label>Nombre *</Form.Label>
                  <Form.Control
                    required
                    value={formData.name}
                    onChange={(e) => setFormData({...formData, name: e.target.value})}
                    placeholder="Nombre del producto"
                  />
                </Form.Group>
              </Col>
            </Row>
            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Categoría *</Form.Label>
                  <Form.Control
                    required
                    value={formData.category}
                    onChange={(e) => setFormData({...formData, category: e.target.value})}
                    placeholder="Categoría"
                  />
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Marca</Form.Label>
                  <Form.Control
                    value={formData.brand}
                    onChange={(e) => setFormData({...formData, brand: e.target.value})}
                    placeholder="Marca"
                  />
                </Form.Group>
              </Col>
            </Row>
            <Form.Group className="mb-3">
              <Form.Label>Descripción</Form.Label>
              <Form.Control
                as="textarea"
                rows={2}
                value={formData.description}
                onChange={(e) => setFormData({...formData, description: e.target.value})}
                placeholder="Descripción del producto"
              />
            </Form.Group>
            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Precio Compra (Bs.)</Form.Label>
                  <Form.Control
                    type="number"
                    step="0.01"
                    value={formData.precio_compra}
                    onChange={(e) => setFormData({...formData, precio_compra: e.target.value})}
                    placeholder="0.00"
                  />
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Precio Venta (Bs.) *</Form.Label>
                  <Form.Control
                    type="number"
                    step="0.01"
                    required
                    value={formData.precio_venta}
                    onChange={(e) => setFormData({...formData, precio_venta: e.target.value})}
                    placeholder="0.00"
                  />
                </Form.Group>
              </Col>
            </Row>
            <Form.Group className="mb-3">
              <Form.Label>URL de Imagen</Form.Label>
              <Form.Control
                value={formData.image_url}
                onChange={(e) => setFormData({...formData, image_url: e.target.value})}
                placeholder="https://..."
              />
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowModal(false)}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              {editingProduct ? 'Guardar Cambios' : 'Crear Producto'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

export default Productos;
