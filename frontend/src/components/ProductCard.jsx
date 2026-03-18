import { Card, Button, Badge } from 'react-bootstrap';
import { CartPlus, Eye } from 'react-bootstrap-icons';
import { useNavigate } from 'react-router-dom';

function ProductCard({ product, onAddToCart }) {
  const navigate = useNavigate();
  
  // Parsear images si es string JSON
  let imagesArray = [];
  if (product.images) {
    try {
      imagesArray = typeof product.images === 'string' ? JSON.parse(product.images) : product.images;
    } catch (e) {
      imagesArray = [];
    }
  }
  
  const mainImage = imagesArray && imagesArray.length > 0 
    ? imagesArray[0] 
    : product.image_url || 'https://via.placeholder.com/300x200';

  return (
    <Card className="h-100 shadow-sm hover-shadow" style={{ cursor: 'pointer' }}>
      <div onClick={() => navigate(`/product/${product.id}`)}>
        <Card.Img 
          variant="top" 
          src={mainImage} 
          style={{ height: '200px', objectFit: 'cover' }}
        />
      </div>
      <Card.Body className="d-flex flex-column">
        <div className="mb-2">
          <Badge bg="info" className="me-2">{product.category}</Badge>
          <Badge bg="secondary">{product.brand}</Badge>
        </div>
        <Card.Title 
          className="h6" 
          onClick={() => navigate(`/product/${product.id}`)}
          style={{ cursor: 'pointer' }}
        >
          {product.name}
        </Card.Title>
        <Card.Text className="text-muted small flex-grow-1">
          {product.description?.substring(0, 80)}
          {product.description?.length > 80 && '...'}
        </Card.Text>
        <div className="mt-auto">
          <div className="d-flex justify-content-between align-items-center mb-2">
            <h5 className="text-primary mb-0">${product.price.toFixed(2)}</h5>
            <small className="text-muted">
              Stock: {product.stock_quantity}
            </small>
          </div>
          <div className="d-flex gap-2">
            <Button 
              variant="outline-primary" 
              onClick={() => navigate(`/product/${product.id}`)}
              style={{ flex: '0 0 auto' }}
            >
              <Eye />
            </Button>
            <Button 
              variant="primary" 
              className="flex-grow-1"
              onClick={() => onAddToCart(product)}
              disabled={product.stock_quantity === 0}
            >
              <CartPlus className="me-2" />
              {product.stock_quantity === 0 ? 'Sin Stock' : 'Agregar'}
            </Button>
          </div>
        </div>
      </Card.Body>
    </Card>
  );
}

export default ProductCard;
