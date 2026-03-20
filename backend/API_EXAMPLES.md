# API Examples - Request & Response

This document provides practical examples of API requests and responses for the Legion Store backend.

---

## Authentication

### Login Request
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@legionstore.com",
    "password": "password123"
  }'
```

### Login Response (200)
```json
{
  "id": 1,
  "name": "Admin User",
  "email": "admin@legionstore.com",
  "role": "admin",
  "sede_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOjEsInJvbGUiOiJhZG1pbiIsImVtYWlsIjoiYWRtaW5AbGVnaW9uc3RvcmUuY29tIn0.xyz123..."
}
```

### Setting Authorization Header (for all subsequent authenticated requests)
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOjEsInJvbGUiOiJhZG1pbiIsImVtYWlsIjoiYWRtaW5AbGVnaW9uc3RvcmUuY29tIn0.xyz123...
```

---

## Products

### Get All Products
```bash
curl -X GET "http://localhost:8080/api/productos?page=1&limit=20" \
  -H "Content-Type: application/json"
```

### Get All Products Response (200)
```json
{
  "data": [
    {
      "id": 1,
      "codigo": "LAPTOP-HP-015",
      "name": "HP Pavilion 15 Laptop",
      "description": "15.6 inch HD display, Intel Core i5, 8GB RAM, 256GB SSD",
      "price": 599.99,
      "category": "Laptops",
      "brand": "HP",
      "stock": 15,
      "activo": true,
      "created_at": "2026-01-15T08:00:00Z",
      "updated_at": "2026-03-10T14:30:00Z"
    },
    {
      "id": 2,
      "codigo": "MOUSE-LOG-001",
      "name": "Logitech MX Master 3S",
      "description": "Advanced ergonomic mouse with precision tracking",
      "price": 99.99,
      "category": "Peripherals",
      "brand": "Logitech",
      "stock": 45,
      "activo": true,
      "created_at": "2026-01-20T08:00:00Z",
      "updated_at": "2026-03-12T10:15:00Z"
    }
  ],
  "total": 2,
  "page": 1,
  "limit": 20,
  "pages": 1
}
```

### Get Single Product
```bash
curl -X GET http://localhost:8080/api/productos/1 \
  -H "Content-Type: application/json"
```

### Get Product Response (200)
```json
{
  "id": 1,
  "codigo": "LAPTOP-HP-015",
  "name": "HP Pavilion 15 Laptop",
  "description": "15.6 inch HD display, Intel Core i5, 8GB RAM, 256GB SSD",
  "price": 599.99,
  "category": "Laptops",
  "brand": "HP",
  "stock": 15,
  "activo": true,
  "created_at": "2026-01-15T08:00:00Z",
  "updated_at": "2026-03-10T14:30:00Z"
}
```

### Create Product
```bash
curl -X POST http://localhost:8080/api/productos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "codigo": "KEYBOARD-MECH-001",
    "name": "Mechanical Gaming Keyboard RGB",
    "description": "Cherry MX switches, programmable keys, RGB lighting",
    "price": 149.99,
    "category": "Peripherals",
    "brand": "Corsair"
  }'
```

### Create Product Response (201)
```json
{
  "id": 3,
  "codigo": "KEYBOARD-MECH-001",
  "name": "Mechanical Gaming Keyboard RGB",
  "description": "Cherry MX switches, programmable keys, RGB lighting",
  "price": 149.99,
  "category": "Peripherals",
  "brand": "Corsair",
  "stock": 0,
  "activo": true,
  "created_at": "2026-03-20T10:30:00Z",
  "updated_at": "2026-03-20T10:30:00Z"
}
```

### Update Product
```bash
curl -X PUT http://localhost:8080/api/productos/3 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "price": 139.99,
    "stock": 10
  }'
```

### Delete Product
```bash
curl -X DELETE http://localhost:8080/api/productos/3 \
  -H "Authorization: Bearer {token}"
```

### Delete Response (200)
```json
{
  "message": "Product deleted successfully"
}
```

---

## Orders

### Get All Orders
```bash
curl -X GET "http://localhost:8080/api/ordenes?page=1&estado=confirmada" \
  -H "Authorization: Bearer {token}"
```

### Orders List Response (200)
```json
{
  "data": [
    {
      "id": 5,
      "numero_orden": "ORD-2026-0005",
      "usuario_id": 2,
      "total": 699.98,
      "estado": "confirmada",
      "fecha_pago": "2026-03-18T15:20:00Z",
      "created_at": "2026-03-18T14:50:00Z",
      "updated_at": "2026-03-18T15:20:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 20,
  "pages": 1
}
```

### Get Order Detail
```bash
curl -X GET http://localhost:8080/api/ordenes/5 \
  -H "Authorization: Bearer {token}"
```

### Order Detail Response (200)
```json
{
  "id": 5,
  "numero_orden": "ORD-2026-0005",
  "usuario_id": 2,
  "items": [
    {
      "id": 8,
      "orden_id": 5,
      "producto_id": 1,
      "cantidad": 1,
      "precio_unit": 599.99,
      "subtotal": 599.99
    },
    {
      "id": 9,
      "orden_id": 5,
      "producto_id": 2,
      "cantidad": 1,
      "precio_unit": 99.99,
      "subtotal": 99.99
    }
  ],
  "total": 699.98,
  "estado": "confirmada",
  "fecha_pago": "2026-03-18T15:20:00Z",
  "direccion_entrega": "123 Customer St, City, State",
  "created_at": "2026-03-18T14:50:00Z",
  "updated_at": "2026-03-18T15:20:00Z"
}
```

### Create Order
```bash
curl -X POST http://localhost:8080/api/ordenes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "items": [
      {
        "producto_id": 1,
        "cantidad": 1
      },
      {
        "producto_id": 2,
        "cantidad": 1
      }
    ]
  }'
```

### Create Order Response (201)
```json
{
  "id": 6,
  "numero_orden": "ORD-2026-0006",
  "usuario_id": 2,
  "items": [
    {
      "id": 10,
      "orden_id": 6,
      "producto_id": 1,
      "cantidad": 1,
      "precio_unit": 599.99,
      "subtotal": 599.99
    },
    {
      "id": 11,
      "orden_id": 6,
      "producto_id": 2,
      "cantidad": 1,
      "precio_unit": 99.99,
      "subtotal": 99.99
    }
  ],
  "total": 699.98,
  "estado": "pendiente",
  "created_at": "2026-03-20T11:00:00Z",
  "updated_at": "2026-03-20T11:00:00Z"
}
```

### Update Order Status
```bash
curl -X PUT http://localhost:8080/api/ordenes/6/estado \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "estado": "entregada"
  }'
```

### Order Status Update Response (200)
```json
{
  "message": "Order status updated successfully",
  "estado": "entregada",
  "updated_at": "2026-03-20T11:05:00Z"
}
```

---

## Quotations

### Create Quotation
```bash
curl -X POST http://localhost:8080/api/cotizaciones \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "cliente_nombre": "Acme Corporation",
    "cliente_telefono": "555-1234",
    "cliente_email": "contact@acme.com",
    "items": [
      {
        "producto_id": 1,
        "cantidad": 5,
        "precio_unit": 599.99
      },
      {
        "producto_id": 2,
        "cantidad": 10,
        "precio_unit": 99.99
      }
    ]
  }'
```

### Create Quotation Response (201)
```json
{
  "id": 3,
  "numero_cotizacion": "COT-2026-0003",
  "cliente_nombre": "Acme Corporation",
  "cliente_telefono": "555-1234",
  "cliente_email": "contact@acme.com",
  "items": [
    {
      "id": 5,
      "cotizacion_id": 3,
      "producto_id": 1,
      "cantidad": 5,
      "precio_unit": 599.99,
      "subtotal": 2999.95
    },
    {
      "id": 6,
      "cotizacion_id": 3,
      "producto_id": 2,
      "cantidad": 10,
      "precio_unit": 99.99,
      "subtotal": 999.90
    }
  ],
  "estado": "pendiente",
  "total": 3999.85,
  "created_at": "2026-03-20T11:15:00Z",
  "updated_at": "2026-03-20T11:15:00Z"
}
```

### Approve Quotation
```bash
curl -X PUT http://localhost:8080/api/cotizaciones/3/estado \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "estado": "aprobada"
  }'
```

### Convert to Sale
```bash
curl -X POST http://localhost:8080/api/cotizaciones/3/convertir-venta \
  -H "Authorization: Bearer {token}"
```

### Convert to Sale Response (201)
```json
{
  "venta_id": 7,
  "numero_venta": "VNT-2026-0007",
  "cotizacion_id": 3,
  "total": 3999.85,
  "estado": "creada",
  "created_at": "2026-03-20T11:20:00Z"
}
```

---

## Users

### Get All Users
```bash
curl -X GET "http://localhost:8080/api/users?page=1&limit=20" \
  -H "Authorization: Bearer {token}"
```

### Users List Response (200)
```json
{
  "data": [
    {
      "id": 1,
      "name": "Admin User",
      "email": "admin@legionstore.com",
      "role": "admin",
      "sede_id": 1,
      "activo": true,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-03-10T14:30:00Z"
    },
    {
      "id": 2,
      "name": "John Seller",
      "email": "john@legionstore.com",
      "role": "vendor",
      "sede_id": 1,
      "activo": true,
      "created_at": "2026-01-15T08:00:00Z",
      "updated_at": "2026-03-15T10:00:00Z"
    }
  ],
  "total": 2,
  "page": 1,
  "limit": 20,
  "pages": 1
}
```

### Create User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "name": "Sarah Manager",
    "email": "sarah@legionstore.com",
    "password": "securePassword@123",
    "role": "manager",
    "sede_id": 1
  }'
```

### Create User Response (201)
```json
{
  "id": 3,
  "name": "Sarah Manager",
  "email": "sarah@legionstore.com",
  "role": "manager",
  "sede_id": 1,
  "activo": true,
  "created_at": "2026-03-20T11:30:00Z",
  "updated_at": "2026-03-20T11:30:00Z"
}
```

---

## Providers

### Get All Providers
```bash
curl -X GET http://localhost:8080/api/proveedores \
  -H "Authorization: Bearer {token}"
```

### Providers List Response (200)
```json
{
  "data": [
    {
      "id": 1,
      "nombre": "Tech Supplies Corp",
      "ruc": "12345678901",
      "direccion": "789 Industrial Blvd, Tech City",
      "telefono": "555-TECH1",
      "email": "sales@techsupplies.com",
      "contacto": "Mike Johnson",
      "activo": true,
      "created_at": "2026-01-10T00:00:00Z",
      "updated_at": "2026-03-15T09:00:00Z"
    }
  ],
  "total": 1
}
```

### Record Provider Debt
```bash
curl -X POST http://localhost:8080/api/proveedores/deudas \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "proveedor_id": 1,
    "num_factura": "INV-2026-0015",
    "monto": 5000.00,
    "fecha_vencimiento": "2026-04-20"
  }'
```

### Record Debt Response (201)
```json
{
  "id": 1,
  "proveedor_id": 1,
  "num_factura": "INV-2026-0015",
  "monto": 5000.00,
  "monto_pagado": 0.00,
  "fecha_vencimiento": "2026-04-20",
  "estado": "pendiente",
  "created_at": "2026-03-20T11:40:00Z",
  "updated_at": "2026-03-20T11:40:00Z"
}
```

---

## RMA (Return & Maintenance)

### Create RMA
```bash
curl -X POST http://localhost:8080/api/rmas \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "orden_id": 5,
    "motivo": "Defective Product",
    "descripcion": "Laptop has hardware failure, does not turn on"
  }'
```

### Create RMA Response (201)
```json
{
  "id": 1,
  "numero_rma": "RMA-2026-0001",
  "orden_id": 5,
  "usuario_id": 2,
  "motivo": "Defective Product",
  "descripcion": "Laptop has hardware failure, does not turn on",
  "estado": "abierta",
  "solucion": null,
  "created_at": "2026-03-20T12:00:00Z",
  "updated_at": "2026-03-20T12:00:00Z"
}
```

### Close RMA with Resolution
```bash
curl -X PUT http://localhost:8080/api/rmas/1/estado \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "estado": "resuelto",
    "solucion": "Replaced with new unit. Original unit returned to warehouse."
  }'
```

### Close RMA Response (200)
```json
{
  "id": 1,
  "numero_rma": "RMA-2026-0001",
  "estado": "resuelto",
  "solucion": "Replaced with new unit. Original unit returned to warehouse.",
  "updated_at": "2026-03-20T12:15:00Z"
}
```

---

## Locations (Sedes)

### Get All Locations
```bash
curl -X GET http://localhost:8080/api/sedes \
  -H "Content-Type: application/json"
```

### Locations List Response (200)
```json
[
  {
    "id": 1,
    "nombre": "Main Office - Lima",
    "direccion": "123 Main Street, Lima",
    "telefono": "555-0001",
    "activa": true,
    "created_at": "2026-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "nombre": "Branch Office - Arequipa",
    "direccion": "456 Branch Avenue, Arequipa",
    "telefono": "555-0002",
    "activa": true,
    "created_at": "2026-02-01T00:00:00Z"
  }
]
```

### Get Stock by Location
```bash
curl -X GET "http://localhost:8080/api/sedes/1/stock" \
  -H "Authorization: Bearer {token}"
```

### Location Stock Response (200)
```json
[
  {
    "sede_id": 1,
    "sede_nombre": "Main Office - Lima",
    "producto_id": 1,
    "nombre_producto": "HP Pavilion 15 Laptop",
    "cantidad": 15,
    "stock_minimo": 5,
    "stock_maximo": 50
  },
  {
    "sede_id": 1,
    "sede_nombre": "Main Office - Lima",
    "producto_id": 2,
    "nombre_producto": "Logitech MX Master 3S",
    "cantidad": 45,
    "stock_minimo": 10,
    "stock_maximo": 100
  }
]
```

---

## Error Examples

### Validation Error (400)
```bash
curl -X POST http://localhost:8080/api/productos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "name": "Product",
    "price": "invalid"
  }'
```

### Error Response (400)
```json
{
  "code": 400,
  "message": "Validation failed",
  "details": "Field 'price' must be a number"
}
```

### Unauthorized (401)
```bash
curl -X GET http://localhost:8080/api/users \
  -H "Content-Type: application/json"
```

### Unauthorized Response (401)
```json
{
  "code": 401,
  "message": "Unauthorized",
  "details": "Missing or invalid authentication token"
}
```

### Not Found (404)
```bash
curl -X GET http://localhost:8080/api/productos/999 \
  -H "Content-Type: application/json"
```

### Not Found Response (404)
```json
{
  "code": 404,
  "message": "Not Found",
  "details": "Product with ID 999 not found"
}
```

### Conflict (409)
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "email": "admin@legionstore.com",
    "password": "password123",
    "name": "Duplicate",
    "role": "admin"
  }'
```

### Conflict Response (409)
```json
{
  "code": 409,
  "message": "Conflict",
  "details": "Email 'admin@legionstore.com' already registered"
}
```

---

## Tips for API Testing

### Using cURL
- Always include `Content-Type: application/json` for POST/PUT
- Include `Authorization: Bearer {token}` for authenticated endpoints
- Use `-d` for request body data
- Use `-H` for headers

### Using Postman
1. Import API_DOCS.md endpoints
2. Set Authorization header globally (Bearer Token)
3. Use environment variables for base URL
4. Save requests in collections for reuse

### Using JavaScript/Fetch
```javascript
const response = await fetch('http://localhost:8080/api/usuarios', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  }
})
const data = await response.json()
```

---

**Last Updated**: March 20, 2026

