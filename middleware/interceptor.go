package middleware

import (
	"log"
	"social-media-app/metrics"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Interceptor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		ctx.Next()

		httpStatusCode := strconv.Itoa(ctx.Writer.Status())
		httpMethod := ctx.Request.Method
		url := ctx.Request.RequestURI

		prometheusMetric(url, httpMethod, httpStatusCode)
	}
}

func prometheusMetric(path string, method string, status string) {
	startTime := time.Now()

	duration := time.Since(startTime).Seconds()
	log.Printf("Request path: %s, method: %s, status: %s, duration: %f seconds\n", path, method, status, duration)
	metrics.RequestHistogram.WithLabelValues(path, method, status).Observe(duration)
}