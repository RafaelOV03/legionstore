import { useState, useEffect } from 'react';
import { Container, Row, Col, Form, Alert, Carousel, Card, Button, Badge } from 'react-bootstrap';
import ProductCard from '../components/ProductCard';
import { getProducts, addToCart } from '../services/resourceApi';
import { useCart } from '../hooks/useCart';
import { ChevronRight } from 'react-bootstrap-icons';

function Home() {
  const [products, setProducts] = useState([]);
  const [filteredProducts, setFilteredProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [searchTerm, setSearchTerm] = useState('');
  const [alert, setAlert] = useState(null);
  const { refreshCart } = useCart();

  useEffect(() => {
    loadProducts();
  }, []);

  useEffect(() => {
    filterProducts();
  }, [products, selectedCategory, searchTerm]);

  const loadProducts = async () => {
    try {
      const data = await getProducts();
      setProducts(data || []);
      
      // Extraer categorías únicas
      const uniqueCategories = [...new Set(data.map(p => p.category))];
      setCategories(uniqueCategories);
    } catch (error) {
      console.error('Error loading products:', error);
      showAlert('Error al cargar los productos', 'danger');
    }
  };

  const filterProducts = () => {
    let filtered = products;

    if (selectedCategory !== 'all') {
      filtered = filtered.filter(p => p.category === selectedCategory);
    }

    if (searchTerm) {
      filtered = filtered.filter(p =>
        p.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        p.brand.toLowerCase().includes(searchTerm.toLowerCase()) ||
        p.description.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    setFilteredProducts(filtered);
  };

  const handleAddToCart = async (product) => {
    try {
      await addToCart(product.id, 1);
      await refreshCart();
      showAlert(`${product.name} agregado al carrito`, 'success');
    } catch (error) {
      console.error('Error adding to cart:', error);
      showAlert('Error al agregar al carrito', 'danger');
    }
  };

  const showAlert = (message, variant) => {
    setAlert({ message, variant });
    setTimeout(() => setAlert(null), 3000);
  };

  return (
    <Container>
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

      {/* Hero Carousel */}
      <Carousel className="mb-5" style={{ 
        borderRadius: '15px', 
        overflow: 'hidden',
        boxShadow: '0 10px 40px rgba(13, 110, 253, 0.3)',
        animation: 'fadeIn 0.8s ease-out'
      }}>
        <Carousel.Item>
          <div style={{
            height: '400px',
            position: 'relative',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
          }}>
            <img
              src="/carrousel1.jpg"
              alt="Bienvenido a Smartech"
              style={{ 
                position: 'absolute',
                width: '100%',
                height: '100%',
                objectFit: 'cover',
                filter: 'brightness(0.7)'
              }}
            />
            <div className="text-center text-white" style={{ position: 'relative', zIndex: 2 }}>
              <h1 className="display-3 fw-bold mb-3" style={{ textShadow: '0 4px 20px rgba(0,0,0,0.8)' }}>
                Bienvenido a Smartech
              </h1>
              <p className="lead mb-4" style={{ fontSize: '1.5rem', textShadow: '0 2px 10px rgba(0,0,0,0.8)' }}>
                La tecnología que necesitas al mejor precio
              </p>
            </div>
          </div>
        </Carousel.Item>
        
        <Carousel.Item>
          <div style={{
            height: '400px',
            position: 'relative',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
          }}>
            <img
              src="/carrousel2.jpg"
              alt="Nuevos Productos"
              style={{ 
                position: 'absolute',
                width: '100%',
                height: '100%',
                objectFit: 'cover',
                filter: 'brightness(0.7)'
              }}
            />
            <div className="text-center text-white" style={{ position: 'relative', zIndex: 2 }}>
              <h2 className="display-4 fw-bold mb-3" style={{ textShadow: '0 4px 20px rgba(0,0,0,0.8)' }}>
                Nuevos Productos
              </h2>
              <p className="lead mb-4" style={{ fontSize: '1.3rem', textShadow: '0 2px 10px rgba(0,0,0,0.8)' }}>
                Descubre las últimas novedades en tecnología
              </p>
            </div>
          </div>
        </Carousel.Item>
        
        <Carousel.Item>
          <div style={{
            height: '400px',
            position: 'relative',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
          }}>
            <img
              src="/carrousel3.jpg"
              alt="Ofertas Especiales"
              style={{ 
                position: 'absolute',
                width: '100%',
                height: '100%',
                objectFit: 'cover',
                filter: 'brightness(0.7)'
              }}
            />
            <div className="text-center text-white" style={{ position: 'relative', zIndex: 2 }}>
              <h2 className="display-4 fw-bold mb-3" style={{ textShadow: '0 4px 20px rgba(0,0,0,0.8)' }}>
                Ofertas Especiales
              </h2>
              <p className="lead mb-4" style={{ fontSize: '1.3rem', textShadow: '0 2px 10px rgba(0,0,0,0.8)' }}>
                Los mejores precios del mercado
              </p>
            </div>
          </div>
        </Carousel.Item>
      </Carousel>

      {/* Search Bar */}
      <Row className="mb-5">
        <Col md={0} className="mx-auto">
          <Card style={{ 
            background: 'var(--bg-card)',
            border: '1px solid var(--border-color)',
            boxShadow: '0 4px 15px rgba(0,0,0,0.2)'
          }}>
            <Card.Body>
              <Row>
                <Col md={9}>
                  <Form.Control
                    type="text"
                    placeholder="Buscar productos..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        e.preventDefault();
                        setSearchTerm(e.target.value);
                      }
                    }}
                    size="lg"
                    style={{
                      background: 'var(--bg-darker)',
                      border: '1px solid var(--border-color)',
                      color: 'var(--text-primary)'
                    }}
                  />
                </Col>
                <Col md={3}>
                  <Form.Select
                    value={selectedCategory}
                    onChange={(e) => setSelectedCategory(e.target.value)}
                    size="lg"
                    style={{
                      border: '1px solid var(--border-color)',
                      color: 'var(--text-primary)'
                    }}
                  >
                    <option value="all">Todas las categorías</option>
                    {categories.map(cat => (
                      <option key={cat} value={cat}>{cat}</option>
                    ))}
                  </Form.Select>
                </Col>
              </Row>
            </Card.Body>
          </Card>
        </Col>
      </Row>

      {/* Products by Category */}
      {selectedCategory === 'all' && !searchTerm ? (
        <>
          {categories.map(category => {
            const categoryProducts = products.filter(p => p.category === category);
            if (categoryProducts.length === 0) return null;
            
            return (
              <div key={category} className="mb-5" style={{ animation: 'fadeInUp 0.6s ease-out' }}>
                <div className="d-flex justify-content-between align-items-center mb-4">
                  <div>
                    <h2 className="fw-bold mb-1" style={{ 
                      background: 'linear-gradient(135deg, #0d6efd 0%, #0dcaf0 100%)',
                      WebkitBackgroundClip: 'text',
                      WebkitTextFillColor: 'transparent',
                      backgroundClip: 'text'
                    }}>
                      {category}
                    </h2>
                    <small style={{ color: 'var(--text-secondary)' }}>
                      {categoryProducts.length} {categoryProducts.length === 1 ? 'producto' : 'productos'}
                    </small>
                  </div>
                  <Button 
                    variant="outline-primary"
                    onClick={() => setSelectedCategory(category)}
                  >
                    Ver todos <ChevronRight />
                  </Button>
                </div>
                <Row xs={1} md={2} lg={3} xl={4} className="mb-5 g-4">
                  {categoryProducts.slice(0, 4).map(product => (
                    <Col key={product.id}>
                      <ProductCard product={product} onAddToCart={handleAddToCart} />
                    </Col>
                  ))}
                </Row>
              </div>
            );
          })}
        </>
      ) : (
        <>
          <div className="d-flex justify-content-between align-items-center mb-4">
            <div>
              <h2 className="fw-bold mb-1" style={{ 
                background: 'linear-gradient(135deg, #0d6efd 0%, #0dcaf0 100%)',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                backgroundClip: 'text'
              }}>
                {searchTerm ? `Resultados de búsqueda: "${searchTerm}"` : selectedCategory}
              </h2>
              <small style={{ color: 'var(--text-secondary)' }}>
                {filteredProducts.length} {filteredProducts.length === 1 ? 'resultado' : 'resultados'}
                {selectedCategory !== 'all' && ` en ${selectedCategory}`}
              </small>
            </div>
            <Button 
              variant="outline-secondary"
              onClick={() => {
                setSelectedCategory('all');
                setSearchTerm('');
              }}
            >
              Limpiar filtros
            </Button>
          </div>
          <Row xs={1} md={2} lg={3} xl={4} className="mb-5 g-4">
            {filteredProducts.map(product => (
              <Col key={product.id}>
                <ProductCard product={product} onAddToCart={handleAddToCart} />
              </Col>
            ))}
          </Row>
        </>
      )}

      {filteredProducts.length === 0 && (selectedCategory !== 'all' || searchTerm) && (
        <div className="text-center py-5">
          <p className="text-muted">
            {searchTerm 
              ? `No se encontraron productos que coincidan con "${searchTerm}"`
              : `No se encontraron productos en esta categoría`
            }
          </p>
          <Button variant="primary" onClick={() => {
            setSelectedCategory('all');
            setSearchTerm('');
          }}>
            Ver todos los productos
          </Button>
        </div>
      )}
    </Container>
  );
}

export default Home;
