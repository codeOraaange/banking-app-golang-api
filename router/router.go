package router

import (
	"social-media-app/controllers"
	"social-media-app/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartApp(DB *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("DB", DB)
		c.Next()
	})

	userAccount := router.Group("/v1/user")
	{
		userAccount.POST("/register", middleware.RegisterValidator(), controllers.UserRegister)
		userAccount.POST("/login", middleware.AuthValidator(), controllers.UserLogin)
	}

	bankAccount := router.Group("/v1/balance")
	{
		bankAccount.POST("/", middleware.BankAccountValidator(), controllers.PostBalance)
		bankAccount.GET("/")
		bankAccount.GET("/history")
	}

	router.POST("/v1/image", middleware.Authentication(), controllers.CreateUploadImage)

	router.GET("/health-check", controllers.ServerCheck)

	return router
}
