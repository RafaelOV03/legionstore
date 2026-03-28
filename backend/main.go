package main

import (
	"log"
	"os"
	"smartech/backend/controllers"
	"smartech/backend/database"
	"smartech/backend/middleware"

	"github.com/gin-gonic/gin"
)



/*
 * ==========================================
 * PUNTO DE ENTRADA PRINCIPAL
 * ==========================================
 */

func main() {
	// Inicializar la base de datos
	log.Println("Initializing database...")
	database.InitDatabase()
	log.Println("Database initialized.")

	// Crear el router de Gin
	log.Println("Creating Gin router...")
	router := gin.Default()
	log.Println("Gin router created.")

	// Usar el middleware de CORS
	log.Println("Using CORS middleware...")
	router.Use(middleware.CorsMiddleware())
	log.Println("CORS middleware used.")

	// Agrupar las rutas de la API
	log.Println("Grouping API routes...")
	api := router.Group("/api")
	{
		// ------------------------------------------
		// MÓDULO 1: GESTIÓN DE ACCESO (Auth)
		// ------------------------------------------
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.POST("/logout", controllers.Logout)
			auth.GET("/me", middleware.AuthMiddleware(), controllers.GetCurrentUser)
		}

		// ==========================================
		// MÓDULO 2: NÚCLEO DE NEGOCIO (Products & Stock)
		// ==========================================
		products := api.Group("/products")
		{
			products.GET("", controllers.GetProducts)
			products.GET("/random", controllers.GetRandomProducts)
			products.GET("/category/:category", controllers.GetProductsByCategory)
			products.GET("/:id", controllers.GetProduct)
			products.GET("/:id/images", controllers.GetProductImages)
			// Rutas protegidas
			products.POST("", middleware.AuthMiddleware(), middleware.RequirePermission("products.create"), controllers.CreateProduct)
			products.PUT("/:id", middleware.AuthMiddleware(), middleware.RequirePermission("products.update"), controllers.UpdateProduct)
			products.PUT("/:id/images", middleware.AuthMiddleware(), middleware.RequirePermission("products.update"), controllers.UpdateProductImages)
			products.DELETE("/:id", middleware.AuthMiddleware(), middleware.RequirePermission("products.delete"), controllers.DeleteProduct)
			// Actualizar precios (administrador)
			products.PATCH("/:id/precios", middleware.AuthMiddleware(), middleware.RequirePermission("precios.update"), controllers.UpdateProduct)
		}

		// ==========================================
		// MÓDULO 3: OPERACIONES Y LOGÍSTICA (RMA, Traspasos, Ordenes)
		// ==========================================
		upload := api.Group("/upload")
		{
			upload.POST("/image", middleware.AuthMiddleware(), middleware.RequirePermission("products.create"), controllers.UploadProductImage)
			upload.DELETE("/image", middleware.AuthMiddleware(), middleware.RequirePermission("products.delete"), controllers.DeleteProductImage)
		}

		// ------------------------------------------
		// MÓDULO 4: VENTAS Y LOGÍSTICA (Cotizaciones, Traspasos, RMA)
		// ------------------------------------------
		sedes := api.Group("/sedes")
		sedes.Use(middleware.AuthMiddleware())
		{
			sedes.GET("", controllers.GetSedes)
			sedes.GET("/:id", controllers.GetSede)
			sedes.POST("", middleware.RequirePermission("sedes.create"), controllers.CreateSede)
			sedes.PUT("/:id", middleware.RequirePermission("sedes.update"), controllers.UpdateSede)
			sedes.DELETE("/:id", middleware.RequirePermission("sedes.delete"), controllers.DeleteSede)
		}

		// Stock multisede (vendedor)
		stock := api.Group("/stock")
		stock.Use(middleware.AuthMiddleware())
		{
			stock.GET("", middleware.RequirePermission("stock.read"), controllers.GetStockMultisede)
			stock.GET("/sede/:sede_id", middleware.RequirePermission("stock.read"), controllers.GetStockBySede)
			stock.PUT("", middleware.RequirePermission("stock.update"), controllers.UpdateStock)
		}

		// ==========================================
		// RUTAS DE RMA/GARANTÍAS (administrador)
		// ==========================================
		rmas := api.Group("/rmas")
		rmas.Use(middleware.AuthMiddleware())
		{
			rmas.GET("", middleware.RequirePermission("rmas.read"), controllers.GetRMAs)
			rmas.GET("/stats", middleware.RequirePermission("rmas.read"), controllers.GetRMAStats)
			rmas.GET("/:id", middleware.RequirePermission("rmas.read"), controllers.GetRMA)
			rmas.POST("", middleware.RequirePermission("rmas.create"), controllers.CreateRMA)
			rmas.PUT("/:id", middleware.RequirePermission("rmas.update"), controllers.UpdateRMA)
			rmas.DELETE("/:id", middleware.RequirePermission("rmas.delete"), controllers.DeleteRMA)
		}

		// ==========================================
		// RUTAS DE COTIZACIONES (vendedor)
		// ==========================================
		cotizaciones := api.Group("/cotizaciones")
		cotizaciones.Use(middleware.AuthMiddleware())
		{
			cotizaciones.GET("", middleware.RequirePermission("cotizaciones.read"), controllers.GetCotizaciones)
			cotizaciones.GET("/:id", middleware.RequirePermission("cotizaciones.read"), controllers.GetCotizacion)
			cotizaciones.GET("/:id/pdf", middleware.RequirePermission("cotizaciones.read"), controllers.GenerarPDFCotizacion)
			cotizaciones.POST("", middleware.RequirePermission("cotizaciones.create"), controllers.CreateCotizacion)
			cotizaciones.PUT("/:id/estado", middleware.RequirePermission("cotizaciones.update"), controllers.UpdateCotizacionEstado)
			cotizaciones.POST("/:id/convertir", middleware.RequirePermission("cotizaciones.update"), controllers.ConvertirCotizacionAVenta)
			cotizaciones.DELETE("/:id", middleware.RequirePermission("cotizaciones.delete"), controllers.DeleteCotizacion)
		}

		// ==========================================
		// RUTAS DE TRASPASOS (administrador)
		// ==========================================
		traspasos := api.Group("/traspasos")
		traspasos.Use(middleware.AuthMiddleware())
		{
			traspasos.GET("", middleware.RequirePermission("traspasos.read"), controllers.GetTraspasos)
			traspasos.GET("/:id", middleware.RequirePermission("traspasos.read"), controllers.GetTraspaso)
			traspasos.POST("", middleware.RequirePermission("traspasos.create"), controllers.CreateTraspaso)
			traspasos.POST("/:id/enviar", middleware.RequirePermission("traspasos.update"), controllers.EnviarTraspaso)
			traspasos.POST("/:id/recibir", middleware.RequirePermission("traspasos.update"), controllers.RecibirTraspaso)
			traspasos.POST("/:id/cancelar", middleware.RequirePermission("traspasos.update"), controllers.CancelarTraspaso)
			traspasos.DELETE("/:id", middleware.RequirePermission("traspasos.delete"), controllers.DeleteTraspaso)
		}

		// ------------------------------------------
		// MÓDULO 5: SERVICIO TÉCNICO E INSUMOS
		// ------------------------------------------
		ordenes := api.Group("/ordenes-trabajo")
		ordenes.Use(middleware.AuthMiddleware())
		{
			ordenes.GET("", middleware.RequirePermission("ordenes.read"), controllers.GetOrdenesTrabajo)
			ordenes.GET("/stats", middleware.RequirePermission("ordenes.read"), controllers.GetOrdenesStats)
			ordenes.GET("/tecnicos", middleware.RequirePermission("ordenes.read"), controllers.GetTecnicos) // Lista de técnicos para asignar
			ordenes.GET("/:id", middleware.RequirePermission("ordenes.read"), controllers.GetOrdenTrabajo)
			ordenes.POST("", middleware.RequirePermission("ordenes.create"), controllers.CreateOrdenTrabajo)
			ordenes.PUT("/:id", middleware.RequirePermission("ordenes.update"), controllers.UpdateOrdenTrabajo)
			ordenes.POST("/:id/asignar", middleware.RequirePermission("ordenes.update"), controllers.AsignarTecnico)
			ordenes.POST("/:id/insumo", middleware.RequirePermission("ordenes.update"), controllers.AgregarInsumo)
			ordenes.POST("/:id/trazabilidad", middleware.RequirePermission("ordenes.update"), controllers.RegistrarTrazabilidad)
			ordenes.DELETE("/:id", middleware.RequirePermission("ordenes.delete"), controllers.DeleteOrdenTrabajo)
		}

		// ------------------------------------------
		// MÓDULO 6: FINANZAS Y PROVEEDORES
		// ------------------------------------------
		proveedores := api.Group("/proveedores")
		proveedores.Use(middleware.AuthMiddleware())
		{
			proveedores.GET("", middleware.RequirePermission("proveedores.read"), controllers.GetProveedores)
			proveedores.GET("/:id", middleware.RequirePermission("proveedores.read"), controllers.GetProveedor)
			proveedores.POST("", middleware.RequirePermission("proveedores.create"), controllers.CreateProveedor)
			proveedores.PUT("/:id", middleware.RequirePermission("proveedores.update"), controllers.UpdateProveedor)
			proveedores.DELETE("/:id", middleware.RequirePermission("proveedores.delete"), controllers.DeleteProveedor)
		}

		deudas := api.Group("/deudas")
		deudas.Use(middleware.AuthMiddleware())
		{
			deudas.GET("", middleware.RequirePermission("deudas.read"), controllers.GetDeudas)
			deudas.GET("/resumen", middleware.RequirePermission("deudas.read"), controllers.GetResumenDeudas)
			deudas.POST("", middleware.RequirePermission("deudas.create"), controllers.CreateDeuda)
			deudas.POST("/:id/pago", middleware.RequirePermission("deudas.update"), controllers.RegistrarPago)
			deudas.GET("/:id/pagos", middleware.RequirePermission("deudas.read"), controllers.GetPagosDeuda)
		}

		// ==========================================
		// RUTAS DE INSUMOS (técnico)
		// ==========================================
		insumos := api.Group("/insumos")
		insumos.Use(middleware.AuthMiddleware())
		{
			insumos.GET("", middleware.RequirePermission("insumos.read"), controllers.GetInsumos)
			insumos.GET("/stats", middleware.RequirePermission("insumos.read"), controllers.GetInsumosStats)
			insumos.GET("/:id", middleware.RequirePermission("insumos.read"), controllers.GetInsumo)
			insumos.POST("", middleware.RequirePermission("insumos.create"), controllers.CreateInsumo)
			insumos.PUT("/:id", middleware.RequirePermission("insumos.update"), controllers.UpdateInsumo)
			insumos.POST("/:id/ajustar", middleware.RequirePermission("insumos.update"), controllers.AjustarStockInsumo)
			insumos.DELETE("/:id", middleware.RequirePermission("insumos.delete"), controllers.DeleteInsumo)
		}

		// ==========================================
		// RUTAS DE COMPATIBILIDAD (vendedor)
		// ==========================================
		compatibilidad := api.Group("/compatibilidad")
		compatibilidad.Use(middleware.AuthMiddleware())
		{
			compatibilidad.GET("", middleware.RequirePermission("compatibilidad.read"), controllers.GetCompatibilidades)
			compatibilidad.GET("/buscar/:id", middleware.RequirePermission("compatibilidad.read"), controllers.BuscarCompatibles)
			compatibilidad.POST("", middleware.RequirePermission("compatibilidad.create"), controllers.CreateCompatibilidad)
			compatibilidad.DELETE("/:id", middleware.RequirePermission("compatibilidad.delete"), controllers.DeleteCompatibilidad)
		}

		// ==========================================
		// RUTAS DE AUDITORÍA Y REPORTES (gerente)
		// ==========================================
		auditoria := api.Group("/auditoria")
		auditoria.Use(middleware.AuthMiddleware())
		{
			auditoria.GET("/logs", middleware.RequirePermission("auditoria.read"), controllers.GetLogs)
			auditoria.GET("/stats", middleware.RequirePermission("auditoria.read"), controllers.GetLogStats)
		}

		reportes := api.Group("/reportes")
		reportes.Use(middleware.AuthMiddleware())
		{
			reportes.GET("/ganancias", middleware.RequirePermission("reportes.read"), controllers.GetReporteGanancias)
		}

		// ==========================================
		// RUTAS DE SEGMENTACIÓN Y PROMOCIONES (gerente)
		// ==========================================
		segmentaciones := api.Group("/segmentaciones")
		segmentaciones.Use(middleware.AuthMiddleware())
		{
			segmentaciones.GET("", middleware.RequirePermission("segmentacion.read"), controllers.GetSegmentaciones)
			segmentaciones.POST("", middleware.RequirePermission("segmentacion.create"), controllers.CreateSegmentacion)
		}

		// ------------------------------------------
		// MÓDULO 7: MARKETING Y AUDITORÍA (Reportes)
		// ------------------------------------------

		promociones := api.Group("/promociones")
		promociones.Use(middleware.AuthMiddleware())
		{
			promociones.GET("", middleware.RequirePermission("promociones.read"), controllers.GetPromociones)
			promociones.POST("", middleware.RequirePermission("promociones.create"), controllers.CreatePromocion)
			promociones.PUT("/:id", middleware.RequirePermission("promociones.update"), controllers.UpdatePromocion)
			promociones.DELETE("/:id", middleware.RequirePermission("promociones.delete"), controllers.DeletePromocion)
		}

		// ------------------------------------------
		// MÓDULO 8: ADMINISTRACIÓN DE SISTEMA (Users & Roles)
		// ------------------------------------------
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("", middleware.RequirePermission("users.read"), controllers.GetUsers)
			users.GET("/:id", middleware.RequirePermission("users.read"), controllers.GetUser)
			users.POST("", middleware.RequirePermission("users.create"), controllers.CreateUser)
			users.PUT("/:id", middleware.RequirePermission("users.update"), controllers.UpdateUser)
			users.DELETE("/:id", middleware.RequirePermission("users.delete"), controllers.DeleteUser)
		}

		// ==========================================
		// RUTAS DE ROLES (administrador)
		// ==========================================
		roles := api.Group("/roles")
		roles.Use(middleware.AuthMiddleware())
		{
			roles.GET("", middleware.RequirePermission("roles.read"), controllers.GetRoles)
			roles.GET("/:id", middleware.RequirePermission("roles.read"), controllers.GetRole)
			roles.POST("", middleware.RequirePermission("roles.create"), controllers.CreateRole)
			roles.PUT("/:id", middleware.RequirePermission("roles.update"), controllers.UpdateRole)
			roles.DELETE("/:id", middleware.RequirePermission("roles.delete"), controllers.DeleteRole)
		}

		// Rutas de permisos (solo lectura)
		api.GET("/permissions", middleware.AuthMiddleware(), middleware.RequirePermission("roles.read"), controllers.GetPermissions)
	}
	log.Println("API routes grouped.")

	// Servir archivos estáticos del directorio uploads
    router.Static("/uploads", "./uploads")

    // --- CONFIGURACIÓN DE ARRANQUE ---
    
    // Obtener el puerto de las variables de entorno o usar 8080 por defecto
    port := os.Getenv("PORT") 
    if port == "" {
        port = "8080"
    }

    // Un banner profesional para identificar el proyecto en la consola
    log.Println("==================================================")
    log.Println("   LEGION STORE - BACKEND SERVICES (UPDS 2026)    ")
    log.Println("   Status: ONLINE                                 ")
    log.Println("   Running on port: " + port                       )
    log.Println("==================================================")
    
    // Iniciar el servidor con el puerto configurado
    router.Run(":" + port)
}