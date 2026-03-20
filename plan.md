# Plan de Refactorización - LegionStore (Nivel Básico)

## Resumen Ejecutivo
LegionStore es un e-commerce full-stack (Go + React) con problemas fundamentales de arquitectura: código duplicado masivo, seguridad débil (credenciales hardcodeadas), y falta de capas de abstracción. Este plan de refactorización BÁSICA prioriza mejoras incrementales sin reescribir el sistema, enfocándose en 5 áreas con ROI más alto: (1) Seguridad, (2) Reducción de duplicación, (3) Mantenibilidad, (4) Configuración, (5) Error handling.

---

## FASE 1: SEGURIDAD E CONFIGURACIÓN (5-7 días)
*Dependencias: Ninguna. Puede ejecutarse inmediatamente.*

### 1.1 Backend - Migrar a Variáables de Entorno

**Archivos afectados:**
- `backend/middleware/auth.go` - JWT secret
- `backend/controllers/order_controller.go` - PayPal credentials
- `backend/database/database.go` - DB path
- `backend/main.go` - Puerto y modo Gin

**Acciones:**
1. Instalar package: `go get github.com/joho/godotenv`
2. Crear `backend/.env.example`:
   ```
   JWT_SECRET=your-secret-key-here
   PAYPAL_CLIENT_ID=your-paypal-id
   PAYPAL_SECRET=your-paypal-secret
   DB_PATH=./legionstore.db
   GIN_MODE=release
   PORT=8080
   CORS_ORIGINS=http://localhost:5173,https://yourdomain.com
   ```
3. Crear `backend/.env` (NO versionar, agregar a .gitignore)
4. Modificar `backend/main.go` init:
   ```go
   import "github.com/joho/godotenv"
   
   func init() {
       if err := godotenv.Load(); err != nil {
           log.Println("No .env file found")
       }
   }
   ```
5. Reemplazar cada hardcode:
   - `auth.go`: `SECRET = os.Getenv("JWT_SECRET")`
   - `order_controller.go`: `os.Getenv("PAYPAL_CLIENT_ID")`
   - `main.go`: Lee puerto y CORS origen de ENV

**Verificación:** `go run main.go` funciona sin modificar código, solo .env

---

### 1.2 Frontend - Configurar Variables de Entorno

**Archivos afectados:**
- `frontend/.env.example`
- `frontend/src/services/api.js`
- `frontend/src/services/inventarioApi.js`
- `frontend/vite.config.js`

**Acciones:**
1. Crear `frontend/.env.local` (NO versionar):
   ```
   VITE_API_URL=http://localhost:8080/api
   VITE_APP_NAME=LegionStore
   ```
2. Crear `frontend/.env.example`:
   ```
   VITE_API_URL=http://localhost:8080/api
   VITE_APP_NAME=LegionStore
   ```
3. Modificar `src/services/api.js`:
   ```javascript
   const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';
   ```
4. Modificar `vite.config.js` para soportar HMR:
   ```javascript
   export default defineConfig({
     plugins: [react()],
     server: {
       port: 5173,
     },
   })
   ```

**Verificación:** Variables de entorno se cargan sin hardcodes en el código

---

## FASE 2: UNIFICACIÓN DE SERVICIOS API (3-4 días)
*Dependencia: Fase 1 completa (usa APIs desde variables)*

### 2.1 Crear Servicio API Unificado

**Archivo nuevo:** `frontend/src/services/apiClient.js`

**Contenido:**
```javascript
// Centralizar toda lógica de fetch, manejo de errores, auth
const API_URL = import.meta.env.VITE_API_URL;

export const apiClient = {
  async request(endpoint, options = {}) {
    const url = `${API_URL}${endpoint}`;
    const headers = {
      'Content-Type': 'application/json',
      ...getAuthHeader(),
      ...options.headers,
    };

    const response = await fetch(url, { ...options, headers });
    
    // Manejar 401 globalmente
    if (response.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
      throw new Error('Session expired. Please login again.');
    }

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || `HTTP ${response.status}`);
    }

    return response.json().catch(() => ({}));
  },

  get: (endpoint) => apiClient.request(endpoint),
  
  post: (endpoint, data) => 
    apiClient.request(endpoint, { method: 'POST', body: JSON.stringify(data) }),
  
  put: (endpoint, data) => 
    apiClient.request(endpoint, { method: 'PUT', body: JSON.stringify(data) }),
  
  delete: (endpoint) => 
    apiClient.request(endpoint, { method: 'DELETE' }),
};

function getAuthHeader() {
  const token = localStorage.getItem('token');
  return token ? { Authorization: `Bearer ${token}` } : {};
}
```

### 2.2 Consolidar `api.js` e `inventarioApi.js`

**Archivo nuevo:** `frontend/src/services/resourceApi.js`

Importar `apiClient` y exportar todas las funciones de negocio (50-60 funciones):
```javascript
import { apiClient } from './apiClient';

// Products
export const getProducts = () => apiClient.get('/products');
export const getProduct = (id) => apiClient.get(`/products/${id}`);
export const createProduct = (data) => apiClient.post('/products', data);
export const updateProduct = (id, data) => apiClient.put(`/products/${id}`, data);
export const deleteProduct = (id) => apiClient.delete(`/products/${id}`);

// Users
export const getUsers = () => apiClient.get('/usuarios');
// ... más funciones
```

### 2.3 Actualizar Importes en Páginas

Buscar y reemplazar:
- `import * as api from '../services/api'` → `import * as api from '../services/resourceApi'`
- `import * as api from '../services/inventarioApi'` → `import * as api from '../services/resourceApi'`

**Verificación:** Todas las páginas funcionan sin cambios lógicos

---

## FASE 3: CREAR COMPONENTE REUTILIZABLE CRUD (4-5 días)
*Dependencia: Fase 2 completa*

### 3.1 Extraer CRUD Component

**Archivo nuevo:** `frontend/src/components/CRUDTable.jsx`

Aceptar props genéricas:
```javascript
export function CRUDTable({
  title,
  columns,              // [{ key: 'id', label: 'ID' }, ...]
  data,
  loading,
  onLoad,               // función para cargar datos
  onAdd,                // (formData) => Promise
  onUpdate,             // (id, formData) => Promise
  onDelete,             // (id) => Promise
  itemShape,            // { nombre: '', email: '' }
  renderForm,           // (formData, setFormData) => JSX
}) {
  // Estado compartido
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState(null);
  const [formData, setFormData] = useState(itemShape);
  const [alert, setAlert] = useState(null);

  // Métodos compartidos
  const load = async () => { /* call onLoad, catch errors */ };
  const handleSubmit = async (e) => { /* call onAdd/onUpdate, handle response */ };
  const handleDelete = async (id) => { /* call onDelete, confirm */ };
  const handleEdit = (item) => { /* set form data, open modal */ };
  const handleCloseModal = () => { /* reset form, close */ };

  return (
    <Container>
      <Button onClick={() => { setEditing(null); setFormData(itemShape); setShowModal(true); }}>
        Add {title}
      </Button>
      
      <Table>
        {/* Render columns dinamicamente */}
      </Table>

      <Modal show={showModal} onHide={handleCloseModal}>
        {renderForm(formData, setFormData)}
        <Button onClick={handleSubmit}>Save</Button>
      </Modal>

      {alert && <Alert variant={alert.type}>{alert.message}</Alert>}
    </Container>
  );
}
```

### 3.2 Refactorizar Páginas CRUD Principales

Iniciar con: `Productos.jsx`, `Users.jsx`, `Sedes.jsx`

**Antes (200 líneas):**
```javascript
const [products, setProducts] = useState([]);
const [loading, setLoading] = useState(true);
// ... 50 líneas de estado
useEffect(() => { loadProducts(); }, []);
const loadProducts = async () => { /* fetch y setState */ };
const handleSubmit = async (e) => { /* validate, post, refetch */ };
// ... más 100 líneas
```

**Después (50 líneas):**
```javascript
const { user } = useAuth();

return (
  <CRUDTable
    title="Productos"
    columns={[
      { key: 'id', label: 'ID' },
      { key: 'nombre', label: 'Nombre' },
      { key: 'precio', label: 'Precio' },
    ]}
    itemShape={{ nombre: '', precio: 0, stock: 0 }}
    onLoad={() => api.getProducts()}
    onAdd={(data) => api.createProduct(data)}
    onUpdate={(id, data) => api.updateProduct(id, data)}
    onDelete={(id) => api.deleteProduct(id)}
    renderForm={(formData, setFormData) => (
      <>
        <Form.Group>
          <Form.Label>Nombre</Form.Label>
          <Form.Control
            value={formData.nombre}
            onChange={(e) => setFormData({ ...formData, nombre: e.target.value })}
          />
        </Form.Group>
        {/* Más campos */}
      </>
    )}
  />
);
```

**Verificación:** Productos, Users, Sedes funcionan idénticamente a antes con 60% menos código

---

## FASE 4: BACKEND - CREAR CAPA REPOSITORY (6-8 días)
*Dependencia: Fase 1*

### 4.1 Crear Paquete de Repositorios

**Estructura nueva:**
```
backend/
├── repository/
│   ├── base.go           # Funciones comunes (scan/map)
│   ├── product_repo.go   # Product CRUD
│   ├── user_repo.go      # User CRUD
│   ├── order_repo.go     # Order CRUD
│   └── ... otros repos
```

### 4.2 Implementar Product Repository

**Archivo nuevo:** `backend/repository/base.go`

Helpers comunes:
```go
package repository

type QueryBuilder struct {
    query  string
    args   []interface{}
}

func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
    if qb.query == "" {
        qb.query = condition
    } else {
        qb.query += " AND " + condition
    }
    qb.args = append(qb.args, args...)
    return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
    return qb.query, qb.args
}

// Función genérica de scan
func ScanProduct(rows *sql.Rows) (Product, error) {
    var p Product
    err := rows.Scan(&p.ID, &p.Name, &p.Price, ...)
    return p, err
}
```

**Archivo nuevo:** `backend/repository/product_repo.go`

```go
package repository

type ProductRepository struct {
    db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
    return &ProductRepository{db: db}
}

// GetAll - Reemplaza GetProducts de controller
func (r *ProductRepository) GetAll(filters map[string]interface{}) ([]Product, error) {
    query := "SELECT id, name, price, stock FROM products WHERE 1=1"
    var args []interface{}

    if name, ok := filters["name"]; ok {
        query += " AND name LIKE ?"
        args = append(args, "%"+name.(string)+"%")
    }

    if category, ok := filters["category"]; ok {
        query += " AND category = ?"
        args = append(args, category)
    }

    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    products := []Product{}
    for rows.Next() {
        p, err := ScanProduct(rows)
        if err != nil {
            return nil, err
        }
        products = append(products, p)
    }

    return products, rows.Err()
}

// GetByID
func (r *ProductRepository) GetByID(id int) (*Product, error) {
    // ...
}

// Create, Update, Delete - similar pattern
```

### 4.3 Refactorizar Controllers para usar Repositories

**Antes** `controllers/product_controller.go`:
```go
func GetProducts(c *gin.Context) {
    rows, err := database.DB.Query("SELECT ...")  // SQL aquí
    // 20 líneas de scan logic
    c.JSON(...)
}
```

**Después:**
```go
var productRepo *repository.ProductRepository

func GetProducts(c *gin.Context) {
    filters := map[string]interface{}{
        "name": c.Query("name"),
        "category": c.Query("category"),
    }
    
    products, err := productRepo.GetAll(filters)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to fetch products"})
        return
    }
    
    c.JSON(200, products)
}
```

### 4.4 Implementar para 5 Entidades Principales

Prioridad:
1. Product → product_repo.go
2. User → user_repo.go
3. Order → order_repo.go
4. Role → role_repo.go
5. Cart → cart_repo.go

**Verificación:** Controllers reducen de 200 líneas a 100 líneas cada uno

---

## FASE 5: MEJORAR ERROR HANDLING Y VALIDACIÓN (3-4 días)
*Dependencia: Fase 4 parcial*

### 5.1 Crear Paquete de Errores Centralizado

**Archivo nuevo:** `backend/errors/errors.go`

```go
package errors

type APIError struct {
    Code    int    `json:"-"`
    Message string `json:"error"`
    Details string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
    return e.Message
}

var (
    ErrNotFound = &APIError{Code: 404, Message: "Resource not found"}
    ErrUnauth   = &APIError{Code: 401, Message: "Unauthorized"}
    ErrForbid   = &APIError{Code: 403, Message: "Forbidden"}
    ErrConflict = &APIError{Code: 409, Message: "Resource already exists"}
    ErrInternal = &APIError{Code: 500, Message: "Internal server error"}
)

func NewValidationError(field, message string) *APIError {
    return &APIError{
        Code:    422,
        Message: "Validation failed",
        Details: fmt.Sprintf("%s: %s", field, message),
    }
}
```

### 5.2 Crear Validador Centralizado

**Archivo nuevo:** `backend/validation/validator.go`

```go
package validation

import "github.com/go-playground/validator/v10"

var v = validator.New()

type ValidationErrors map[string]string

func ValidateStruct(data interface{}) ValidationErrors {
    errs := ValidationErrors{}
    err := v.Struct(data)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            errs[err.Field()] = err.Error()
        }
    }
    return errs
}
```

Usar en modelos:
```go
type Product struct {
    ID    int     `validate:"required,number"`
    Name  string  `validate:"required,min=3,max=100"`
    Price float64 `validate:"required,gt=0"`
}
```

### 5.3 Crear Middleware de Error Handling

**Archivo nuevo:** `backend/middleware/errorhandler.go`

```go
func ErrorHandlingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic: %v", err)
                c.JSON(500, gin.H{"error": "Internal server error"})
            }
        }()
        
        c.Next()
    }
}
```

### 5.4 Actualizar Cotizaciones para usar Transacciones Seguras

**Archivo:** `backend/controllers/cotizacion_controller.go`

Buscar líneas con `tx.Rollback()` y reemplazar:
```go
// Antes
if err != nil {
    tx.Rollback()  // Riesgoso
}

// Después
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

tx, err := database.DB.Begin()
if err != nil {
    return err
}
defer func() {
    if err != nil {
        tx.Rollback()
    } else {
        tx.Commit()
    }
}()
```

---

## FASE 6: CONSOLIDACIÓN DE NAVEGACIÓN FRONT (2-3 días)
*Dependencia: Ninguna*

### 6.1 Unificar Navigation Components

Mantener `Navigation.jsx` (simple) como componente principal actual y eliminar `InventarioNavigation.jsx`.

**Opción A - Simple:** Usar solo Navigation.jsx y agregar condicionales según permisos:
```javascript
export function Navigation({ user, onLogout }) {
  const { hasRole } = useAuth();
  
  return (
    <nav>
      {hasRole('admin') && <NavLink to="/admin">Admin</NavLink>}
      {hasRole('gerente') && <NavLink to="/reportes">Reportes</NavLink>}
      {/* ... */}
    </nav>
  );
}
```

**Opción B - Flexible:** Crear `Navigation.jsx` mejorado que acepte menuItems como prop

**Verificación:** Solo un componente Navigation usado en App.jsx

---

## FASE 7: DOCUMENTACIÓN Y LIMPIEZA (2-3 días)
*Dependencia: Fases 1-6*

### 7.1 Crear API Documentation

**Archivo nuevo:** `backend/API_DOCS.md`

```markdown
# LegionStore API Documentation

## Endpoints

### Products
GET /api/products              - Get all products
GET /api/products/:id          - Get product by ID
POST /api/products             - Create product (admin)
PUT /api/products/:id          - Update product (admin)
DELETE /api/products/:id       - Delete product (admin)

### Users
GET /api/usuarios              - Get all users (admin)
GET /api/usuarios/:id          - Get user by ID
POST /api/usuarios             - Create user (admin)
PUT /api/usuarios/:id          - Update user
DELETE /api/usuarios/:id       - Delete user (admin)

[... más endpoints]

## Error Responses
400 - Bad Request
401 - Unauthorized
403 - Forbidden
404 - Not Found
422 - Validation Error
500 - Internal Server Error
```

### 7.2 Crear README de Backend

**Archivo:** `backend/README.md`

```markdown
# LegionStore Backend

## Setup

1. Copiar .env.example a .env
2. Configurar variables (JWT_SECRET, PAYPAL_*, etc)
3. go run main.go

## Architecture

- `controllers/` - HTTP handlers
- `repository/` - Data access layer
- `models/` - Data structures
- `middleware/` - Auth, CORS, error handling
- `database/` - DB init and seed

## Repositories

Cada repository (product_repo, user_repo) encapsula toda la lógica SQL
para su entidad. El controller solo maneja HTTP.
```

### 7.3 Crear README de Frontend

**Archivo:** `frontend/README-DEV.md`

```markdown
# LegionStore Frontend

## Setup

1. npm install
2. Copiar .env.example a .env.local
3. npm run dev

## Architecture

- `services/apiClient.js` - Centralizado fetch client
- `services/resourceApi.js` - Todas las API calls
- `components/` - Reutilizables (Navigation, CRUDTable, etc)
- `pages/` - Páginas conectadas a routes
- `context/AuthContext.jsx` - Gestión de auth global

## Component Patterns

Use `<CRUDTable>` para operaciones CRUD genéricas.
```

### 7.4 Agregar JSDoc Básico

En archivos principales:
```go
// GetProducts retrieves all products with optional filters
// Parameters: filters (map[string]interface{}) - Optional filters for name, category
// Returns: ([]Product, error)
func (r *ProductRepository) GetAll(filters map[string]interface{}) ([]Product, error) {
```

---

## RESUMEN DE CAMBIOS POR FASE

| Fase | Tiempo | Cambios | Archivos | Impacto |
|------|--------|---------|----------|---------|
| 1 | 5-7d | ENV vars centralizadas | 4 backend, 2 frontend | 🔴 Seguridad Crítica |
| 2 | 3-4d | API client unificado | 2 nuevos, ~15 actualizaciones | ⚠️ Mantenibilidad |
| 3 | 4-5d | CRUDTable reutilizable | 1 nuevo, 3-5 páginas reducidas | ⚠️ Duplication -60% |
| 4 | 6-8d | Repository layer backend | 6-8 nuevos, 15 controllers simplificados | 🟡 Arquitectura |
| 5 | 3-4d | Error handling + validación | 3 nuevos, Controllers mejorados | ⚠️ Robustez |
| 6 | 2-3d | Consolidar Navigation | 1 eliminado, 1 mejorado | 🟢 Consistencia |
| 7 | 2-3d | Documentación | 4 nuevos .md, JSDoc | 🟢 Mantenibilidad |

**Total Completo:** 25-35 días

**MVP Recomendado (PLAN SELECCIONADO):** 10-14 días (Fases 1-3)
- Máxima seguridad inmediata
- -60% código duplicado
- ROI 85% con bajo riesgo
- Fácil de hacer incremental sin bloquear feature work

---

## BENEFICIOS ESPERADOS

✅ **Seguridad:** Credenciales fuera del código  
✅ **Mantenibilidad:** -60% código duplicado  
✅ **Testabilidad:** Repository layer facilita unit tests  
✅ **Escalabilidad:** Arquitectura separada de capas  
✅ **Consistencia:** Patrones y convenciones unificadas  
✅ **Error Handling:** Manejo centralizado de errores y validación  

---

## DECISIONES CLAVE

1. **No reescribir todo** - Refactorización incremental
2. **SQLite → PostgreSQL** se recomienda para producción pero NO está en este plan
3. **TypeScript** se recomienda para futuro pero NO está en este plan
4. **Testing** no incluido pero será más fácil después de estas refactorizaciones
5. **Presupuesto:** Si solo se hacen Fases 1-3 (10-14 días), se obtiene 85% del ROI

