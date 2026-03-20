import { Container, Navbar, Nav, NavDropdown, Badge } from 'react-bootstrap';
import { Link, useLocation } from 'react-router-dom';
import { 
  BoxSeam, House, PersonCircle, BoxArrowRight, Cart3,
  Building, Truck, ClipboardCheck, Tools, People,
  FileEarmarkText, GraphUp, Shield, Receipt, Gear
} from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../hooks/useCart';

function Navigation() {
  const { user, logout, hasPermission } = useAuth();
  const { cartItemsCount } = useCart();
  const location = useLocation();

  const handleLogout = () => {
    logout();
    window.location.href = '/login';
  };

  const isActive = (path) => location.pathname.startsWith(path);

  // Determinar roles basado en permisos
  const isGerente = hasPermission('auditoria.read') && hasPermission('reportes.read');
  const isVendedor = hasPermission('cotizaciones.read') && hasPermission('stock.read');
  const isTecnico = hasPermission('ordenes.read') && hasPermission('insumos.read');
  const isAdmin = hasPermission('rmas.read') && hasPermission('traspasos.read');

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
      <Container fluid>
        <Navbar.Brand as={Link} to="/" className="fw-bold" style={{
          fontSize: '1.4rem',
          background: 'linear-gradient(135deg, #0d6efd 0%, #0dcaf0 100%)',
          WebkitBackgroundClip: 'text',
          WebkitTextFillColor: 'transparent',
          backgroundClip: 'text',
        }}>
          <BoxSeam className="me-2" style={{ color: '#0d6efd' }} />
          LegionStore
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="navbar-nav" />
        <Navbar.Collapse id="navbar-nav">
          <Nav className="me-auto">
            <Nav.Link as={Link} to="/" active={location.pathname === '/'}>
              <House className="me-1" /> Inicio
            </Nav.Link>

            {/* INVENTARIO - Productos y Stock */}
            {(hasPermission('products.read') || hasPermission('stock.read')) && (
              <NavDropdown 
                title={<><BoxSeam className="me-1" /> Inventario</>} 
                active={isActive('/productos') || isActive('/stock')}
              >
                {hasPermission('products.read') && (
                  <NavDropdown.Item as={Link} to="/productos">
                    <BoxSeam className="me-2" /> Productos
                  </NavDropdown.Item>
                )}
                {hasPermission('stock.read') && (
                  <NavDropdown.Item as={Link} to="/stock">
                    <Building className="me-2" /> Stock por Sede
                  </NavDropdown.Item>
                )}
                {hasPermission('sedes.read') && (
                  <NavDropdown.Item as={Link} to="/sedes">
                    <Building className="me-2" /> Sedes
                  </NavDropdown.Item>
                )}
              </NavDropdown>
            )}

            {/* VENTAS - Cotizaciones y Compatibilidad (Vendedor) */}
            {isVendedor && (
              <NavDropdown 
                title={<><Receipt className="me-1" /> Ventas</>}
                active={isActive('/cotizaciones') || isActive('/compatibilidad')}
              >
                <NavDropdown.Item as={Link} to="/cotizaciones">
                  <FileEarmarkText className="me-2" /> Cotizaciones
                </NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/compatibilidad">
                  <Gear className="me-2" /> Compatibilidad
                </NavDropdown.Item>
              </NavDropdown>
            )}

            {/* SERVICIO TÉCNICO - Órdenes e Insumos (Técnico) */}
            {isTecnico && (
              <NavDropdown 
                title={<><Tools className="me-1" /> Servicio Técnico</>}
                active={isActive('/ordenes-trabajo') || isActive('/insumos')}
              >
                <NavDropdown.Item as={Link} to="/ordenes-trabajo">
                  <ClipboardCheck className="me-2" /> Órdenes de Trabajo
                </NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/insumos">
                  <Gear className="me-2" /> Insumos
                </NavDropdown.Item>
              </NavDropdown>
            )}

            {/* ADMINISTRACIÓN - RMA, Traspasos, Proveedores, Deudas (Administrador) */}
            {isAdmin && (
              <NavDropdown 
                title={<><Shield className="me-1" /> Administración</>}
                active={isActive('/rma') || isActive('/traspasos') || isActive('/proveedores') || isActive('/deudas')}
              >
                <NavDropdown.Item as={Link} to="/rma">
                  <ClipboardCheck className="me-2" /> RMA / Garantías
                </NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/traspasos">
                  <Truck className="me-2" /> Traspasos
                </NavDropdown.Item>
                <NavDropdown.Divider />
                <NavDropdown.Item as={Link} to="/proveedores">
                  <People className="me-2" /> Proveedores
                </NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/deudas">
                  <Receipt className="me-2" /> Deudas
                </NavDropdown.Item>
              </NavDropdown>
            )}

            {/* GERENCIA - Reportes, Auditoría, Promociones (Gerente) */}
            {isGerente && (
              <NavDropdown 
                title={<><GraphUp className="me-1" /> Gerencia</>}
                active={isActive('/reportes') || isActive('/auditoria') || isActive('/promociones')}
              >
                <NavDropdown.Item as={Link} to="/reportes">
                  <GraphUp className="me-2" /> Reportes
                </NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/auditoria">
                  <Shield className="me-2" /> Auditoría
                </NavDropdown.Item>
                <NavDropdown.Divider />
                <NavDropdown.Item as={Link} to="/promociones">
                  <FileEarmarkText className="me-2" /> Promociones
                </NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/segmentacion">
                  <People className="me-2" /> Segmentación
                </NavDropdown.Item>
              </NavDropdown>
            )}

            {/* CONFIGURACIÓN - Usuarios y Roles */}
            {hasPermission('users.read') && (
              <NavDropdown 
                title={<><Gear className="me-1" /> Config</>}
                active={isActive('/usuarios') || isActive('/roles')}
              >
                <NavDropdown.Item as={Link} to="/usuarios">
                  <People className="me-2" /> Usuarios
                </NavDropdown.Item>
                {hasPermission('roles.read') && (
                  <NavDropdown.Item as={Link} to="/roles">
                    <Shield className="me-2" /> Roles
                  </NavDropdown.Item>
                )}
              </NavDropdown>
            )}
          </Nav>

          {/* Carrito + Usuario */}
          <Nav>
            <Nav.Link as={Link} to="/cart" className="position-relative me-3">
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
                  <span>
                    <PersonCircle className="me-1" />
                    {user.name}
                    <Badge 
                      bg={
                        user.role?.name === 'administrador' ? 'danger' :
                        user.role?.name === 'gerente' ? 'warning' :
                        user.role?.name === 'vendedor' ? 'info' :
                        user.role?.name === 'tecnico' ? 'success' :
                        'secondary'
                      } 
                      className="ms-2"
                      style={{ fontSize: '0.7em' }}
                    >
                      {user.role?.name || 'Usuario'}
                    </Badge>
                  </span>
                }
                id="user-dropdown"
                align="end"
              >
                <NavDropdown.Item disabled>
                  <small className="text-muted">{user.email}</small>
                </NavDropdown.Item>
                <NavDropdown.Item disabled>
                  <small className="text-muted">Sede: {user.sede?.nombre || 'No asignada'}</small>
                </NavDropdown.Item>
                <NavDropdown.Divider />
                <NavDropdown.Item onClick={handleLogout} className="text-danger">
                  <BoxArrowRight className="me-2" />
                  Cerrar Sesión
                </NavDropdown.Item>
              </NavDropdown>
            ) : (
              <Nav.Link as={Link} to="/login">
                <PersonCircle className="me-1" />
                Iniciar Sesión
              </Nav.Link>
            )}
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

export default Navigation;
