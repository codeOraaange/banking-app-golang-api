package middleware

import (
	// "fmt"
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models/bankAccount"

	"github.com/gin-gonic/gin"
)

func BankAccountValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var bankAccountRequest bankAccount.BankAccountRequest

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

		// Check if the currency code is valid
		if !isValidCurrency(bankAccountRequest.Currency) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid currency code",
				"message": "Failed to validate",
			})
			return
		}

		context.Set("request", bankAccountRequest)
		context.Next()
	}
}

// isValidCurrency checks if the provided currency code is valid
func isValidCurrency(currencyCode string) bool {
	// For simplicity, let's assume only a few currency codes are valid
	validCurrencyCodes := map[string]bool{
		"AFN": true, "EUR": true, "ALL": true, "DZD": true, "USD": true,
		"AOA": true, "XCD": true, "ARS": true, "AMD": true, "AWG": true,
		"AUD": true, "AZN": true, "BSD": true, "BHD": true, "BDT": true,
		"BBD": true, "BYN": true, "BZD": true, "XOF": true, "BMD": true,
		"INR": true, "BTN": true, "BOB": true, "BOV": true, "BAM": true,
		"BWP": true, "NOK": true, "BRL": true, "BND": true, "BGN": true,
		"BIF": true, "CVE": true, "KHR": true, "XAF": true, "CAD": true,
		"KYD": true, "CLP": true, "CLF": true, "CNY": true, "COP": true,
		"COU": true, "KMF": true, "CDF": true, "NZD": true, "CRC": true,
		"HRK": true, "CUP": true, "CUC": true, "ANG": true, "CZK": true,
		"DKK": true, "DJF": true, "DOP": true, "EGP": true, "SVC": true,
		"ERN": true, "SZL": true, "ETB": true, "FKP": true, "FJD": true,
		"XPF": true, "GMD": true, "GEL": true, "GHS": true, "GIP": true,
		"GTQ": true, "GBP": true, "GNF": true, "GYD": true, "HTG": true,
		"HNL": true, "HKD": true, "HUF": true, "ISK": true, "IDR": true,
		"XDR": true, "IRR": true, "IQD": true, "ILS": true, "JMD": true,
		"JPY": true, "JOD": true, "KZT": true, "KES": true, "KPW": true,
		"KRW": true, "KWD": true, "KGS": true, "LAK": true, "LBP": true,
		"LSL": true, "ZAR": true, "LRD": true, "LYD": true, "CHF": true,
		"MOP": true, "MKD": true, "MGA": true, "MWK": true, "MYR": true,
		"MVR": true, "MRU": true, "MUR": true, "XUA": true, "MXN": true,
		"MXV": true, "MDL": true, "MNT": true, "MAD": true, "MZN": true,
		"MMK": true, "NAD": true, "NPR": true, "NIO": true, "NGN": true,
		"OMR": true, "PKR": true, "PAB": true, "PGK": true, "PYG": true,
		"PEN": true, "PHP": true, "PLN": true, "QAR": true, "RON": true,
		"RUB": true, "RWF": true, "SHP": true, "WST": true, "STN": true,
		"SAR": true, "RSD": true, "SCR": true, "SLL": true, "SGD": true,
		"XSU": true, "SBD": true, "SOS": true, "SSP": true, "LKR": true,
		"SDG": true, "SRD": true, "SEK": true, "CHE": true, "CHW": true,
		"SYP": true, "TWD": true, "TJS": true, "TZS": true, "THB": true,
		"TOP": true, "TTD": true, "TND": true, "TRY": true, "TMT": true,
		"UGX": true, "UAH": true, "AED": true, "USN": true, "UYU": true,
		"UYI": true, "UYW": true, "UZS": true, "VUV": true, "VES": true,
		"VND": true, "YER": true, "ZMW": true, "ZWL": true, "XBA": true,
		"XBB": true, "XBC": true, "XBD": true, "XTS": true, "XXX": true,
		"XAU": true, "XPD": true, "XPT": true, "XAG": true,
	}
	return validCurrencyCodes[currencyCode]
}
