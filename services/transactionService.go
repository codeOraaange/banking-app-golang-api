package services

import (
	"context"
	"errors"
	"fmt"
	"social-media-app/models/transaction"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
    ErrInsufficientBalance = errors.New("insufficient balance")
)

func PostTransactionService(DB *pgxpool.Pool, userID int, transactionRequest transaction.TransactionRequest) (*transaction.TransactionResponse, error) {
    ctx := context.Background()
    tx, err := DB.Begin(ctx)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback(ctx)

    var transactionResponse transaction.TransactionResponse

    // Get the sender's bank account ID and balance
    var senderBankAccountID int
    var senderBalance int
    err = tx.QueryRow(ctx, `SELECT id, balance FROM bank_accounts WHERE user_id = $1`, userID).Scan(&senderBankAccountID, &senderBalance)
    if err != nil {
        return nil, err
    }

    // Check if the sender has enough balance
    if senderBalance < transactionRequest.Balance {
        return nil, ErrInsufficientBalance
    }

    // Get the recipient's bank account ID
    var recipientBankAccountID int
    err = tx.QueryRow(ctx, `SELECT id FROM bank_accounts WHERE account_number = $1 AND bank_name = $2`, transactionRequest.RecipientBankAccountNumber, transactionRequest.RecipientBankName).Scan(&recipientBankAccountID)
    if err != nil {
        return nil, err
    }

    // Insert the transaction record
    err = tx.QueryRow(ctx, `
        INSERT INTO balance_outcome (sender_id, recipient_id, credited_amount, from_currency, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, sender_id, recipient_id, credited_amount, from_currency, created_at, updated_at;
    `, senderBankAccountID, recipientBankAccountID, transactionRequest.Balance, transactionRequest.Currency, time.Now(), time.Now()).Scan(
        &transactionResponse.ID,
        &transactionResponse.SenderId,
        &transactionResponse.RecipientId,
        &transactionResponse.CreditedAmount,
        &transactionResponse.FromCurrency,
        &transactionResponse.CreatedAt,
        &transactionResponse.CreatedAt,
    )
    if err != nil {
        fmt.Println("Error inserting transaction record:", err)
        return nil, err
    }

    // Log the SQL queries
    senderQuery := fmt.Sprintf(`UPDATE bank_accounts SET balance = balance - %d WHERE id = %d`, transactionRequest.Balance, senderBankAccountID)
    recipientQuery := fmt.Sprintf(`UPDATE bank_accounts SET balance = balance + %d WHERE id = %d`, transactionRequest.Balance, recipientBankAccountID)
    insertQuery := `INSERT INTO balance_outcome (sender_id, recipient_id, credited_amount, from_currency, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ...`

    fmt.Println("SQL Query:", senderQuery)
    fmt.Println("SQL Query:", recipientQuery)
    fmt.Println("SQL Query:", insertQuery)

    // Update the sender's balance
    _, err = tx.Exec(ctx, `UPDATE bank_accounts SET balance = balance - $1 WHERE id = $2`, transactionRequest.Balance, senderBankAccountID)
    if err != nil {
        return nil, err
    }

    // Update the recipient's balance
    _, err = tx.Exec(ctx, `UPDATE bank_accounts SET balance = balance + $1 WHERE id = $2`, transactionRequest.Balance, recipientBankAccountID)
    if err != nil {
        return nil, err
    }

	if err := tx.Commit(ctx); err != nil {
        // Log the error
        fmt.Println("Error committing transaction:", err)
        return nil, err
    }

	fmt.Println("Transaction committed successfully")

    return &transactionResponse, nil
}