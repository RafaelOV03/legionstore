# Sistema de Roles y Permisos - Smartech

## Descripción General

Se ha implementado un sistema completo de roles y permisos basado en RBAC (Role-Based Access Control) que permite gestionar el acceso a diferentes funcionalidades del sistema.

## Estructura del Sistema

### Roles Predefinidos

El sistema incluye 3 roles predefinidos:

#### 1. **Usuario** (usuario)
- Rol básico para clientes
- **Permisos:**
  - `products.read` - Ver productos
  - `orders.create` - Crear órdenes
  - `orders.read` - Ver sus propias órdenes

#### 2. **Empleado** (empleado)
- Rol para personal de la tienda
- **Permisos:**
  - `products.read` - Ver productos
  - `products.create` - Crear productos
  - `products.update` - Actualizar productos
  - `orders.read` - Ver órdenes
  - `orders.read_all` - Ver todas las órdenes
  - `orders.create` - Crear órdenes

#### 3. **Administrador** (administrador)
- Acceso total al sistema
- **Permisos:** Todos los permisos disponibles
  - **Productos:** create, read, update, delete
  - **Órdenes:** create, read, read_all, update, delete
  - **Usuarios:** create, read, update, delete
  - **Roles:** create, read, update, delete

### Permisos Disponibles

Los permisos están organizados por recurso y acción:

| Recurso | Acciones | Descripción |
|---------|----------|-------------|
| `products` | create, read, update, delete | Gestión de productos |
| `orders` | create, read, read_all, update, delete | Gestión de órdenes/ventas |
| `users` | create, read, update, delete | Gestión de usuarios |
| `roles` | create, read, update, delete | Gestión de roles |

## Credenciales por Defecto

### Administrador
- **Email:** admin@smartech.com
- **Password:** admin123
- **Rol:** administrador

## Endpoints del Backend

### Usuarios
```
GET    /api/users          - Listar usuarios (requiere: users.read)
GET    /api/users/:id      - Obtener usuario (requiere: users.read)
POST   /api/users          - Crear usuario (requiere: users.create)
PUT    /api/users/:id      - Actualizar usuario (requiere: users.update)
DELETE /api/users/:id      - Eliminar usuario (requiere: users.delete)
```

### Roles
```
GET    /api/roles          - Listar roles (requiere: roles.read)
GET    /api/roles/:id      - Obtener rol (requiere: roles.read)
POST   /api/roles          - Crear rol (requiere: roles.create)
PUT    /api/roles/:id      - Actualizar rol (requiere: roles.update)
DELETE /api/roles/:id      - Eliminar rol (requiere: roles.delete)
```

### Permisos
```
GET    /api/permissions    - Listar permisos (requiere: roles.read)
```

### Órdenes (Actualizadas)
```
PUT    /api/orders/:id     - Actualizar orden (requiere: orders.update)
DELETE /api/orders/:id     - Eliminar orden (requiere: orders.delete)
```

## Frontend

### Nuevas Páginas

1. **Gestión de Usuarios** (`/users`)
   - Lista de todos los usuarios
   - Crear, editar y eliminar usuarios
   - Asignar roles a usuarios
   - Solo accesible con permiso `users.read`

2. **Gestión de Roles** (`/roles`)
   - Lista de roles con tarjetas
   - Crear, editar y eliminar roles personalizados
   - Asignar permisos a roles
   - Ver permisos de cada rol
   - Los roles del sistema no pueden ser eliminados
   - Solo accesible con permiso `roles.read`

### Actualización del Contexto de Autenticación

Se agregaron nuevas funciones al `AuthContext`:

```javascript
hasPermission(permission)  // Verifica si el usuario tiene un permiso específico
hasRole(roleName)         // Verifica si el usuario tiene un rol específico
isAdmin()                 // Verifica si es administrador
isEmployee()              // Verifica si es empleado
isUser()                  // Verifica si es usuario regular
```

### Navegación

La navegación ahora muestra enlaces dinámicos basados en permisos:
- **Admin Panel:** visible con `products.update` o `orders.read_all`
- **Usuarios:** visible con `users.read`
- **Roles:** visible con `roles.read`

## Características de Seguridad

### Backend
1. **JWT con Permisos:** El token JWT incluye los permisos del usuario
2. **Middleware de Permisos:** `RequirePermission(permission)` valida acceso
3. **Roles del Sistema:** Roles predefinidos marcados como `is_system` no pueden ser eliminados
4. **Validación de Usuarios:** No se puede eliminar el propio usuario
5. **Validación de Roles:** No se puede eliminar un rol con usuarios asignados

### Frontend
1. **Protección de Rutas:** ProtectedRoute valida autenticación
2. **UI Condicional:** Botones y enlaces se muestran según permisos
3. **Validación de Acceso:** Redirección automática si no tiene permisos

## Uso del Sistema

### Crear un Nuevo Rol

1. Ir a `/roles`
2. Clic en "Nuevo Rol"
3. Ingresar nombre y descripción
4. Seleccionar los permisos deseados
5. Guardar

### Crear un Nuevo Usuario

1. Ir a `/users`
2. Clic en "Nuevo Usuario"
3. Completar el formulario:
   - Nombre
   - Email
   - Contraseña
   - Rol
4. Guardar

### Modificar Permisos de un Rol

1. Ir a `/roles`
2. Clic en el botón de editar (lápiz) del rol
3. Marcar/desmarcar permisos
4. Guardar cambios

**Nota:** Los roles del sistema (usuario, empleado, administrador) no pueden ser editados ni eliminados.

## Migraciones

Al iniciar el backend por primera vez después de esta actualización:

1. Se eliminarán las tablas antiguas
2. Se crearán las nuevas tablas:
   - `permissions`
   - `roles`
   - `role_permissions` (tabla de relación muchos a muchos)
   - `users` (actualizada con `role_id`)
3. Se poblarán los permisos predefinidos
4. Se crearán los 3 roles del sistema
5. Se creará el usuario administrador por defecto

## Modelo de Datos

### User
```go
type User struct {
    id        uint
    Name      string
    Email     string
    Password  string
    Roleid    uint
    Role      Role
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Role
```go
type Role struct {
    id          uint
    Name        string
    Description string
    Permissions []Permission
    IsSystem    bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Permission
```go
type Permission struct {
    id          uint
    Name        string
    Description string
    Resource    string
    Action      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

## Testing

Para probar el sistema:

1. **Como Administrador:**
   - Login con admin@smartech.com
   - Acceder a todas las páginas
   - Crear roles y usuarios

2. **Como Empleado:**
   - Crear un usuario con rol "empleado"
   - Login con ese usuario
   - Verificar acceso a productos y órdenes
   - Verificar NO acceso a usuarios y roles

3. **Como Usuario:**
   - Registrar un nuevo usuario (se asigna rol "usuario" automáticamente)
   - Login con ese usuario
   - Verificar acceso solo a productos y sus propias órdenes

## Próximas Mejoras

- [ ] Auditoría de cambios en usuarios y roles
- [ ] Historial de permisos asignados
- [ ] Exportar/importar configuración de roles
- [ ] Permisos granulares por producto o categoría
- [ ] Multi-tenancy para múltiples tiendas
