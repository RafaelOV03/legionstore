// Usar URL relativa para aprovechar el proxy de Vite
const API_URL = '/api';

// Obtener el token de localStorage
const getAuthHeader = () => {
  const token = localStorage.getItem('token');
  if (!token) {
    console.warn('No authentication token found in localStorage');
  }
  return token ? { 'Authorization': `Bearer ${token}` } : {};
};

// Generar o obtener un id de sesión único
export const getSessionId = () => {
  let sessionId = localStorage.getItem('sessionId');
  if (!sessionId) {
    sessionId = 'session_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
    localStorage.setItem('sessionId', sessionId);
  }
  return sessionId;
};

// Productos
export const getProducts = async () => {
  const response = await fetch(`${API_URL}/products`);
  return response.json();
};

export const getProduct = async (id) => {
  const response = await fetch(`${API_URL}/products/${id}`);
  return response.json();
};

export const createProduct = async (product) => {
  const response = await fetch(`${API_URL}/products`, {
    method: 'POST',
    headers: { 
      'Content-Type': 'application/json',
      ...getAuthHeader()
    },
    body: JSON.stringify(product),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to create product');
  }
  
  return response.json();
};

export const updateProduct = async (id, product) => {
  const authHeaders = getAuthHeader();
  console.log('Updating product with headers:', authHeaders); // Debug
  
  const response = await fetch(`${API_URL}/products/${id}`, {
    method: 'PUT',
    headers: { 
      'Content-Type': 'application/json',
      ...authHeaders
    },
    body: JSON.stringify(product),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to update product');
  }
  
  return response.json();
};

export const deleteProduct = async (id) => {
  const response = await fetch(`${API_URL}/products/${id}`, {
    method: 'DELETE',
    headers: getAuthHeader(),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to delete product');
  }
  
  return response.json();
};

export const getRelatedProducts = async (id) => {
  const response = await fetch(`${API_URL}/products/${id}/related`);
  
  if (!response.ok) {
    throw new Error('Failed to fetch related products');
  }
  
  return response.json();
};

export const getRandomProducts = async (limit = 8) => {
  const response = await fetch(`${API_URL}/products/random?limit=${limit}`);
  
  if (!response.ok) {
    throw new Error('Failed to fetch random products');
  }
  
  return response.json();
};

export const getBestSellingProducts = async (limit = 8) => {
  const response = await fetch(`${API_URL}/products/bestsellers?limit=${limit}`);
  
  if (!response.ok) {
    throw new Error('Failed to fetch best selling products');
  }
  
  return response.json();
};

export const getNewProducts = async (limit = 8) => {
  const response = await fetch(`${API_URL}/products/new?limit=${limit}`);
  
  if (!response.ok) {
    throw new Error('Failed to fetch new products');
  }
  
  return response.json();
};

export const getFeaturedProducts = async (limit = 8) => {
  const response = await fetch(`${API_URL}/products/featured?limit=${limit}`);
  
  if (!response.ok) {
    throw new Error('Failed to fetch featured products');
  }
  
  return response.json();
};

export const getProductsByCategory = async (category) => {
  const response = await fetch(`${API_URL}/products/category/${encodeURIComponent(category)}`);
  
  if (!response.ok) {
    throw new Error('Failed to fetch products by category');
  }
  
  return response.json();
};

// Upload de imágenes
export const uploadImage = async (file) => {
  const formData = new FormData();
  formData.append('image', file);
  
  const response = await fetch(`${API_URL}/upload/image`, {
    method: 'POST',
    headers: getAuthHeader(),
    body: formData,
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to upload image');
  }
  
  return response.json();
};

export const deleteImage = async (url) => {
  const response = await fetch(`${API_URL}/upload/image`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeader()
    },
    body: JSON.stringify({ url }),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to delete image');
  }
  
  return response.json();
};

// Carrito
export const getCart = async () => {
  const sessionId = getSessionId();
  const response = await fetch(`${API_URL}/cart?session_id=${sessionId}`);
  return response.json();
};

export const addToCart = async (productId, quantity) => {
  const sessionId = getSessionId();
  const response = await fetch(`${API_URL}/cart/items?session_id=${sessionId}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ product_id: productId, quantity }),
  });
  return response.json();
};

export const updateCartItem = async (itemId, quantity) => {
  const sessionId = getSessionId();
  const response = await fetch(`${API_URL}/cart/items/${itemId}?session_id=${sessionId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ quantity }),
  });
  return response.json();
};

export const removeFromCart = async (itemId) => {
  const sessionId = getSessionId();
  const response = await fetch(`${API_URL}/cart/items/${itemId}?session_id=${sessionId}`, {
    method: 'DELETE',
  });
  return response.json();
};

export const clearCart = async () => {
  const sessionId = getSessionId();
  const response = await fetch(`${API_URL}/cart/clear?session_id=${sessionId}`, {
    method: 'DELETE',
  });
  return response.json();
};

// Órdenes y PayPal
export const getPayPalConfig = async () => {
  const response = await fetch(`${API_URL}/orders/config`);
  return response.json();
};

export const createOrder = async () => {
  const sessionId = getSessionId();
  const response = await fetch(`${API_URL}/orders`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ session_id: sessionId }),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to create order');
  }
  
  return response.json();
};

export const captureOrder = async (orderId) => {
  const response = await fetch(`${API_URL}/orders/${orderId}/capture`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to capture order');
  }
  
  return response.json();
};

export const getOrder = async (orderId) => {
  const response = await fetch(`${API_URL}/orders/${orderId}`);
  return response.json();
};

export const getOrders = async () => {
  const sessionId = getSessionId();
  const response = await fetch(`${API_URL}/orders?session_id=${sessionId}`);
  return response.json();
};

export const getAllOrders = async () => {
  const response = await fetch(`${API_URL}/orders/all`, {
    headers: getAuthHeader()
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch all orders');
  }
  
  return response.json();
};

export const finalizeOrder = async (orderId) => {
  const response = await fetch(`${API_URL}/orders/${orderId}/finalize`, {
    method: 'PUT',
    headers: getAuthHeader()
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to finalize order');
  }
  
  return response.json();
};

export const updateOrder = async (orderId, data) => {
  const response = await fetch(`${API_URL}/orders/${orderId}`, {
    method: 'PUT',
    headers: { 
      'Content-Type': 'application/json',
      ...getAuthHeader()
    },
    body: JSON.stringify(data),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to update order');
  }
  
  return response.json();
};

export const deleteOrder = async (orderId) => {
  const response = await fetch(`${API_URL}/orders/${orderId}`, {
    method: 'DELETE',
    headers: getAuthHeader(),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to delete order');
  }
  
  return response.json();
};

// Usuarios
export const getUsers = async () => {
  const response = await fetch(`${API_URL}/users`, {
    headers: getAuthHeader(),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch users');
  }
  
  return response.json();
};

export const getUser = async (userId) => {
  const response = await fetch(`${API_URL}/users/${userId}`, {
    headers: getAuthHeader(),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch user');
  }
  
  return response.json();
};

export const createUser = async (userData) => {
  const response = await fetch(`${API_URL}/users`, {
    method: 'POST',
    headers: { 
      'Content-Type': 'application/json',
      ...getAuthHeader()
    },
    body: JSON.stringify(userData),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to create user');
  }
  
  return response.json();
};

export const updateUser = async (userId, userData) => {
  const response = await fetch(`${API_URL}/users/${userId}`, {
    method: 'PUT',
    headers: { 
      'Content-Type': 'application/json',
      ...getAuthHeader()
    },
    body: JSON.stringify(userData),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to update user');
  }
  
  return response.json();
};

export const deleteUser = async (userId) => {
  const response = await fetch(`${API_URL}/users/${userId}`, {
    method: 'DELETE',
    headers: getAuthHeader(),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to delete user');
  }
  
  return response.json();
};

// Roles
export const getRoles = async () => {
  const response = await fetch(`${API_URL}/roles`, {
    headers: getAuthHeader(),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch roles');
  }
  
  return response.json();
};

export const getRole = async (roleId) => {
  const response = await fetch(`${API_URL}/roles/${roleId}`, {
    headers: getAuthHeader(),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch role');
  }
  
  return response.json();
};

export const createRole = async (roleData) => {
  const response = await fetch(`${API_URL}/roles`, {
    method: 'POST',
    headers: { 
      'Content-Type': 'application/json',
      ...getAuthHeader()
    },
    body: JSON.stringify(roleData),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to create role');
  }
  
  return response.json();
};

export const updateRole = async (roleId, roleData) => {
  const response = await fetch(`${API_URL}/roles/${roleId}`, {
    method: 'PUT',
    headers: { 
      'Content-Type': 'application/json',
      ...getAuthHeader()
    },
    body: JSON.stringify(roleData),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to update role');
  }
  
  return response.json();
};

export const deleteRole = async (roleId) => {
  const response = await fetch(`${API_URL}/roles/${roleId}`, {
    method: 'DELETE',
    headers: getAuthHeader(),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to delete role');
  }
  
  return response.json();
};

// Permisos
export const getPermissions = async () => {
  const response = await fetch(`${API_URL}/permissions`, {
    headers: getAuthHeader(),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch permissions');
  }
  
  return response.json();
};
