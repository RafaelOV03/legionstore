/**
 * apiClient.js - Centralizado API client con manejo global de errores y autenticación
 * Proporciona una interfaz consistente para todas las llamadas HTTP
 */

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

/**
 * Obtener header de autenticación desde localStorage
 */
function getAuthHeader() {
  const token = localStorage.getItem('token');
  return token ? { 'Authorization': `Bearer ${token}` } : {};
}

/**
 * Cliente API centralizado con manejo de errores y autenticación
 */
export const apiClient = {
  /**
   * Realiza una petición HTTP genérica
   * @param {string} endpoint - Ruta del endpoint (ej: '/products')
   * @param {object} options - Opciones adicionales (method, body, headers, etc)
   * @returns {Promise<object>} Respuesta parseada como JSON
   * @throws {Error} Si la petición falla o retorna error
   */
  async request(endpoint, options = {}) {
    const url = `${API_URL}${endpoint}`;
    
    const headers = {
      'Content-Type': 'application/json',
      ...getAuthHeader(),
      ...options.headers,
    };

    const config = {
      ...options,
      headers,
    };

    // Si el body es un objeto, convertir a JSON (a menos que sea FormData)
    if (config.body && typeof config.body === 'object' && !(config.body instanceof FormData)) {
      config.body = JSON.stringify(config.body);
    }

    try {
      const response = await fetch(url, config);

      // Manejar 401 Unauthorized globalmente
      if (response.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
        throw new Error('Session expired. Please login again.');
      }

      // Manejar respuestas no OK
      if (!response.ok) {
        let error;
        try {
          error = await response.json();
        } catch {
          error = { error: `HTTP ${response.status}` };
        }
        throw new Error(error.error || `HTTP ${response.status}`);
      }

      // Parsear respuesta como JSON (o retornar objeto vacío si está vacío)
      try {
        return await response.json();
      } catch {
        return {};
      }
    } catch (error) {
      // Re-lanzar error para que el componente lo maneje
      throw error;
    }
  },

  /**
   * GET request
   */
  get(endpoint, options = {}) {
    return this.request(endpoint, { ...options, method: 'GET' });
  },

  /**
   * POST request
   */
  post(endpoint, data, options = {}) {
    return this.request(endpoint, { 
      ...options, 
      method: 'POST', 
      body: data 
    });
  },

  /**
   * PUT request
   */
  put(endpoint, data, options = {}) {
    return this.request(endpoint, { 
      ...options, 
      method: 'PUT', 
      body: data 
    });
  },

  /**
   * DELETE request
   */
  delete(endpoint, options = {}) {
    return this.request(endpoint, { ...options, method: 'DELETE' });
  },

  /**
   * PATCH request
   */
  patch(endpoint, data, options = {}) {
    return this.request(endpoint, { 
      ...options, 
      method: 'PATCH', 
      body: data 
    });
  },

  /**
   * Petición sin headers JSON (para FormData en uploads)
   */
  async requestFormData(endpoint, formData, options = {}) {
    const url = `${API_URL}${endpoint}`;
    
    const config = {
      method: options.method || 'POST',
      headers: getAuthHeader(),
      body: formData,
      ...options,
    };

    try {
      const response = await fetch(url, config);

      if (response.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
        throw new Error('Session expired. Please login again.');
      }

      if (!response.ok) {
        let error;
        try {
          error = await response.json();
        } catch {
          error = { error: `HTTP ${response.status}` };
        }
        throw new Error(error.error || `HTTP ${response.status}`);
      }

      try {
        return await response.json();
      } catch {
        return {};
      }
    } catch (error) {
      throw error;
    }
  },
};

export default apiClient;
