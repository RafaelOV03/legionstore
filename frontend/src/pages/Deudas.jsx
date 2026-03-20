import { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Table, Button, Modal, Form, Badge, Spinner, Alert, InputGroup, ProgressBar } from 'react-bootstrap';
import { Plus, CashStack, Check, Clock, ExclamationTriangle, Receipt } from 'react-bootstrap-icons';
import { useAuth } from '../context/AuthContext';
import * as api from '../services/resourceApi';

function Deudas() {
  const { hasPermission, user } = useAuth();
  const [deudas, setDeudas] = useState([]);
  const [proveedores, setProveedores] = useState([]);
  const [resumen, setResumen] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [statusFilter, setStatusFilter] = useState('');
  
  const [showModal, setShowModal] = useState(false);
  const [showPagoModal, setShowPagoModal] = useState(false);
  const [selectedDeuda, setSelectedDeuda] = useState(null);
  const [formData, setFormData] = useState({
    proveedor_id: '',
    monto_total: '',
    fecha_factura: new Date().toISOString().split('T')[0],
    fecha_vencimiento: '',
    notas: '',
    numero_factura: ''
  });
  const [pagoData, setPagoData] = useState({ monto: '', metodo_pago: 'efectivo', referencia: '' });

  const canCreate = hasPermission('deudas.create');
  const canPay = hasPermission('deudas.update');

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [deudasData, proveedoresData, resumenData] = await Promise.all([
        api.getDeudas(),
        api.getProveedores(),
        api.getResumenDeudas()
      ]);
      setDeudas(deudasData || []);
      setProveedores(proveedoresData || []);
      setResumen(resumenData);
    } catch (err) {
      setError('Error cargando datos');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        proveedor_id: parseInt(formData.proveedor_id),
        numero_factura: formData.numero_factura,
        monto_total: parseFloat(formData.monto_total),
        fecha_factura: formData.fecha_factura,
        fecha_vencimiento: formData.fecha_vencimiento || null,
        notas: formData.notas
      };

      await api.createDeuda(payload);
      setShowModal(false);
      resetForm();
      loadData();
      setSuccess('Deuda registrada exitosamente');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error registrando deuda');
      console.error(err);
    }
  };

  const handlePago = async (e) => {
    e.preventDefault();
    try {
      await api.registrarPago(selectedDeuda.id, {
        monto: parseFloat(pagoData.monto),
        metodo_pago: pagoData.metodo_pago,
        referencia: pagoData.referencia,
        usuario_id: user?.id
      });
      setShowPagoModal(false);
      setPagoData({ monto: '', metodo_pago: 'efectivo', referencia: '' });
      loadData();
      setSuccess('Pago registrado exitosamente');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('Error registrando pago');
      console.error(err);
    }
  };

  const resetForm = () => {
    setFormData({
      proveedor_id: '',
      monto_total: '',
      fecha_factura: new Date().toISOString().split('T')[0],
      fecha_vencimiento: '',
      notas: '',
      numero_factura: ''
    });
  };

  const getEstadoBadge = (estado) => {
    const colors = {
      pendiente: 'warning',
      parcial: 'info',
      pagado: 'success',
      vencido: 'danger'
    };
    return <Badge bg={colors[estado] || 'secondary'}>{estado}</Badge>;
  };

  const calcularProgreso = (deuda) => {
    if (!deuda.monto_total) return 0;
    const pagado = deuda.monto_total - (deuda.monto_pendiente || 0);
    return (pagado / deuda.monto_total) * 100;
  };

  const isVencida = (fecha) => {
    if (!fecha) return false;
    return new Date(fecha) < new Date();
  };

  const filteredDeudas = deudas.filter(d => 
    !statusFilter || d.estado === statusFilter
  );

  if (loading) {
    return (
      <Container className="py-5 text-center">
        <Spinner animation="border" variant="primary" />
      </Container>
    );
  }

  return (
    <Container fluid className="py-4">
      <Row className="mb-4 align-items-center">
        <Col>
          <h2><CashStack className="me-2" />Deudas a Proveedores</h2>
          <p className="text-muted">Control de cuentas por pagar</p>
        </Col>
        <Col xs="auto">
          {canCreate && (
            <Button variant="primary" onClick={() => { resetForm(); setShowModal(true); }}>
              <Plus className="me-1" /> Nueva Deuda
            </Button>
          )}
        </Col>
      </Row>

      {error && <Alert variant="danger" onClose={() => setError(null)} dismissible>{error}</Alert>}
      {success && <Alert variant="success" onClose={() => setSuccess(null)} dismissible>{success}</Alert>}

      {/* Resumen */}
      {resumen && (
        <Row className="mb-4 g-3">
          <Col xs={6} md={3}>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-3">
                <h4 className="mb-0 text-primary">Bs. {resumen.total_deuda?.toFixed(2) || '0.00'}</h4>
                <small className="text-muted">Total Deuda</small>
              </Card.Body>
            </Card>
          </Col>
          <Col xs={6} md={3}>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-3">
                <h4 className="mb-0 text-warning">{resumen.pendientes || 0}</h4>
                <small className="text-muted">Pendientes</small>
              </Card.Body>
            </Card>
          </Col>
          <Col xs={6} md={3}>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-3">
                <h4 className="mb-0 text-danger">{resumen.vencidas || 0}</h4>
                <small className="text-muted">Vencidas</small>
              </Card.Body>
            </Card>
          </Col>
          <Col xs={6} md={3}>
            <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
              <Card.Body className="text-center py-3">
                <h4 className="mb-0 text-success">{resumen.pagadas || 0}</h4>
                <small className="text-muted">Pagadas</small>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      )}

      {/* Filtros */}
      <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }} className="mb-4">
        <Card.Body className="py-2">
          <div className="d-flex gap-2 flex-wrap">
            <Button 
              variant={statusFilter === '' ? 'primary' : 'outline-secondary'} 
              size="sm"
              onClick={() => setStatusFilter('')}
            >
              Todas
            </Button>
            <Button 
              variant={statusFilter === 'pendiente' ? 'warning' : 'outline-warning'} 
              size="sm"
              onClick={() => setStatusFilter('pendiente')}
            >
              <Clock className="me-1" /> Pendientes
            </Button>
            <Button 
              variant={statusFilter === 'vencido' ? 'danger' : 'outline-danger'} 
              size="sm"
              onClick={() => setStatusFilter('vencido')}
            >
              <ExclamationTriangle className="me-1" /> Vencidas
            </Button>
            <Button 
              variant={statusFilter === 'pagado' ? 'success' : 'outline-success'} 
              size="sm"
              onClick={() => setStatusFilter('pagado')}
            >
              <Check className="me-1" /> Pagadas
            </Button>
          </div>
        </Card.Body>
      </Card>

      <Card style={{ background: 'var(--card-bg)', border: '1px solid var(--border-color)' }}>
        <Card.Body className="p-0">
          <Table responsive hover className="mb-0">
            <thead>
              <tr>
                <th>Proveedor</th>
                <th>Concepto</th>
                <th>N° Factura</th>
                <th className="text-end">Monto Total</th>
                <th className="text-end">Pendiente</th>
                <th style={{ width: 150 }}>Progreso</th>
                <th>Vencimiento</th>
                <th>Estado</th>
                <th className="text-center">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {filteredDeudas.length === 0 ? (
                <tr>
                  <td colSpan="9" className="text-center py-4 text-muted">
                    No hay deudas para mostrar
                  </td>
                </tr>
              ) : (
                filteredDeudas.map(deuda => {
                  const progreso = calcularProgreso(deuda);
                  const vencida = isVencida(deuda.fecha_vencimiento) && deuda.estado !== 'pagado';
                  
                  return (
                    <tr key={deuda.id} className={vencida ? 'table-danger' : ''}>
                      <td>
                        <strong>{deuda.proveedor?.nombre || 'N/A'}</strong>
                      </td>
                      <td>{deuda.notas || '-'}</td>
                      <td>{deuda.numero_factura || '-'}</td>
                      <td className="text-end">Bs. {deuda.monto_total?.toFixed(2) || '0.00'}</td>
                      <td className="text-end">
                        <strong className={deuda.monto_pendiente > 0 ? 'text-danger' : 'text-success'}>
                          Bs. {deuda.monto_pendiente?.toFixed(2) || '0.00'}
                        </strong>
                      </td>
                      <td>
                        <ProgressBar 
                          now={progreso} 
                          variant={progreso >= 100 ? 'success' : progreso > 0 ? 'info' : 'secondary'}
                          style={{ height: 8 }}
                        />
                        <small className="text-muted">{progreso.toFixed(0)}%</small>
                      </td>
                      <td>
                        {deuda.fecha_vencimiento ? (
                          <span className={vencida ? 'text-danger fw-bold' : ''}>
                            {new Date(deuda.fecha_vencimiento).toLocaleDateString()}
                            {vencida && <ExclamationTriangle className="ms-1" size={12} />}
                          </span>
                        ) : '-'}
                      </td>
                      <td>{getEstadoBadge(vencida && deuda.estado !== 'pagado' ? 'vencido' : deuda.estado)}</td>
                      <td className="text-center">
                        {canPay && deuda.estado !== 'pagado' && (
                          <Button 
                            variant="outline-success" 
                            size="sm"
                            onClick={() => { 
                              setSelectedDeuda(deuda); 
                              setPagoData({ 
                                monto: deuda.monto_pendiente?.toString() || '', 
                                metodo_pago: 'efectivo', 
                                referencia: '' 
                              });
                              setShowPagoModal(true); 
                            }}
                            title="Registrar pago"
                          >
                            <CashStack />
                          </Button>
                        )}
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </Table>
        </Card.Body>
        <Card.Footer className="bg-transparent text-muted">
          Total: {filteredDeudas.length} deuda(s)
        </Card.Footer>
      </Card>

      {/* Modal de nueva deuda */}
      <Modal show={showModal} onHide={() => setShowModal(false)} size="lg">
        <Form onSubmit={handleSubmit}>
          <Modal.Header closeButton>
            <Modal.Title>Nueva Deuda</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Row>
              <Col md={8}>
                <Form.Group className="mb-3">
                  <Form.Label>Proveedor *</Form.Label>
                  <Form.Select
                    required
                    value={formData.proveedor_id}
                    onChange={(e) => setFormData({...formData, proveedor_id: e.target.value})}
                  >
                    <option value="">Seleccionar proveedor</option>
                    {proveedores.map(p => (
                      <option key={p.id} value={p.id}>{p.nombre}</option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>N° Factura *</Form.Label>
                  <Form.Control
                    required
                    value={formData.numero_factura}
                    onChange={(e) => setFormData({...formData, numero_factura: e.target.value})}
                    placeholder="Número de factura"
                  />
                </Form.Group>
              </Col>
            </Row>
            <Form.Group className="mb-3">
              <Form.Label>Notas</Form.Label>
              <Form.Control
                value={formData.notas}
                onChange={(e) => setFormData({...formData, notas: e.target.value})}
                placeholder="Descripción de la compra"
              />
            </Form.Group>
            <Row>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Monto Total (Bs.) *</Form.Label>
                  <Form.Control
                    type="number"
                    step="0.01"
                    required
                    value={formData.monto_total}
                    onChange={(e) => setFormData({...formData, monto_total: e.target.value})}
                    placeholder="0.00"
                  />
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Fecha Factura *</Form.Label>
                  <Form.Control
                    type="date"
                    required
                    value={formData.fecha_factura}
                    onChange={(e) => setFormData({...formData, fecha_factura: e.target.value})}
                  />
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Fecha de Vencimiento</Form.Label>
                  <Form.Control
                    type="date"
                    value={formData.fecha_vencimiento}
                    onChange={(e) => setFormData({...formData, fecha_vencimiento: e.target.value})}
                  />
                </Form.Group>
              </Col>
            </Row>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowModal(false)}>
              Cancelar
            </Button>
            <Button variant="primary" type="submit">
              Registrar Deuda
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>

      {/* Modal de pago */}
      <Modal show={showPagoModal} onHide={() => setShowPagoModal(false)}>
        <Form onSubmit={handlePago}>
          <Modal.Header closeButton>
            <Modal.Title>Registrar Pago</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            {selectedDeuda && (
              <Alert variant="info">
                <strong>Proveedor:</strong> {selectedDeuda.proveedor?.nombre}<br />
                <strong>Pendiente:</strong> Bs. {selectedDeuda.monto_pendiente?.toFixed(2)}
              </Alert>
            )}
            <Form.Group className="mb-3">
              <Form.Label>Monto a Pagar (Bs.) *</Form.Label>
              <Form.Control
                type="number"
                step="0.01"
                required
                max={selectedDeuda?.monto_pendiente}
                value={pagoData.monto}
                onChange={(e) => setPagoData({...pagoData, monto: e.target.value})}
                placeholder="0.00"
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Método de Pago</Form.Label>
              <Form.Select
                value={pagoData.metodo_pago}
                onChange={(e) => setPagoData({...pagoData, metodo_pago: e.target.value})}
              >
                <option value="efectivo">Efectivo</option>
                <option value="transferencia">Transferencia</option>
                <option value="cheque">Cheque</option>
                <option value="tarjeta">Tarjeta</option>
              </Form.Select>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Referencia / Comprobante</Form.Label>
              <Form.Control
                value={pagoData.referencia}
                onChange={(e) => setPagoData({...pagoData, referencia: e.target.value})}
                placeholder="Número de referencia o comprobante"
              />
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => setShowPagoModal(false)}>
              Cancelar
            </Button>
            <Button variant="success" type="submit">
              Registrar Pago
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Container>
  );
}

export default Deudas;
