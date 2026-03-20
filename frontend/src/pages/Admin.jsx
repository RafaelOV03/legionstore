import { useState, useEffect } from 'react';
import { Container, Table, Button, Modal, Form, Alert, Badge, Card, Row, Col, InputGroup, Tabs, Tab, Nav, Spinner } from 'react-bootstrap';
import { PencilSquare, Trash, PlusCircle, Search, ExclamationTriangle, BoxSeam, CurrencyDollar, Diagram3, Receipt, Eye, XCircle, Upload, BellFill } from 'react-bootstrap-icons';
import { getProducts, createProduct, updateProduct, deleteProduct, getOrders, getAllOrders, uploadImage, deleteImage, finalizeOrder } from '../services/resourceApi';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

function Admin() {
  const { hasPermission } = useAuth();
  const navigate = useNavigate();

  // Verificar si tiene permisos para ver esta página
  useEffect(() => {
    if (!hasPermission('products.read')) {
      navigate('/');
    }
  }, [hasPermission, navigate]);
  const [products, setProducts] = useState([]);
  const [filteredProducts, setFilteredProducts] = useState([]);
  const [orders, setOrders] = useState([]);
  const [filteredOrders, setFilteredOrders] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [showOrderModal, setShowOrderModal] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState(null);
  const [editingProduct, setEditingProduct] = useState(null);
  const [alert, setAlert] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterCategory, setFilterCategory] = useState('');
  const [activeTab, setActiveTab] = useState('all');
  const [activeSection, setActiveSection] = useState('inventory');
  const [orderSearchTerm, setOrderSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [currentOrderPage, setCurrentOrderPage] = useState(1);
  const [itemsPerPage] = useState(10);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    price: '',
    stock_quantity: '',
    category: '',
    brand: '',
    image_url: '',
    images: []
  });
  const [uploadingImage, setUploadingImage] = useState(false);
  const [newOrdersCount, setNewOrdersCount] = useState(0);
  const [lastCheckedTime, setLastCheckedTime] = useState(null);

  useEffect(() => {
    loadProducts();
    loadOrders();
    
    // Cargar el último tiempo de verificación desde localStorage
    const savedTime = localStorage.getItem('lastOrderCheckTime');
    if (savedTime) {
      setLastCheckedTime(new Date(savedTime));
    }
    
    // Verificar nuevas órdenes cada 30 segundos si el usuario tiene permiso
    if (hasPermission('orders.read')) {
      const interval = setInterval(() => {
        checkForNewOrders();
      }, 30000); // 30 segundos
      
      return () => clearInterval(interval);
    }
  }, []);

  const checkForNewOrders = async () => {
    if (!hasPermission('orders.read')) return;
    
    try {
      const data = hasPermission('orders.read') 
        ? await getAllOrders() 
        : await getOrders();
      
      const savedTime = localStorage.getItem('lastOrderCheckTime');
      if (savedTime) {
        const lastCheck = new Date(savedTime);
        const newOrders = data.filter(order => {
          const orderDate = new Date(order.created_at);
          return orderDate > lastCheck && order.status === 'COMPLETED' && !order.finalized;
        });
        
        if (newOrders.length > 0) {
          setNewOrdersCount(newOrders.length);
          // Mostrar alerta de nuevas órdenes
          showAlert(`¡${newOrders.length} nueva(s) venta(s) registrada(s)!`, 'info');
        }
      }
    } catch (error) {
      console.error('Error checking for new orders:', error);
    }
  };

  const markOrdersAsChecked = () => {
    const now = new Date().toISOString();
    localStorage.setItem('lastOrderCheckTime', now);
    setLastCheckedTime(new Date(now));
    setNewOrdersCount(0);
  };

  useEffect(() => {
    filterProducts();
    setCurrentPage(1); // Reset a la primera página cuando cambian los filtros
  }, [products, searchTerm, filterCategory, activeTab]);

  useEffect(() => {
    filterOrders();
    setCurrentOrderPage(1); // Reset a la primera página cuando cambia el filtro
  }, [orders, orderSearchTerm]);

  const loadProducts = async () => {
    try {
      const data = await getProducts();
      setProducts(data || []);
    } catch (error) {
      console.error('Error loading products:', error);
      showAlert('Error al cargar los productos', 'danger');
    }
  };

  const loadOrders = async () => {
    try {
      // Usar getAllOrders si el usuario tiene el permiso orders.read_all
      const data = hasPermission('orders.read') 
        ? await getAllOrders() 
        : await getOrders();
      
      // Ordenar por fecha más reciente
      const sortedOrders = (data || []).sort((a, b) => 
        new Date(b.created_at) - new Date(a.created_at)
      );
      setOrders(sortedOrders);
      
      // Verificar si hay nuevas órdenes desde la última verificación
      const savedTime = localStorage.getItem('lastOrderCheckTime');
      if (savedTime && (hasPermission('orders.read'))) {
        const lastCheck = new Date(savedTime);
        const newOrders = sortedOrders.filter(order => {
          const orderDate = new Date(order.created_at);
          return orderDate > lastCheck && order.status === 'COMPLETED' && !order.finalized;
        });
        setNewOrdersCount(newOrders.length);
      }
    } catch (error) {
      console.error('Error loading orders:', error);
      showAlert('Error al cargar las órdenes: ' + error.message, 'danger');
    }
  };

  const filterOrders = () => {
    let filtered = [...orders];
    
    if (orderSearchTerm) {
      filtered = filtered.filter(o =>
        o.payer_email?.toLowerCase().includes(orderSearchTerm.toLowerCase()) ||
        o.payer_name?.toLowerCase().includes(orderSearchTerm.toLowerCase()) ||
        o.paypal_order_id?.toLowerCase().includes(orderSearchTerm.toLowerCase())
      );
    }
    
    setFilteredOrders(filtered);
  };

  const filterProducts = () => {
    let filtered = [...products];

    // Filtrar por búsqueda
    if (searchTerm) {
      filtered = filtered.filter(p =>
        p.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        p.brand?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        p.category?.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Filtrar por categoría
    if (filterCategory) {
      filtered = filtered.filter(p => p.category === filterCategory);
    }

    // Filtrar por tab activo
    if (activeTab === 'low-stock') {
      filtered = filtered.filter(p => p.stock_quantity <= 10);
    } else if (activeTab === 'out-of-stock') {
      filtered = filtered.filter(p => p.stock_quantity === 0);
    }

    setFilteredProducts(filtered);
  };

  const getCategories = () => {
    const categories = [...new Set(products.map(p => p.category))];
    return categories.filter(Boolean).sort();
  };

  const getStatistics = () => {
    const totalProducts = products.length;
    const totalValue = products.reduce((sum, p) => sum + (p.price * p.stock_quantity), 0);
    const lowStock = products.filter(p => p.stock_quantity <= 10 && p.stock_quantity > 0).length;
    const outOfStock = products.filter(p => p.stock_quantity === 0).length;
    const totalStock = products.reduce((sum, p) => sum + p.stock_quantity, 0);

    return { totalProducts, totalValue, lowStock, outOfStock, totalStock };
  };

  const getOrderStatistics = () => {
    const totalOrders = orders.length;
    const completedOrders = orders.filter(o => o.status === 'COMPLETED').length;
    const totalRevenue = orders
      .filter(o => o.status === 'COMPLETED')
      .reduce((sum, o) => sum + o.total_amount, 0);
    const averageOrderValue = completedOrders > 0 ? totalRevenue / completedOrders : 0;

    return { totalOrders, completedOrders, totalRevenue, averageOrderValue };
  };

  const stats = getStatistics();
  const orderStats = getOrderStatistics();

  const handleViewOrder = (order) => {
    console.log('Order selected:', order);
    console.log('Order finalized status:', order.finalized);
    console.log('Has orders.read permission:', hasPermission('orders.read'));
    setSelectedOrder(order);
    setShowOrderModal(true);
  };

  const handleShowModal = (product = null) => {
    if (product) {
      setEditingProduct(product);
      // Parsear images si es string JSON
      let imagesArray = [];
      if (product.images) {
        try {
          imagesArray = typeof product.images === 'string' ? JSON.parse(product.images) : product.images;
        } catch (e) {
          imagesArray = [];
        }
      }
      setFormData({
        name: product.name,
        description: product.description,
        price: product.price,
        stock_quantity: product.stock_quantity,
        category: product.category,
        brand: product.brand,
        image_url: product.image_url || '',
        images: imagesArray
      });
    } else {
      setEditingProduct(null);
      setFormData({
        name: '',
        description: '',
        price: '',
        stock_quantity: '',
        category: '',
        brand: '',
        image_url: '',
        images: []
      });
    }
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingProduct(null);
  };

  const handleImageUpload = async (e) => {
    const files = Array.from(e.target.files);
    if (files.length === 0) return;

    setUploadingImage(true);
    try {
      const uploadPromises = files.map(file => uploadImage(file));
      const results = await Promise.all(uploadPromises);
      const newImageUrls = results.map(r => 'http://localhost:8080' + r.url);
      
      setFormData(prev => ({
        ...prev,
        images: [...prev.images, ...newImageUrls],
        image_url: prev.image_url || newImageUrls[0] // Si no hay imagen principal, usar la primera
      }));
      
      showAlert(`${files.length} imagen(es) subida(s) exitosamente`, 'success');
    } catch (error) {
      console.error('Error uploading images:', error);
      showAlert('Error al subir las imágenes', 'danger');
    } finally {
      setUploadingImage(false);
    }
  };

  const handleRemoveImage = async (imageUrl, index) => {
    try {
      await deleteImage(imageUrl);
      setFormData(prev => {
        const newImages = prev.images.filter((_, i) => i !== index);
        // Si la imagen eliminada era la imagen principal, usar la siguiente disponible
        const newImageUrl = prev.image_url === imageUrl ? (newImages[0] || '') : prev.image_url;
        return {
          ...prev,
          images: newImages,
          image_url: newImageUrl
        };
      });
      showAlert('Imagen eliminada', 'success');
    } catch (error) {
      console.error('Error deleting image:', error);
      showAlert('Error al eliminar la imagen', 'danger');
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const productData = {
        ...formData,
        price: parseFloat(formData.price),
        stock_quantity: parseInt(formData.stock_quantity),
        images: formData.images
      };

      if (editingProduct) {
        await updateProduct(editingProduct.id, productData);
        showAlert('Producto actualizado exitosamente', 'success');
      } else {
        await createProduct(productData);
        showAlert('Producto creado exitosamente', 'success');
      }

      handleCloseModal();
      loadProducts();
    } catch (error) {
      console.error('Error saving product:', error);
      if (error.message.includes('Invalid token') || error.message.includes('Authorization')) {
        showAlert('Sesión expirada. Por favor inicia sesión nuevamente.', 'danger');
        setTimeout(() => {
          window.location.href = '/login';
        }, 2000);
      } else {
        showAlert('Error al guardar el producto: ' + error.message, 'danger');
      }
    }
  };

  const handleDelete = async (id, name) => {
    if (window.confirm(`¿Estás seguro de eliminar "${name}"?`)) {
      try {
        await deleteProduct(id);
        showAlert('Producto eliminado exitosamente', 'success');
        loadProducts();
      } catch (error) {
        console.error('Error deleting product:', error);
        showAlert('Error al eliminar el producto', 'danger');
      }
    }
  };

  const handleFinalizeOrder = async (orderId) => {
    try {
      await finalizeOrder(orderId);
      showAlert('Orden finalizada exitosamente', 'success');
      setShowOrderModal(false);
      await loadOrders(); // Recargar las órdenes
    } catch (error) {
      console.error('Error finalizing order:', error);
      showAlert('Error al finalizar la orden: ' + error.message, 'danger');
    }
  };

  const isNewOrder = (order) => {
    if (order.finalized) return false;
    
    const savedTime = localStorage.getItem('lastOrderCheckTime');
    if (!savedTime) return false;
    
    const lastCheck = new Date(savedTime);
    const orderDate = new Date(order.created_at);
    return orderDate > lastCheck && order.status === 'COMPLETED';
  };

  const showAlert = (message, variant) => {
    setAlert({ message, variant });
    setTimeout(() => setAlert(null), 3000);
  };

  const renderProductTable = (productsToShow) => {
    if (productsToShow.length === 0) {
      return (
        <div className="text-center py-5">
          <p className="text-muted">No hay productos que mostrar</p>
        </div>
      );
    }

    // Calcular paginación
    const indexOfLastItem = currentPage * itemsPerPage;
    const indexOfFirstItem = indexOfLastItem - itemsPerPage;
    const currentItems = productsToShow.slice(indexOfFirstItem, indexOfLastItem);
    const totalPages = Math.ceil(productsToShow.length / itemsPerPage);

    return (
      <>
      <Table responsive striped bordered hover>
        <thead className="table-dark">
          <tr>
            <th style={{ width: '80px' }}>Imagen</th>
            <th>id</th>
            <th>Nombre</th>
            <th>Descripción</th>
            <th>Categoría</th>
            <th>Marca</th>
            <th>Precio</th>
            <th>Stock</th>
            <th>Valor</th>
            <th style={{ width: '150px' }}>Acciones</th>
          </tr>
        </thead>
        <tbody>
          {currentItems.map(product => (
            <tr key={product.id} className={product.stock_quantity === 0 ? 'table-danger' : ''}>
              <td>
                <img
                  src={product.image_url || 'https://via.placeholder.com/50'}
                  alt={product.name}
                  style={{ width: '50px', height: '50px', objectFit: 'cover', borderRadius: '4px' }}
                />
              </td>
              <td>{product.id}</td>
              <td><strong>{product.name}</strong></td>
              <td>
                <small className="text-muted">
                  {product.description?.substring(0, 50)}
                  {product.description?.length > 50 ? '...' : ''}
                </small>
              </td>
              <td><Badge bg="info">{product.category}</Badge></td>
              <td>{product.brand}</td>
              <td><strong>${product.price.toFixed(2)}</strong></td>
              <td>
                <Badge bg={
                  product.stock_quantity === 0 ? 'danger' :
                  product.stock_quantity <= 10 ? 'warning' : 'success'
                }>
                  {product.stock_quantity}
                </Badge>
              </td>
              <td>
                <small className="text-muted">
                  ${(product.price * product.stock_quantity).toFixed(2)}
                </small>
              </td>
              <td>
                <div className="d-flex gap-2">
                  {hasPermission('products.update') && (
                    <Button
                      variant="outline-primary"
                      size="sm"
                      onClick={() => handleShowModal(product)}
                      title="Editar"
                    >
                      <PencilSquare />
                    </Button>
                  )}
                  {hasPermission('products.delete') && (
                    <Button
                      variant="outline-danger"
                      size="sm"
                      onClick={() => handleDelete(product.id, product.name)}
                      title="Eliminar"
                    >
                      <Trash />
                    </Button>
                  )}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>

      {/* Paginación */}
      {totalPages > 1 && (
        <div className="d-flex justify-content-between align-items-center my-3">
          <div className="text-muted">
            Mostrando {indexOfFirstItem + 1} - {Math.min(indexOfLastItem, productsToShow.length)} de {productsToShow.length} productos
          </div>
          <nav>
            <ul className="pagination mb-0">
              <li className={`page-item ${currentPage === 1 ? 'disabled' : ''}`}>
                <button 
                  className="page-link" 
                  onClick={() => setCurrentPage(currentPage - 1)}
                  disabled={currentPage === 1}
                >
                  Anterior
                </button>
              </li>
              {[...Array(totalPages)].map((_, index) => (
                <li key={index} className={`page-item ${currentPage === index + 1 ? 'active' : ''}`}>
                  <button 
                    className="page-link" 
                    onClick={() => setCurrentPage(index + 1)}
                  >
                    {index + 1}
                  </button>
                </li>
              ))}
              <li className={`page-item ${currentPage === totalPages ? 'disabled' : ''}`}>
                <button 
                  className="page-link" 
                  onClick={() => setCurrentPage(currentPage + 1)}
                  disabled={currentPage === totalPages}
                >
                  Siguiente
                </button>
              </li>
            </ul>
          </nav>
        </div>
      )}
      </>
    );
  };

  return (
    <Container fluid>
      {/* Toast Alert en esquina superior derecha */}
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

      <div className="d-flex justify-content-between align-items-center mb-4">
        <h2>Panel de Administración</h2>
        {activeSection === 'inventory' && hasPermission('products.create') && (
          <Button variant="primary" onClick={() => handleShowModal()}>
            <PlusCircle className="me-2" />
            Nuevo Producto
          </Button>
        )}
      </div>

      {/* Navigation Tabs */}
      <Nav variant="tabs" className="mb-4">
        {hasPermission('products.read') && (
          <Nav.Item>
            <Nav.Link 
              active={activeSection === 'inventory'}
              onClick={() => setActiveSection('inventory')}
            >
              <BoxSeam className="me-2" />
              Inventario
            </Nav.Link>
          </Nav.Item>
        )}
        {(hasPermission('orders.read') || hasPermission('orders.read')) && (
          <Nav.Item>
            <Nav.Link 
              active={activeSection === 'sales'}
              onClick={() => {
                setActiveSection('sales');
                markOrdersAsChecked();
              }}
              className="position-relative"
            >
              <Receipt className="me-2" />
              Registro de Ventas
              {newOrdersCount > 0 && (
                <Badge 
                  bg="danger" 
                  pill 
                  className="position-absolute top-0 start-100 translate-middle"
                  style={{ fontSize: '0.7rem' }}
                >
                  {newOrdersCount}
                  <span className="visually-hidden">nuevas órdenes</span>
                </Badge>
              )}
            </Nav.Link>
          </Nav.Item>
        )}
      </Nav>

      {/* Estadísticas */}
      {activeSection === 'inventory' && (
        <Row className="mb-4">
          <Col md={3}>
            <Card className="text-center shadow-sm">
              <Card.Body>
                <BoxSeam size={32} className="text-primary mb-2" />
                <h3 className="mb-0">{stats.totalProducts}</h3>
                <small className="text-muted">Total Productos</small>
              </Card.Body>
            </Card>
          </Col>
          <Col md={3}>
            <Card className="text-center shadow-sm">
              <Card.Body>
                <CurrencyDollar size={32} className="text-success mb-2" />
                <h3 className="mb-0">${stats.totalValue.toFixed(2)}</h3>
                <small className="text-muted">Valor Total Inventario</small>
              </Card.Body>
            </Card>
          </Col>
          <Col md={3}>
            <Card className="text-center shadow-sm">
              <Card.Body>
                <ExclamationTriangle size={32} className="text-warning mb-2" />
                <h3 className="mb-0">{stats.lowStock}</h3>
                <small className="text-muted">Stock Bajo</small>
              </Card.Body>
            </Card>
          </Col>
          <Col md={3}>
            <Card className="text-center shadow-sm">
              <Card.Body>
                <Diagram3 size={32} className="text-info mb-2" />
                <h3 className="mb-0">{stats.totalStock}</h3>
                <small className="text-muted">Unidades Totales</small>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      )}

      {activeSection === 'sales' && (
        <Row className="mb-4">
          <Col md={3}>
            <Card className="text-center shadow-sm">
              <Card.Body>
                <Receipt size={32} className="text-primary mb-2" />
                <h3 className="mb-0">{orderStats.totalOrders}</h3>
                <small className="text-muted">Total Órdenes</small>
              </Card.Body>
            </Card>
          </Col>
          <Col md={3}>
            <Card className="text-center shadow-sm">
              <Card.Body>
                <Badge bg="success" className="mb-2" style={{ fontSize: '1.5rem' }}>✓</Badge>
                <h3 className="mb-0">{orderStats.completedOrders}</h3>
                <small className="text-muted">Completadas</small>
              </Card.Body>
            </Card>
          </Col>
          <Col md={3}>
            <Card className="text-center shadow-sm">
              <Card.Body>
                <CurrencyDollar size={32} className="text-success mb-2" />
                <h3 className="mb-0">${orderStats.totalRevenue.toFixed(2)}</h3>
                <small className="text-muted">Ingresos Totales</small>
              </Card.Body>
            </Card>
          </Col>
          <Col md={3}>
            <Card className="text-center shadow-sm position-relative">
              <Card.Body>
                <Diagram3 size={32} className="text-info mb-2" />
                <h3 className="mb-0">${orderStats.averageOrderValue.toFixed(2)}</h3>
                <small className="text-muted">Ticket Promedio</small>
                {newOrdersCount > 0 && (
                  <Badge 
                    bg="danger" 
                    className="position-absolute top-0 end-0 m-2"
                    style={{ animation: 'pulse 2s infinite' }}
                  >
                    <BellFill className="me-1" />
                    {newOrdersCount} nueva(s)
                  </Badge>
                )}
              </Card.Body>
            </Card>
          </Col>
        </Row>
      )}

      <style>
        {`
          @keyframes pulse {
            0%, 100% {
              opacity: 1;
              transform: scale(1);
            }
            50% {
              opacity: 0.8;
              transform: scale(1.05);
            }
          }
        `}
      </style>

      {/* Filtros Inventario */}
      {activeSection === 'inventory' && (
        <Card className="mb-4 shadow-sm">
          <Card.Body>
            <Row>
              <Col md={6}>
                <InputGroup>
                  <InputGroup.Text>
                    <Search />
                  </InputGroup.Text>
                  <Form.Control
                    placeholder="Buscar por nombre, marca o categoría..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                  />
                </InputGroup>
              </Col>
              <Col md={4}>
                <Form.Select
                  value={filterCategory}
                  onChange={(e) => setFilterCategory(e.target.value)}
                >
                  <option value="">Todas las categorías</option>
                  {getCategories().map(cat => (
                    <option key={cat} value={cat}>{cat}</option>
                  ))}
                </Form.Select>
              </Col>
              <Col md={2}>
                <Button 
                  variant="outline-secondary" 
                  className="w-100"
                  onClick={() => { setSearchTerm(''); setFilterCategory(''); }}
                >
                  Limpiar Filtros
                </Button>
              </Col>
            </Row>
          </Card.Body>
        </Card>
      )}

      {/* Filtros Ventas */}
      {activeSection === 'sales' && (
        <Card className="mb-4 shadow-sm">
          <Card.Body>
            <Row>
              <Col md={10}>
                <InputGroup>
                  <InputGroup.Text>
                    <Search />
                  </InputGroup.Text>
                  <Form.Control
                    placeholder="Buscar por email, nombre del cliente o id de orden..."
                    value={orderSearchTerm}
                    onChange={(e) => setOrderSearchTerm(e.target.value)}
                  />
                </InputGroup>
              </Col>
              <Col md={2}>
                <Button 
                  variant="outline-secondary" 
                  className="w-100"
                  onClick={() => setOrderSearchTerm('')}
                >
                  Limpiar
                </Button>
              </Col>
            </Row>
          </Card.Body>
        </Card>
      )}

      {/* Tabs Inventario */}
      {activeSection === 'inventory' && (
        <Tabs activeKey={activeTab} onSelect={(k) => setActiveTab(k)} className="mb-3">
          <Tab eventKey="all" title={`Todos (${products.length})`}>
            {renderProductTable(filteredProducts)}
          </Tab>
          <Tab 
            eventKey="low-stock" 
            title={
              <span>
                Stock Bajo <Badge bg="warning" className="ms-2">{stats.lowStock}</Badge>
              </span>
            }
          >
            {renderProductTable(filteredProducts)}
          </Tab>
          <Tab 
            eventKey="out-of-stock" 
            title={
              <span>
                Sin Stock <Badge bg="danger" className="ms-2">{stats.outOfStock}</Badge>
              </span>
            }
          >
            {renderProductTable(filteredProducts)}
          </Tab>
        </Tabs>
      )}

      {/* Tabla de Ventas */}
      {activeSection === 'sales' && (
        <>
          {filteredOrders.length === 0 ? (
            <div className="text-center py-5">
              <p className="text-muted">No hay ventas registradas</p>
            </div>
          ) : (
            <>
            {(() => {
              // Calcular paginación para órdenes
              const indexOfLastOrder = currentOrderPage * itemsPerPage;
              const indexOfFirstOrder = indexOfLastOrder - itemsPerPage;
              const currentOrders = filteredOrders.slice(indexOfFirstOrder, indexOfLastOrder);
              const totalOrderPages = Math.ceil(filteredOrders.length / itemsPerPage);

              return (
                <>
            <Table responsive striped bordered hover>
              <thead className="table-dark">
                <tr>
                  <th>id Orden</th>
                  <th>Fecha</th>
                  <th>Cliente</th>
                  <th>Email</th>
                  <th>Total</th>
                  <th>Estado</th>
                  <th>Items</th>
                  <th>Acciones</th>
                </tr>
              </thead>
              <tbody>
                {currentOrders.map(order => (
                  <tr 
                    key={order.id}
                    style={{
                      backgroundColor: isNewOrder(order) ? '#fff3cd' : 'transparent',
                      fontWeight: isNewOrder(order) ? 'bold' : 'normal'
                    }}
                  >
                    <td>
                      <small className="text-muted">
                        {order.paypal_order_id?.substring(0, 15)}...
                      </small>
                    </td>
                    <td>
                      {new Date(order.completed_at || order.created_at).toLocaleDateString('es-ES', {
                        day: '2-digit',
                        month: '2-digit',
                        year: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit'
                      })}
                    </td>
                    <td>{order.payer_name || 'N/A'}</td>
                    <td>{order.payer_email || 'N/A'}</td>
                    <td>
                      <strong>${order.total_amount?.toFixed(2)}</strong>
                    </td>
                    <td>
                      <Badge bg={order.status === 'COMPLETED' ? 'success' : 'warning'}>
                        {order.status}
                      </Badge>
                    </td>
                    <td>
                      <Badge bg="info">{order.items?.length || 0}</Badge>
                    </td>
                    <td>
                      <Button
                        variant="outline-primary"
                        size="sm"
                        onClick={() => handleViewOrder(order)}
                        title="Ver detalles"
                      >
                        <Eye />
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </Table>

            {/* Paginación Órdenes */}
            {totalOrderPages > 1 && (
              <div className="d-flex justify-content-between align-items-center mt-3">
                <div className="text-muted">
                  Mostrando {indexOfFirstOrder + 1} - {Math.min(indexOfLastOrder, filteredOrders.length)} de {filteredOrders.length} órdenes
                </div>
                <nav>
                  <ul className="pagination mb-0">
                    <li className={`page-item ${currentOrderPage === 1 ? 'disabled' : ''}`}>
                      <button 
                        className="page-link" 
                        onClick={() => setCurrentOrderPage(currentOrderPage - 1)}
                        disabled={currentOrderPage === 1}
                      >
                        Anterior
                      </button>
                    </li>
                    {[...Array(totalOrderPages)].map((_, index) => (
                      <li key={index} className={`page-item ${currentOrderPage === index + 1 ? 'active' : ''}`}>
                        <button 
                          className="page-link" 
                          onClick={() => setCurrentOrderPage(index + 1)}
                        >
                          {index + 1}
                        </button>
                      </li>
                    ))}
                    <li className={`page-item ${currentOrderPage === totalOrderPages ? 'disabled' : ''}`}>
                      <button 
                        className="page-link" 
                        onClick={() => setCurrentOrderPage(currentOrderPage + 1)}
                        disabled={currentOrderPage === totalOrderPages}
                      >
                        Siguiente
                      </button>
                    </li>
                  </ul>
                </nav>
              </div>
            )}
                </>
              );
            })()}
            </>
          )}
        </>
      )}

      <Modal show={showModal} onHide={handleCloseModal} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>
            {editingProduct ? 'Editar Producto' : 'Nuevo Producto'}
          </Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form onSubmit={handleSubmit}>
            <Form.Group className="mb-3">
              <Form.Label>Nombre *</Form.Label>
              <Form.Control
                type="text"
                required
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Descripción</Form.Label>
              <Form.Control
                as="textarea"
                rows={3}
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Categoría *</Form.Label>
              <Form.Control
                type="text"
                required
                placeholder="Ej: Smartphones, Laptops, Tablets"
                value={formData.category}
                onChange={(e) => setFormData({ ...formData, category: e.target.value })}
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Marca</Form.Label>
              <Form.Control
                type="text"
                placeholder="Ej: Apple, Samsung, Dell"
                value={formData.brand}
                onChange={(e) => setFormData({ ...formData, brand: e.target.value })}
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Precio *</Form.Label>
              <Form.Control
                type="number"
                step="0.01"
                required
                value={formData.price}
                onChange={(e) => setFormData({ ...formData, price: e.target.value })}
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Stock *</Form.Label>
              <Form.Control
                type="number"
                required
                value={formData.stock_quantity}
                onChange={(e) => setFormData({ ...formData, stock_quantity: e.target.value })}
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Imágenes del Producto</Form.Label>
              
              {/* Galería de imágenes actuales */}
              {formData.images && formData.images.length > 0 && (
                <div className="mb-3">
                  <div className="d-flex flex-wrap gap-2">
                    {formData.images.map((imageUrl, index) => (
                      <div key={index} className="position-relative" style={{ width: '100px', height: '100px' }}>
                        <img
                          src={imageUrl}
                          alt={`Producto ${index + 1}`}
                          className="img-thumbnail"
                          style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                        />
                        <Button
                          variant="danger"
                          size="sm"
                          className="position-absolute top-0 end-0 p-1"
                          style={{ lineHeight: 1 }}
                          onClick={() => handleRemoveImage(imageUrl, index)}
                          disabled={uploadingImage}
                        >
                          <XCircle size={16} />
                        </Button>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Input para subir nuevas imágenes */}
              <div className="d-flex align-items-center gap-2">
                <Form.Control
                  type="file"
                  accept="image/*"
                  multiple
                  onChange={handleImageUpload}
                  disabled={uploadingImage}
                />
                {uploadingImage && (
                  <div className="d-flex align-items-center gap-2">
                    <Spinner animation="border" size="sm" />
                    <small className="text-muted">Subiendo...</small>
                  </div>
                )}
              </div>
              <Form.Text className="text-muted">
                Puedes subir múltiples imágenes. Formatos aceptados: JPG, PNG, GIF, WEBP
              </Form.Text>
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>URL de Imagen (Opcional)</Form.Label>
              <Form.Control
                type="url"
                placeholder="https://example.com/image.jpg"
                value={formData.image_url}
                onChange={(e) => setFormData({ ...formData, image_url: e.target.value })}
              />
              <Form.Text className="text-muted">
                También puedes usar una URL externa si lo prefieres
              </Form.Text>
            </Form.Group>

            <div className="d-flex justify-content-end gap-2">
              <Button variant="secondary" onClick={handleCloseModal}>
                Cancelar
              </Button>
              <Button variant="primary" type="submit">
                {editingProduct ? 'Actualizar' : 'Crear'}
              </Button>
            </div>
          </Form>
        </Modal.Body>
      </Modal>

      {/* Modal de Detalles de Orden */}
      <Modal show={showOrderModal} onHide={() => setShowOrderModal(false)} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>
            Detalles de la Orden
            {selectedOrder && isNewOrder(selectedOrder) && (
              <Badge bg="warning" text="dark" className="ms-2">Nueva</Badge>
            )}
          </Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {selectedOrder && (
            <>
              <Row className="mb-3">
                <Col md={6}>
                  <p><strong>id de Orden:</strong><br />
                    <small className="text-muted">{selectedOrder.paypal_order_id}</small>
                  </p>
                  <p><strong>Cliente:</strong> {selectedOrder.payer_name || 'N/A'}</p>
                  <p><strong>Email:</strong> {selectedOrder.payer_email || 'N/A'}</p>
                </Col>
                <Col md={6}>
                  <p><strong>Fecha:</strong> {new Date(selectedOrder.completed_at || selectedOrder.created_at).toLocaleString('es-ES')}</p>
                  <p><strong>Estado:</strong> <Badge bg={selectedOrder.status === 'COMPLETED' ? 'success' : 'warning'}>{selectedOrder.status}</Badge></p>
                  <p><strong>Total:</strong> <span className="h4 text-success">${selectedOrder.total_amount?.toFixed(2)}</span></p>
                  <p><small className="text-muted">Finalized: {selectedOrder.finalized ? 'true' : 'false'}</small></p>
                </Col>
              </Row>

              <h5 className="mb-3">Productos Comprados</h5>
              <Table striped bordered>
                <thead>
                  <tr>
                    <th>Producto</th>
                    <th>Precio Unit.</th>
                    <th>Cantidad</th>
                    <th>Subtotal</th>
                  </tr>
                </thead>
                <tbody>
                  {selectedOrder.items?.map((item, idx) => (
                    <tr key={idx}>
                      <td>
                        <strong>{item.product_name}</strong>
                        {item.product?.brand && (
                          <><br /><small className="text-muted">{item.product.brand}</small></>
                        )}
                      </td>
                      <td>${item.price?.toFixed(2)}</td>
                      <td>{item.quantity}</td>
                      <td><strong>${(item.price * item.quantity).toFixed(2)}</strong></td>
                    </tr>
                  ))}
                </tbody>
                <tfoot>
                  <tr>
                    <td colSpan="3" className="text-end"><strong>Total:</strong></td>
                    <td><strong className="text-success">${selectedOrder.total_amount?.toFixed(2)}</strong></td>
                  </tr>
                </tfoot>
              </Table>
            </>
          )}
        </Modal.Body>
        <Modal.Footer>
          {selectedOrder && (
            <>
              {selectedOrder.finalized ? (
                <Badge bg="success" className="me-auto">
                  ✓ Orden Finalizada
                </Badge>
              ) : (
                hasPermission('orders.read') && (
                  <Button 
                    variant="success" 
                    onClick={() => handleFinalizeOrder(selectedOrder.id)}
                  >
                    Finalizar Orden
                  </Button>
                )
              )}
            </>
          )}
          <Button variant="secondary" onClick={() => setShowOrderModal(false)}>
            Cerrar
          </Button>
        </Modal.Footer>
      </Modal>
    </Container>
  );
}

export default Admin;
