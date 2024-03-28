package middleware

import (
	"net/http"
	"banking-app-golang-api/helpers"
	"banking-app-golang-api/models"

	"github.com/gin-gonic/gin"
)

func RegisterValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var user models.UsersForAuth

		if payloadValidationError := context.ShouldBindJSON(&user); payloadValidationError != nil {
			var errors []string

			if payloadValidationError.Error() == "EOF" {
				errors = append(errors, "Request body is missing")
			} else {
				errors = helpers.GeneralValidator(payloadValidationError)
			}

			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   errors,
				"message": "Failed to validate",
			})
			return
		}

		context.Set("request", user)
		context.Next()
	}
}
