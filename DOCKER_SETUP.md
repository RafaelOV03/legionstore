# Guía de Docker Compose para LegionStore

## Requisitos Previos

- Docker (versión 20.10+): https://docs.docker.com/install/
- Docker Compose (versión 2.0+): https://docs.docker.com/compose/install/

## Estructura de Contenedores

```
legionstore/
├── backend/
│   ├── Dockerfile
│   ├── .dockerignore
│   └── (código Go)
├── frontend/
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── nginx.conf
│   └── (código React)
├── docker-compose.yml
└── .env.example
```

## Instalación y Ejecución

### 1. Clonar o descargar el proyecto

```bash
cd /ruta/del/proyecto
```

### 2. Configurar variables de entorno (opcional)

```bash
cp .env.example .env
# Editar .env según tus necesidades
```

### 3. Construir y ejecutar con Docker Compose

```bash
# Construir las imágenes Docker
docker-compose build

# Iniciar los servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Ver logs de un servicio específico
docker-compose logs -f backend
docker-compose logs -f frontend
```

### 4. Acceder a la aplicación

- **Frontend**: http://localhost
- **Backend API**: http://localhost:8080/api

## Comandos Útiles

```bash
# Detener servicios
docker-compose down

# Detener y eliminar volúmenes
docker-compose down -v

# Reconstruir servicios
docker-compose build --no-cache

# Ejecutar comando en un contenedor
docker-compose exec backend sh
docker-compose exec frontend sh

# Ver estado de los servicios
docker-compose ps

# Reiniciar un servicio
docker-compose restart backend
docker-compose restart frontend

# Ver variables de entorno
docker-compose config
```

## Estructura de Puertos

| Servicio | Puerto Interno | Puerto Externo | URL |
|----------|---|---|---|
| Backend  | 8080 | 8080 | http://localhost:8080 |
| Frontend | 80 | 80 | http://localhost |

## Volúmenes

- `backend-data`: Almacena datos de SQLite del backend

## Red

Los servicios se comunican entre sí usando la red Docker `legionstore-network`:
- Frontend se conecta a Backend usando `http://backend:8080`
- Backend está disponible en `http://localhost:8080` desde el exterior

## Health Check

El backend incluye un health check que verifica la conectividad cada 30 segundos. El frontend solo inicia después de que el backend esté listo.

## Desarrollo Local

Para cambios en desarrollo, puedes:

```bash
# Construir solo un servicio
docker-compose build backend
docker-compose build frontend

# Recrear solo un servicio
docker-compose up -d backend
docker-compose up -d frontend

# Inspeccionar contenedor
docker-compose exec backend ls -la
```

## Solución de Problemas

### El frontend no se conecta al backend
- Verifica que la variable `VITE_API_URL` en el nginx.conf sea correcta
- El frontend debe acceder al backend como `http://backend:8080` desde el contenedor

### Puerto 80 ya está en uso
Cambia en `docker-compose.yml`:
```yaml
frontend:
  ports:
    - "3000:80"
```
Luego accede a http://localhost:3000

### Limpiar todo y empezar de cero

```bash
docker-compose down -v
docker system prune -a
docker-compose up -d --build
```

## Notas Importantes

1. **SQLite**: El backend usa SQLite. La base de datos se almacena en el volumen `backend-data`
2. **CORS**: El backend tiene CORS habilitado para todas las rutas
3. **Static Files**: Nginx sirve archivos estáticos con caché de 1 año
4. **API Proxy**: El frontend proxyfiere peticiones `/api/*` al backend

## Variables de Entorno Disponibles

```bash
# Backend
GIN_MODE=release  # Modo de Gin (release, debug)
PORT=8080         # Puerto del backend

# Frontend
VITE_API_URL      # URL del API (usado en tiempo de compilación)
```
