/**
 * resourceApi.js - Consolidación de todas las operaciones de API
 * Unifica api.js e inventarioApi.js usando apiClient centralizado
 */

import { apiClient } from './apiClient';

// ==========================================
// UTILIDADES
// ==========================================

/**
 * Generar o obtener un id de sesión único del carrito
 */
export const getSessionId = () => {
  let sessionId = localStorage.getItem('sessionId');
  if (!sessionId) {
    sessionId = 'session_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
    localStorage.setItem('sessionId', sessionId);
  }
  return sessionId;
};

// ==========================================
// AUTENTICACIÓN
// ==========================================

export const login = (email, password) => 
  apiClient.post('/auth/login', { email, password });

export const register = (userData) => 
  apiClient.post('/auth/register', userData);

export const getCurrentUser = () => 
  apiClient.get('/auth/me');

// ==========================================
// PRODUCTOS
// ==========================================

export const getProducts = () => 
  apiClient.get('/products');

export const getProduct = (id) => 
  apiClient.get(`/products/${id}`);

export const createProduct = (product) => 
  apiClient.post('/products', product);

export const updateProduct = (id, product) => 
  apiClient.put(`/products/${id}`, product);

export const deleteProduct = (id) => 
  apiClient.delete(`/products/${id}`);

export const getRelatedProducts = (id) => 
  apiClient.get(`/products/${id}/related`);

export const getRandomProducts = (limit = 8) => 
  apiClient.get(`/products/random?limit=${limit}`);

export const getBestSellingProducts = (limit = 8) => 
  apiClient.get(`/products/bestsellers?limit=${limit}`);

export const getNewProducts = (limit = 8) => 
  apiClient.get(`/products/new?limit=${limit}`);

export const getFeaturedProducts = (limit = 8) => 
  apiClient.get(`/products/featured?limit=${limit}`);

export const getProductsByCategory = (category) => 
  apiClient.get(`/products/category/${encodeURIComponent(category)}`);

// ==========================================
// CARRITO
// ==========================================

export const getCart = () => {
  const sessionId = getSessionId();
  return apiClient.get(`/cart?session_id=${sessionId}`);
};

export const addToCart = (productId, quantity) => {
  const sessionId = getSessionId();
  return apiClient.post(`/cart/items?session_id=${sessionId}`, { 
    product_id: productId, 
    quantity 
  });
};

export const updateCartItem = (itemId, quantity) => {
  const sessionId = getSessionId();
  return apiClient.put(`/cart/items/${itemId}?session_id=${sessionId}`, { quantity });
};

export const removeFromCart = (itemId) => {
  const sessionId = getSessionId();
  return apiClient.delete(`/cart/items/${itemId}?session_id=${sessionId}`);
};

export const clearCart = () => {
  const sessionId = getSessionId();
  return apiClient.delete(`/cart/clear?session_id=${sessionId}`);
};

// ==========================================
// ÓRDENES Y PAYPAL
// ==========================================

export const getPayPalConfig = () => 
  apiClient.get('/orders/config');

export const createOrder = () => {
  const sessionId = getSessionId();
  return apiClient.post('/orders', { session_id: sessionId });
};

export const captureOrder = (orderId) => 
  apiClient.post(`/orders/${orderId}/capture`, {});

export const getOrder = (orderId) => 
  apiClient.get(`/orders/${orderId}`);

export const getOrders = () => {
  const sessionId = getSessionId();
  return apiClient.get(`/orders?session_id=${sessionId}`);
};

export const getAllOrders = () => 
  apiClient.get(`/orders/all`);

export const finalizeOrder = (orderId) => 
  apiClient.put(`/orders/${orderId}/finalize`, {});

export const updateOrder = (orderId, data) => 
  apiClient.put(`/orders/${orderId}`, data);

export const deleteOrder = (orderId) => 
  apiClient.delete(`/orders/${orderId}`);

// ==========================================
// UPLOAD DE IMÁGENES
// ==========================================

export const uploadImage = async (file) => {
  const formData = new FormData();
  formData.append('image', file);
  return apiClient.requestFormData('/upload/image', formData, { 
    method: 'POST' 
  });
};

export const deleteImage = (imageUrl) => 
  apiClient.delete('/upload/image', { 
    body: { image_url: imageUrl } 
  });

// ==========================================
// SEDES Y STOCK
// ==========================================

export const getSedes = () => 
  apiClient.get('/sedes');

export const getSede = (id) => 
  apiClient.get(`/sedes/${id}`);

export const createSede = (sede) => 
  apiClient.post('/sedes', sede);

export const updateSede = (id, sede) => 
  apiClient.put(`/sedes/${id}`, sede);

export const deleteSede = (id) => 
  apiClient.delete(`/sedes/${id}`);

// Stock multisede
export const getStockMultisede = () => 
  apiClient.get('/stock');

export const getStockBySede = (sedeId) => 
  apiClient.get(`/stock/sede/${sedeId}`);

export const updateStock = (stockData) => 
  apiClient.put('/stock', stockData);

// ==========================================
// RMA / GARANTÍAS
// ==========================================

export const getRMAs = (params = {}) => {
  const query = new URLSearchParams(params).toString();
  return apiClient.get(`/rmas${query ? '?' + query : ''}`);
};

export const getRMA = (id) => 
  apiClient.get(`/rmas/${id}`);

export const getRMAStats = () => 
  apiClient.get('/rmas/stats');

export const createRMA = (rma) => 
  apiClient.post('/rmas', rma);

export const updateRMA = (id, rma) => 
  apiClient.put(`/rmas/${id}`, rma);

export const deleteRMA = (id) => 
  apiClient.delete(`/rmas/${id}`);

// ==========================================
// COTIZACIONES
// ==========================================

export const getCotizaciones = (params = {}) => {
  const query = new URLSearchParams(params).toString();
  return apiClient.get(`/cotizaciones${query ? '?' + query : ''}`);
};

export const getCotizacion = (id) => 
  apiClient.get(`/cotizaciones/${id}`);

export const getCotizacionPDF = (id) => 
  apiClient.get(`/cotizaciones/${id}/pdf`);

export const createCotizacion = (cotizacion) => 
  apiClient.post('/cotizaciones', cotizacion);

export const updateCotizacionEstado = (id, estado) => 
  apiClient.put(`/cotizaciones/${id}/estado`, { estado });

export const convertirCotizacionAVenta = (id) => 
  apiClient.post(`/cotizaciones/${id}/convertir`, {});

export const deleteCotizacion = (id) => 
  apiClient.delete(`/cotizaciones/${id}`);

// Alias para compatibilidad
export const convertirCotizacion = convertirCotizacionAVenta;

// ==========================================
// TRASPASOS
// ==========================================

export const getTraspasos = (params = {}) => {
  const query = new URLSearchParams(params).toString();
  return apiClient.get(`/traspasos${query ? '?' + query : ''}`);
};

export const getTraspaso = (id) => 
  apiClient.get(`/traspasos/${id}`);

export const createTraspaso = (traspaso) => 
  apiClient.post('/traspasos', traspaso);

export const enviarTraspaso = (id) => 
  apiClient.post(`/traspasos/${id}/enviar`, {});

export const recibirTraspaso = (id, data) => 
  apiClient.post(`/traspasos/${id}/recibir`, data);

export const cancelarTraspaso = (id) => 
  apiClient.post(`/traspasos/${id}/cancelar`, {});

export const deleteTraspaso = (id) => 
  apiClient.delete(`/traspasos/${id}`);

// ==========================================
// ÓRDENES DE TRABAJO
// ==========================================

export const getOrdenesTrabajo = (params = {}) => {
  const query = new URLSearchParams(params).toString();
  return apiClient.get(`/ordenes-trabajo${query ? '?' + query : ''}`);
};

export const getOrdenTrabajo = (id) => 
  apiClient.get(`/ordenes-trabajo/${id}`);

export const getOrdenesStats = () => 
  apiClient.get('/ordenes-trabajo/stats');

export const getTecnicos = () => 
  apiClient.get('/ordenes-trabajo/tecnicos');

export const createOrdenTrabajo = (orden) => 
  apiClient.post('/ordenes-trabajo', orden);

export const updateOrdenTrabajo = (id, orden) => 
  apiClient.put(`/ordenes-trabajo/${id}`, orden);

export const asignarTecnico = (id, tecnicoId) => 
  apiClient.post(`/ordenes-trabajo/${id}/asignar`, { tecnico_id: tecnicoId });

export const agregarInsumo = (id, insumoData) => 
  apiClient.post(`/ordenes-trabajo/${id}/insumo`, insumoData);

// Alias para compatibilidad
export const agregarInsumoOrden = agregarInsumo;

export const registrarTrazabilidad = (id, data) => 
  apiClient.post(`/ordenes-trabajo/${id}/trazabilidad`, data);

export const deleteOrdenTrabajo = (id) => 
  apiClient.delete(`/ordenes-trabajo/${id}`);

// ==========================================
// PROVEEDORES
// ==========================================

export const getProveedores = () => 
  apiClient.get('/proveedores');

export const getProveedor = (id) => 
  apiClient.get(`/proveedores/${id}`);

export const createProveedor = (proveedor) => 
  apiClient.post('/proveedores', proveedor);

export const updateProveedor = (id, proveedor) => 
  apiClient.put(`/proveedores/${id}`, proveedor);

export const deleteProveedor = (id) => 
  apiClient.delete(`/proveedores/${id}`);

// ==========================================
// DEUDAS Y PAGOS
// ==========================================

export const getDeudas = (params = {}) => {
  const query = new URLSearchParams(params).toString();
  return apiClient.get(`/deudas${query ? '?' + query : ''}`);
};

export const getResumenDeudas = () => 
  apiClient.get('/deudas/resumen');

export const createDeuda = (deuda) => 
  apiClient.post('/deudas', deuda);

export const registrarPago = (deudaId, pago) => 
  apiClient.post(`/deudas/${deudaId}/pago`, pago);

export const getPagosDeuda = (deudaId) => 
  apiClient.get(`/deudas/${deudaId}/pagos`);

// ==========================================
// INSUMOS
// ==========================================

export const getInsumos = (params = {}) => {
  const query = new URLSearchParams(params).toString();
  return apiClient.get(`/insumos${query ? '?' + query : ''}`);
};

export const getInsumo = (id) => 
  apiClient.get(`/insumos/${id}`);

export const getInsumosStats = () => 
  apiClient.get('/insumos/stats');

export const createInsumo = (insumo) => 
  apiClient.post('/insumos', insumo);

export const updateInsumo = (id, insumo) => 
  apiClient.put(`/insumos/${id}`, insumo);

export const ajustarStockInsumo = (id, data) => 
  apiClient.post(`/insumos/${id}/ajustar`, data);

export const deleteInsumo = (id) => 
  apiClient.delete(`/insumos/${id}`);

// ==========================================
// COMPATIBILIDAD
// ==========================================

export const getCompatibilidades = (params = {}) => {
  const query = new URLSearchParams(params).toString();
  return apiClient.get(`/compatibilidad${query ? '?' + query : ''}`);
};

export const buscarCompatibles = (productoId) => 
  apiClient.get(`/compatibilidad/buscar/${productoId}`);

export const createCompatibilidad = (data) => 
  apiClient.post('/compatibilidad', data);

export const deleteCompatibilidad = (id) => 
  apiClient.delete(`/compatibilidad/${id}`);

// ==========================================
// AUDITORÍA Y REPORTES
// ==========================================

export const getLogs = (params = {}) => {
  const query = new URLSearchParams(params).toString();
  return apiClient.get(`/auditoria/logs${query ? '?' + query : ''}`);
};

// Alias para compatibilidad
export const getAuditoriaLogs = getLogs;

export const getLogStats = () => 
  apiClient.get('/auditoria/stats');

export const getReporteGanancias = (params = {}) => {
  const query = new URLSearchParams(params).toString();
  return apiClient.get(`/reportes/ganancias${query ? '?' + query : ''}`);
};

// ==========================================
// SEGMENTACIÓN Y PROMOCIONES
// ==========================================

export const getSegmentaciones = () => 
  apiClient.get('/segmentaciones');

export const createSegmentacion = (segmentacion) => 
  apiClient.post('/segmentaciones', segmentacion);

export const getPromociones = (params = {}) => {
  const query = new URLSearchParams(params).toString();
  return apiClient.get(`/promociones${query ? '?' + query : ''}`);
};

export const createPromocion = (promocion) => 
  apiClient.post('/promociones', promocion);

export const updatePromocion = (id, promocion) => 
  apiClient.put(`/promociones/${id}`, promocion);

export const deletePromocion = (id) => 
  apiClient.delete(`/promociones/${id}`);

// ==========================================
// USUARIOS Y ROLES
// ==========================================

export const getUsers = () => 
  apiClient.get('/users');

export const getUser = (id) => 
  apiClient.get(`/users/${id}`);

export const createUser = (user) => 
  apiClient.post('/users', user);

export const updateUser = (id, user) => 
  apiClient.put(`/users/${id}`, user);

export const deleteUser = (id) => 
  apiClient.delete(`/users/${id}`);

export const getRoles = () => 
  apiClient.get('/roles');

export const getRole = (id) => 
  apiClient.get(`/roles/${id}`);

export const createRole = (role) => 
  apiClient.post('/roles', role);

export const updateRole = (id, role) => 
  apiClient.put(`/roles/${id}`, role);

export const deleteRole = (id) => 
  apiClient.delete(`/roles/${id}`);

export const getPermissions = () => 
  apiClient.get('/permissions');
