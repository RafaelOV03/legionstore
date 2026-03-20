# Dockerfile para desarrollo con go run
FROM golang:1.24-alpine

# Instalar dependencias necesarias
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Exponer puerto
EXPOSE 8080

# Comando para ejecutar en modo desarrollo
CMD ["go", "run", "./"]
