package services

import (
	"database/sql"
	"errors"
	"smartech/backend/models"
	"smartech/backend/repositories"
)

var (
	ErrCartNotFound      = errors.New("cart not found")
	ErrCartItemNotFound  = errors.New("cart item not found")
	ErrCartUnauthorized  = errors.New("unauthorized cart access")
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type AddToCartInput struct {
	SessionID string
	ProductID int64
	Quantity  int
}

type UpdateCartItemInput struct {
	SessionID string
	ItemID    int64
	Quantity  int
}

type CartService struct {
	repo *repositories.CartRepository
}

func NewCartService(repo *repositories.CartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) getOrCreateCart(sessionID string) (models.Cart, error) {
	cart, err := s.repo.GetCartBySession(sessionID)
	if err == sql.ErrNoRows {
		cartID, createErr := s.repo.CreateCart(sessionID)
		if createErr != nil {
			return models.Cart{}, createErr
		}
		cart.ID = cartID
		cart.SessionID = sessionID
		cart.CartItems = []models.CartItem{}
		return cart, nil
	}
	if err != nil {
		return models.Cart{}, err
	}
	return cart, nil
}

func (s *CartService) GetCart(sessionID string) (models.Cart, error) {
	cart, err := s.getOrCreateCart(sessionID)
	if err != nil {
		return models.Cart{}, err
	}

	items, err := s.repo.GetCartItemsWithProducts(cart.ID)
	if err != nil {
		cart.CartItems = []models.CartItem{}
		return cart, nil
	}
	cart.CartItems = items
	if cart.CartItems == nil {
		cart.CartItems = []models.CartItem{}
	}

	return cart, nil
}

func (s *CartService) AddToCart(input AddToCartInput) (models.Cart, error) {
	stockQuantity, err := s.repo.GetProductStock(input.ProductID)
	if err == sql.ErrNoRows {
		return models.Cart{}, ErrProductNotFound
	}
	if err != nil {
		return models.Cart{}, err
	}
	if stockQuantity < input.Quantity {
		return models.Cart{}, ErrInsufficientStock
	}

	cart, err := s.getOrCreateCart(input.SessionID)
	if err != nil {
		return models.Cart{}, err
	}

	itemID, existingQty, err := s.repo.GetCartItemByCartAndProduct(cart.ID, input.ProductID)
	if err == nil {
		if err := s.repo.UpdateCartItemQuantity(itemID, existingQty+input.Quantity); err != nil {
			return models.Cart{}, err
		}
	} else if err == sql.ErrNoRows {
		if err := s.repo.InsertCartItem(cart.ID, input.ProductID, input.Quantity); err != nil {
			return models.Cart{}, err
		}
	} else {
		return models.Cart{}, err
	}

	items, _ := s.repo.GetCartItemsWithProducts(cart.ID)
	cart.CartItems = items
	if cart.CartItems == nil {
		cart.CartItems = []models.CartItem{}
	}
	return cart, nil
}

func (s *CartService) UpdateCartItem(input UpdateCartItemInput) (models.CartItem, error) {
	item, err := s.repo.GetCartItemWithProduct(input.ItemID)
	if err == sql.ErrNoRows {
		return models.CartItem{}, ErrCartItemNotFound
	}
	if err != nil {
		return models.CartItem{}, err
	}

	cartSessionID, err := s.repo.GetSessionByCartID(item.CartID)
	if err != nil || cartSessionID != input.SessionID {
		return models.CartItem{}, ErrCartUnauthorized
	}

	if err := s.repo.UpdateCartItemQuantity(input.ItemID, input.Quantity); err != nil {
		return models.CartItem{}, err
	}

	item.Quantity = input.Quantity
	return item, nil
}

func (s *CartService) RemoveFromCart(sessionID string, itemID int64) error {
	cartID, err := s.repo.GetCartIDByItem(itemID)
	if err == sql.ErrNoRows {
		return ErrCartItemNotFound
	}
	if err != nil {
		return err
	}

	cartSessionID, err := s.repo.GetSessionByCartID(cartID)
	if err != nil || cartSessionID != sessionID {
		return ErrCartUnauthorized
	}

	return s.repo.DeleteCartItem(itemID)
}

func (s *CartService) ClearCart(sessionID string) error {
	cartID, err := s.repo.GetCartIDBySession(sessionID)
	if err == sql.ErrNoRows {
		return ErrCartNotFound
	}
	if err != nil {
		return err
	}

	return s.repo.DeleteCartItemsByCartID(cartID)
}
