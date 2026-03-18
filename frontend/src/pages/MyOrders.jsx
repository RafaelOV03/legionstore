import { useState, useEffect } from 'react';
import { Container, Card, Row, Col, Badge, Table, Spinner, Alert, Button } from 'react-bootstrap';
import { BoxSeam, Calendar, CurrencyDollar, Receipt, ChevronDown, ChevronUp } from 'react-bootstrap-icons';
import { getOrders } from '../services/api';
import { useNavigate } from 'react-router-dom';

function MyOrders() {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [expandedOrder, setExpandedOrder] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    loadOrders();
  }, []);

  const loadOrders = async () => {
    try {
      setLoading(true);
      const data = await getOrders();
      setOrders(data || []);
    } catch (err) {
      console.error('Error loading orders:', err);
      setError('Error al cargar tus órdenes');
    } finally {
      setLoading(false);
    }
  };

  const toggleOrderDetails = (orderID) => {
    setExpandedOrder(expandedOrder === orderID ? null : orderID);
  };

  const getOrderStats = () => {
    const totalOrders = orders.length;
    const completedOrders = orders.filter(o => o.status === 'COMPLETED').length;
    const totalSpent = orders
      .filter(o => o.status === 'COMPLETED')
      .reduce((sum, o) => sum + o.total_amount, 0);

    return { totalOrders, completedOrders, totalSpent };
  };

  const stats = getOrderStats();

  if (loading) {
    return (
      <Container className="py-5 text-center">
        <Spinner animation="border" variant="primary" />
        <p className="mt-3">Cargando tus órdenes...</p>
      </Container>
    );
  }

  if (error) {
    return (
      <Container className="py-5">
        <Alert variant="danger">{error}</Alert>
      </Container>
    );
  }

  return (
    <Container className="py-4">
      <h2 className="mb-4">Mis Compras</h2>

      {/* Estadísticas */}
      <Row className="mb-4">
        <Col md={4}>
          <Card className="text-center shadow-sm">
            <Card.Body>
              <Receipt size={32} className="text-primary mb-2" />
              <h3 className="mb-0">{stats.totalOrders}</h3>
              <small className="text-muted">Total de Órdenes</small>
            </Card.Body>
          </Card>
        </Col>
        <Col md={4}>
          <Card className="text-center shadow-sm">
            <Card.Body>
              <Badge bg="success" className="mb-2" style={{ fontSize: '1.5rem' }}>✓</Badge>
              <h3 className="mb-0">{stats.completedOrders}</h3>
              <small className="text-muted">Completadas</small>
            </Card.Body>
          </Card>
        </Col>
        <Col md={4}>
          <Card className="text-center shadow-sm">
            <Card.Body>
              <CurrencyDollar size={32} className="text-success mb-2" />
              <h3 className="mb-0">${stats.totalSpent.toFixed(2)}</h3>
              <small className="text-muted">Total Gastado</small>
            </Card.Body>
          </Card>
        </Col>
      </Row>

      {/* Lista de Órdenes */}
      {orders.length === 0 ? (
        <Card className="text-center py-5">
          <Card.Body>
            <BoxSeam size={64} className="text-muted mb-3" />
            <h4 className="text-muted">No tienes órdenes aún</h4>
            <p className="text-muted">¡Empieza a comprar y tus órdenes aparecerán aquí!</p>
            <Button variant="primary" onClick={() => navigate('/')}>
              Ir a la tienda
            </Button>
          </Card.Body>
        </Card>
      ) : (
        <div className="d-flex flex-column gap-3">
          {orders.map(order => (
            <Card key={order.id} className="shadow-sm hover-shadow">
              <Card.Header 
                className="d-flex justify-content-between align-items-center cursor-pointer"
                onClick={() => toggleOrderDetails(order.id)}
                style={{ cursor: 'pointer' }}
              >
                <div className="d-flex align-items-center gap-3">
                  <Receipt size={24} className="text-primary" />
                  <div>
                    <h6 className="mb-0">
                      Orden #{order.id}
                      <Badge bg={order.status === 'COMPLETED' ? 'success' : 'warning'} className="ms-2">
                        {order.status}
                      </Badge>
                    </h6>
                    <small className="text-muted">
                      <Calendar className="me-1" />
                      {new Date(order.completed_at || order.created_at).toLocaleDateString('es-ES', {
                        day: '2-digit',
                        month: 'long',
                        year: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit'
                      })}
                    </small>
                  </div>
                </div>
                <div className="d-flex align-items-center gap-3">
                  <div className="text-end">
                    <h5 className="mb-0 text-success">${order.total_amount?.toFixed(2)}</h5>
                    <small className="text-muted">{order.items?.length || 0} productos</small>
                  </div>
                  {expandedOrder === order.id ? <ChevronUp /> : <ChevronDown />}
                </div>
              </Card.Header>

              {expandedOrder === order.id && (
                <Card.Body>
                  <Row className="mb-3">
                    <Col md={6}>
                      <p className="mb-1">
                        <strong>id PayPal:</strong><br />
                        <small className="text-muted">{order.paypal_order_id}</small>
                      </p>
                    </Col>
                    <Col md={6}>
                      <p className="mb-1">
                        <strong>Cliente:</strong> {order.payer_name || 'N/A'}
                      </p>
                      <p className="mb-0">
                        <strong>Email:</strong> {order.payer_email || 'N/A'}
                      </p>
                    </Col>
                  </Row>

                  <h6 className="mb-3">Productos Comprados</h6>
                  <Table responsive striped bordered hover size="sm">
                    <thead>
                      <tr>
                        <th style={{ width: '80px' }}>Imagen</th>
                        <th>Producto</th>
                        <th>Marca</th>
                        <th>Precio</th>
                        <th>Cantidad</th>
                        <th>Subtotal</th>
                      </tr>
                    </thead>
                    <tbody>
                      {order.items?.map((item, idx) => (
                        <tr key={idx}>
                          <td>
                            <img
                              src={item.product?.image_url || 'https://via.placeholder.com/50'}
                              alt={item.product_name}
                              style={{ width: '50px', height: '50px', objectFit: 'cover', borderRadius: '4px' }}
                            />
                          </td>
                          <td>
                            <strong>{item.product_name}</strong>
                          </td>
                          <td>{item.product?.brand || 'N/A'}</td>
                          <td>${item.price?.toFixed(2)}</td>
                          <td>
                            <Badge bg="secondary">{item.quantity}</Badge>
                          </td>
                          <td>
                            <strong>${(item.price * item.quantity).toFixed(2)}</strong>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                    <tfoot>
                      <tr>
                        <td colSpan="5" className="text-end">
                          <strong>Total:</strong>
                        </td>
                        <td>
                          <strong className="text-success">${order.total_amount?.toFixed(2)}</strong>
                        </td>
                      </tr>
                    </tfoot>
                  </Table>
                </Card.Body>
              )}
            </Card>
          ))}
        </div>
      )}
    </Container>
  );
}

export default MyOrders;
