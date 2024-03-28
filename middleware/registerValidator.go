package middleware

import (
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models"

	"github.com/gin-gonic/gin"
)

func RegisterValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var userRegister user.UserRegisterRequest

		if payloadValidationError := context.ShouldBindJSON(&userRegister); payloadValidationError != nil {
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

		context.Set("request", userRegister)
		context.Next()
	}
}