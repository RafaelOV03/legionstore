package controllers

// GUÍA DE REFACTORIZACIÓN DE CONTROLADORES
// ==========================================
// Este documento proporciona patrones consistentes para actualizar
// controladores a usar el nuevo error handling centralizado.

/*
PASO 1: Actualizar Imports
==========================

Cambiar de:
    import (
        "net/http"
        "smartech/backend/database"
        "smartech/backend/models"
    )

A:
    import (
        "smartech/backend/database"
        "smartech/backend/errors"
        "smartech/backend/models"
        "smartech/backend/validation"
    )

NOTA: Remover "net/http" y agregar "errors", "validation"


PASO 2: Reemplazar Patrones de Error
=====================================

Patrón 1 - Bad Request:
    OLD: c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
    NEW: apiErr := errors.NewBadRequest("Invalid id format")
         c.JSON(apiErr.Code, apiErr)

Patrón 2 - Not Found:
    OLD: c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
    NEW: apiErr := errors.NewNotFound("User", id)
         c.JSON(apiErr.Code, apiErr)

Patrón 3 - Conflict:
    OLD: c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
    NEW: apiErr := errors.NewConflict("Email already registered")
         c.JSON(apiErr.Code, apiErr)

Patrón 4 - Database Error:
    OLD: c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
    NEW: apiErr := errors.NewDatabaseError("Fetch users", err)
         c.JSON(apiErr.Code, apiErr)

Patrón 5 - Internal Server Error:
    OLD: c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process"})
    NEW: apiErr := errors.NewInternal("Failed to process request")
         c.JSON(apiErr.Code, apiErr)

Patrón 6 - Forbidden:
    OLD: c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
    NEW: apiErr := errors.ErrForbidden
         c.JSON(apiErr.Code, apiErr)

Patrón 7 - Unauthorized:
    OLD: c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
    NEW: apiErr := errors.ErrUnauthorized
         c.JSON(apiErr.Code, apiErr)


PASO 3: Agregar Validación de Structs
======================================

Patrón - Request Validation:
    var req struct {
        Name  string `json:"name" validate:"required,min=3"`
        Email string `json:"email" validate:"required,email"`
    }

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


PASO 4: Códigos HTTP Numéricos
===============================

Cambiar de:
    c.JSON(http.StatusOK, data)
    c.JSON(http.StatusCreated, data)

A:
    c.JSON(200, data)
    c.JSON(201, data)

Mapeo Completo:
    http.StatusOK             → 200
    http.StatusCreated        → 201
    http.StatusBadRequest     → 400 (usa errors.NewBadRequest)
    http.StatusUnauthorized   → 401 (usa errors.ErrUnauthorized)
    http.StatusForbidden      → 403 (usa errors.ErrForbidden)
    http.StatusNotFound       → 404 (usa errors.NewNotFound)
    http.StatusConflict       → 409 (usa errors.NewConflict)
    http.StatusUnprocessable  → 422 (usa ValidateStruct.ToAPIError)
    http.StatusInternalServer → 500 (usa errors.NewInternal)


PASO 5: Orden de Reemplazo Recomendado
=======================================

1. Actualizar imports (top of file)
2. Reemplazar en orden de aparición:
   - ShouldBindJSON errors
   - ValidateStruct agregaciones
   - Validaciones manuales
   - Database errors
   - 404/Not Found
   - 409/Conflict
   - Status codes numéricos
3. Compilar y verificar


CONTROLADORES YA REFACTORIZADOS
================================
✅ auth_controller.go (Register, Login, GetCurrentUser)
✅ product_controller.go (GetProducts, GetProduct, CreateProduct, UpdateProduct, DeleteProduct, GetRandomProducts, GetProductsByCategory)
✅ user_controller.go (GetUsers, GetUser, CreateUser, UpdateUser, DeleteUser)
✅ order_controller.go (GetOrder, GetOrders, GetAllOrders, CreateOrder)


SIGUIENTES CONTROLADORES A REFACTORIZAR
========================================
⏳ Otros controladores pueden seguir los mismos patrones

Palabras clave para búsqueda rápida:
    - c.JSON(http.Status → reemplazar con códigos numéricos
    - gin.H{"error" → reemplazar con errors.NewXXX
    - binding:"required" → cambiar a validate:"required"
    - ShouldBindJSON → agregar ValidateStruct después
*/
