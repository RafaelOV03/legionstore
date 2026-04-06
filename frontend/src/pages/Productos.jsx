import { useState, useEffect } from 'react';
import { Form, Row, Col } from 'react-bootstrap';
import { useAuth } from '../context/AuthContext';
import { CRUDTable } from '../components/CRUDTable';
import * as api from '../services/resourceApi';

function Productos() {
  const { hasPermission } = useAuth();

  const columns = [
    { key: 'codigo', label: 'Código' },
    { key: 'name', label: 'Nombre' },
    { key: 'category', label: 'Categoría' },
    { key: 'stock_total', label: 'Stock' },
    { key: 'precio_compra', label: 'Costo' },
    { key: 'precio_venta', label: 'Precio' },
    { key: 'margen', label: 'Margen' }
  ];

  const itemShape = {
    codigo: '',
    name: '',
    description: '',
    precio_venta: '',
    precio_compra: '',
    category: '',
    brand: '',
    image_url: ''
  };

  const renderForm = (formData, setFormData) => (
    <>
      <Row>
        <Col md={4}>
          <Form.Group className="mb-3">
            <Form.Label>Código *</Form.Label>
            <Form.Control
              required
              value={formData.codigo}
              onChange={(e) => setFormData({...formData, codigo: e.target.value})}
              placeholder="SKU-001"
            />
          </Form.Group>
        </Col>
        <Col md={8}>
          <Form.Group className="mb-3">
            <Form.Label>Nombre *</Form.Label>
            <Form.Control
              required
              value={formData.name}
              onChange={(e) => setFormData({...formData, name: e.target.value})}
              placeholder="Nombre del producto"
            />
          </Form.Group>
        </Col>
      </Row>
      <Row>
        <Col md={6}>
          <Form.Group className="mb-3">
            <Form.Label>Categoría *</Form.Label>
            <Form.Control
              required
              value={formData.category}
              onChange={(e) => setFormData({...formData, category: e.target.value})}
              placeholder="Categoría"
            />
          </Form.Group>
        </Col>
        <Col md={6}>
          <Form.Group className="mb-3">
            <Form.Label>Marca</Form.Label>
            <Form.Control
              value={formData.brand}
              onChange={(e) => setFormData({...formData, brand: e.target.value})}
              placeholder="Marca"
            />
          </Form.Group>
        </Col>
      </Row>
      <Row>
        <Col md={6}>
          <Form.Group className="mb-3">
            <Form.Label>Precio Compra (Bs.) *</Form.Label>
            <Form.Control
              required
              type="number"
              step="0.01"
              value={formData.precio_compra}
              onChange={(e) => setFormData({...formData, precio_compra: e.target.value})}
              placeholder="0.00"
            />
          </Form.Group>
        </Col>
        <Col md={6}>
          <Form.Group className="mb-3">
            <Form.Label>Precio Venta (Bs.) *</Form.Label>
            <Form.Control
              required
              type="number"
              step="0.01"
              value={formData.precio_venta}
              onChange={(e) => setFormData({...formData, precio_venta: e.target.value})}
              placeholder="0.00"
            />
          </Form.Group>
        </Col>
      </Row>
      <Form.Group className="mb-3">
        <Form.Label>Descripción</Form.Label>
        <Form.Control
          as="textarea"
          rows={3}
          value={formData.description}
          onChange={(e) => setFormData({...formData, description: e.target.value})}
          placeholder="Descripción del producto"
        />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>URL Imagen</Form.Label>
        <Form.Control
          value={formData.image_url}
          onChange={(e) => setFormData({...formData, image_url: e.target.value})}
          placeholder="https://..."
        />
      </Form.Group>
    </>
  );

  const renderCustomCell = (key, value, item) => {
    if (key === 'margen') {
      const margen = item.precio_venta && item.precio_compra 
        ? (((item.precio_venta - item.precio_compra) / item.precio_compra) * 100).toFixed(1) 
        : 0;
      return `${margen}%`;
    }
    if (key === 'precio_compra' || key === 'precio_venta') {
      return `Bs. ${parseFloat(value || 0).toFixed(2)}`;
    }
    return value;
  };

  return (
    <CRUDTable
      title="Productos"
      columns={columns}
      onLoad={api.getProducts}
      onAdd={(data) => api.createProduct(transformProductData(data))}
      onUpdate={(id, data) => api.updateProduct(id, transformProductData(data))}
      onDelete={api.deleteProduct}
      itemShape={itemShape}
      renderForm={renderForm}
      renderCustomCell={renderCustomCell}
      canEdit={[
        hasPermission('products.create'),
        hasPermission('products.update'),
        hasPermission('products.delete')
      ]}
    />
  );
}

function transformProductData(data) {
  return {
    ...data,
    precio_venta: parseFloat(data.precio_venta) || 0,
    precio_compra: parseFloat(data.precio_compra) || 0
  };
}

export default Productos;
