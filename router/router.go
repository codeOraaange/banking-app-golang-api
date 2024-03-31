package router

import (
	"log"
	"social-media-app/controllers"
	"social-media-app/middleware"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Histogram of request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10),
	}, []string{"path", "method", "status"})
)

func StartApp(DB *pgxpool.Pool) *gin.Engine {
	prometheus.Register(requestHistogram)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("DB", DB)
		c.Next()
	})

	userAccount := router.Group("/v1/user")
	{
		userAccount.POST("/register", wrapHandlerWithMetrics("v1/user/register", "POST", controllers.UserRegister, middleware.RegisterValidator()))
		userAccount.POST("/login", wrapHandlerWithMetrics("v1/user/login", "POST", controllers.UserLogin, middleware.AuthValidator()))
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

func wrapHandlerWithMetrics(path, method string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc) gin.HandlerFunc {
  return func(c *gin.Context) {
    startTime := time.Now()

    // Chain middleware functions
    for _, m := range middleware {
      m(c)
    }

    handler(c)

    duration := time.Since(startTime).Seconds()
    log.Printf("Request path: %s, method: %s, status: %d, duration: %f seconds\n", path, method, c.Writer.Status(), duration)
    requestHistogram.WithLabelValues(path, method, strconv.Itoa(c.Writer.Status())).Observe(duration)
  }
}
