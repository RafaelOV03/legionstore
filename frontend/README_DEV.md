# Legion Store - Frontend Development Guide

**Framework**: React 18  
**Build Tool**: Vite  
**CSS Framework**: Bootstrap 5  
**State Management**: React Context API  
**Routing**: React Router v6

---

## Table of Contents

1. [Setup & Installation](#setup--installation)
2. [Project Structure](#project-structure)
3. [Development Workflow](#development-workflow)
4. [Architecture Overview](#architecture-overview)
5. [Key Components](#key-components)
6. [Custom Hooks](#custom-hooks)
7. [Context API Usage](#context-api-usage)
8. [Styling](#styling)
9. [API Integration](#api-integration)
10. [Error Handling](#error-handling)

---

## Setup & Installation

### Prerequisites
- Node.js 18+ or 22+
- npm 9+
- Git

### Installation Steps

```bash
# 1. Navigate to frontend directory
cd frontend

# 2. Install dependencies
npm install

# 3. Start development server
npm run dev

# 4. Open in browser
# http://localhost:5173
```

### Environment Variables

**File**: `.env` (create if doesn't exist)

```env
# Backend API
VITE_API_URL=http://localhost:8080/api

# Optional: Build optimization
VITE_BUILD_MODE=development
```

### Build for Production

```bash
npm run build

# Preview production build
npm run preview
```

---

## Project Structure

```
frontend/
├── src/
│   ├── assets/              # Static assets
│   │   └── bootstrap-custom.css
│   ├── components/          # Reusable components
│   │   ├── Navigation.jsx          # Consolidated navbar
│   │   ├── ProtectedRoute.jsx      # Route protection wrapper
│   │   └── ProductCard.jsx         # Product display card
│   ├── context/             # React Context
│   │   └── AuthContext.jsx         # Authentication state
│   ├── hooks/               # Custom React hooks
│   │   └── useCart.js              # Shopping cart logic
│   ├── pages/               # Full-page components
│   │   ├── Admin.jsx               # Admin dashboard
│   │   ├── Dashboard.jsx           # User dashboard
│   │   ├── Productos.jsx           # Product listing
│   │   ├── Login.jsx               # Login page
│   │   ├── Cart.jsx                # Shopping cart
│   │   └── ...                     # Other pages
│   ├── services/            # API integration
│   │   ├── api.js                  # Base API client
│   │   └── inventarioApi.js        # Inventory API
│   ├── App.jsx              # Root component
│   ├── App.css              # Global styles
│   ├── custom-bootstrap.scss # Bootstrap customization
│   ├── index.css            # Global CSS
│   └── main.jsx             # Entry point
├── package.json
├── vite.config.js
└── index.html
```

---

## Development Workflow

### Running the Development Server

```bash
npm run dev
```

**Features**:
- Hot Module Replacement (HMR) enabled
- Automatic reload on code changes
- Debug mode enabled
- localhost:5173 auto-open

### Directory Structure Tips

1. **Pages** go in `src/pages/` (route-level components)
2. **Reusable components** go in `src/components/`
3. **Shared logic** goes in `src/hooks/`
4. **API calls** go in `src/services/`
5. **Global state** in `src/context/`

### Linting & Code Style

```bash
# Run ESLint
npm run lint

# Fix linting issues
npm run lint -- --fix
```

**ESLint Config**: [eslint.config.js](eslint.config.js)

---

## Architecture Overview

### Component Hierarchy

```
App.jsx
├── AuthContext.Provider
├── BrowserRouter
│   ├── Routes
│   │   ├── ProtectedRoute
│   │   │   ├── Admin Dashboard
│   │   │   ├── User Pages
│   │   │   └── ...
│   │   └── Public Routes
│   │       ├── Login
│   │       └── Public Products
│   └── Navigation.jsx
```

### Data Flow

```
1. Page Component
   ↓
2. API Service (api.js or inventarioApi.js)
   ↓
3. Backend (http://localhost:8080/api)
   ↓
4. Response → State/Context
   ↓
5. Re-render Component
```

### Authentication Flow

```
1. User enters email/password
2. Login.jsx sends to /api/auth/login
3. Backend returns JWT token
4. Store token in localStorage (AuthContext)
5. Add token to all API requests
6. ProtectedRoute checks auth status
7. Redirect if unauthorized
```

---

## Key Components

### Navigation.jsx
**Consolidated navigation component** for both admin and user views.

**Props**:
- None (uses AuthContext for user info)

**Features**:
- User role-based menu items
- Logout functionality
- Navigation to all sections

**Usage**:
```jsx
import Navigation from './components/Navigation'

// In App.jsx
<Navigation />
```

---

### ProtectedRoute.jsx
**Wrapper for protecting routes** that require authentication.

**Props**:
```jsx
{
  children: ReactNode,
  requiredRole?: string  // Optional: restrict to role
}
```

**Features**:
- Redirects to login if not authenticated
- Checks user roles
- Preserves intended route

**Usage**:
```jsx
<ProtectedRoute requiredRole="admin">
  <Admin />
</ProtectedRoute>
```

---

### ProductCard.jsx
**Reusable card for displaying products**.

**Props**:
```jsx
{
  product: {
    id: number,
    name: string,
    price: number,
    description: string,
    stock: number
  },
  onAddToCart?: function,
  showActions?: boolean
}
```

**Usage**:
```jsx
<ProductCard 
  product={product} 
  onAddToCart={handleAdd}
  showActions={true}
/>
```

---

## Custom Hooks

### useCart
**Shopping cart state management hook**.

**Returns**:
```jsx
{
  cart: Array<CartItem>,
  addToCart: (product, quantity) => void,
  removeFromCart: (productId) => void,
  clearCart: () => void,
  getTotalPrice: () => number,
  getItemCount: () => number
}
```

**Usage**:
```jsx
import useCart from '../hooks/useCart'

function ProductDetail() {
  const { addToCart, cart } = useCart()
  
  return (
    <button onClick={() => addToCart(product, 1)}>
      Add to Cart
    </button>
  )
}
```

---

## Context API Usage

### AuthContext
**Global authentication state**.

**Provides**:
```jsx
{
  user: {
    id: number,
    name: string,
    email: string,
    role: string,
    token: string
  },
  login: (email, password) => Promise,
  logout: () => void,
  isAuthenticated: boolean,
  isLoading: boolean
}
```

**Usage**:
```jsx
import { useContext } from 'react'
import { AuthContext } from '../context/AuthContext'

function MyComponent() {
  const { user, logout, isAuthenticated } = useContext(AuthContext)
  
  if (!isAuthenticated) {
    return <Navigate to="/login" />
  }
  
  return <div>Welcome, {user.name}!</div>
}
```

**Location**: [src/context/AuthContext.jsx](src/context/AuthContext.jsx)

---

## Styling

### Bootstrap Integration

**Custom Bootstrap Setup**:
- File: [src/custom-bootstrap.scss](src/custom-bootstrap.scss)
- Imports: Bootstrap 5 from node_modules
- Custom variables override available

**Usage**:
```jsx
// In components
<div className="container">
  <div className="row">
    <div className="col-md-6">
      <button className="btn btn-primary">Click me</button>
    </div>
  </div>
</div>
```

### Global CSS

**Files**:
- [src/index.css](src/index.css) - Global styles
- [src/App.css](src/App.css) - App-specific styles

### Component CSS

Use `import "./Component.css"` for component-scoped styles:

```jsx
// Component.jsx
import './Component.css'

export default function Component() {
  return <div className="component">...</div>
}
```

---

## API Integration

### Base API Client

**File**: [src/services/api.js](src/services/api.js)

**Features**:
- Axios-based HTTP client
- Automatic JWT token attachment
- Error handling
- Base URL configuration

**Methods**:
```jsx
api.get(endpoint, config)
api.post(endpoint, data, config)
api.put(endpoint, data, config)
api.delete(endpoint, config)
api.patch(endpoint, data, config)
```

**Usage**:
```jsx
import api from './services/api'

// GET request
const products = await api.get('/productos')

// POST request
const newProduct = await api.post('/productos', {
  name: 'Product Name',
  price: 99.99
})

// With auth token (automatic)
const orders = await api.get('/ordenes')  // Token auto-attached
```

### Error Handling

```jsx
try {
  const response = await api.get('/productos')
  console.log(response.data)
} catch (error) {
  console.error('API Error:', error.response.data)
  // error.response.data = { code, message, details }
}
```

### Inventory API Client

**File**: [src/services/inventarioApi.js](src/services/inventarioApi.js)

Specialized API client for inventory operations.

---

## Error Handling

### Standard Error Response

All API errors return this format:

```json
{
  "code": 400,
  "message": "Error description",
  "details": "Additional context"
}
```

### Frontend Error Handling Pattern

```jsx
import { useState } from 'react'
import api from './services/api'

function MyComponent() {
  const [error, setError] = useState(null)
  const [loading, setLoading] = useState(false)

  const fetchData = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await api.get('/endpoint')
      // Handle success
    } catch (err) {
      setError(err.response?.data?.message || 'Unknown error')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      {error && <div className="alert alert-danger">{error}</div>}
      {loading && <div>Loading...</div>}
      <button onClick={fetchData}>Load Data</button>
    </div>
  )
}
```

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| 401 Unauthorized | Invalid token | Re-login user |
| 403 Forbidden | Insufficient permissions | Show error message |
| 404 Not Found | Resource deleted | Redirect to list |
| 422 Unprocessable | Validation failed | Show field errors |
| 500 Server Error | Backend error | Retry or show message |

---

## Page Overview

### Login.jsx
- Email/password form
- JWT token storage
- Redirect to dashboard on success
- Remember "next" URL for redirect

### Dashboard.jsx
- User-specific orders and info
- Cart view and checkout
- Profile management

### Admin.jsx
- Product management
- User management
- Order management
- Audit logs
- Role-based access

### Productos.jsx
- Product listing with filters
- Search functionality
- Category filtering
- Pagination

### Cart.jsx
- Review cart items
- Adjust quantities
- Proceed to checkout
- PayPal integration (if configured)

### Cotizaciones.jsx
- Create quotations
- View quotation history
- Convert to sales

### Proveedores.jsx
- Supplier management
- Debt tracking
- Contact information

---

## Development Tips

### Debugging

1. **React DevTools Browser Extension**
   - Inspect component hierarchy
   - Check props and state changes

2. **Network Tab**
   - Monitor API calls
   - Check request/response headers
   - Verify JWT token

3. **Console Logging**
   ```jsx
   console.log('Component rendered', { props, state })
   ```

### Common Issues

**Issue**: "Cannot find 'api' module"  
**Solution**: Check import path - should be relative from file location

**Issue**: "401 Unauthorized on all requests"  
**Solution**: Check if token exists in localStorage - verify login success

**Issue**: "CORS error"  
**Solution**: 
- Backend CORS_ORIGINS must include `http://localhost:5173`
- Check docker-compose.yml environment variables

**Issue**: "Bootstrap styles not loading"  
**Solution**: 
- Verify `custom-bootstrap.scss` is imported
- Check import path is correct (`./node_modules` in container)

### Performance Optimization

1. **Code Splitting**
   ```jsx
   const Admin = React.lazy(() => import('./pages/Admin'))
   <Suspense fallback={<div>Loading...</div>}>
     <Admin />
   </Suspense>
   ```

2. **Memoization**
   ```jsx
   const ProductCard = React.memo(({ product }) => ...)
   ```

3. **useCallback for handlers**
   ```jsx
   const handleClick = useCallback(() => {...}, [dep1, dep2])
   ```

---

## Production Deployment

### Build Optimization

```bash
# Generate optimized production build
npm run build

# Output: dist/ directory with minified files
```

### Environment Variables for Production

```env
VITE_API_URL=https://api.legionstore.com
```

### Deployment Checklist

- [ ] Remove console.log statements
- [ ] Test all features in production build
- [ ] Verify API endpoints (HTTPS)
- [ ] Check environment variables
- [ ] Test PayPal integration (production keys)
- [ ] Set up CORS correctly on backend
- [ ] Monitor analytics and errors

---

## Future Enhancements

- [ ] State management with Redux or Zustand
- [ ] Component library (Storybook)
- [ ] E2E testing (Cypress/Playwright)
- [ ] Unit testing (Vitest)
- [ ] Dark mode support
- [ ] Internationalization (i18n)
- [ ] Accessibility (a11y) improvements
- [ ] Service Worker for offline support
- [ ] Real-time notifications (WebSocket)

---

**Last Updated**: March 20, 2026  
**Version**: 1.0  
**Status**: Production Ready
