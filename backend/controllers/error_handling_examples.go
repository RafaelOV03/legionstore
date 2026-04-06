package controllers

import (
	"smartech/backend/errors"
	"smartech/backend/validation"

	"github.com/gin-gonic/gin"
)

// ==========================================
// EJEMPLOS DE USO DEL NUEVO ERROR HANDLING
// ==========================================

// Ejemplo 1: Usar ValidationErrors en controladores
func exampleValidationUsage(c *gin.Context) {
	type CreateUserRequest struct {
		Name     string `json:"name" validate:"required,min=3,max=100"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Validar la estructura
	validationErrors := validation.ValidateStruct(req)
	if len(validationErrors) > 0 {
		c.JSON(422, validationErrors.ToAPIError())
		return
	}

	// Procesar la solicitud
	c.JSON(200, gin.H{"message": "User created successfully"})
}

// Ejemplo 2: Usar errores predefinidos
func exampleNotFoundUsage(c *gin.Context) {
	userID := c.Param("id")

	// Simular búsqueda de usuario
	user := findUserByID(userID)
	if user == nil {
		err := errors.NewNotFound("User", userID)
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, user)
}

// Ejemplo 3: Manejo de errores de base de datos
func exampleDatabaseErrorUsage(c *gin.Context) {
	// result, err := someDatabase()
	// if err != nil {
	//     apiErr := errors.NewDatabaseError("Query products", err)
	//     c.JSON(apiErr.Code, apiErr)
	//     return
	// }

	c.JSON(200, gin.H{"message": "Success"})
}

// Ejemplo 4: Validar variable individual
func exampleVariableValidation(c *gin.Context) {
	email := c.Query("email")

	if err := validation.ValidateVar(email, "required,email"); err != nil {
		apiErr := errors.NewValidationError("email", "Invalid email format")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, gin.H{"message": "Email is valid"})
}

// ==========================================
// RESUMEN DE PATRONES
// ==========================================

/*
PATRÓN 1: Validación de Request Body
==================================
var req CreateUserRequest
if err := c.ShouldBindJSON(&req); err != nil {
    apiErr := errors.NewBadRequest(err.Error())
    c.JSON(apiErr.Code, apiErr)
    return
}

validationErrors := validation.ValidateStruct(req)
if len(validationErrors) > 0 {
    c.JSON(422, validationErrors.ToAPIError())
    return
}

PATRÓN 2: Buscar recurso y retornar NotFound
=============================================
resource := repo.GetByID(id)
if resource == nil {
    err := errors.NewNotFound("Resource", id)
    c.JSON(err.Code, err)
    return
}

PATRÓN 3: Manejo de errores de BD
================================
result, err := db.Exec(query, args...)
if err != nil {
    apiErr := errors.NewDatabaseError("operación", err)
    c.JSON(apiErr.Code, apiErr)
    return
}

PATRÓN 4: Conflicto (recurso duplicado)
=======================================
if userExists {
    err := errors.NewConflict("Email already registered")
    c.JSON(err.Code, err)
    return
}

PATRÓN 5: Error interno con detalles
====================================
if unexpectedCondition {
    err := errors.NewInternal("Unable to process request: " + details)
    c.JSON(err.Code, err)
    return
}

PATRÓN 6: Validación simple de variable
=====================================
if err := validation.ValidateVar(email, "email"); err != nil {
    apiErr := errors.NewValidationError("email", "Invalid format")
    c.JSON(apiErr.Code, apiErr)
    return
}

PATRÓN 7: Verificar validez completa
===================================
if !validation.IsValid(req) {
    // Hacer algo
}

*/

// Funciones auxiliares para ejemplos
func findUserByID(id string) interface{} {
	return nil // Simulado
}
