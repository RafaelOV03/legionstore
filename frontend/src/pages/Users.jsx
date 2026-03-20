import { useState, useEffect } from 'react';
import { Form, Badge } from 'react-bootstrap';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { CRUDTable } from '../components/CRUDTable';
import { getUsers, getRoles, createUser, updateUser, deleteUser } from '../services/resourceApi';

function Users() {
  const { hasPermission } = useAuth();
  const navigate = useNavigate();
  const [roles, setRoles] = useState([]);

  useEffect(() => {
    if (!hasPermission('users.read')) {
      navigate('/');
      return;
    }
    loadRoles();
  }, []);

  const loadRoles = async () => {
    try {
      const data = await getRoles();
      setRoles(data);
    } catch (error) {
      console.error(error);
    }
  };

  const columns = [
    { key: 'id', label: 'ID' },
    { key: 'name', label: 'Nombre' },
    { key: 'email', label: 'Email' },
    { key: 'role', label: 'Rol' },
    { key: 'created_at', label: 'Fecha Registro' }
  ];

  const itemShape = {
    name: '',
    email: '',
    password: '',
    role_id: ''
  };

  const renderForm = (formData, setFormData) => (
    <>
      <Form.Group className="mb-3">
        <Form.Label>Nombre *</Form.Label>
        <Form.Control
          type="text"
          value={formData.name}
          onChange={(e) => setFormData({...formData, name: e.target.value})}
          placeholder="Nombre completo"
          required
        />
      </Form.Group>

      <Form.Group className="mb-3">
        <Form.Label>Email *</Form.Label>
        <Form.Control
          type="email"
          value={formData.email}
          onChange={(e) => setFormData({...formData, email: e.target.value})}
          placeholder="usuario@ejemplo.com"
          required
        />
      </Form.Group>

      <Form.Group className="mb-3">
        <Form.Label>Contraseña {formData.id && '(dejar en blanco para no cambiar)'}</Form.Label>
        <Form.Control
          type="password"
          value={formData.password || ''}
          onChange={(e) => setFormData({...formData, password: e.target.value})}
          required={!formData.id}
          minLength={6}
        />
      </Form.Group>

      <Form.Group className="mb-3">
        <Form.Label>Rol *</Form.Label>
        <Form.Select
          value={formData.role_id || ''}
          onChange={(e) => setFormData({...formData, role_id: parseInt(e.target.value)})}
          required
        >
          <option value="">Seleccione un rol</option>
          {roles.map((role) => (
            <option key={role.id} value={role.id}>
              {role.name}
            </option>
          ))}
        </Form.Select>
      </Form.Group>
    </>
  );

  const renderCustomCell = (key, value, item) => {
    if (key === 'role') {
      const role = item.role || {};
      const variants = {
        'administrador': 'danger',
        'empleado': 'primary',
        'usuario': 'secondary'
      };
      return <Badge bg={variants[role.name] || 'secondary'}>{role.name || 'Sin rol'}</Badge>;
    }
    if (key === 'created_at') {
      return value ? new Date(value).toLocaleDateString() : '-';
    }
    return value;
  };

  const handleAddUser = async (data) => {
    const submitData = { ...data };
    if (!submitData.password) {
      throw new Error('La contraseña es requerida para nuevos usuarios');
    }
    return createUser(submitData);
  };

  const handleUpdateUser = async (id, data) => {
    const submitData = { ...data };
    if (!submitData.password) {
      delete submitData.password;
    }
    return updateUser(id, submitData);
  };

  return (
    <CRUDTable
      title="Usuarios"
      columns={columns}
      onLoad={getUsers}
      onAdd={handleAddUser}
      onUpdate={handleUpdateUser}
      onDelete={deleteUser}
      itemShape={itemShape}
      renderForm={renderForm}
      renderCustomCell={renderCustomCell}
      canEdit={[
        hasPermission('users.create'),
        hasPermission('users.update'),
        hasPermission('users.delete')
      ]}
    />
  );
}

export default Users;

