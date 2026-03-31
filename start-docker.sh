#!/bin/bash

# Script de inicio rápido para LegionStore con Docker

set -e

echo "🐳 LegionStore - Docker Compose Startup"
echo "========================================"

# Verificar si Docker está instalado
if ! command -v docker &> /dev/null; then
    echo "❌ Docker no está instalado. Por favor instala Docker."
    exit 1
fi

# Verificar si Docker Compose está instalado
if ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose no está instalado. Por favor instala Docker Compose."
    exit 1
fi

# Crear .env si no existe
if [ ! -f .env ]; then
    echo "📝 Creando archivo .env..."
    cp .env.example .env
    echo "✅ Archivo .env creado"
fi

# Construir imágenes
echo ""
echo "🔨 Construyendo imágenes Docker..."
docker compose build

# Iniciar servicios
echo ""
echo "🚀 Iniciando servicios..."
docker compose up -d

# Esperar a que los servicios estén listos
echo ""
echo "⏳ Esperando a que los servicios estén listos..."
sleep 5

# Mostrar estado
echo ""
echo "📊 Estado de los servicios:"
docker compose ps

# Mostrar logs iniciales
echo ""
echo "📋 Logs recientes:"
docker-compose logs --tail=20

echo ""
echo "✅ ¡LegionStore está listo!"
echo ""
echo "🌐 URLs de acceso:"
echo "   Frontend: http://localhost"
echo "   Backend API: http://localhost:8080/api"
echo ""
echo "💡 Para ver logs en vivo:"
echo "   docker-compose logs -f"
echo ""
echo "🛑 Para detener:"
echo "   docker-compose down"
