@echo off
REM Script de inicio rápido para LegionStore con Docker (Windows)

setlocal enabledelayedexpansion

echo.
echo 🐳 LegionStore - Docker Compose Startup
echo ========================================

REM Verificar si Docker está instalado
where docker >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Docker no está instalado. Por favor instala Docker Desktop.
    pause
    exit /b 1
)

REM Verificar si Docker Compose está instalado
where docker-compose >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Docker Compose no está instalado.
    pause
    exit /b 1
)

REM Crear .env si no existe
if not exist .env (
    echo 📝 Creando archivo .env...
    copy .env.example .env
    echo ✅ Archivo .env creado
)

REM Construir imágenes
echo.
echo 🔨 Construyendo imágenes Docker...
docker-compose build

REM Iniciar servicios
echo.
echo 🚀 Iniciando servicios...
docker-compose up -d

REM Esperar a que los servicios estén listos
echo.
echo ⏳ Esperando a que los servicios estén listos...
timeout /t 5

REM Mostrar estado
echo.
echo 📊 Estado de los servicios:
docker-compose ps

REM Mostrar logs iniciales
echo.
echo 📋 Logs recientes:
docker-compose logs --tail=20

echo.
echo ✅ ¡LegionStore está listo!
echo.
echo 🌐 URLs de acceso:
echo    Frontend: http://localhost
echo    Backend API: http://localhost:8080/api
echo.
echo 💡 Para ver logs en vivo:
echo    docker-compose logs -f
echo.
echo 🛑 Para detener:
echo    docker-compose down
echo.
pause
