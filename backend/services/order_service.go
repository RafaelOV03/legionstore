package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"smartech/backend/models"
	"smartech/backend/repositories"
	"time"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderCartNotFound  = errors.New("cart not found")
	ErrOrderCartEmpty     = errors.New("cart is empty")
	ErrOrderPayPalToken   = errors.New("paypal token error")
	ErrOrderPayPalCreate  = errors.New("paypal order create error")
	ErrOrderPayPalCapture = errors.New("paypal order capture error")
)

type PayPalConfig struct {
	Clientid string
	Secret   string
	BaseURL  string
}

type OrderService struct {
	repo         *repositories.OrderRepository
	paypalConfig PayPalConfig
}

func NewOrderService(repo *repositories.OrderRepository) *OrderService {
	cfg := PayPalConfig{
		Clientid: os.Getenv("PAYPAL_CLIENT_id"),
		Secret:   os.Getenv("PAYPAL_SECRET"),
		BaseURL:  "https://api-m.sandbox.paypal.com",
	}
	if cfg.Clientid == "" {
		cfg.Clientid = "Aec0f6y57ztt6tvhADUXj1GAZqlH_vF0hGyK89cWblFEs3Gq3kunXw1iaqGeKe32CjPMuxy3PMuzLz6A"
	}
	if cfg.Secret == "" {
		cfg.Secret = "EAF3VXSIZfDhEr522LM9UAyVOGn0kVvyiJBaJKPah4iRweg6MpAv2IZR8TDNhqj-9GUZuS17P066shVf"
	}

	return &OrderService{repo: repo, paypalConfig: cfg}
}

func (s *OrderService) GetPayPalClientID() string {
	return s.paypalConfig.Clientid
}

func (s *OrderService) getPayPalAccessToken() (string, error) {
	url := fmt.Sprintf("%s/v1/oauth2/token", s.paypalConfig.BaseURL)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(s.paypalConfig.Clientid, s.paypalConfig.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", ErrOrderPayPalToken
	}

	return token, nil
}

func (s *OrderService) CreateOrder(sessionID string) (models.Order, map[string]interface{}, error) {
	cart, err := s.repo.GetCartBySession(sessionID)
	if err == sql.ErrNoRows {
		return models.Order{}, nil, ErrOrderCartNotFound
	}
	if err != nil {
		return models.Order{}, nil, err
	}

	cartItems, err := s.repo.GetCartItemsWithProducts(cart.ID)
	if err != nil {
		return models.Order{}, nil, err
	}
	if len(cartItems) == 0 {
		return models.Order{}, nil, ErrOrderCartEmpty
	}

	var total float64
	items := make([]map[string]interface{}, 0, len(cartItems))
	for _, item := range cartItems {
		itemTotal := float64(item.Quantity) * item.Product.PrecioVenta
		total += itemTotal
		items = append(items, map[string]interface{}{
			"name": item.Product.Name,
			"unit_amount": map[string]interface{}{
				"currency_code": "USD",
				"value":         fmt.Sprintf("%.2f", item.Product.PrecioVenta),
			},
			"quantity": fmt.Sprintf("%d", item.Quantity),
		})
	}

	token, err := s.getPayPalAccessToken()
	if err != nil {
		return models.Order{}, nil, ErrOrderPayPalToken
	}

	orderData := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"amount": map[string]interface{}{
					"currency_code": "USD",
					"value":         fmt.Sprintf("%.2f", total),
					"breakdown": map[string]interface{}{
						"item_total": map[string]interface{}{
							"currency_code": "USD",
							"value":         fmt.Sprintf("%.2f", total),
						},
					},
				},
				"items": items,
			},
		},
	}

	orderJSON, _ := json.Marshal(orderData)
	url := fmt.Sprintf("%s/v2/checkout/orders", s.paypalConfig.BaseURL)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(orderJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Order{}, nil, ErrOrderPayPalCreate
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	paypalOrder := make(map[string]interface{})
	if err := json.Unmarshal(body, &paypalOrder); err != nil {
		return models.Order{}, nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return models.Order{}, paypalOrder, ErrOrderPayPalCreate
	}

	paypalOrderID, ok := paypalOrder["id"].(string)
	if !ok {
		return models.Order{}, paypalOrder, ErrOrderPayPalCreate
	}

	order, err := s.repo.CreateOrderWithItems(sessionID, paypalOrderID, total, cartItems)
	if err != nil {
		return models.Order{}, paypalOrder, err
	}

	return order, paypalOrder, nil
}

func (s *OrderService) CaptureOrder(paypalOrderID string) (models.Order, map[string]interface{}, error) {
	token, err := s.getPayPalAccessToken()
	if err != nil {
		return models.Order{}, nil, ErrOrderPayPalToken
	}

	url := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", s.paypalConfig.BaseURL, paypalOrderID)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Order{}, nil, ErrOrderPayPalCapture
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	captureResult := make(map[string]interface{})
	if err := json.Unmarshal(body, &captureResult); err != nil {
		return models.Order{}, nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return models.Order{}, captureResult, ErrOrderPayPalCapture
	}

	order, err := s.repo.GetOrderByPayPalID(paypalOrderID)
	if err == sql.ErrNoRows {
		return models.Order{}, captureResult, ErrOrderNotFound
	}
	if err != nil {
		return models.Order{}, captureResult, err
	}

	payerEmail := ""
	payerName := ""
	if payer, ok := captureResult["payer"].(map[string]interface{}); ok {
		if email, ok := payer["email_address"].(string); ok {
			payerEmail = email
		}
		if name, ok := payer["name"].(map[string]interface{}); ok {
			if givenName, ok := name["given_name"].(string); ok {
				payerName = givenName
				if surname, ok := name["surname"].(string); ok {
					payerName += " " + surname
				}
			}
		}
	}

	now := time.Now()
	if err := s.repo.MarkOrderCompleted(order.ID, now, payerEmail, payerName); err != nil {
		return models.Order{}, captureResult, err
	}

	order.Status = "COMPLETED"
	order.CompletedAt = &now
	order.PayerEmail = payerEmail
	order.PayerName = payerName

	_ = s.repo.DecreaseStockByOrder(order.ID)
	_ = s.repo.ClearCartBySession(order.SessionID)

	return order, captureResult, nil
}

func (s *OrderService) GetOrder(paypalOrderID string) (models.Order, error) {
	order, err := s.repo.GetOrderByPayPalID(paypalOrderID)
	if err == sql.ErrNoRows {
		return models.Order{}, ErrOrderNotFound
	}
	if err != nil {
		return models.Order{}, err
	}

	items, _ := s.repo.GetOrderItems(order.ID)
	order.OrderItems = items
	return order, nil
}

func (s *OrderService) GetOrdersBySession(sessionID string) ([]models.Order, error) {
	orders, err := s.repo.ListOrdersBySession(sessionID)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		items, _ := s.repo.GetOrderItems(orders[i].ID)
		orders[i].OrderItems = items
	}

	return orders, nil
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	orders, err := s.repo.ListAllOrders()
	if err != nil {
		return nil, err
	}

	for i := range orders {
		items, _ := s.repo.GetOrderItems(orders[i].ID)
		orders[i].OrderItems = items
	}

	return orders, nil
}

func (s *OrderService) FinalizeOrder(orderID int64) (models.Order, error) {
	count, err := s.repo.CountOrderByID(orderID)
	if err != nil || count == 0 {
		return models.Order{}, ErrOrderNotFound
	}

	if err := s.repo.FinalizeOrder(orderID); err != nil {
		return models.Order{}, err
	}

	return s.repo.GetOrderByID(orderID)
}

func (s *OrderService) UpdateOrder(orderID int64, status string, totalAmount float64) (models.Order, error) {
	count, err := s.repo.CountOrderByID(orderID)
	if err != nil || count == 0 {
		return models.Order{}, ErrOrderNotFound
	}

	if err := s.repo.UpdateOrderFields(orderID, status, totalAmount); err != nil {
		return models.Order{}, err
	}

	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return models.Order{}, err
	}
	items, _ := s.repo.GetOrderItems(order.ID)
	order.OrderItems = items
	return order, nil
}

func (s *OrderService) DeleteOrder(orderID int64) error {
	count, err := s.repo.CountOrderByID(orderID)
	if err != nil || count == 0 {
		return ErrOrderNotFound
	}
	return s.repo.DeleteOrder(orderID)
}
