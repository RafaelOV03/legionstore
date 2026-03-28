package controllers

import (
	"smartech/backend/database"
	"smartech/backend/models"
	"smartech/backend/repositories"
)

func getCartItemsWithProducts(cartID int64) []models.CartItem {
	repo := repositories.NewCartRepository(database.DB)
	items, err := repo.GetCartItemsWithProducts(cartID)
	if err != nil {
		return []models.CartItem{}
	}
	return items
}
