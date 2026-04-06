# API Documentation - LegionStore Backend

**Base URL**: `http://localhost:8080/api`  
**Default Port**: 8080  
**Framework**: Gin (Go)

---

## Table of Contents

1. [Authentication](#authentication)
2. [Error Handling](#error-handling)
3. [Endpoints by Module](#endpoints-by-module)
   - [Health Check](#health-check)
   - [Authentication](#auth-endpoints)
   - [Users & Roles](#users--roles)
   - [Products](#products)
   - [Orders](#orders)
   - [Quotations](#quotations)
   - [Providers](#providers)
   - [RMA](#rma)
   - [Locations](#locations)
   - [Audit](#audit)

---

## Authentication

### JWT Token
All protected endpoints require a Bearer token in the `Authorization` header:

```bash
Authorization: Bearer {token}
```

### Header Example
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Token Payload
```json
{
  "userid": 1,
  "role": "admin",
  "email": "user@example.com"
}
```

---

## Error Handling

### Standard API Error Response

All error responses follow this structure:

```json
{
  "code": 400,
  "message": "Error description",
  "details": "Additional context about the error"
}
```

### Common HTTP Status Codes

| Code | Meaning | Example |
|------|---------|---------|
| 200 | Success | Request processed successfully |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid parameters or validation failed |
| 401 | Unauthorized | Missing or invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource does not exist |
| 409 | Conflict | Resource already exists or state conflict |
| 422 | Unprocessable | Validation error |
| 500 | Internal Server Error | Server processing error |

### Error Response Examples

**Bad Request (400)**
```json
{
  "code": 400,
  "message": "Validation failed",
  "details": "Field 'email' is required"
}
```

**Unauthorized (401)**
```json
{
  "code": 401,
  "message": "Unauthorized",
  "details": "Invalid or missing authentication token"
}
```

**Not Found (404)**
```json
{
  "code": 404,
  "message": "Resource not found",
  "details": "Product with ID 999 not found"
}
```

---

## Endpoints by Module

### Health Check

#### GET /api/health
Health check endpoint for Docker health checks and monitoring.

**Authentication**: Not required  
**Status**: ✅ Available

**Response** (200):
```json
{
  "service": "legionstore-backend",
  "status": "ok",
  "timestamp": 1774011042
}
```

---

### Auth Endpoints

#### POST /api/auth/login
Authenticate user with email and password.

**Authentication**: Not required  
**Status**: ✅ Available

**Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response** (200):
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "role": "admin",
  "sede_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors**:
- 401: Invalid credentials
- 400: Missing email or password

---

#### POST /api/auth/logout
Logout current user (invalidate token on client side).

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Response** (200):
```json
{
  "message": "Logged out successfully"
}
```

---

### Users & Roles

#### GET /api/users
List all users with pagination.

**Authentication**: Required ✔️  
**Status**: ✅ Available  
**Permissions**: Admin required

**Query Parameters**:
- `page` (integer): Page number (default: 1)
- `limit` (integer): Items per page (default: 20)
- `role` (string): Filter by role

**Response** (200):
```json
{
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "admin",
      "sede_id": 1,
      "activo": true,
      "created_at": "2026-03-20T10:30:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "limit": 20
}
```

---

#### GET /api/users/:id
Get user details by ID.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Path Parameters**:
- `id` (integer): User ID

**Response** (200): User object (see GET /api/users)

**Errors**:
- 404: User not found

---

#### POST /api/users
Create a new user.

**Authentication**: Required ✔️  
**Status**: ✅ Available  
**Permissions**: Admin required

**Body**:
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "securepassword",
  "role": "vendor",
  "sede_id": 1
}
```

**Response** (201): User object

**Errors**:
- 400: Validation failed (missing fields)
- 409: Email already exists

---

#### PUT /api/users/:id
Update user information.

**Authentication**: Required ✔️  
**Status**: ✅ Available  
**Permissions**: Admin or self

**Body**:
```json
{
  "name": "Jane Smith",
  "email": "jane.smith@example.com",
  "role": "vendor"
}
```

**Response** (200): Updated user object

---

#### DELETE /api/users/:id
Delete a user.

**Authentication**: Required ✔️  
**Status**: ✅ Available  
**Permissions**: Admin required

**Response** (200):
```json
{
  "message": "User deleted successfully"
}
```

---

#### GET /api/roles
List all available roles.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Response** (200):
```json
[
  {
    "id": 1,
    "name": "admin",
    "description": "Administrator with full access"
  },
  {
    "id": 2,
    "name": "vendor",
    "description": "Vendor user"
  }
]
```

---

### Products

#### GET /api/productos
List all active products.

**Authentication**: Not required  
**Status**: ✅ Available

**Query Parameters**:
- `category` (string): Filter by category
- `page` (integer): Pagination
- `search` (string): Search by name or code

**Response** (200):
```json
{
  "data": [
    {
      "id": 1,
      "codigo": "PROD001",
      "name": "Product Name",
      "description": "Product description",
      "price": 99.99,
      "category": "Electronics",
      "brand": "Brand Name",
      "stock": 50,
      "activo": true
    }
  ],
  "total": 100,
  "page": 1
}
```

---

#### GET /api/productos/:id
Get product details.

**Authentication**: Not required  
**Status**: ✅ Available

**Response** (200): Product object

---

#### POST /api/productos
Create a new product.

**Authentication**: Required ✔️  
**Status**: ✅ Available  
**Permissions**: Admin required

**Body**:
```json
{
  "codigo": "PROD002",
  "name": "New Product",
  "description": "Product description",
  "price": 149.99,
  "category": "Electronics",
  "brand": "Brand"
}
```

**Response** (201): Product object

---

#### PUT /api/productos/:id
Update product information.

**Authentication**: Required ✔️  
**Status**: ✅ Available  
**Permissions**: Admin required

**Body**: Same as POST

**Response** (200): Updated product object

---

#### DELETE /api/productos/:id
Delete a product.

**Authentication**: Required ✔️  
**Status**: ✅ Available  
**Permissions**: Admin required

**Response** (200):
```json
{
  "message": "Product deleted successfully"
}
```

---

### Orders

#### GET /api/ordenes
List orders with filters.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Query Parameters**:
- `estado` (string): Filter by status (pendiente, confirmada, entregada, cancelada)
- `usuario_id` (integer): Filter by user
- `page` (integer): Pagination

**Response** (200):
```json
{
  "data": [
    {
      "id": 1,
      "numero_orden": "ORD-001",
      "usuario_id": 1,
      "total": 299.99,
      "estado": "confirmada",
      "created_at": "2026-03-20T10:30:00Z"
    }
  ],
  "total": 25
}
```

---

#### GET /api/ordenes/:id
Get order details with items.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Response** (200):
```json
{
  "id": 1,
  "numero_orden": "ORD-001",
  "usuario_id": 1,
  "items": [
    {
      "producto_id": 1,
      "cantidad": 2,
      "precio_unit": 99.99
    }
  ],
  "total": 299.99,
  "estado": "confirmada"
}
```

---

#### POST /api/ordenes
Create a new order.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Body**:
```json
{
  "items": [
    {
      "producto_id": 1,
      "cantidad": 2
    }
  ]
}
```

**Response** (201): Order object

---

#### PUT /api/ordenes/:id/estado
Update order status.

**Authentication**: Required ✔️  
**Status**: ✅ Available  
**Permissions**: Admin or owner

**Body**:
```json
{
  "estado": "entregada"
}
```

**Response** (200):
```json
{
  "message": "Order status updated successfully",
  "estado": "entregada"
}
```

---

### Quotations

#### GET /api/cotizaciones
List quotations.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Response** (200): Array of quotation objects

---

#### GET /api/cotizaciones/:id
Get quotation details.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Response** (200): Quotation object

---

#### POST /api/cotizaciones
Create a new quotation.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Body**:
```json
{
  "cliente_nombre": "Client Name",
  "cliente_telefono": "555-1234",
  "cliente_email": "client@example.com",
  "items": [
    {
      "producto_id": 1,
      "cantidad": 5,
      "precio_unit": 99.99
    }
  ]
}
```

**Response** (201): Quotation object

---

#### PUT /api/cotizaciones/:id/estado
Update quotation status.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Body**:
```json
{
  "estado": "aprobada"
}
```

**Response** (200): Updated quotation

---

#### POST /api/cotizaciones/:id/convertir-venta
Convert approved quotation to sale.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Response** (201):
```json
{
  "venta_id": 5,
  "numero_venta": "VNT-2026-0001"
}
```

---

### Providers

#### GET /api/proveedores
List all providers.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Response** (200):
```json
{
  "data": [
    {
      "id": 1,
      "nombre": "Provider Name",
      "ruc": "12345678901",
      "direccion": "Address",
      "telefono": "555-1234",
      "email": "provider@example.com",
      "contacto": "Contact Person",
      "activo": true
    }
  ]
}
```

---

#### POST /api/proveedores
Create provider.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Body**:
```json
{
  "nombre": "New Provider",
  "ruc": "87654321098",
  "direccion": "Provider Address",
  "telefono": "555-5678",
  "email": "newprovider@example.com"
}
```

**Response** (201): Provider object

---

#### GET /api/proveedores/deudas
List provider debts.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Query Parameters**:
- `estado` (string): pendiente, parcial, pagada
- `proveedor_id` (integer): Filter by provider

**Response** (200): Array of debt objects

---

#### POST /api/proveedores/deudas
Create provider debt.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Body**:
```json
{
  "proveedor_id": 1,
  "num_factura": "INV-001",
  "monto": 5000.00,
  "fecha_vencimiento": "2026-04-20"
}
```

**Response** (201): Debt object

---

### RMA

#### GET /api/rmas
List RMA requests.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Response** (200): Array of RMA objects

---

#### POST /api/rmas
Create RMA request.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Body**:
```json
{
  "orden_id": 1,
  "motivo": "Defective product",
  "descripcion": "Product not working"
}
```

**Response** (201): RMA object

---

#### PUT /api/rmas/:id/estado
Update RMA status.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Body**:
```json
{
  "estado": "resuelto",
  "solucion": "Product replaced"
}
```

**Response** (200): Updated RMA

---

### Locations

#### GET /api/sedes
List all locations.

**Authentication**: Not required  
**Status**: ✅ Available

**Response** (200):
```json
[
  {
    "id": 1,
    "nombre": "Main Office",
    "direccion": "123 Main St",
    "telefono": "555-0000",
    "activa": true
  }
]
```

---

#### POST /api/sedes
Create location.

**Authentication**: Required ✔️  
**Status**: ✅ Available  
**Permissions**: Admin required

**Body**:
```json
{
  "nombre": "New Location",
  "direccion": "456 New St",
  "telefono": "555-1111"
}
```

**Response** (201): Location object

---

#### GET /api/sedes/stock
Get stock across all locations.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Response** (200):
```json
[
  {
    "sede_id": 1,
    "sede_nombre": "Main Office",
    "producto_id": 1,
    "cantidad": 50,
    "stock_minimo": 10,
    "stock_maximo": 100
  }
]
```

---

#### GET /api/sedes/:sede_id/stock
Get stock for specific location.

**Authentication**: Required ✔️  
**Status**: ✅ Available

**Response** (200): Array of stock items for that location

---

### Audit

#### GET /api/auditoria
List audit log entries.

**Authentication**: Required ✔️  
**Status**: ✅ Available  
**Permissions**: Admin required

**Query Parameters**:
- `usuario_id` (integer): Filter by user
- `accion` (string): Filter by action
- `fecha_inicio` (string): Start date
- `fecha_fin` (string): End date

**Response** (200):
```json
{
  "data": [
    {
      "id": 1,
      "usuario_id": 1,
      "accion": "crear",
      "tabla": "productos",
      "registro_id": 1,
      "valores_antiguos": null,
      "valores_nuevos": "{}",
      "created_at": "2026-03-20T10:30:00Z"
    }
  ],
  "total": 100
}
```

---

## Rate Limiting

No rate limiting is currently implemented. Please implement appropriate rate limiting for production deployment.

---

## Pagination

Endpoints that return multiple items support pagination:

**Query Parameters**:
- `page` (integer): Page number (default: 1)
- `limit` (integer): Items per page (default: 20, max: 100)

**Response Format**:
```json
{
  "data": [...],
  "total": 150,
  "page": 1,
  "limit": 20,
  "pages": 8
}
```

---

## Timestamps

All timestamps are in ISO 8601 format (UTC):
```
2026-03-20T10:30:45Z
```

---

## Notes for Developers

1. **CORS**: Configured for `http://localhost:5173` (frontend dev server) and Docker networking
2. **Database**: SQLite3 at `./legionstore.db`
3. **JWT Secret**: Set via `JWT_SECRET` environment variable
4. **PayPal Integration**: Uses sandbox credentials (configure for production)
5. **Health Check**: Available without authentication at GET /api/health

---

## Future Enhancements

- [ ] API versioning (v1, v2, etc.)
- [ ] OpenAPI/Swagger documentation
- [ ] GraphQL support
- [ ] WebSocket support for real-time updates
- [ ] Request validation middleware
- [ ] Rate limiting per user/IP
- [ ] Caching layer (Redis)
- [ ] Comprehensive API metrics

---

**Last Updated**: March 20, 2026  
**Status**: Production Ready (v1.0)
