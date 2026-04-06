import { useState, useEffect } from 'react';
import { Container, Table, Button, Modal, Form, Alert, Badge, Card, Row, Col, ListGroup } from 'react-bootstrap';
import { PencilSquare, Trash, PlusCircle, ShieldCheck, CheckCircleFill } from 'react-bootstrap-icons';
import { getRoles, getPermissions, createRole, updateRole, deleteRole } from '../services/resourceApi';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

function Roles() {
  const [roles, setRoles] = useState([]);
  const [permissions, setPermissions] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [showPermissionsModal, setShowPermissionsModal] = useState(false);
  const [selectedRole, setSelectedRole] = useState(null);
  const [editingRole, setEditingRole] = useState(null);
  const [alert, setAlert] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    permission_ids: []
  });
  
  const { hasPermission } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!hasPermission('roles.read')) {
      navigate('/');
      return;
    }
    loadRoles();
    loadPermissions();
  }, []);

  const loadRoles = async () => {
    try {
      const data = await getRoles();
      setRoles(data);
    } catch (error) {
      showAlert('Error al cargar roles: ' + error.message, 'danger');
    }
  };

  const loadPermissions = async () => {
    try {
      const data = await getPermissions();
      setPermissions(data);
    } catch (error) {
      showAlert('Error al cargar permisos: ' + error.message, 'danger');
    }
  };

  const showAlert = (message, variant) => {
    setAlert({ message, variant });
    setTimeout(() => setAlert(null), 3000);
  };

  const handleOpenModal = (role = null) => {
    if (role) {
      setEditingRole(role);
      setFormData({
        name: role.name,
        description: role.description,
        permission_ids: role.permissions?.map(p => p.id) || []
      });
    } else {
      setEditingRole(null);
      setFormData({
        name: '',
        description: '',
        permission_ids: []
      });
    }
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingRole(null);
    setFormData({
      name: '',
      description: '',
      permission_ids: []
    });
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handlePermissionToggle = (permissionId) => {
    setFormData(prev => {
      const ids = [...prev.permission_ids];
      const index = ids.indexOf(permissionId);
      if (index > -1) {
        ids.splice(index, 1);
      } else {
        ids.push(permissionId);
      }
      return { ...prev, permission_ids: ids };
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      if (editingRole) {
        await updateRole(editingRole.id, formData);
        showAlert('Rol actualizado exitosamente', 'success');
      } else {
        await createRole(formData);
        showAlert('Rol creado exitosamente', 'success');
      }

      handleCloseModal();
      loadRoles();
    } catch (error) {
      showAlert('Error: ' + error.message, 'danger');
    }
  };

  const handleDelete = async (roleId) => {
    if (window.confirm('¿Está seguro de eliminar este rol?')) {
      try {
        await deleteRole(roleId);
        showAlert('Rol eliminado exitosamente', 'success');
        loadRoles();
      } catch (error) {
        showAlert('Error al eliminar rol: ' + error.message, 'danger');
      }
    }
  };

  const handleShowPermissions = (role) => {
    setSelectedRole(role);
    setShowPermissionsModal(true);
  };

  const groupPermissionsByResource = (perms) => {
    const grouped = {};
    perms.forEach(perm => {
      if (!grouped[perm.resource]) {
        grouped[perm.resource] = [];
      }
      grouped[perm.resource].push(perm);
    });
    return grouped;
  };

  const getActionBadge = (action) => {
    const variants = {
      'create': 'success',
      'read': 'info',
      'update': 'warning',
      'delete': 'danger'
    };
    return variants[action] || 'secondary';
  };

  const getName = (resource) => {
    const icons = {
      'products': 'Productos',
      'orders': 'Ventas',
      'users': 'Usuarios',
      'roles': 'Roles'
    };
    return icons[resource];
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
          <h2><ShieldCheck className="me-2" />Gestión de Roles y Permisos</h2>
        </Col>
        <Col className="text-end">
          {hasPermission('roles.create') && (
            <Button variant="primary" onClick={() => handleOpenModal()}>
              <PlusCircle className="me-2" />Nuevo Rol
            </Button>
          )}
        </Col>
      </Row>

      <Row>
        {roles.map(role => (
          <Col md={6} lg={4} key={role.id} className="mb-4">
            <Card>
              <Card.Header className="d-flex justify-content-between align-items-center">
                <div>
                  <h5 className="mb-0">{role.name}</h5>
                  {role.is_system && <Badge bg="info" className="mt-1">Sistema</Badge>}
                </div>
                <div>
                  {!role.is_system && hasPermission('roles.update') && (
                    <Button 
                      variant="warning" 
                      size="sm" 
                      className="me-2"
                      onClick={() => handleOpenModal(role)}
                    >
                      <PencilSquare />
                    </Button>
                  )}
                  {!role.is_system && hasPermission('roles.delete') && (
                    <Button 
                      variant="danger" 
                      size="sm"
                      onClick={() => handleDelete(role.id)}
                    >
                      <Trash />
                    </Button>
                  )}
                </div>
              </Card.Header>
              <Card.Body>
                <Card.Text className="text-muted small mb-3">
                  {role.description || 'Sin descripción'}
                </Card.Text>
                <div className="d-flex justify-content-between align-items-center">
                  <span className="text-muted">
                    {role.permissions?.length || 0} permisos
                  </span>
                  <Button 
                    variant="outline-primary" 
                    size="sm"
                    onClick={() => handleShowPermissions(role)}
                  >
                    Ver Permisos
                  </Button>
                </div>
              </Card.Body>
            </Card>
          </Col>
        ))}
      </Row>

      {/* Modal para crear/editar rol */}
      <Modal show={showModal} onHide={handleCloseModal} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>{editingRole ? 'Editar Rol' : 'Nuevo Rol'}</Modal.Title>
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
              <Form.Label>Descripción</Form.Label>
              <Form.Control
                as="textarea"
                rows={2}
                name="description"
                value={formData.description}
                onChange={handleChange}
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Permisos</Form.Label>
              <div style={{ maxHeight: '400px', overflowY: 'auto', border: '1px solid #dee2e6', borderRadius: '4px', padding: '10px' }}>
                {Object.entries(groupPermissionsByResource(permissions)).map(([resource, perms]) => (
                  <div key={resource} className="mb-3">
                    <h6 className="mb-2">{getName(resource)}</h6>
                    {perms.map(perm => (
                      <Form.Check
                        key={perm.id}
                        type="checkbox"
                        id={`perm-${perm.id}`}
                        label={
                          <span>
                            <Badge bg={getActionBadge(perm.action)} className="me-2">
                              {perm.action}
                            </Badge>
                            {perm.description}
                          </span>
                        }
                        checked={formData.permission_ids.includes(perm.id)}
                        onChange={() => handlePermissionToggle(perm.id)}
                        className="mb-2"
                      />
                    ))}
                  </div>
                ))}
              </div>
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={handleCloseModal}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              {editingRole ? 'Actualizar' : 'Crear'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>

      {/* Modal para ver permisos de un rol */}
      <Modal show={showPermissionsModal} onHide={() => setShowPermissionsModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Permisos de {selectedRole?.name}</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {selectedRole && (
            <>
              <p className="text-muted mb-3">{selectedRole.description}</p>
              {Object.entries(groupPermissionsByResource(selectedRole.permissions || [])).map(([resource, perms]) => (
                <div key={resource} className="mb-3">
                  <h6>{getName(resource)}</h6>
                  <ListGroup variant="flush">
                    {perms.map(perm => (
                      <ListGroup.Item key={perm.id} className="d-flex align-items-center">
                        <CheckCircleFill className="text-success me-2" />
                        <Badge bg={getActionBadge(perm.action)} className="me-2">
                          {perm.action}
                        </Badge>
                        <span>{perm.description}</span>
                      </ListGroup.Item>
                    ))}
                  </ListGroup>
                </div>
              ))}
            </>
          )}
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => setShowPermissionsModal(false)}>
            Cerrar
          </Button>
        </Modal.Footer>
      </Modal>
    </Container>
  );
}

export default Roles;
