import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Button, Modal, Form, Badge, Spinner, Alert, InputGroup } from 'react-bootstrap';
import { Plus, Pencil, Trash, Search, Tools, ExclamationTriangle } from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import * as api from '../services/inventarioApi';

function Insumos() {
  const { hasPermission } = useAuth();
  const [insumos, setInsumos] = useState([]);
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [search, setSearch] = useState('');
  
  const [showModal, setShowModal] = useState(false);
  const [showAjusteModal, setShowAjusteModal] = useState(false);
  const [editingInsumo, setEditingInsumo] = useState(null);
  const [selectedInsumo, setSelectedInsumo] = useState(null);
  const [sedes, setSedes] = useState([]);
  const [formData, setFormData] = useState({
    codigo: '',
    nombre: '',
    descripcion: '',
    categoria: '',
    stock: 0,
    stock_minimo: 5,
    costo: '',
    unidad_medida: 'unidad',
    sede_id: ''
  });
  const [ajusteData, setAjusteData] = useState({ cantidad: 0, tipo: 'entrada', motivo: '' });

  const canCreate = hasPermission('insumos.create');
  const canEdit = hasPermission('insumos.update');
  const canDelete = hasPermission('insumos.delete');

  const categorias = [
    'Repuestos',
    'Componentes Electrónicos',
    'Cables y Conectores',
    'Herramientas',
    'Consumibles',
    'Otros'
  ];

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [insumosData, statsData, sedesData] = await Promise.all([
        api.getInsumos(),
        api.getInsumosStats(),
        api.getSedes()
      ]);
      setInsumos(insumosData || []);
      setStats(statsData);
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
      const payload = {
        codigo: formData.codigo,
        nombre: formData.nombre,
        descripcion: formData.descripcion,
        categoria: formData.categoria,
        unidad_medida: formData.unidad_medida,
        stock: parseInt(formData.stock) || 0,
        stock_minimo: parseInt(formData.stock_minimo) || 5,
        costo: parseFloat(formData.costo) || 0,
        sede_id: parseInt(formData.sede_id)
      };

      if (editingInsumo) {
        await api.updateInsumo(editingInsumo.id, payload);
      } else {
        await api.createInsumo(payload);
      }
      
      setShowModal(false);
      resetForm();
      loadData();
    } catch (err) {
      setError('Error guardando insumo');
      console.error(err);
    }
  };

  const handleAjusteStock = async (e) => {
    e.preventDefault();
    try {
      await api.ajustarStockInsumo(selectedInsumo.id, {
        cantidad: parseInt(ajusteData.cantidad),
        tipo: ajusteData.tipo,
        motivo: ajusteData.motivo
      });
      setShowAjusteModal(false);
      setAjusteData({ cantidad: 0, tipo: 'entrada', motivo: '' });
      loadData();
    } catch (err) {
      setError('Error ajustando stock');
      console.error(err);
    }
  };

  const handleEdit = (insumo) => {
    setEditingInsumo(insumo);
    setFormData({
      codigo: insumo.codigo || '',
      nombre: insumo.nombre || '',
      descripcion: insumo.descripcion || '',
      categoria: insumo.categoria || '',
      stock: insumo.stock || 0,
      stock_minimo: insumo.stock_minimo || 5,
      costo: insumo.costo?.toString() || '',
      unidad_medida: insumo.unidad_medida || 'unidad',
      sede_id: insumo.sede_id?.toString() || ''
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('¿Está seguro de eliminar este insumo?')) return;
    try {
      await api.deleteInsumo(id);
      loadData();
    } catch (err) {
      setError('Error eliminando insumo');
      console.error(err);
    }
  };

  const resetForm = () => {
    setEditingInsumo(null);
    setFormData({
      codigo: '',
      nombre: '',
      descripcion: '',
      categoria: '',
      stock: 0,
      stock_minimo: 5,
      costo: '',
      unidad_medida: 'unidad',
      sede_id: sedes.length > 0 ? sedes[0].id.toString() : ''
    });
  };

  const filteredInsumos = insumos.filter(i =>
    i.nombre?.toLowerCase().includes(search.toLowerCase()) ||
    i.categoria?.toLowerCase().includes(search.toLowerCase())
  );

  const bajoStock = filteredInsumos.filter(i => i.stock <= (i.stock_minimo || 5));

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
          <h2><Tools className="me-2" />Insumos</h2>
          <p className="text-muted">Gestión de repuestos e insumos para servicio técnico</p>
        </Col>
        <Col xs="auto">
          {canCreate && (
            <Button variant="primary" onClick={() => { resetForm(); setShowModal(true); }}>
              <Plus className="me-1" /> Nuevo Insumo
            </Button>
          )}
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}

      {/* Estadísticas */}
      {stats && (
        <Row className="mb-4 g-3">
          <Col xs={6} md={3}>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-3">
                <h3 className="mb-0 text-primary">{stats.total_insumos || 0}</h3>
                <small className="text-muted">Total Insumos</small>
              </Card.Body>
            </Card>
          </Col>
          <Col xs={6} md={3}>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-3">
                <h3 className="mb-0 text-warning">{stats.bajo_stock || 0}</h3>
                <small className="text-muted">Bajo Stock</small>
              </Card.Body>
            </Card>
          </Col>
          <Col xs={6} md={3}>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-3">
                <h3 className="mb-0 text-danger">{stats.sin_stock || 0}</h3>
                <small className="text-muted">Sin Stock</small>
              </Card.Body>
            </Card>
          </Col>
          <Col xs={6} md={3}>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-3">
                <h3 className="mb-0 text-success">Bs. {(stats.valor_total || 0).toFixed(2)}</h3>
                <small className="text-muted">Valor Total</small>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      )}

      {bajoStock.length > 0 && (
        <Alert variant="warning" className="mb-4">
          <ExclamationTriangle className="me-2" />
          <strong>{bajoStock.length}</strong> insumo(s) con stock bajo o agotado:
          <span className="ms-2">
            {bajoStock.slice(0, 5).map(i => i.nombre).join(', ')}
            {bajoStock.length > 5 && '...'}
          </span>
        </Alert>
      )}

      <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
        <Card.Header className="bg-transparent">
          <InputGroup>
            <InputGroup.Text><Search /></InputGroup.Text>
            <Form.Control
              placeholder="Buscar por nombre o categoría..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </InputGroup>
        </Card.Header>
        <Card.Body className="p-0">
          <Table responsive hover className="mb-0">
            <thead>
              <tr>
                <th>Insumo</th>
                <th>Categoría</th>
                <th className="text-center">Stock</th>
                <th className="text-center">Mínimo</th>
                <th className="text-end">Costo Unit.</th>
                <th>Unidad</th>
                <th className="text-center">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {filteredInsumos.length === 0 ? (
                <tr>
                  <td colSpan="7" className="text-center py-4 text-muted">
                    No hay insumos para mostrar
                  </td>
                </tr>
              ) : (
                filteredInsumos.map(insumo => {
                  const stockBajo = insumo.stock <= (insumo.stock_minimo || 5);
                  return (
                    <tr key={insumo.id} className={stockBajo ? 'table-warning' : ''}>
                      <td>
                        <div>
                          <strong>{insumo.nombre}</strong>
                          {stockBajo && <ExclamationTriangle className="ms-2 text-warning" size={14} />}
                        </div>
                        {insumo.descripcion && (
                          <small className="text-muted">{insumo.descripcion}</small>
                        )}
                      </td>
                      <td>
                        <Badge bg="secondary">{insumo.categoria || 'Sin categoría'}</Badge>
                      </td>
                      <td className="text-center">
                        <Badge 
                          bg={insumo.stock > (insumo.stock_minimo || 5) ? 'success' : insumo.stock > 0 ? 'warning' : 'danger'}
                          style={{ fontSize: '0.9em' }}
                        >
                          {insumo.stock || 0}
                        </Badge>
                      </td>
                      <td className="text-center text-muted">{insumo.stock_minimo || 5}</td>
                      <td className="text-end">Bs. {insumo.costo?.toFixed(2) || '0.00'}</td>
                      <td>{insumo.unidad_medida || 'unidad'}</td>
                      <td className="text-center">
                        {canEdit && (
                          <>
                            <Button 
                              variant="outline-success" 
                              size="sm" 
                              className="me-1"
                              onClick={() => { 
                                setSelectedInsumo(insumo); 
                                setShowAjusteModal(true); 
                              }}
                              title="Ajustar stock"
                            >
                              <Plus />
                            </Button>
                            <Button 
                              variant="outline-primary" 
                              size="sm" 
                              className="me-1"
                              onClick={() => handleEdit(insumo)}
                              title="Editar"
                            >
                              <Pencil />
                            </Button>
                          </>
                        )}
                        {canDelete && (
                          <Button 
                            variant="outline-danger" 
                            size="sm"
                            onClick={() => handleDelete(insumo.id)}
                            title="Eliminar"
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
          Total: {filteredInsumos.length} insumo(s)
        </Card.Footer>
      </Card>

      {/* Modal de creación/edición */}
      <Modal show={showModal} onHide={() => setShowModal(false)} size="lg">
        <Form onSubmit={handleSubmit}>
          <Modal.Header closeButton>
            <Modal.Title>
              {editingInsumo ? 'Editar Insumo' : 'Nuevo Insumo'}
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
                    placeholder="INS-001"
                  />
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Nombre *</Form.Label>
                  <Form.Control
                    required
                    value={formData.nombre}
                    onChange={(e) => setFormData({...formData, nombre: e.target.value})}
                    placeholder="Nombre del insumo"
                  />
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Sede *</Form.Label>
                  <Form.Select
                    required
                    value={formData.sede_id}
                    onChange={(e) => setFormData({...formData, sede_id: e.target.value})}
                  >
                    <option value="">Seleccionar...</option>
                    {sedes.map(sede => (
                      <option key={sede.id} value={sede.id}>{sede.nombre}</option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
            </Row>
            <Row>
              <Col md={8}>
                <Form.Group className="mb-3">
                  <Form.Label>Descripción</Form.Label>
                  <Form.Control
                    as="textarea"
                    rows={2}
                    value={formData.descripcion}
                    onChange={(e) => setFormData({...formData, descripcion: e.target.value})}
                    placeholder="Descripción del insumo"
                  />
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Categoría</Form.Label>
                  <Form.Select
                    value={formData.categoria}
                    onChange={(e) => setFormData({...formData, categoria: e.target.value})}
                  >
                    <option value="">Seleccionar...</option>
                    {categorias.map(cat => (
                      <option key={cat} value={cat}>{cat}</option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
            </Row>
            <Row>
              <Col md={3}>
                <Form.Group className="mb-3">
                  <Form.Label>Stock Inicial</Form.Label>
                  <Form.Control
                    type="number"
                    min="0"
                    value={formData.stock}
                    onChange={(e) => setFormData({...formData, stock: e.target.value})}
                  />
                </Form.Group>
              </Col>
              <Col md={3}>
                <Form.Group className="mb-3">
                  <Form.Label>Stock Mínimo</Form.Label>
                  <Form.Control
                    type="number"
                    min="0"
                    value={formData.stock_minimo}
                    onChange={(e) => setFormData({...formData, stock_minimo: e.target.value})}
                  />
                </Form.Group>
              </Col>
              <Col md={3}>
                <Form.Group className="mb-3">
                  <Form.Label>Costo Unitario</Form.Label>
                  <Form.Control
                    type="number"
                    step="0.01"
                    value={formData.costo}
                    onChange={(e) => setFormData({...formData, costo: e.target.value})}
                    placeholder="0.00"
                  />
                </Form.Group>
              </Col>
              <Col md={3}>
                <Form.Group className="mb-3">
                  <Form.Label>Unidad de Medida</Form.Label>
                  <Form.Select
                    value={formData.unidad_medida}
                    onChange={(e) => setFormData({...formData, unidad_medida: e.target.value})}
                  >
                    <option value="unidad">Unidad</option>
                    <option value="pieza">Pieza</option>
                    <option value="metro">Metro</option>
                    <option value="litro">Litro</option>
                    <option value="kg">Kilogramo</option>
                  </Form.Select>
                </Form.Group>
              </Col>
            </Row>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowModal(false)}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              {editingInsumo ? 'Guardar Cambios' : 'Crear Insumo'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>

      {/* Modal de ajuste de stock */}
      <Modal show={showAjusteModal} onHide={() => setShowAjusteModal(false)}>
        <Form onSubmit={handleAjusteStock}>
          <Modal.Header closeButton>
            <Modal.Title>Ajustar Stock</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            {selectedInsumo && (
              <Alert variant="info">
                <strong>{selectedInsumo.nombre}</strong><br />
                Stock actual: <Badge bg="primary">{selectedInsumo.stock}</Badge>
              </Alert>
            )}
            <Form.Group className="mb-3">
              <Form.Label>Tipo de Movimiento</Form.Label>
              <Form.Select
                value={ajusteData.tipo}
                onChange={(e) => setAjusteData({...ajusteData, tipo: e.target.value})}
              >
                <option value="entrada">Entrada (aumentar)</option>
                <option value="salida">Salida (disminuir)</option>
              </Form.Select>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Cantidad</Form.Label>
              <Form.Control
                type="number"
                min="1"
                required
                value={ajusteData.cantidad}
                onChange={(e) => setAjusteData({...ajusteData, cantidad: e.target.value})}
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Motivo</Form.Label>
              <Form.Control
                as="textarea"
                rows={2}
                value={ajusteData.motivo}
                onChange={(e) => setAjusteData({...ajusteData, motivo: e.target.value})}
                placeholder="Motivo del ajuste"
              />
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowAjusteModal(false)}>
              Cancelar
            </Button>
            <Button variant="success" type="submit">
              Aplicar Ajuste
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

export default Insumos;
