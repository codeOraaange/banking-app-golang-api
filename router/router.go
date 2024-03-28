package router

import (
	"banking-app-golang-api/controllers"
	// "banking-app-golang-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartApp(DB *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("DB", DB)
		c.Next()
	})

	// userAccount := router.Group("/v1/user")
	// {
	// 	userAccount.POST("/register", middleware.RegisterValidator(), controllers.UserRegister)
	// 	userAccount.POST("/login", middleware.AuthValidator(), controllers.UserLogin)
	// }

	router.GET("/health-check", controllers.ServerCheck)

	return router
}
