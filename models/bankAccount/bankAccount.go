package bankAccount

type BankAccountRequest struct {
	SenderBankAccountNumber string  `json:"senderBankAccountNumber" binding:"required,min=5,max=30" validate:"required,min=5,max=30"`
	SenderBankName          string  `json:"senderBankName" binding:"required,min=5,max=30" validate:"required,min=5,max=30"`
	AddedBalance            float64 `json:"addedBalance" binding:"required,min=0" validate:"required,min=0"`
	Currency                string  `json:"currency" binding:"required,len=3" validate:"required,len=3"`
	TransferProofImg        string  `json:"transferProofImg" binding:"required,url"`
}

type BankAccountResponse struct {
	ID                int    `json:"id"`
	UserID            int    `json:"userId"`
	BankAccountNumber string `json:"bankAccountNumber"`
	BankName          string `json:"bankName"`
	Balance           int    `json:"balance"`
	Currency          string `json:"currency"`
}
