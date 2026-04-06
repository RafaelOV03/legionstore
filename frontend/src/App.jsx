import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import Navigation from './components/Navigation';
import ProtectedRoute from './components/ProtectedRoute';
import { lazy, Suspense } from 'react';
import { Spinner, Container } from 'react-bootstrap';

// Pages
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import Users from './pages/Users';
import Roles from './pages/Roles';

// Lazy loaded pages
const Productos = lazy(() => import('./pages/Productos'));
const Stock = lazy(() => import('./pages/Stock'));
const Sedes = lazy(() => import('./pages/Sedes'));
const RMA = lazy(() => import('./pages/RMA'));
const Cotizaciones = lazy(() => import('./pages/Cotizaciones'));
const Traspasos = lazy(() => import('./pages/Traspasos'));
const OrdenesTrabajoPage = lazy(() => import('./pages/OrdenesTrabajoPage'));
const Proveedores = lazy(() => import('./pages/Proveedores'));
const Deudas = lazy(() => import('./pages/Deudas'));
const Insumos = lazy(() => import('./pages/Insumos'));
const Compatibilidad = lazy(() => import('./pages/Compatibilidad'));
const Auditoria = lazy(() => import('./pages/Auditoria'));
const Reportes = lazy(() => import('./pages/Reportes'));
const Promociones = lazy(() => import('./pages/Promociones'));
const Segmentacion = lazy(() => import('./pages/Segmentacion'));

// Loading component
const PageLoader = () => (
  <Container className="d-flex justify-content-center align-items-center" style={{ minHeight: '50vh' }}>
    <Spinner animation="border" variant="primary" />
  </Container>
);

function App() {
  return (
    <AuthProvider>
      <BrowserRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
        <div className="min-vh-100">
          <Navigation />
          <Suspense fallback={<PageLoader />}>
            <Routes>
              {/* Dashboard */}
              <Route path="/" element={<Dashboard />} />
              <Route path="/login" element={<Login />} />

              {/* Inventario */}
              <Route
                path="/productos"
                element={
                  <ProtectedRoute requiredPermission="products.read">
                    <Productos />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/stock"
                element={
                  <ProtectedRoute requiredPermission="stock.read">
                    <Stock />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/sedes"
                element={
                  <ProtectedRoute requiredPermission="sedes.read">
                    <Sedes />
                  </ProtectedRoute>
                }
              />

              {/* Ventas / Vendedor */}
              <Route
                path="/cotizaciones"
                element={
                  <ProtectedRoute requiredPermission="cotizaciones.read">
                    <Cotizaciones />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/compatibilidad"
                element={
                  <ProtectedRoute requiredPermission="stock.read">
                    <Compatibilidad />
                  </ProtectedRoute>
                }
              />

              {/* Servicio Técnico */}
              <Route
                path="/ordenes-trabajo"
                element={
                  <ProtectedRoute requiredPermission="ordenes.read">
                    <OrdenesTrabajoPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/insumos"
                element={
                  <ProtectedRoute requiredPermission="insumos.read">
                    <Insumos />
                  </ProtectedRoute>
                }
              />

              {/* Administración */}
              <Route
                path="/rma"
                element={
                  <ProtectedRoute requiredPermission="rmas.read">
                    <RMA />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/traspasos"
                element={
                  <ProtectedRoute requiredPermission="traspasos.read">
                    <Traspasos />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/proveedores"
                element={
                  <ProtectedRoute requiredPermission="proveedores.read">
                    <Proveedores />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/deudas"
                element={
                  <ProtectedRoute requiredPermission="deudas.read">
                    <Deudas />
                  </ProtectedRoute>
                }
              />

              {/* Gerencia */}
              <Route
                path="/reportes"
                element={
                  <ProtectedRoute requiredPermission="reportes.read">
                    <Reportes />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/auditoria"
                element={
                  <ProtectedRoute requiredPermission="auditoria.read">
                    <Auditoria />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/promociones"
                element={
                  <ProtectedRoute requiredPermission="promociones.read">
                    <Promociones />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/segmentacion"
                element={
                  <ProtectedRoute requiredPermission="segmentacion.read">
                    <Segmentacion />
                  </ProtectedRoute>
                }
              />

              {/* Configuración */}
              <Route
                path="/usuarios"
                element={
                  <ProtectedRoute requiredPermission="users.read">
                    <Users />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/roles"
                element={
                  <ProtectedRoute requiredPermission="roles.read">
                    <Roles />
                  </ProtectedRoute>
                }
              />
            </Routes>
          </Suspense>
        </div>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
