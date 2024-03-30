package models

// import "time"

type BalanceHistoryResponse struct {
    TransactionID    int     `json:"transactionId"`
    Balance          int     `json:"balance"`
    Currency         string  `json:"currency"`
    TransferProofImg string  `json:"transferProofImg"`
    Source           Source  `json:"source"`
    CreatedAt        int64  `json:"createdAt"`
    // CreatedAtMillis  int64   `json:"createdAtMillis"` // Add CreatedAtMillis field
}

type Source struct {
    BankAccountNumber string `json:"bankAccountNumber"`
    BankName          string `json:"bankName"`
}

type BalanceHistoryData struct {
	Message string                   `json:"message"`
	Data    []BalanceHistoryResponse `json:"data"`
	Meta    BalanceHistoryMeta       `json:"meta"`
}

type BalanceHistoryMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
