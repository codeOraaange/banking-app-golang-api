package middleware

import (
	"fmt"
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models/bankAccount"

	"github.com/gin-gonic/gin"
)

func BankAccountValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var bankAccountRequest bankAccount.BankAccountRequest
		fmt.Println("line 14")
		if payloadValidationError := context.ShouldBindJSON(&bankAccountRequest); payloadValidationError != nil {
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
		fmt.Println("line 31")
		context.Set("request", bankAccountRequest)
		context.Next()
	}
}
