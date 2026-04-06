import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Button, Modal, Form, Alert, Spinner, Badge } from 'react-bootstrap';
import { Plus, Pencil, Trash } from 'react-bootstrap-icons';

/**
 * CRUDTable - Componente genérico reutilizable para operaciones CRUD
 * 
 * @param {Object} props
 * @param {string} props.title - Título de la tabla (ej: "Productos")
 * @param {Array} props.columns - Definición de columnas: [{ key: 'id', label: 'ID' }, ...]
 * @param {Function} props.onLoad - Función async que carga los datos: () => Promise<array>
 * @param {Function} props.onAdd - Función async que crea: (formData) => Promise
 * @param {Function} props.onUpdate - Función async que actualiza: (id, formData) => Promise
 * @param {Function} props.onDelete - Función async que elimina: (id) => Promise
 * @param {Object} props.itemShape - Forma del objeto: { campo: valor, ... }
 * @param {Function} props.renderForm - Función que renderea el formulario: (formData, setFormData) => JSX
 * @param {Array} props.canEdit - Permisos: [canCreate, canEdit, canDelete]
 * @param {Function} props.renderCustomCell - (opcional) Renderizar celdas personalizadas: (key, value, item) => JSX
 */
export function CRUDTable({
  title,
  columns,
  onLoad,
  onAdd,
  onUpdate,
  onDelete,
  itemShape,
  renderForm,
  canEdit = [true, true, true], // [canCreate, canEdit, canDelete]
  renderCustomCell = null
}) {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState(null);
  const [formData, setFormData] = useState({ ...itemShape });
  const [alert, setAlert] = useState(null);

  const [canCreate, canEditItem, canDelete] = canEdit;

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const data = await onLoad();
      setItems(data || []);
      setAlert(null);
    } catch (err) {
      showAlert(`Error cargando ${title}`, 'danger');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const showAlert = (message, variant = 'success') => {
    setAlert({ message, variant });
    setTimeout(() => setAlert(null), 3000);
  };

  const handleOpenModal = (item = null) => {
    if (item) {
      setEditing(item);
      setFormData({ ...item });
    } else {
      setEditing(null);
      setFormData({ ...itemShape });
    }
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditing(null);
    setFormData({ ...itemShape });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (editing) {
        await onUpdate(editing.id, formData);
        showAlert(`${title} actualizado correctamente`);
      } else {
        await onAdd(formData);
        showAlert(`${title} creado correctamente`);
      }
      handleCloseModal();
      loadData();
    } catch (err) {
      showAlert(`Error guardando ${title}`, 'danger');
      console.error(err);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm(`¿Está seguro de eliminar este ${title}?`)) return;
    try {
      await onDelete(id);
      showAlert(`${title} eliminado correctamente`);
      loadData();
    } catch (err) {
      showAlert(`Error eliminando ${title}`, 'danger');
      console.error(err);
    }
  };

  if (loading) {
    return (
      <Container className="py-5 text-center">
        <Spinner animation="border" variant="primary" />
      </Container>
    );
  }

  return (
    <Container fluid className="py-4">
      <Row className="mb-4">
        <Col>
          <h2>{title}</h2>
        </Col>
        {canCreate && (
          <Col md="auto">
            <Button
              variant="success"
              onClick={() => handleOpenModal()}
              className="d-flex align-items-center gap-2"
            >
              <Plus size={18} /> Agregar {title}
            </Button>
          </Col>
        )}
      </Row>

      {alert && <Alert variant={alert.variant} onClose={() => setAlert(null)} dismissible>{alert.message}</Alert>}

      <Card>
        <Card.Body className="p-0">
          <Table hover responsive className="mb-0">
            <thead className="table-light">
              <tr>
                {columns.map((col) => (
                  <th key={col.key}>{col.label}</th>
                ))}
                {(canEditItem || canDelete) && <th>Acciones</th>}
              </tr>
            </thead>
            <tbody>
              {items.length === 0 ? (
                <tr>
                  <td colSpan={columns.length + (canEditItem || canDelete ? 1 : 0)} className="text-center text-muted py-4">
                    No hay {title.toLowerCase()} disponibles
                  </td>
                </tr>
              ) : (
                items.map((item) => (
                  <tr key={item.id}>
                    {columns.map((col) => (
                      <td key={`${item.id}-${col.key}`}>
                        {renderCustomCell ? (
                          renderCustomCell(col.key, item[col.key], item)
                        ) : (
                          renderCellValue(item[col.key])
                        )}
                      </td>
                    ))}
                    {(canEditItem || canDelete) && (
                      <td>
                        <div className="d-flex gap-2">
                          {canEditItem && (
                            <Button
                              size="sm"
                              variant="primary"
                              onClick={() => handleOpenModal(item)}
                              title="Editar"
                            >
                              <Pencil size={16} />
                            </Button>
                          )}
                          {canDelete && (
                            <Button
                              size="sm"
                              variant="danger"
                              onClick={() => handleDelete(item.id)}
                              title="Eliminar"
                            >
                              <Trash size={16} />
                            </Button>
                          )}
                        </div>
                      </td>
                    )}
                  </tr>
                ))
              )}
            </tbody>
          </Table>
        </Card.Body>
      </Card>

      <Modal show={showModal} onHide={handleCloseModal} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>
            {editing ? `Editar ${title}` : `Nuevo ${title}`}
          </Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleSubmit}>
          <Modal.Body>
            {renderForm(formData, setFormData)}
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={handleCloseModal}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              {editing ? 'Actualizar' : 'Crear'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

/**
 * Renderiza el valor de una celda según su tipo
 */
function renderCellValue(value) {
  if (value === null || value === undefined) return '—';
  if (typeof value === 'boolean') return value ? <Badge bg="success">Sí</Badge> : <Badge bg="secondary">No</Badge>;
  if (typeof value === 'number') return value.toFixed(2);
  if (typeof value === 'object') return JSON.stringify(value);
  return String(value).substring(0, 50);
}
