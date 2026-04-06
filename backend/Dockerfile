# Dockerfile para desarrollo con go run
FROM golang:1.25-alpine

# Instalar dependencias necesarias del sistema
RUN apk add --no-cache gcc musl-dev sqlite-dev wget

WORKDIR /app

# Copiar archivos de módulos
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Exponer puerto
EXPOSE 8080

# Comando para ejecutar en modo desarrollo
CMD ["go", "run", "./"]
