import { useState } from 'react';
import { Container, Row, Col, Card, Button, Carousel, ListGroup, Alert } from 'react-bootstrap';
import { useParams, useNavigate } from 'react-router-dom';
import { useEffect } from 'react';
import { getProduct, addToCart, getRelatedProducts } from '../services/api';
import { Cart, CheckCircle, ExclamationTriangle, ArrowLeft } from 'react-bootstrap-icons';
import ProductCard from '../components/ProductCard';

function ProductDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [product, setProduct] = useState(null);
  const [relatedProducts, setRelatedProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [quantity, setQuantity] = useState(1);
  const [alert, setAlert] = useState(null);
  const [selectedImageIndex, setSelectedImageIndex] = useState(0);

  useEffect(() => {
    loadProduct();
    loadRelatedProducts();
    window.scrollTo(0, 0);
  }, [id]);

  const loadProduct = async () => {
    try {
      setLoading(true);
      const data = await getProduct(id);
      setProduct(data);
    } catch (error) {
      console.error('Error loading product:', error);
      showAlert('Error al cargar el producto', 'danger');
    } finally {
      setLoading(false);
    }
  };

  const loadRelatedProducts = async () => {
    try {
      const data = await getRelatedProducts(id);
      setRelatedProducts(data || []);
    } catch (error) {
      console.error('Error loading related products:', error);
    }
  };

  const handleAddToCart = async () => {
    try {
      await addToCart(product.id, quantity);
      showAlert('Producto agregado al carrito', 'success');
    } catch (error) {
      console.error('Error adding to cart:', error);
      showAlert('Error al agregar al carrito', 'danger');
    }
  };

  const showAlert = (message, variant) => {
    setAlert({ message, variant });
    setTimeout(() => setAlert(null), 3000);
  };

  if (loading) {
    return (
      <>
        {alert && (
          <Alert 
            variant={alert.variant} 
            dismissible 
            onClose={() => setAlert(null)}
            style={{
              position: 'fixed',
              top: '20px',
              right: '20px',
              zIndex: 9999,
              minWidth: '300px',
              boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
              animation: 'slideIn 0.3s ease-out'
            }}
          >
            {alert.message}
          </Alert>
        )}
        <style>
          {`
            @keyframes slideIn {
              from {
                transform: translateX(400px);
                opacity: 0;
              }
              to {
                transform: translateX(0);
                opacity: 1;
              }
            }
          `}
        </style>
        <Container className="py-5">
          <div className="text-center">
            <p>Cargando producto...</p>
          </div>
        </Container>
      </>
    );
  }

  if (!product) {
    return (
      <Container className="py-5">
        <Alert variant="danger">Producto no encontrado</Alert>
        <Button variant="primary" onClick={() => navigate('/')}>
          Volver a la tienda
        </Button>
      </Container>
    );
  }

  // Parsear images si es string JSON
  let imagesArray = [];
  if (product.images) {
    try {
      imagesArray = typeof product.images === 'string' ? JSON.parse(product.images) : product.images;
    } catch (e) {
      imagesArray = [];
    }
  }

  const images = imagesArray && imagesArray.length > 0 
    ? imagesArray 
    : product.image_url 
    ? [product.image_url] 
    : ['https://via.placeholder.com/600x400'];

  return (
    <>
      {alert && (
        <Alert 
          variant={alert.variant} 
          dismissible 
          onClose={() => setAlert(null)}
          style={{
            position: 'fixed',
            top: '20px',
            right: '20px',
            zIndex: 9999,
            minWidth: '300px',
            boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
            animation: 'slideIn 0.3s ease-out'
          }}
        >
          {alert.message}
        </Alert>
      )}

      <style>
        {`
          @keyframes slideIn {
            from {
              transform: translateX(400px);
              opacity: 0;
            }
            to {
              transform: translateX(0);
              opacity: 1;
            }
          }
        `}
      </style>

      <Container className="py-4" style={{ animation: 'fadeIn 0.5s ease-out' }}>
      <Button 
        variant="outline-primary" 
        className="mb-4" 
        onClick={() => navigate('/')}
        style={{ 
          transition: 'all 0.3s ease',
          animation: 'slideInLeft 0.5s ease-out'
        }}
      >
        <ArrowLeft className="me-2" />
        Volver a la tienda
      </Button>

      <Row>
        {/* Galería de Imágenes */}
        <Col lg={6} className="mb-4 d-flex flex-column justify-content-center" style={{ animation: 'fadeInUp 0.6s ease-out' }}>
            <h1 className="mb-auto mx-auto" style={{ 
                color: 'var(--text-primary)',
                fontWeight: 'bold',
                animation: 'fadeInUp 0.6s ease-out 0.3s both'
            }}>
                {product.name}
          </h1>
          <Card style={{ 
            background: 'var(--bg-card)',
            border: '1px solid var(--border-color)',
            borderRadius: '15px',
            overflow: 'hidden',
            boxShadow: '0 10px 30px rgba(0,0,0,0.3)'
          }}>
            <Carousel 
              activeIndex={selectedImageIndex}
              onSelect={(index) => setSelectedImageIndex(index)}
              interval={null}
            >
              {images.map((image, index) => (
                <Carousel.Item key={index}>
                  <div style={{
                    height: '500px',
                    background: 'linear-gradient(135deg, var(--bg-darker) 0%, var(--bg-dark) 100%)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    padding: '20px'
                  }}>
                    <img
                      src={image}
                      alt={`${product.name} - Imagen ${index + 1}`}
                      style={{ 
                        maxHeight: '100%',
                        maxWidth: '100%',
                        objectFit: 'contain',
                        transition: 'transform 0.3s ease'
                      }}
                      onMouseOver={(e) => e.target.style.transform = 'scale(1.05)'}
                      onMouseOut={(e) => e.target.style.transform = 'scale(1)'}
                    />
                  </div>
                </Carousel.Item>
              ))}
            </Carousel>
          </Card>

          {/* Miniaturas */}
          {images.length > 1 && (
            <Row className="mt-3 g-2">
              {images.map((image, index) => (
                <Col xs={3} key={index}>
                  <div
                    onClick={() => setSelectedImageIndex(index)}
                    style={{
                      cursor: 'pointer',
                      border: selectedImageIndex === index 
                        ? '3px solid var(--primary-color)' 
                        : '2px solid var(--border-color)',
                      borderRadius: '10px',
                      overflow: 'hidden',
                      transition: 'all 0.3s ease',
                      transform: selectedImageIndex === index ? 'scale(1.05)' : 'scale(1)',
                      boxShadow: selectedImageIndex === index 
                        ? '0 5px 15px rgba(13, 110, 253, 0.4)' 
                        : 'none'
                    }}
                  >
                    <img
                      src={image}
                      alt={`Miniatura ${index + 1}`}
                      style={{ 
                        height: '80px',
                        width: '100%',
                        objectFit: 'cover',
                        display: 'block'
                      }}
                    />
                  </div>
                </Col>
              ))}
            </Row>
          )}
        </Col>

        {/* Información del Producto */}
        <Col lg={6} style={{ animation: 'fadeInUp 0.6s ease-out 0.2s both' }}>
          
          <div className="mb-3" style={{ animation: 'fadeIn 0.6s ease-out 0.4s both' }}>
            <span style={{
              background: 'rgba(13, 202, 240, 0.2)',
              color: '#0dcaf0',
              padding: '0.5em 1em',
              borderRadius: '20px',
              fontSize: '0.9rem',
              fontWeight: '600',
              marginRight: '10px',
              border: '1px solid rgba(13, 202, 240, 0.3)'
            }}>
              {product.category}
            </span>
            {product.brand && (
              <span style={{
                background: 'rgba(108, 117, 125, 0.2)',
                color: '#9fa8da',
                padding: '0.5em 1em',
                borderRadius: '20px',
                fontSize: '0.9rem',
                fontWeight: '600',
                border: '1px solid rgba(108, 117, 125, 0.3)'
              }}>
                {product.brand}
              </span>
            )}
          </div>

          <div className="mb-4" style={{ animation: 'fadeInUp 0.6s ease-out 0.5s both' }}>
            <h3 style={{ 
              color: 'var(--primary-color)',
              fontSize: '2.5rem',
              fontWeight: 'bold'
            }}>
              ${product.price.toFixed(2)}
            </h3>
          </div>

          {/* Stock */}
          <Card className="mb-4" style={{
            background: 'var(--bg-card)',
            border: '1px solid var(--border-color)',
            animation: 'fadeInUp 0.6s ease-out 0.6s both'
          }}>
            <Card.Body>
              <Row className="align-items-center">
                <Col>
                  <strong style={{ color: 'var(--text-primary)' }}>Disponibilidad:</strong>
                </Col>
                <Col className="text-end">
                  {product.stock_quantity > 0 ? (
                    <span style={{
                      background: 'rgba(25, 135, 84, 0.2)',
                      color: '#75b798',
                      padding: '0.5em 1em',
                      borderRadius: '20px',
                      fontSize: '0.9rem',
                      fontWeight: '600',
                      border: '1px solid rgba(25, 135, 84, 0.3)',
                      display: 'inline-block'
                    }}>
                      <CheckCircle className="me-1" />
                      {product.stock_quantity} en stock
                    </span>
                  ) : (
                    <span style={{
                      background: 'rgba(220, 53, 69, 0.2)',
                      color: '#ea868f',
                      padding: '0.5em 1em',
                      borderRadius: '20px',
                      fontSize: '0.9rem',
                      fontWeight: '600',
                      border: '1px solid rgba(220, 53, 69, 0.3)',
                      display: 'inline-block'
                    }}>
                      <ExclamationTriangle className="me-1" />
                      Agotado
                    </span>
                  )}
                </Col>
              </Row>
            </Card.Body>
          </Card>

          {/* Descripción */}
          {product.description && (
            <Card className="mb-4" style={{
              background: 'var(--bg-card)',
              border: '1px solid var(--border-color)',
              animation: 'fadeInUp 0.6s ease-out 0.7s both'
            }}>
              <Card.Header style={{
                background: 'transparent',
                borderBottom: '1px solid var(--border-color)',
                color: 'var(--text-primary)'
              }}>
                <strong>Descripción</strong>
              </Card.Header>
              <Card.Body>
                <p className="mb-0" style={{ 
                  color: 'var(--text-secondary)',
                  lineHeight: '1.8'
                }}>
                  {product.description}
                </p>
              </Card.Body>
            </Card>
          )}

          {/* Agregar al carrito */}
          {product.stock_quantity > 0 && (
            <Card style={{
              background: 'var(--bg-card)',
              border: '2px solid var(--primary-color)',
              boxShadow: '0 5px 20px rgba(13, 110, 253, 0.3)',
              animation: 'fadeInUp 0.6s ease-out 0.8s both'
            }}>
              <Card.Body>
                <Row className="align-items-end">
                  <Col md={4}>
                    <label className="mb-2" style={{ 
                      color: 'var(--text-primary)',
                      fontWeight: '600'
                    }}>
                      Cantidad:
                    </label>
                    <input
                      type="number"
                      className="form-control"
                      min="1"
                      max={product.stock_quantity}
                      value={quantity}
                      onChange={(e) => setQuantity(Math.max(1, Math.min(product.stock_quantity, parseInt(e.target.value) || 1)))}
                      style={{
                        background: 'var(--bg-darker)',
                        border: '1px solid var(--border-color)',
                        color: 'var(--text-primary)',
                        fontSize: '1.1rem',
                        fontWeight: '600',
                        textAlign: 'center'
                      }}
                    />
                  </Col>
                  <Col md={8}>
                    <Button 
                      variant="primary" 
                      size="lg" 
                      className="w-100"
                      onClick={handleAddToCart}
                      style={{
                        background: 'linear-gradient(135deg, #0d6efd 0%, #0a58ca 100%)',
                        border: 'none',
                        padding: '12px',
                        fontWeight: '600',
                        fontSize: '1.1rem',
                        transition: 'all 0.3s ease'
                      }}
                      onMouseOver={(e) => {
                        e.target.style.transform = 'translateY(-3px)';
                        e.target.style.boxShadow = '0 10px 25px rgba(13, 110, 253, 0.4)';
                      }}
                      onMouseOut={(e) => {
                        e.target.style.transform = 'translateY(0)';
                        e.target.style.boxShadow = 'none';
                      }}
                    >
                      <Cart className="me-2" />
                      Agregar al Carrito
                    </Button>
                  </Col>
                </Row>
              </Card.Body>
            </Card>
          )}

          {/* Especificaciones */}
          <Card className="mt-4" style={{
            background: 'var(--bg-card)',
            border: '1px solid var(--border-color)',
            animation: 'fadeInUp 0.6s ease-out 0.9s both'
          }}>
            <Card.Header style={{
              background: 'transparent',
              borderBottom: '1px solid var(--border-color)',
              color: 'var(--text-primary)'
            }}>
              <strong>Especificaciones</strong>
            </Card.Header>
            <ListGroup variant="flush">
              <ListGroup.Item style={{
                background: 'transparent',
                borderColor: 'var(--border-color)',
                color: 'var(--text-primary)'
              }}>
                <Row>
                  <Col><strong>id del Producto:</strong></Col>
                  <Col className="text-end" style={{ color: 'var(--text-secondary)' }}>
                    {product.id}
                  </Col>
                </Row>
              </ListGroup.Item>
              <ListGroup.Item style={{
                background: 'transparent',
                borderColor: 'var(--border-color)',
                color: 'var(--text-primary)'
              }}>
                <Row>
                  <Col><strong>Categoría:</strong></Col>
                  <Col className="text-end" style={{ color: 'var(--text-secondary)' }}>
                    {product.category}
                  </Col>
                </Row>
              </ListGroup.Item>
              {product.brand && (
                <ListGroup.Item style={{
                  background: 'transparent',
                  borderColor: 'var(--border-color)',
                  color: 'var(--text-primary)'
                }}>
                  <Row>
                    <Col><strong>Marca:</strong></Col>
                    <Col className="text-end" style={{ color: 'var(--text-secondary)' }}>
                      {product.brand}
                    </Col>
                  </Row>
                </ListGroup.Item>
              )}
              <ListGroup.Item style={{
                background: 'transparent',
                borderColor: 'var(--border-color)',
                color: 'var(--text-primary)'
              }}>
                <Row>
                  <Col><strong>Precio:</strong></Col>
                  <Col className="text-end">
                    <strong style={{ color: 'var(--primary-color)', fontSize: '1.2rem' }}>
                      ${product.price.toFixed(2)}
                    </strong>
                  </Col>
                </Row>
              </ListGroup.Item>
            </ListGroup>
          </Card>
        </Col>
      </Row>

      {/* Productos Relacionados */}
      {relatedProducts.length > 0 && (
        <div className="mt-5" style={{ animation: 'fadeInUp 0.6s ease-out 1s both' }}>
          <h3 className="mb-4" style={{
            color: 'var(--text-primary)',
            fontWeight: 'bold'
          }}>
            Productos Relacionados
          </h3>
          <Row xs={1} md={2} lg={4} className="g-4">
            {relatedProducts.map(relatedProduct => (
              <Col key={relatedProduct.id}>
                <ProductCard 
                  product={relatedProduct} 
                  onAddToCart={async (p) => {
                    try {
                      await addToCart(p.id, 1);
                      showAlert(`${p.name} agregado al carrito`, 'success');
                    } catch (error) {
                      showAlert('Error al agregar al carrito', 'danger');
                    }
                  }} 
                />
              </Col>
            ))}
          </Row>
        </div>
      )}
    </Container>
    </>
  );
}

export default ProductDetail;
