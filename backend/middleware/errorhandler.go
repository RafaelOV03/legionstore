package middleware

import (
	"log"
	"smartech/backend/errors"

	"github.com/gin-gonic/gin"
)

// ErrorHandlingMiddleware es un middleware que maneja panics y errores globales
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Usar defer para capturar panics
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC: %v\n", r)

				// Convertir el panic a un error API
				err := errors.NewInternal("An unexpected error occurred")

				c.JSON(err.Code, err)
			}
		}()

		// Continuar con la siguiente ruta
		c.Next()

		// Opcional: Manejo de errores después de la ruta
		// (Por ejemplo, si algún handler asignó un error al contexto)
		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last()

			// Si es un APIError, usarlo directamente
			if apiErr, ok := lastErr.Err.(*errors.APIError); ok {
				c.JSON(apiErr.Code, apiErr)
				return
			}

			// Si no es APIError, envolver en un error interno
			err := errors.NewInternal(lastErr.Error())
			c.JSON(err.Code, err)
		}
	}
}

// JSONErrorMiddleware retorna errores en formato JSON consistente
func JSONErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}
