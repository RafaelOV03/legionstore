# Sistema E-Commerce "Smartech"

Sistema de comercio electrónico full-stack para la venta de dispositivos electrónicos con autenticación de usuarios.

## Tecnologías

### Backend
- **Go 1.24** con framework **Gin**
- **GORM** para ORM
- **SQLite** como base de datos
- **JWT (JSON Web Tokens)** para autenticación
- **Bcrypt** para encriptación de contraseñas

### Frontend
- **React** con **Vite**
- **React Router** para navegación
- **Bootstrap 5** con **React Bootstrap** para el diseño
- **Context API** para gestión de estado global

## Características

### Para Clientes
- Registro e inicio de sesión de usuarios
- Catálogo de productos con filtros por categoría
- Búsqueda de productos
- Carrito de compras funcional
- Visualización de productos con imágenes

### Para Administradores
- Acceso protegido al panel de administración
- Gestión completa de productos (CRUD)
- Control de inventario
- Actualización de precios y stock

## Instalación y Ejecución

### Backend

```bash
cd backend
go mod download
go run main.go
```

El servidor estará disponible en `http://localhost:8080`

**Usuario administrador por defecto:**
- Email: `admin@inventario.com`
- Contraseña: `admin123`

### Frontend

```bash
cd frontend
npm install
npm run dev
```

La aplicación web estará disponible en `http://localhost:5173`

## API Endpoints

### Autenticación
- `POST /api/auth/register` - Registrar nuevo usuario
- `POST /api/auth/login` - Iniciar sesión
- `POST /api/auth/logout` - Cerrar sesión
- `GET /api/auth/me` - Obtener usuario actual (requiere autenticación)

### Productos
- `GET /api/products` - Obtener todos los productos (público)
- `GET /api/products/:id` - Obtener un producto (público)
- `POST /api/products` - Crear producto (requiere rol admin)
- `PUT /api/products/:id` - Actualizar producto (requiere rol admin)
- `DELETE /api/products/:id` - Eliminar producto (requiere rol admin)

### Carrito
- `GET /api/cart?session_id={id}` - Obtener carrito
- `POST /api/cart/items?session_id={id}` - Agregar al carrito
- `PUT /api/cart/items/:id?session_id={id}` - Actualizar cantidad
- `DELETE /api/cart/items/:id?session_id={id}` - Eliminar del carrito
- `DELETE /api/cart/clear?session_id={id}` - Vaciar carrito

## Autenticación y Autorización

El sistema implementa autenticación basada en JWT:

1. **Registro/Login**: Los usuarios se registran o inician sesión y reciben un token JWT
2. **Almacenamiento**: El token se guarda en localStorage del navegador
3. **Autorización**: Las rutas protegidas requieren el token en el header `Authorization: Bearer {token}`
4. **Roles**: 
   - `customer`: Usuario normal con acceso al catálogo y carrito
   - `admin`: Acceso completo incluyendo gestión de productos

## Estructura del Proyecto

```
Proyecto/
├── backend/
│   ├── controllers/
│   │   ├── product_controller.go
│   │   ├── cart_controller.go
│   │   └── auth_controller.go
│   ├── middleware/
│   │   └── auth.go
│   ├── database/
│   │   ├── database.go
│   │   └── seed.go
│   ├── models/
│   │   ├── product.go
│   │   ├── cart.go
│   │   └── user.go
│   ├── main.go
│   └── go.mod
└── frontend/
    ├── src/
    │   ├── components/
    │   │   ├── Navigation.jsx
    │   │   ├── ProductCard.jsx
    │   │   └── ProtectedRoute.jsx
    │   ├── context/
    │   │   └── AuthContext.jsx
    │   ├── pages/
    │   │   ├── Home.jsx
    │   │   ├── Admin.jsx
    │   │   ├── Cart.jsx
    │   │   └── Login.jsx
    │   ├── services/
    │   │   └── api.js
    │   ├── hooks/
    │   │   └── useCart.js
    │   ├── App.jsx
    │   └── main.jsx
    └── package.json
```

## Diseño

- Paleta de colores: Azules, grises y blancos (tema tecnológico)
- Diseño responsive (Mobile-first)
- Principios de Material Design
- Interfaz moderna y minimalista

## Datos de Ejemplo

El sistema incluye 12 productos de ejemplo en las siguientes categorías:
- Smartphones
- Laptops
- Tablets
- Accesorios
- Wearables
- Monitores

Marcas incluidas: Apple, Samsung, Dell, Sony, Logitech, LG

## Seguridad

- Contraseñas encriptadas con bcrypt
- Tokens JWT con expiración de 24 horas
- Rutas protegidas con middleware de autenticación
- Validación de roles para acceso administrativo
- Headers CORS configurados correctamente
