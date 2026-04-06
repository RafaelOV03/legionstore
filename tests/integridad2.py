import subprocess
import sys
import os

CONTAINER_NAME = "legionstore-backend"

def run_integrity_test():
    print("==================================================")
    print("    TEST DE INTEGRIDAD Y CONDICIONES DE CARRERA   ")
    print("==================================================")
    print(f"[*] Ejecutando con Docker en el contenedor '{CONTAINER_NAME}'...\n")
    print("> go test -race -v ./\n")

    # Comando para ejecutar las pruebas con el detector de carreras (race detector) en todos los paquetes
    os.chdir("backend")
    command = [
        "go", "test", "-race", "-v"
    ]

    try:
        # Ejecutamos el comando capturando la salida en vivo para ver el progreso
        process = subprocess.Popen(
            command,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            bufsize=1
        )

        # Imprimir salida directamente conforme se generan los resultados
        for line in process.stdout:
            print(line, end="")

        process.wait()

        print("\n==================================================")
        if process.returncode == 0:
            print("[✓] PRUEBA SUPERADA: No se encontraron condiciones de carrera (Race Conditions) ni errores de integridad.")
        else:
            print(f"[𐄂] PRUEBA FALLIDA: Se detectaron problemas de integridad o el contenedor no está corriendo (Código {process.returncode}).")
            sys.exit(process.returncode)

    except FileNotFoundError:
        print("[𐄂] ERROR: El comando 'docker' no fue encontrado. Asegurate de tener Docker instalado.")
        sys.exit(1)
    except Exception as e:
        print(f"[𐄂] ERROR INESPERADO al ejecutar Docker: {e}")
        sys.exit(1)

if __name__ == "__main__":
    run_integrity_test()
