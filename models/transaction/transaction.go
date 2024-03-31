package transaction

import "time"

type TransactionRequest struct {
	RecipientBankAccountNumber string `json:"recipientBankAccountNumber" binding:"required,min=5,max=30" validate:"required,min=5,max=30"`
	RecipientBankName          string `json:"recipientBankName" binding:"required,min=5,max=30" validate:"required,min=5,max=30"`
	Currency                   string `json:"fromCurrency" binding:"required,len=3" validate:"required,len=3,iso4217"`
	Balance                    int    `json:"balances" binding:"required,min=0" validate:"required,gte=0"`
}

type TransactionResponse struct {
	ID             int       `json:"id"`
	SenderId       int       `json:"senderId"`
	RecipientId    int       `json:"recipientId"`
	CreditedAmount int       `json:"creditedAmount"`
	FromCurrency   string    `json:"fromCurrency"`
	CreatedAt      time.Time `json:"createdAt"`
}
