import urllib.request
import urllib.error
import time
import concurrent.futures
import sys

BASE_URL = "http://localhost:8080/api"
ENDPOINT = "/products"
TOTAL_REQUESTS = 5000
CONCURRENT_THREADS = 100

def make_request(req_id):
    url = f"{BASE_URL}{ENDPOINT}"
    req = urllib.request.Request(url, method="GET")
    start_time = time.time()
    try:
        with urllib.request.urlopen(req, timeout=10) as response:
            status = response.getcode()
            response.read() # Consumir el body
            end_time = time.time()
            return status, end_time - start_time, None
    except urllib.error.HTTPError as e:
        end_time = time.time()
        return e.code, end_time - start_time, None
    except Exception as e:
        end_time = time.time()
        return 0, end_time - start_time, str(e)

def run_stress_test():
    print("==================================================")
    print("           TEST DE ESTRÉS / RENDIMIENTO           ")
    print("==================================================")
    print(f"URL Objetivo      : {BASE_URL}{ENDPOINT}")
    print(f"Total de Peticiones: {TOTAL_REQUESTS}")
    print(f"Hilos Concurrentes: {CONCURRENT_THREADS}")
    print("Iniciando prueba...\n")

    start_total_time = time.time()

    results = []
    with concurrent.futures.ThreadPoolExecutor(max_workers=CONCURRENT_THREADS) as executor:
        futures = [executor.submit(make_request, i) for i in range(TOTAL_REQUESTS)]
        
        completed = 0
        for future in concurrent.futures.as_completed(futures):
            results.append(future.result())
            completed += 1
            if completed % 50 == 0:
                print(f"Progreso: {completed}/{TOTAL_REQUESTS} peticiones completadas...")

    end_total_time = time.time()

    # Analizar resultados
    successful = 0
    failed = 0
    total_time = 0
    min_time = float('inf')
    max_time = 0
    errors = {}

    for status, duration, err in results:
        total_time += duration
        if duration < min_time:
            min_time = duration
        if duration > max_time:
            max_time = duration
        
        if status in [200, 201, 204]:
            successful += 1
        else:
            failed += 1
            if err:
                errors[err] = errors.get(err, 0) + 1
            else:
                key = f"HTTP {status}"
                errors[key] = errors.get(key, 0) + 1

    avg_time = (total_time / TOTAL_REQUESTS) * 1000
    total_duration = end_total_time - start_total_time
    req_per_sec = TOTAL_REQUESTS / total_duration

    print("\n==================================================")
    print("                 RESULTADOS FINALES               ")
    print("==================================================")
    print(f"Tiempo Total de la Prueba : {total_duration:.2f} segundos")
    print(f"Troughput (Req/sec)       : {req_per_sec:.2f} peticiones/seg")
    print(f"Peticiones Exitosas       : {successful}")
    print(f"Peticiones Fallidas       : {failed}")
    print("-" * 50)
    print("Tiempos de Respuesta:")
    print(f"  - Mínimo    : {min_time * 1000:.2f} ms")
    print(f"  - Máximo    : {max_time * 1000:.2f} ms")
    print(f"  - Promedio  : {avg_time:.2f} ms")
    
    if failed > 0:
        print("-" * 50)
        print("Resumen de Errores:")
        for err_msg, count in errors.items():
            print(f"  - {err_msg}: {count} veces")
    print("==================================================")

if __name__ == "__main__":
    run_stress_test()
