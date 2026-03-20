# Refactor del Backend (Producto, Sede y Stock)

## Objetivo
Reducir la complejidad de controladores y separar responsabilidades para acercar el backend a una arquitectura por capas:

- Controller: HTTP (request/response, status codes, validaciones de entrada)
- Service: reglas de negocio y orquestacion
- Repository: acceso a datos SQL

## Cambios realizados

### 1) Modulo de Productos
- Se extrajo logica de acceso a datos de `product_controller.go`.
- Se agrego `backend/repositories/product_repository.go` para consultas SQL.
- Se agrego `backend/services/product_service.go` para reglas de negocio.
- `product_controller.go` quedo enfocado en flujo HTTP.

### 2) Modulo de Sedes
- Se extrajo logica de datos de `sede_controller.go`.
- Se agrego `backend/repositories/sede_repository.go`.
- Se agrego `backend/services/sede_service.go`.
- `sede_controller.go` ahora maneja principalmente validacion de parametros y respuestas.

### 3) Modulo de Stock (separado de Sedes)
- Se separo la logica de stock a un controlador dedicado: `backend/controllers/stock_controller.go`.
- Se agrego `backend/repositories/stock_repository.go`.
- Se agrego `backend/services/stock_service.go`.
- Las rutas de stock en `main.go` apuntan al controlador de stock.

## Por que Repository + Service (y no solo Service)
Si solo hubiera Service, ese archivo terminaria mezclando:
- SQL y detalles de persistencia
- reglas de negocio
- decisiones de aplicacion

Con Repository + Service:
- se desacopla la base de datos del dominio
- los controladores quedan delgados y mas mantenibles
- es mas facil probar reglas de negocio en aislamiento
- se simplifica cambiar estrategia de persistencia en el futuro

## Beneficios obtenidos
- Menor acoplamiento entre HTTP y SQL.
- Codigo mas legible y mantenible.
- Base preparada para continuar el mismo patron en otros modulos.

## Alcance actual
Refactor aplicado a:
- productos
- sedes
- stock

No se cambiaron contratos de API de forma intencional fuera de lo necesario para la separacion interna.

## Proximo paso recomendado
Aplicar el mismo patron por etapas al resto de controladores con mayor complejidad:
- ordenes
- cotizaciones
- rma
- usuarios/roles

Hacerlo modulo por modulo permite validar sin romper funcionalidades existentes.
