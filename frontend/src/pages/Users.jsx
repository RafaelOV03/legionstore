import { useState, useEffect } from 'react';
import { Container, Table, Button, Modal, Form, Alert, Badge, Card, Row, Col } from 'react-bootstrap';
import { PencilSquare, Trash, PlusCircle, PersonFill } from 'react-bootstrap-icons';
import { getUsers, getRoles, createUser, updateUser, deleteUser } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

function Users() {
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [alert, setAlert] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    role_id: ''
  });
  
  const { hasPermission } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!hasPermission('users.read')) {
      navigate('/');
      return;
    }
    loadUsers();
    loadRoles();
  }, []);

  const loadUsers = async () => {
    try {
      const data = await getUsers();
      setUsers(data);
    } catch (error) {
      showAlert('Error al cargar usuarios: ' + error.message, 'danger');
    }
  };

  const loadRoles = async () => {
    try {
      const data = await getRoles();
      setRoles(data);
    } catch (error) {
      showAlert('Error al cargar roles: ' + error.message, 'danger');
    }
  };

  const showAlert = (message, variant) => {
    setAlert({ message, variant });
    setTimeout(() => setAlert(null), 3000);
  };

  const handleOpenModal = (user = null) => {
    if (user) {
      setEditingUser(user);
      setFormData({
        name: user.name,
        email: user.email,
        password: '',
        role_id: user.role.id
      });
    } else {
      setEditingUser(null);
      setFormData({
        name: '',
        email: '',
        password: '',
        role_id: ''
      });
    }
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingUser(null);
    setFormData({
      name: '',
      email: '',
      password: '',
      role_id: ''
    });
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'role_id' ? parseInt(value) : value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const submitData = { ...formData };
      
      // Si estamos editando y no se cambió la contraseña, no enviarla
      if (editingUser && !submitData.password) {
        delete submitData.password;
      }

      if (editingUser) {
        await updateUser(editingUser.id, submitData);
        showAlert('Usuario actualizado exitosamente', 'success');
      } else {
        if (!submitData.password) {
          showAlert('La contraseña es requerida para nuevos usuarios', 'danger');
          return;
        }
        await createUser(submitData);
        showAlert('Usuario creado exitosamente', 'success');
      }

      handleCloseModal();
      loadUsers();
    } catch (error) {
      showAlert('Error: ' + error.message, 'danger');
    }
  };

  const handleDelete = async (userId) => {
    if (window.confirm('¿Está seguro de eliminar este usuario?')) {
      try {
        await deleteUser(userId);
        showAlert('Usuario eliminado exitosamente', 'success');
        loadUsers();
      } catch (error) {
        showAlert('Error al eliminar usuario: ' + error.message, 'danger');
      }
    }
  };

  const getRoleBadge = (roleName) => {
    const variants = {
      'administrador': 'danger',
      'empleado': 'primary',
      'usuario': 'secondary'
    };
    return variants[roleName] || 'secondary';
  };

  return (
    <Container className="py-4">
      {alert && (
        <Alert 
          variant={alert.variant} 
          dismissible 
          onClose={() => setAlert(null)}
          style={{
            position: 'fixed',
            top: '20px',
            right: '20px',
            zIndex: 9999,
            minWidth: '300px',
            boxShadow: '0 4px 12px rgba(0,0,0,0.15)'
          }}
        >
          {alert.message}
        </Alert>
      )}

      <Row className="mb-4">
        <Col>
          <h2><PersonFill className="me-2" />Gestión de Usuarios</h2>
        </Col>
        <Col className="text-end">
          {hasPermission('users.create') && (
            <Button variant="primary" onClick={() => handleOpenModal()}>
              <PlusCircle className="me-2" />Nuevo Usuario
            </Button>
          )}
        </Col>
      </Row>

      <Card>
        <Card.Body>
          <Table striped bordered hover responsive>
            <thead>
              <tr>
                <th>id</th>
                <th>Nombre</th>
                <th>Email</th>
                <th>Rol</th>
                <th>Fecha Registro</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {users && users.length > 0 ? users.map((user, index) => (
                <tr key={user.id || `user-${index}`}>
                  <td>{user.id}</td>
                  <td>{user.name}</td>
                  <td>{user.email}</td>
                  <td>
                    <Badge bg={getRoleBadge(user.role?.name)}>
                      {user.role?.name || 'Sin rol'}
                    </Badge>
                  </td>
                  <td>{user.created_at ? new Date(user.created_at).toLocaleDateString() : '-'}</td>
                  <td>
                    {hasPermission('users.update') && (
                      <Button 
                        variant="warning" 
                        size="sm" 
                        className="me-2"
                        onClick={() => handleOpenModal(user)}
                      >
                        <PencilSquare />
                      </Button>
                    )}
                    {hasPermission('users.delete') && (
                      <Button 
                        variant="danger" 
                        size="sm"
                        onClick={() => handleDelete(user.id)}
                      >
                        <Trash />
                      </Button>
                    )}
                  </td>
                </tr>
              )) : (
                <tr>
                  <td colSpan="6" className="text-center text-muted">No hay usuarios</td>
                </tr>
              )}
            </tbody>
          </Table>
        </Card.Body>
      </Card>

      {/* Modal para crear/editar usuario */}
      <Modal show={showModal} onHide={handleCloseModal}>
        <Modal.Header closeButton>
          <Modal.Title>{editingUser ? 'Editar Usuario' : 'Nuevo Usuario'}</Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleSubmit}>
          <Modal.Body>
            <Form.Group className="mb-3">
              <Form.Label>Nombre</Form.Label>
              <Form.Control
                type="text"
                name="name"
                value={formData.name}
                onChange={handleChange}
                required
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Email</Form.Label>
              <Form.Control
                type="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                required
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>
                Contraseña {editingUser && '(dejar en blanco para no cambiar)'}
              </Form.Label>
              <Form.Control
                type="password"
                name="password"
                value={formData.password}
                onChange={handleChange}
                required={!editingUser}
                minLength={6}
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Rol</Form.Label>
              <Form.Select
                name="role_id"
                value={formData.role_id || ''}
                onChange={handleChange}
                required
              >
                <option value="">Seleccione un rol</option>
                {roles && roles.length > 0 && roles.map((role, index) => (
                  <option key={role.id || `role-${index}`} value={role.id}>
                    {role.name || 'Sin nombre'}
                  </option>
                ))}
              </Form.Select>
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={handleCloseModal}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              {editingUser ? 'Actualizar' : 'Crear'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

export default Users;
