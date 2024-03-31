package middleware

import (
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models/transaction"

	"github.com/gin-gonic/gin"
)

func TransactionValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var transactionRequest transaction.TransactionRequest
		if payloadValidationError := context.ShouldBindJSON(&transactionRequest); payloadValidationError != nil {
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

		// Check if the currency code is valid
		if !isValidCurrency(transactionRequest.Currency) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid currency code",
				"message": "Failed to validate",
			})
			return
		}

		context.Set("request", transactionRequest)
		context.Next()
	}

}