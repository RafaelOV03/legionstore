import { Container, Row, Col, ListGroup, Button, Alert, Card, Spinner } from 'react-bootstrap';
import { Trash, Plus, Dash } from 'react-bootstrap-icons';
import { useCart } from '../hooks/useCart';
import { updateCartItem, removeFromCart, clearCart, getPayPalConfig, createOrder, captureOrder } from '../services/resourceApi';
import { useState, useEffect } from 'react';
import { PayPalScriptProvider, PayPalButtons } from '@paypal/react-paypal-js';
import { useNavigate } from 'react-router-dom';

function Cart() {
  const { cart, loading, refreshCart, cartTotal } = useCart();
  const [alert, setAlert] = useState(null);
  const [paypalClientId, setPaypalClientId] = useState(null);
  const [processingPayment, setProcessingPayment] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchPayPalConfig = async () => {
      try {
        const config = await getPayPalConfig();
        console.log('PayPal Client id loaded:', config.clientId);
        setPaypalClientId(config.clientId);
      } catch (error) {
        console.error('Error fetching PayPal config:', error);
        setAlert({ message: 'Error al cargar PayPal. Verifica el servidor.', variant: 'warning' });
      }
    };
    fetchPayPalConfig();
  }, []);

  const handleUpdateQuantity = async (itemId, newQuantity) => {
    if (newQuantity < 1) return;
    
    try {
      await updateCartItem(itemId, newQuantity);
      await refreshCart();
    } catch (error) {
      console.error('Error updating cart:', error);
      showAlert('Error al actualizar el carrito', 'danger');
    }
  };

  const handleRemoveItem = async (itemId) => {
    try {
      await removeFromCart(itemId);
      await refreshCart();
      showAlert('Producto eliminado del carrito', 'success');
    } catch (error) {
      console.error('Error removing item:', error);
      showAlert('Error al eliminar el producto', 'danger');
    }
  };

  const handleClearCart = async () => {
    if (window.confirm('¿Estás seguro de vaciar el carrito?')) {
      try {
        await clearCart();
        await refreshCart();
        showAlert('Carrito vaciado', 'success');
      } catch (error) {
        console.error('Error clearing cart:', error);
        showAlert('Error al vaciar el carrito', 'danger');
      }
    }
  };

  const showAlert = (message, variant) => {
    setAlert({ message, variant });
    setTimeout(() => setAlert(null), 3000);
  };

  if (loading) {
    return (
      <Container>
        <div className="text-center py-5">
          <p>Cargando carrito...</p>
        </div>
      </Container>
    );
  }

  return (
    <Container>
      {alert && (
        <Alert variant={alert.variant} dismissible onClose={() => setAlert(null)}>
          {alert.message}
        </Alert>
      )}

      <h2 className="mb-4">Mi Carrito de Compras</h2>

      {!cart.items || cart.items.length === 0 ? (
        <Card className="text-center py-5">
          <Card.Body>
            <h4 className="text-muted">Tu carrito está vacío</h4>
            <p className="text-muted">¡Agrega algunos productos para comenzar!</p>
            <Button variant="primary" href="/">
              Ir a la tienda
            </Button>
          </Card.Body>
        </Card>
      ) : (
        <Row>
          <Col lg={8}>
            <ListGroup>
              {cart.items.map(item => (
                <ListGroup.Item key={item.id} className="py-3">
                  <Row className="align-items-center">
                    <Col md={2}>
                      <img
                        src={item.product?.image_url || 'https://via.placeholder.com/100'}
                        alt={item.product?.name}
                        className="img-fluid rounded"
                      />
                    </Col>
                    <Col md={4}>
                      <h6 className="mb-1">{item.product?.name}</h6>
                      <small className="text-muted">{item.product?.brand}</small>
                    </Col>
                    <Col md={2}>
                      <strong>${item.product?.price.toFixed(2)}</strong>
                    </Col>
                    <Col md={3}>
                      <div className="d-flex align-items-center gap-2">
                        <Button
                          variant="outline-secondary"
                          size="sm"
                          onClick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
                          disabled={item.quantity <= 1}
                        >
                          <Dash />
                        </Button>
                        <span className="mx-2">{item.quantity}</span>
                        <Button
                          variant="outline-secondary"
                          size="sm"
                          onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                          disabled={item.quantity >= item.product?.stock_quantity}
                        >
                          <Plus />
                        </Button>
                      </div>
                    </Col>
                    <Col md={1}>
                      <Button
                        variant="outline-danger"
                        size="sm"
                        onClick={() => handleRemoveItem(item.id)}
                      >
                        <Trash />
                      </Button>
                    </Col>
                  </Row>
                </ListGroup.Item>
              ))}
            </ListGroup>

            <div className="mt-3">
              <Button variant="outline-danger" onClick={handleClearCart}>
                Vaciar Carrito
              </Button>
            </div>
          </Col>

          <Col lg={4}>
            <Card className="sticky-top" style={{ top: '20px' }}>
              <Card.Header className="bg-primary text-white">
                <h5 className="mb-0">Resumen del Pedido</h5>
              </Card.Header>
              <Card.Body>
                <div className="d-flex justify-content-between mb-2">
                  <span>Subtotal:</span>
                  <strong>${cartTotal.toFixed(2)}</strong>
                </div>
                <div className="d-flex justify-content-between mb-2">
                  <span>Envío:</span>
                  <strong>Gratis</strong>
                </div>
                <hr />
                <div className="d-flex justify-content-between mb-3">
                  <h5>Total:</h5>
                  <h5 className="text-primary">${cartTotal.toFixed(2)}</h5>
                </div>
                
                {processingPayment && (
                  <div className="text-center mb-3">
                    <Spinner animation="border" size="sm" /> Procesando pago...
                  </div>
                )}

                {paypalClientId ? (
                  <div style={{
                    background: 'white',
                    padding: '15px',
                    borderRadius: '12px',
                    boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)'
                  }}>
                    <PayPalScriptProvider options={{ 
                      clientId: paypalClientId,
                      currency: "USD"
                    }}>
                      <PayPalButtons
                        disabled={processingPayment}
                        style={{ 
                          layout: "vertical",
                          color: "gold",
                          shape: "rect",
                          label: "paypal"
                        }}
                        createOrder={async () => {
                          try {
                            const response = await createOrder();
                            return response.id;
                          } catch (error) {
                            showAlert('Error al crear la orden: ' + error.message, 'danger');
                            throw error;
                          }
                        }}
                        onApprove={async (data) => {
                          try {
                            setProcessingPayment(true);
                            const response = await captureOrder(data.orderID);
                            showAlert('¡Pago completado exitosamente!', 'success');
                            await refreshCart();
                            setTimeout(() => {
                              navigate(`/order/${data.orderID}`);
                            }, 1500);
                          } catch (error) {
                            showAlert('Error al procesar el pago: ' + error.message, 'danger');
                            setProcessingPayment(false);
                          }
                        }}
                        onError={(err) => {
                          console.error('PayPal error:', err);
                          showAlert('Error con PayPal. Por favor intenta de nuevo.', 'danger');
                          setProcessingPayment(false);
                        }}
                      />
                    </PayPalScriptProvider>
                  </div>
                ) : (
                  <div className="text-center py-3">
                    <Spinner animation="border" size="sm" />
                  </div>
                )}
                
                <Button variant="outline-secondary" className="w-100 mt-2" href="/">
                  Continuar Comprando
                </Button>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      )}
    </Container>
  );
}

export default Cart;
