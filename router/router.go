package router

import (
	"social-media-app/controllers"
	"social-media-app/metrics"
	"social-media-app/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartApp(DB *pgxpool.Pool) *gin.Engine {
	prometheus.Register(metrics.RequestHistogram)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("DB", DB)
		c.Next()
	})
	router.Use(middleware.Interceptor())

	userAccount := router.Group("/v1/user")
	{
		userAccount.POST("/register", middleware.RegisterValidator(), controllers.UserRegister)
		userAccount.POST("/login", middleware.AuthValidator(), controllers.UserLogin)
	}

	bankAccount := router.Group("/v1/balance")
	{
		bankAccount.POST("/", middleware.Authentication(), middleware.BankAccountValidator(), controllers.PostBalance)
		bankAccount.GET("/", middleware.Authentication(), controllers.GetBalance)
		bankAccount.GET("/history", middleware.Authentication(), controllers.GetBalanceByHistory)
	}

	transaction := router.Group("/v1/transaction")
	{
		transaction.POST("/", middleware.Authentication(), middleware.TransactionValidator(), controllers.PostTransaction)
	}

	router.POST("/v1/image", middleware.Authentication(), controllers.CreateUploadImage)

	router.GET("/health-check", controllers.ServerCheck)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}
