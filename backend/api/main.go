package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.POST("/contact", func(c *gin.Context) {})
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", func(c *gin.Context) {
		})
		authRoutes.POST("/signup", func(c *gin.Context) {
		})
	}
	serviceRoutes := router.Group("/services")
	{
		serviceRoutes.POST("/sync", func(c *gin.Context) {})
		calendarRoutes := serviceRoutes.Group("/calendar")
		{
			// Calendar routes
			calendarRoutes.GET("/", func(c *gin.Context) {
			})
			calendarRoutes.POST("/sync", func(c *gin.Context) {
			})
			calendarRoutes.GET("/:id", func(c *gin.Context) {
			})
			calendarRoutes.POST("/", func(c *gin.Context) {
			})
			calendarRoutes.PUT("/:id/connect", func(c *gin.Context) {
			})
			calendarRoutes.POST("/:id/sync", func(c *gin.Context) {
			})
			calendarRoutes.DELETE("/:id", func(c *gin.Context) {
			})
			calendarRoutes.PUT("/:id", func(c *gin.Context) {
			})
			calendarRoutes.PUT("/:id/disconnect", func(c *gin.Context) {
			})
		}
		musicRoutes := serviceRoutes.Group("/music")
		{
			// Music routes
			musicRoutes.GET("/", func(c *gin.Context) {
			})
			musicRoutes.POST("/sync", func(c *gin.Context) {
			})
			musicRoutes.GET("/:id", func(c *gin.Context) {
			})
			musicRoutes.POST("/", func(c *gin.Context) {
			})
			musicRoutes.PUT("/:id/connect", func(c *gin.Context) {
			})
			musicRoutes.POST("/:id/sync", func(c *gin.Context) {
			})
			musicRoutes.DELETE("/:id", func(c *gin.Context) {
			})
			musicRoutes.PUT("/:id", func(c *gin.Context) {
			})
			musicRoutes.PUT("/:id/disconnect", func(c *gin.Context) {
			})
		}
	}
	router.Run(":8080")
}
