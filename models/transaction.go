package models

type TransactionRequest struct {
	RecipientBankAccountNumber string `json:"recipientBankAccountNumber" binding:"required" validate:""`
	RecipientBankName          string `json:"recipientBankName" binding:"required" validate:""`
	Currency                   string `json:"fromCurrency" binding:"required" validate:""`
	Balance                    int    `json:"balances" binding:"required" validate:""`
}

type TransactionResponse struct {
	ID               int    `json:"id"`
	TransferProofImg string `json:"transferProofImg"`
	SenderId         int    `json:"senderId"`
	RecipientId      int    `json:"recipientId"`
	Amount           int    `json:"amount"`
	Type             string `json:"type"`
}
