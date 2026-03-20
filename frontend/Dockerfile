FROM node:22-alpine

WORKDIR /app

# Exponer puerto para Vite dev server
EXPOSE 5173

# Comando para ejecutar en modo desarrollo
CMD sh -c "npm install && npm run dev -- --host 0.0.0.0"
