import { Container, Navbar, Nav, Badge, NavDropdown } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { Cart3, BoxSeam, House, PersonCircle, BoxArrowRight } from 'react-bootstrap-icons';
import { useCart } from '../hooks/useCart';
import { useAuth } from '../context/AuthContext';

function Navigation() {
  const { cartItemsCount } = useCart();
  const { user, logout, isAdmin, hasPermission } = useAuth();

  const handleLogout = () => {
    logout();
    window.location.href = '/';
  };

  return (
    <Navbar variant="dark" expand="lg" className="mb-4" style={{ 
      background: 'rgba(10, 14, 39, 0.95)',
      backdropFilter: 'blur(10px)',
      boxShadow: '0 4px 20px rgba(0, 0, 0, 0.5)',
      borderBottom: '1px solid var(--border-color)',
      position: 'sticky',
      top: 0,
      zIndex: 1030
    }}>
      <Container>
        <Navbar.Brand as={Link} to="/" className="fw-bold" style={{
          fontSize: '1.5rem',
          background: 'linear-gradient(135deg, #0d6efd 0%, #0dcaf0 100%)',
          WebkitBackgroundClip: 'text',
          WebkitTextFillColor: 'transparent',
          backgroundClip: 'text',
          transition: 'all 0.3s ease'
        }}>
          <BoxSeam className="me-2" style={{ color: '#0d6efd' }} />
          Smartech
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="ms-auto">
            <Nav.Link as={Link} to="/">
              <House className="me-1" />
              Inicio
            </Nav.Link>
            {(hasPermission('products.read') || hasPermission('orders.read')) && (
              <Nav.Link as={Link} to="/admin">
                <BoxSeam className="me-1" />
                Admin
              </Nav.Link>
            )}
            {hasPermission('users.read') && (
              <Nav.Link as={Link} to="/users">
                <PersonCircle className="me-1" />
                Usuarios
              </Nav.Link>
            )}
            {hasPermission('roles.read') && (
              <Nav.Link as={Link} to="/roles">
                <BoxSeam className="me-1" />
                Roles
              </Nav.Link>
            )}
            <Nav.Link as={Link} to="/cart" className="position-relative">
              <Cart3 className="me-1" />
              Carrito
              {cartItemsCount > 0 && (
                <Badge bg="danger" pill className="ms-1">
                  {cartItemsCount}
                </Badge>
              )}
            </Nav.Link>
            
            {user ? (
              <NavDropdown
                title={
                  <>
                    <PersonCircle className="me-1" />
                    {user.name}
                  </>
                }
                id="user-dropdown"
                align="end"
                style={{ zIndex: 9999 }}
                className="position-relative"
              >
                <NavDropdown.Item disabled>
                  <small className="text-muted">{user.email}</small>
                </NavDropdown.Item>
                <NavDropdown.Item disabled>
                  <Badge bg={
                    user.role?.name === 'administrador' ? 'danger' : 
                    user.role?.name === 'empleado' ? 'primary' : 
                    'info'
                  }>
                    {user.role?.name || 'Usuario'}
                  </Badge>
                </NavDropdown.Item>
                <NavDropdown.Divider />
                <NavDropdown.Item as={Link} to="/my-orders">
                  <BoxSeam className="me-2" />
                  Mis Compras
                </NavDropdown.Item>
                <NavDropdown.Divider />
                <NavDropdown.Item onClick={handleLogout}>
                  <BoxArrowRight className="me-2" />
                  Cerrar Sesión
                </NavDropdown.Item>
              </NavDropdown>
            ) : (
              <>
                <Nav.Link as={Link} to="/my-orders">
                  <BoxSeam className="me-1" />
                  Mis Compras
                </Nav.Link>
                <Nav.Link as={Link} to="/login">
                  <PersonCircle className="me-1" />
                  Iniciar Sesión
                </Nav.Link>
              </>
            )}
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

export default Navigation;
