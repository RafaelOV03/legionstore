import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Badge, Spinner, Alert } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { 
  BoxSeam, Building, ClipboardCheck, Truck, Tools, 
  Receipt, GraphUp, People, ExclamationTriangle
} from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import * as api from '../services/resourceApi';

function Dashboard() {
  const { user, hasPermission } = useAuth();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [stats, setStats] = useState({
    ordenes: null,
    rmas: null,
    insumos: null,
  });

  useEffect(() => {
    const loadStats = async () => {
      try {
        const promises = [];
        
        if (hasPermission('ordenes.read')) {
          promises.push(api.getOrdenesStats().then(data => ({ ordenes: data })));
        }
        if (hasPermission('rmas.read')) {
          promises.push(api.getRMAStats().then(data => ({ rmas: data })));
        }
        if (hasPermission('insumos.read')) {
          promises.push(api.getInsumosStats().then(data => ({ insumos: data })));
        }

        const results = await Promise.all(promises);
        const newStats = results.reduce((acc, curr) => ({ ...acc, ...curr }), {});
        setStats(newStats);
      } catch (err) {
        setError('Error cargando estadísticas');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    if (user) {
      loadStats();
    } else {
      setLoading(false);
    }
  }, [user, hasPermission]);

  if (!user) {
    return (
      <Container className="py-5">
        <Row className="justify-content-center">
          <Col md={8} lg={6}>
            <Card className="text-center" style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="py-5">
                <BoxSeam size={80} className="text-primary mb-4" />
                <h2 className="mb-3">Sistema de Gestión de Inventario</h2>
                <p className="text-muted mb-4">
                  Inicia sesión para acceder al sistema de inventario, gestión de órdenes de trabajo, 
                  RMA, cotizaciones y más.
                </p>
                <Link to="/login" className="btn btn-primary btn-lg">
                  Iniciar Sesión
                </Link>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    );
  }

  const QuickAccessCard = ({ icon: Icon, title, description, to, color, stats: cardStats }) => (
    <Card 
      as={Link} 
      to={to} 
      className="h-100 text-decoration-none" 
      style={{ 
        background: 'var(--card-bg)', 
        border: '1px solid var(--border-color)',
        transition: 'all 0.3s ease'
      }}
      onMouseOver={(e) => {
        e.currentTarget.style.transform = 'translateY(-5px)';
        e.currentTarget.style.boxShadow = `0 10px 30px ${color}30`;
      }}
      onMouseOut={(e) => {
        e.currentTarget.style.transform = 'translateY(0)';
        e.currentTarget.style.boxShadow = 'none';
      }}
    >
      <Card.Body>
        <div className="d-flex align-items-center mb-3">
          <div 
            className="rounded-circle p-3 me-3" 
            style={{ background: `${color}20` }}
          >
            <Icon size={24} style={{ color }} />
          </div>
          <div>
            <h5 className="mb-0" style={{ color: 'var(--text-primary)' }}>{title}</h5>
            <small className="text-muted">{description}</small>
          </div>
        </div>
        {cardStats && (
          <div className="d-flex gap-2 flex-wrap">
            {cardStats.map((stat, idx) => (
              <Badge 
                key={idx} 
                bg={stat.variant || 'secondary'}
                className="px-2 py-1"
              >
                {stat.label}: {stat.value}
              </Badge>
            ))}
          </div>
        )}
      </Card.Body>
    </Card>
  );

  return (
    <Container fluid className="py-4">
      {/* Encabezado */}
      <Row className="mb-4">
        <Col>
          <h2>
            Bienvenido, {user.name}
            <Badge 
              bg={
                user.role?.name === 'administrador' ? 'danger' :
                user.role?.name === 'gerente' ? 'warning' :
                user.role?.name === 'vendedor' ? 'info' :
                user.role?.name === 'tecnico' ? 'success' :
                'secondary'
              } 
              className="ms-3"
              style={{ fontSize: '0.5em', verticalAlign: 'middle' }}
            >
              {user.role?.name || 'Usuario'}
            </Badge>
          </h2>
          <p className="text-muted">
            Panel de control del Sistema de Gestión de Inventario
          </p>
        </Col>
      </Row>

      {error && <Alert variant="danger">{error}</Alert>}
      {loading && (
        <div className="text-center py-5">
          <Spinner animation="border" variant="primary" />
        </div>
      )}

      {!loading && (
        <>
          {/* Alertas urgentes */}
          {stats.ordenes?.urgentes > 0 && (
            <Alert variant="warning" className="mb-4">
              <ExclamationTriangle className="me-2" />
              Hay <strong>{stats.ordenes.urgentes}</strong> órdenes de trabajo urgentes pendientes
            </Alert>
          )}
          {stats.insumos?.bajo_stock > 0 && (
            <Alert variant="warning" className="mb-4">
              <ExclamationTriangle className="me-2" />
              Hay <strong>{stats.insumos.bajo_stock}</strong> insumos con stock bajo
            </Alert>
          )}

          {/* Accesos rápidos según rol */}
          <h5 className="mb-3">Accesos Rápidos</h5>
          <Row className="g-3 mb-4">
            {/* Inventario - Todos */}
            {hasPermission('products.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={BoxSeam}
                  title="Productos"
                  description="Gestionar inventario"
                  to="/productos"
                  color="#0d6efd"
                />
              </Col>
            )}

            {hasPermission('stock.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={Building}
                  title="Stock Multisede"
                  description="Ver stock por sede"
                  to="/stock"
                  color="#6f42c1"
                />
              </Col>
            )}

            {/* Cotizaciones */}
            {hasPermission('cotizaciones.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={Receipt}
                  title="Cotizaciones"
                  description="Emitir cotizaciones"
                  to="/cotizaciones"
                  color="#20c997"
                />
              </Col>
            )}

            {/* Compatibilidad */}
            {hasPermission('compatibilidad.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={Tools}
                  title="Compatibilidad"
                  description="Buscar compatibles"
                  to="/compatibilidad"
                  color="#fd7e14"
                />
              </Col>
            )}

            {/* Órdenes de Trabajo */}
            {hasPermission('ordenes.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={ClipboardCheck}
                  title="Órdenes de Trabajo"
                  description="Servicio técnico"
                  to="/ordenes-trabajo"
                  color="#198754"
                  stats={stats.ordenes ? [
                    { label: 'Pendientes', value: stats.ordenes.recibidos, variant: 'warning' },
                    { label: 'En proceso', value: stats.ordenes.en_reparacion, variant: 'info' },
                  ] : null}
                />
              </Col>
            )}

            {/* Insumos */}
            {hasPermission('insumos.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={Tools}
                  title="Insumos"
                  description="Gestionar insumos"
                  to="/insumos"
                  color="#fd7e14"
                  stats={stats.insumos ? [
                    { label: 'Total', value: stats.insumos.total_insumos, variant: 'primary' },
                    { label: 'Bajo stock', value: stats.insumos.bajo_stock, variant: 'danger' },
                  ] : null}
                />
              </Col>
            )}

            {/* RMA / Garantías */}
            {hasPermission('rmas.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={ClipboardCheck}
                  title="RMA / Garantías"
                  description="Gestionar devoluciones"
                  to="/rma"
                  color="#dc3545"
                  stats={stats.rmas ? [
                    { label: 'Recibidos', value: stats.rmas.recibidos, variant: 'warning' },
                    { label: 'En revisión', value: stats.rmas.en_revision, variant: 'info' },
                  ] : null}
                />
              </Col>
            )}

            {/* Traspasos */}
            {hasPermission('traspasos.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={Truck}
                  title="Traspasos"
                  description="Mover inventario"
                  to="/traspasos"
                  color="#0dcaf0"
                />
              </Col>
            )}

            {/* Proveedores */}
            {hasPermission('proveedores.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={People}
                  title="Proveedores"
                  description="Gestionar proveedores"
                  to="/proveedores"
                  color="#6c757d"
                />
              </Col>
            )}

            {/* Deudas */}
            {hasPermission('deudas.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={Receipt}
                  title="Deudas"
                  description="Control de deudas"
                  to="/deudas"
                  color="#dc3545"
                />
              </Col>
            )}

            {/* Reportes */}
            {hasPermission('reportes.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={GraphUp}
                  title="Reportes"
                  description="Análisis de ganancias"
                  to="/reportes"
                  color="#ffc107"
                />
              </Col>
            )}

            {/* Auditoría */}
            {hasPermission('auditoria.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={ClipboardCheck}
                  title="Auditoría"
                  description="Logs del sistema"
                  to="/auditoria"
                  color="#6c757d"
                />
              </Col>
            )}

            {/* Promociones */}
            {hasPermission('promociones.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={Receipt}
                  title="Promociones"
                  description="Gestionar promociones"
                  to="/promociones"
                  color="#e83e8c"
                />
              </Col>
            )}

            {/* Usuarios */}
            {hasPermission('users.read') && (
              <Col md={6} lg={4} xl={3}>
                <QuickAccessCard
                  icon={People}
                  title="Usuarios"
                  description="Gestionar usuarios"
                  to="/usuarios"
                  color="#6c757d"
                />
              </Col>
            )}
          </Row>

          {/* Información del usuario */}
          <Row className="mt-4">
            <Col md={6} lg={4}>
              <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
                <Card.Header>
                  <h6 className="mb-0">Tu información</h6>
                </Card.Header>
                <Card.Body>
                  <p className="mb-1"><strong>Nombre:</strong> {user.name}</p>
                  <p className="mb-1"><strong>Email:</strong> {user.email}</p>
                  <p className="mb-1"><strong>Rol:</strong> {user.role?.name || 'Sin rol'}</p>
                  <p className="mb-0"><strong>Sede:</strong> {user.sede?.nombre || 'No asignada'}</p>
                </Card.Body>
              </Card>
            </Col>
          </Row>
        </>
      )}
    </Container>
  );
}

export default Dashboard;
