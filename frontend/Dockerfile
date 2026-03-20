FROM node:22-alpine

WORKDIR /app

# Copiar archivos de dependencias
COPY package*.json ./

# Instalar dependencias
RUN npm ci

# Copiar código fuente
COPY . .

# Exponer puerto para Vite dev server
EXPOSE 5173

# Comando para ejecutar en modo desarrollo
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
