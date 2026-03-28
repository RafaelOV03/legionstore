package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/models"
	"smartech/backend/repositories"
	"smartech/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getSedeService() *services.SedeService {
	repo := repositories.NewSedeRepository(database.DB)
	return services.NewSedeService(repo)
}

// GetSedes obtiene todas las sedes
func GetSedes(c *gin.Context) {
	sedes, err := getSedeService().ListSedes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sedes"})
		return
	}

	c.JSON(http.StatusOK, sedes)
}

// GetSede obtiene una sede por ID
func GetSede(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sede ID"})
		return
	}

	sede, err := getSedeService().GetSede(id)
	if err == repositories.ErrSedeNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sede not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sede"})
		return
	}
	c.JSON(http.StatusOK, sede)
}

// CreateSede crea una nueva sede
func CreateSede(c *gin.Context) {
	var sede models.Sede
	if err := c.ShouldBindJSON(&sede); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdSede, err := getSedeService().CreateSede(sede)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sede"})
		return
	}

	// Log de auditoría
	logAuditoria(c, "crear", "sede", createdSede.ID, "", createdSede.Nombre)

	c.JSON(http.StatusCreated, createdSede)
}

// UpdateSede actualiza una sede
func UpdateSede(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sede ID"})
		return
	}

	var sede models.Sede
	if err := c.ShouldBindJSON(&sede); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSede, err := getSedeService().UpdateSede(id, sede)
	if err == repositories.ErrSedeNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sede not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sede"})
		return
	}
	logAuditoria(c, "editar", "sede", id, "", updatedSede.Nombre)

	c.JSON(http.StatusOK, updatedSede)
}

// DeleteSede elimina una sede
func DeleteSede(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sede ID"})
		return
	}

	err = getSedeService().DeleteSede(id)
	if err == repositories.ErrSedeHasAssociatedUsers {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede eliminar la sede, tiene usuarios asociados"})
		return
	}
	if err == repositories.ErrSedeNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sede not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sede"})
		return
	}

	logAuditoria(c, "eliminar", "sede", id, "", "")
	c.JSON(http.StatusOK, gin.H{"message": "Sede deleted successfully"})
}
