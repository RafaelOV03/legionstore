import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Badge, Spinner, Alert, Form, InputGroup } from 'react-bootstrap';
import { Building, Search, BoxSeam } from 'react-bootstrap-icons';
import * as api from '../services/inventarioApi';

function Stock() {
  const [productos, setProductos] = useState([]);
  const [sedes, setSedes] = useState([]);
  const [stockData, setStockData] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [search, setSearch] = useState('');
  const [sedeFilter, setSedeFilter] = useState('');

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [productosData, sedesData] = await Promise.all([
        api.getProducts(),
        api.getSedes()
      ]);
      
      setProductos(productosData || []);
      setSedes(sedesData || []);
      
      // Cargar stock de cada sede
      const stockBySede = {};
      for (const sede of (sedesData || [])) {
        try {
          const sedeStock = await api.getStockBySede(sede.id);
          stockBySede[sede.id] = sedeStock || [];
        } catch (e) {
          stockBySede[sede.id] = [];
        }
      }
      setStockData(stockBySede);
    } catch (err) {
      setError('Error cargando datos de stock');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const getStockForProduct = (productoId, sedeId) => {
    const sedeStock = stockData[sedeId] || [];
    const item = sedeStock.find(s => s.product_id === productoId);
    return item?.cantidad || 0;
  };

  const getTotalStock = (productoId) => {
    let total = 0;
    for (const sedeId of Object.keys(stockData)) {
      total += getStockForProduct(productoId, parseInt(sedeId));
    }
    return total;
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
          <h2><Building className="me-2" />Stock Multisede</h2>
          <p className="text-muted">Consulta de stock por sede</p>
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}

      <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
        <Card.Header className="bg-transparent">
          <Row className="g-2">
            <Col md={6}>
              <InputGroup>
                <InputGroup.Text><Search /></InputGroup.Text>
                <Form.Control
                  placeholder="Buscar productos..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
              </InputGroup>
            </Col>
            <Col md={6}>
              <Form.Select
                value={sedeFilter}
                onChange={(e) => setSedeFilter(e.target.value)}
              >
                <option value="">Todas las sedes</option>
                {sedes.map(sede => (
                  <option key={sede.id} value={sede.id}>{sede.nombre}</option>
                ))}
              </Form.Select>
            </Col>
          </Row>
        </Card.Header>
        <Card.Body className="p-0">
          <div style={{ overflowX: 'auto' }}>
            <Table hover className="mb-0">
              <thead>
                <tr>
                  <th style={{ minWidth: 250 }}>Producto</th>
                  <th className="text-center" style={{ minWidth: 80 }}>Total</th>
                  {sedes
                    .filter(s => !sedeFilter || s.id.toString() === sedeFilter)
                    .map(sede => (
                      <th key={sede.id} className="text-center" style={{ minWidth: 100 }}>
                        {sede.nombre}
                      </th>
                    ))
                  }
                </tr>
              </thead>
              <tbody>
                {filteredProducts.length === 0 ? (
                  <tr>
                    <td colSpan={sedes.length + 2} className="text-center py-4 text-muted">
                      No hay productos para mostrar
                    </td>
                  </tr>
                ) : (
                  filteredProducts.map(producto => {
                    const totalStock = getTotalStock(producto.id);
                    return (
                      <tr key={producto.id}>
                        <td>
                          <div className="d-flex align-items-center">
                            <BoxSeam className="text-muted me-2" />
                            <div>
                              <strong>{producto.name}</strong>
                              <div className="small text-muted">{producto.category}</div>
                            </div>
                          </div>
                        </td>
                        <td className="text-center">
                          <Badge 
                            bg={totalStock > 10 ? 'success' : totalStock > 0 ? 'warning' : 'danger'}
                            style={{ fontSize: '0.9em' }}
                          >
                            {totalStock}
                          </Badge>
                        </td>
                        {sedes
                          .filter(s => !sedeFilter || s.id.toString() === sedeFilter)
                          .map(sede => {
                            const stock = getStockForProduct(producto.id, sede.id);
                            return (
                              <td key={sede.id} className="text-center">
                                <Badge 
                                  bg={stock > 10 ? 'secondary' : stock > 0 ? 'warning' : 'dark'}
                                  text={stock === 0 ? 'muted' : undefined}
                                >
                                  {stock}
                                </Badge>
                              </td>
                            );
                          })
                        }
                      </tr>
                    );
                  })
                )}
              </tbody>
            </Table>
          </div>
        </Card.Body>
        <Card.Footer className="bg-transparent text-muted">
          {filteredProducts.length} producto(s) | {sedes.length} sede(s)
        </Card.Footer>
      </Card>

      {/* Resumen por sede */}
      <Row className="mt-4 g-3">
        {sedes.map(sede => {
          const sedeStock = stockData[sede.id] || [];
          const totalItems = sedeStock.reduce((sum, s) => sum + (s.cantidad || 0), 0);
          const productsWithStock = sedeStock.filter(s => s.cantidad > 0).length;
          const lowStock = sedeStock.filter(s => s.cantidad > 0 && s.cantidad <= 5).length;
          
          return (
            <Col md={6} lg={4} xl={3} key={sede.id}>
              <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
                <Card.Body>
                  <div className="d-flex align-items-center mb-3">
                    <Building size={24} className="text-primary me-2" />
                    <h6 className="mb-0">{sede.nombre}</h6>
                  </div>
                  <Row className="g-2">
                    <Col xs={6}>
                      <div className="text-muted small">Total items</div>
                      <h5 className="mb-0">{totalItems}</h5>
                    </Col>
                    <Col xs={6}>
                      <div className="text-muted small">Productos</div>
                      <h5 className="mb-0">{productsWithStock}</h5>
                    </Col>
                    <Col xs={12} className="mt-2">
                      {lowStock > 0 && (
                        <Badge bg="warning" className="me-1">
                          {lowStock} con bajo stock
                        </Badge>
                      )}
                    </Col>
                  </Row>
                </Card.Body>
              </Card>
            </Col>
          );
        })}
      </Row>
    </Container>
  );
}

export default Stock;
