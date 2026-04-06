import { useState, useEffect } from 'react';
import { Form, Badge } from 'react-bootstrap';
import { useAuth } from '../context/AuthContext';
import { CRUDTable } from '../components/CRUDTable';
import * as api from '../services/resourceApi';

function Sedes() {
  const { hasPermission } = useAuth();

  const columns = [
    { key: 'nombre', label: 'Nombre' },
    { key: 'direccion', label: 'Dirección' },
    { key: 'telefono', label: 'Teléfono' },
    { key: 'activa', label: 'Estado' }
  ];

  const itemShape = {
    nombre: '',
    direccion: '',
    telefono: '',
    activa: true
  };

  const renderForm = (formData, setFormData) => (
    <>
      <Form.Group className="mb-3">
        <Form.Label>Nombre *</Form.Label>
        <Form.Control
          required
          value={formData.nombre}
          onChange={(e) => setFormData({...formData, nombre: e.target.value})}
          placeholder="Nombre de la sede"
        />
      </Form.Group>

      <Form.Group className="mb-3">
        <Form.Label>Dirección</Form.Label>
        <Form.Control
          value={formData.direccion}
          onChange={(e) => setFormData({...formData, direccion: e.target.value})}
          placeholder="Dirección completa"
        />
      </Form.Group>

      <Form.Group className="mb-3">
        <Form.Label>Teléfono</Form.Label>
        <Form.Control
          value={formData.telefono}
          onChange={(e) => setFormData({...formData, telefono: e.target.value})}
          placeholder="Número telefónico"
        />
      </Form.Group>

      <Form.Group className="mb-3">
        <Form.Check
          type="switch"
          label="Sede activa"
          checked={formData.activa}
          onChange={(e) => setFormData({...formData, activa: e.target.checked})}
        />
      </Form.Group>
    </>
  );

  const renderCustomCell = (key, value, item) => {
    if (key === 'activa') {
      return value ? <Badge bg="success">Activa</Badge> : <Badge bg="secondary">Inactiva</Badge>;
    }
    return value;
  };

  return (
    <CRUDTable
      title="Sedes"
      columns={columns}
      onLoad={api.getSedes}
      onAdd={api.createSede}
      onUpdate={api.updateSede}
      onDelete={api.deleteSede}
      itemShape={itemShape}
      renderForm={renderForm}
      renderCustomCell={renderCustomCell}
      canEdit={[
        hasPermission('sedes.create'),
        hasPermission('sedes.update'),
        hasPermission('sedes.delete')
      ]}
    />
  );
}

export default Sedes;
