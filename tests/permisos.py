import urllib.request
import urllib.parse
import urllib.error
import json
import getpass
import sys

BASE_URL = "http://localhost:8080/api"

ROUTES_TO_TEST = [
    # Auth
    {"method": "GET", "path": "/auth/me", "name": "Mi Perfil"},
    # Products
    {"method": "GET", "path": "/products", "name": "Listar Productos"},
    {"method": "POST", "path": "/products", "name": "Crear Producto"},
    {"method": "PUT", "path": "/products/1", "name": "Editar Producto"},
    {"method": "DELETE", "path": "/products/1", "name": "Eliminar Producto"},
    {"method": "PATCH", "path": "/products/1/precios", "name": "Actualizar Precios"},
    # Sedes
    {"method": "GET", "path": "/sedes", "name": "Listar Sedes"},
    {"method": "POST", "path": "/sedes", "name": "Crear Sede"},
    # Stock
    {"method": "GET", "path": "/stock", "name": "Listar Stock"},
    {"method": "PUT", "path": "/stock", "name": "Actualizar Stock"},
    # RMAs
    {"method": "GET", "path": "/rmas", "name": "Listar RMAs"},
    {"method": "POST", "path": "/rmas", "name": "Crear RMA"},
    # Cotizaciones
    {"method": "GET", "path": "/cotizaciones", "name": "Listar Cotizaciones"},
    {"method": "POST", "path": "/cotizaciones", "name": "Crear Cotizacion"},
    {"method": "PUT", "path": "/cotizaciones/1/estado", "name": "Actualizar Cotizacion"},
    # Traspasos
    {"method": "GET", "path": "/traspasos", "name": "Listar Traspasos"},
    {"method": "POST", "path": "/traspasos", "name": "Crear Traspaso"},
    # Ordenes de Trabajo
    {"method": "GET", "path": "/ordenes-trabajo", "name": "Listar Ordenes Trabajo"},
    {"method": "POST", "path": "/ordenes-trabajo", "name": "Crear Orden Trabajo"},
    # Proveedores
    {"method": "GET", "path": "/proveedores", "name": "Listar Proveedores"},
    {"method": "POST", "path": "/proveedores", "name": "Crear Proveedor"},
    # Deudas
    {"method": "GET", "path": "/deudas", "name": "Listar Deudas"},
    {"method": "POST", "path": "/deudas", "name": "Crear Deuda"},
    {"method": "POST", "path": "/deudas/1/pago", "name": "Pagar Deuda"},
    # Insumos
    {"method": "GET", "path": "/insumos", "name": "Listar Insumos"},
    {"method": "POST", "path": "/insumos", "name": "Crear Insumo"},
    # Compatibilidad
    {"method": "GET", "path": "/compatibilidad", "name": "Listar Compatibilidad"},
    {"method": "POST", "path": "/compatibilidad", "name": "Crear Compatibilidad"},
    # Auditoria y Reportes
    {"method": "GET", "path": "/auditoria/logs", "name": "Ver Auditoria"},
    {"method": "GET", "path": "/reportes/ganancias", "name": "Ver Reporte de Ganancias"},
    # Segmentaciones
    {"method": "GET", "path": "/segmentaciones", "name": "Listar Segmentaciones"},
    {"method": "POST", "path": "/segmentaciones", "name": "Crear Segmentacion"},
    # Promociones
    {"method": "GET", "path": "/promociones", "name": "Listar Promociones"},
    {"method": "POST", "path": "/promociones", "name": "Crear Promocion"},
    # Users & Roles
    {"method": "GET", "path": "/users", "name": "Listar Usuarios"},
    {"method": "POST", "path": "/users", "name": "Crear Usuario"},
    {"method": "GET", "path": "/roles", "name": "Listar Roles"},
    {"method": "GET", "path": "/permissions", "name": "Listar Permisos"},
]

def login(email, password):
    url = f"{BASE_URL}/auth/login"
    data = json.dumps({"email": email, "password": password}).encode('utf-8')
    req = urllib.request.Request(url, data=data, headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req) as response:
            res = json.loads(response.read().decode('utf-8'))
            return res.get('token')
    except urllib.error.HTTPError as e:
        if e.code == 401:
            print("[𐄂] Credenciales incorrectas.")
            return None
        print(f"Error HTTP {e.code}")
        return None
    except Exception as e:
        print(f"Error de conexion al login: {e}")
        return None

def test_route(method, path, token):
    url = f"{BASE_URL}{path}"
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    
    data = None
    if method in ["POST", "PUT", "PATCH"]:
        data = b"{}"

    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req) as response:
            return response.getcode(), "[✓] PERMITIDO (Autorizado)"
    except urllib.error.HTTPError as e:
        if e.code == 403:
            return e.code, "[𐄂] DENEGADO (Falta Permiso)"
        elif e.code == 401:
            return e.code, "[?] NO AUTORIZADO (Token Invalido)"
        else:
            return e.code, f"[✓] PERMITIDO (Status {e.code}, pasó validación de permisos)"
    except Exception as e:
        return 0, f"Error de conexión: {e}"

def run_tests():
    print("==================================================")
    print("      TEST INTERACTIVO DE RUTAS Y PERMISOS        ")
    print("==================================================")
    print("\nCuentas de prueba disponibles en la base de datos (Seed):")
    print("  - Administrador : admin@inventario.com     / admin123")
    print("  - Gerente       : gerente@inventario.com   / gerente123")
    print("  - Vendedor      : vendedor@inventario.com  / vendedor123")
    print("  - Tecnico       : tecnico@inventario.com   / tecnico123")
    print("--------------------------------------------------\n")
    email = input("Ingrese el correo electrónico: ")
    password = getpass.getpass("Ingrese la contraseña: ")
    
    print(f"\n[*] Iniciando sesión como {email}...")
    token = login(email, password)
    if not token:
        print("No se pudo obtener el token. Saliendo...")
        sys.exit(1)
        
    print("[*] Inicio de sesión exitoso. Iniciando escaneo de rutas...\n")
    print(f"{'MÉTODO':<8} | {'RUTA':<35} | {'ESTADO':<10} | {'ACCESO'}")
    print("-" * 85)
    
    for route in ROUTES_TO_TEST:
        status, msg = test_route(route['method'], route['path'], token)
        status_str = f"HTTP {status}" if status > 0 else "ERROR"
        print(f"{route['method']:<8} | {route['path']:<35} | {status_str:<10} | {msg}")

if __name__ == "__main__":
    run_tests()
