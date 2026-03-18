import { Container, Card, Row, Col, ListGroup, Badge, Spinner, Alert } from 'react-bootstrap';
import { CheckCircle } from 'react-bootstrap-icons';
import { useParams } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { getOrder } from '../services/api';

function OrderConfirmation() {
  const { orderID } = useParams();
  const [order, setOrder] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchOrder = async () => {
      try {
        setLoading(true);
        const data = await getOrder(orderID);
        setOrder(data);
      } catch (err) {
        setError('Error al cargar los detalles de la orden');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    if (orderID) {
      fetchOrder();
    }
  }, [orderID]);

  if (loading) {
    return (
      <Container className="py-5 text-center">
        <Spinner animation="border" />
        <p className="mt-3">Cargando detalles de la orden...</p>
      </Container>
    );
  }

  if (error || !order) {
    return (
      <Container className="py-5">
        <Alert variant="danger">{error || 'Orden no encontrada'}</Alert>
      </Container>
    );
  }

  return (
    <Container className="py-5">
      <Card className="shadow-sm">
        <Card.Body className="text-center py-5">
          <CheckCircle size={64} className="text-success mb-3" />
          <h2 className="mb-3">¡Pago Completado!</h2>
          <p className="text-muted mb-4">
            Gracias por tu compra. Tu orden ha sido procesada exitosamente.
          </p>
          <Badge bg="success" className="mb-3">
            {order.status}
          </Badge>
        </Card.Body>
      </Card>

      <Card className="mt-4 shadow-sm">
        <Card.Header className="bg-primary text-white">
          <h5 className="mb-0">Detalles de la Orden</h5>
        </Card.Header>
        <Card.Body>
          <Row className="mb-3">
            <Col md={6}>
              <p className="mb-2">
                <strong>id de Orden:</strong><br />
                <code>{order.paypal_order_id}</code>
              </p>
            </Col>
            <Col md={6}>
              <p className="mb-2">
                <strong>Fecha:</strong><br />
                {new Date(order.completed_at || order.created_at).toLocaleString('es-ES')}
              </p>
            </Col>
          </Row>

          {order.payer_email && (
            <Row className="mb-3">
              <Col md={6}>
                <p className="mb-2">
                  <strong>Email:</strong><br />
                  {order.payer_email}
                </p>
              </Col>
              {order.payer_name && (
                <Col md={6}>
                  <p className="mb-2">
                    <strong>Nombre:</strong><br />
                    {order.payer_name}
                  </p>
                </Col>
              )}
            </Row>
          )}

          <Row>
            <Col md={6}>
              <p className="mb-2">
                <strong>Total:</strong><br />
                <span className="h4 text-primary">
                  ${order.total_amount?.toFixed(2)} {order.currency}
                </span>
              </p>
            </Col>
          </Row>
        </Card.Body>
      </Card>

      {order.items && order.items.length > 0 && (
        <Card className="mt-4 shadow-sm">
          <Card.Header>
            <h5 className="mb-0">Productos</h5>
          </Card.Header>
          <Card.Body>
            <ListGroup variant="flush">
              {order.items.map((item) => (
                <ListGroup.Item key={item.id}>
                  <Row className="align-items-center">
                    <Col md={2}>
                      <img
                        src={item.product?.image_url || 'https://via.placeholder.com/80'}
                        alt={item.product_name}
                        className="img-fluid rounded"
                      />
                    </Col>
                    <Col md={5}>
                      <h6 className="mb-1">{item.product_name}</h6>
                      {item.product?.brand && (
                        <small className="text-muted">{item.product.brand}</small>
                      )}
                    </Col>
                    <Col md={2}>
                      <span>${item.price?.toFixed(2)}</span>
                    </Col>
                    <Col md={2}>
                      <span>Cant: {item.quantity}</span>
                    </Col>
                    <Col md={1} className="text-end">
                      <strong>${(item.price * item.quantity).toFixed(2)}</strong>
                    </Col>
                  </Row>
                </ListGroup.Item>
              ))}
            </ListGroup>
          </Card.Body>
        </Card>
      )}

      <div className="text-center mt-4 d-flex gap-2 justify-content-center">
        <a href="/my-orders" className="btn btn-outline-primary">
          Ver Mis Compras
        </a>
        <a href="/" className="btn btn-primary">
          Volver a la Tienda
        </a>
      </div>
    </Container>
  );
}

export default OrderConfirmation;
