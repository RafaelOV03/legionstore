import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Form, Button, Table, Badge, Spinner, Alert, InputGroup } from 'react-bootstrap';
import { Search, Tools, Check, BoxSeam } from 'react-bootstrap-icons';
import * as api from '../services/resourceApi';

function Compatibilidad() {
  const [productos, setProductos] = useState([]);
  const [insumos, setInsumos] = useState([]);
  const [compatibilidades, setCompatibilidades] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchType, setSearchType] = useState('producto'); // producto o insumo
  const [selectedItem, setSelectedItem] = useState(null);
  const [results, setResults] = useState([]);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [productosData, insumosData] = await Promise.all([
        api.getProducts(),
        api.getInsumos()
      ]);
      setProductos(productosData || []);
      setInsumos(insumosData || []);
    } catch (err) {
      setError('Error cargando datos');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async () => {
    if (!selectedItem) return;
    
    try {
      setLoading(true);
      setError(null);
      // Buscar compatibilidades del producto seleccionado
      const data = await api.buscarCompatibles(selectedItem);
      setResults(data?.compatibles || []);
    } catch (err) {
      setError('Error buscando compatibles');
      console.error(err);
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  if (loading && !productos.length) {
    return (
      <Container className="py-5 text-center">
        <Spinner animation="border" variant="primary" />
      </Container>
    );
  }

  return (
    <Container fluid className="py-4">
      <Row className="mb-4">
        <Col>
          <h2><Tools className="me-2" />Asistente de Compatibilidad</h2>
          <p className="text-muted">Buscar productos e insumos compatibles</p>
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}

      <Row>
        <Col lg={4}>
          <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
            <Card.Header className="bg-transparent">
              <h6 className="mb-0">Buscar Compatibilidad</h6>
            </Card.Header>
            <Card.Body>
              <Form.Group className="mb-3">
                <Form.Label>Seleccionar Producto</Form.Label>
                <Form.Select
                  value={selectedItem || ''}
                  onChange={(e) => setSelectedItem(e.target.value)}
                >
                  <option value="">Seleccionar...</option>
                  {productos.map(p => (
                    <option key={p.id} value={p.id}>
                      {p.name} - {p.category}
                    </option>
                  ))}
                </Form.Select>
              </Form.Group>

              <Button 
                variant="primary" 
                className="w-100"
                onClick={handleSearch}
                disabled={!selectedItem || loading}
              >
                <Search className="me-1" /> Buscar Compatibles
              </Button>
            </Card.Body>
          </Card>

          {/* Info del item seleccionado */}
          {selectedItem && (
            <Card className="mt-3" style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body>
                {(() => {
                  const prod = productos.find(p => p.id.toString() === selectedItem);
                  return prod ? (
                    <>
                      <div className="d-flex align-items-center mb-2">
                        <BoxSeam className="text-primary me-2" />
                        <strong>{prod.name}</strong>
                      </div>
                      <Badge bg="secondary" className="me-2">{prod.category}</Badge>
                      <Badge bg="info">{prod.brand}</Badge>
                      {prod.description && (
                        <p className="text-muted small mt-2 mb-0">{prod.description}</p>
                      )}
                    </>
                  ) : null;
                })()}
              </Card.Body>
            </Card>
          )}
        </Col>

        <Col lg={8}>
          <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
            <Card.Header className="bg-transparent">
              <h6 className="mb-0">
                {searchType === 'producto' 
                  ? 'Insumos Compatibles' 
                  : 'Productos Compatibles'
                }
                {results.length > 0 && ` (${results.length})`}
              </h6>
            </Card.Header>
            <Card.Body className="p-0">
              {loading ? (
                <div className="text-center py-5">
                  <Spinner animation="border" variant="primary" size="sm" />
                </div>
              ) : results.length === 0 ? (
                <div className="text-center py-5 text-muted">
                  {selectedItem 
                    ? 'No se encontraron compatibilidades'
                    : 'Selecciona un elemento y busca compatibilidades'
                  }
                </div>
              ) : (
                <Table hover className="mb-0">
                  <thead>
                    <tr>
                      <th>Producto Compatible</th>
                      <th>Marca</th>
                      <th>Categoría</th>
                      <th>Precio</th>
                      <th>Tipo Match</th>
                      <th>Notas</th>
                    </tr>
                  </thead>
                  <tbody>
                    {results.map((item, index) => (
                      <tr key={index}>
                        <td>
                          <div className="d-flex align-items-center">
                            <BoxSeam className="text-muted me-2" size={16} />
                            <strong>{item.name}</strong>
                          </div>
                        </td>
                        <td>{item.brand || '-'}</td>
                        <td>
                          <Badge bg="secondary">
                            {item.category || 'N/A'}
                          </Badge>
                        </td>
                        <td>Bs. {item.precio?.toFixed(2) || '0.00'}</td>
                        <td>
                          <Badge 
                            bg={item.tipo_match === 'directo' ? 'primary' : 
                               item.tipo_match === 'misma_categoria' ? 'info' : 'secondary'}
                          >
                            {item.tipo_match === 'directo' ? 'Directo' :
                             item.tipo_match === 'misma_categoria' ? 'Misma Categoría' : 'Mismo Fabricante'}
                          </Badge>
                        </td>
                        <td>
                          <small className="text-muted">{item.notas || '-'}</small>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </Table>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Container>
  );
}

export default Compatibilidad;
