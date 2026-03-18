# Configuración de PayPal

Para habilitar los pagos con PayPal en tu aplicación:

## 1. Crear cuenta de desarrollador de PayPal

1. Ve a https://developer.paypal.com
2. Inicia sesión o crea una cuenta
3. Ve a "Dashboard" → "Apps & Credentials"
4. En "Sandbox", crea una nueva aplicación
5. Copia el **Client id** y **Secret**

## 2. Configurar credenciales en el backend

En `/backend/controllers/order_controller.go`, reemplaza:

```go
Clientid: "YOUR_PAYPAL_CLIENT_id"
Secret:   "YOUR_PAYPAL_SECRET"
```

Con tus credenciales de Sandbox.

O mejor aún, usa variables de entorno:

```bash
export PAYPAL_CLIENT_id="tu_client_id_aqui"
export PAYPAL_SECRET="tu_secret_aqui"
```

## 3. Crear cuentas de prueba

En el Dashboard de PayPal → "Sandbox" → "Accounts":
- Crea una cuenta de **comprador** (buyer)
- Usa estas credenciales para hacer pagos de prueba

## 4. Probar la integración

1. Reinicia el servidor backend
2. Agrega productos al carrito
3. En la página del carrito verás los botones de PayPal
4. Usa la cuenta de comprador de Sandbox para completar el pago

## 5. Pasar a producción

Cuando estés listo para producción:

1. Cambia `BaseURL` en `order_controller.go`:
   ```go
   BaseURL: "https://api-m.paypal.com" // Producción
   ```

2. Obtén credenciales de **Live** en lugar de Sandbox

3. Asegúrate de que tu servidor tenga HTTPS

## Notas

- Los pagos en Sandbox usan dinero ficticio
- Las cuentas de prueba tienen $5000 USD por defecto
- Puedes crear múltiples cuentas de prueba para diferentes escenarios
