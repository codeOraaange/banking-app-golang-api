package services

import (
	"context"
	// "errors"
	"log"
	"social-media-app/models"
	"time"

	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PostBalanceService(DB *pgxpool.Pool, balanceRequest models.BankAccountRequest, userID int) (*models.BankAccountResponse, error) {
	ctx := context.Background()
    
    // Begin a transaction
    tx, err := DB.Begin(ctx)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback(ctx) // Rollback transaction if not committed

    var bankAccountResponse models.BankAccountResponse
    err = tx.QueryRow(ctx, `
        INSERT INTO bank_accounts (user_id, account_number, bank_name, balance, currency, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, user_id, account_number, bank_name, balance, currency
    `, userID, balanceRequest.SenderBankAccountNumber, balanceRequest.SenderBankName, balanceRequest.AddedBalance, balanceRequest.Currency, time.Now(), time.Now()).
        Scan(&bankAccountResponse.ID, &bankAccountResponse.UserID, &bankAccountResponse.BankAccountNumber, &bankAccountResponse.BankName, &bankAccountResponse.Balance, &bankAccountResponse.Currency)
    if err != nil {
        tx.Rollback(ctx)
        log.Println("Failed to insert balance:", err)
        return nil, err
    }

    // Insert data into balance_income table
	
    _, err = tx.Exec(ctx, `
        INSERT INTO balance_income (bank_id, transfer_proof_img, deposited_amount, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `, bankAccountResponse.ID, balanceRequest.TransferProofImg, balanceRequest.AddedBalance, time.Now(), time.Now())
    if err != nil {
        tx.Rollback(ctx)
        log.Println("Failed to insert balance income:", err)
        return nil, err
    }

    // Commit transaction
    err = tx.Commit(ctx)
    if err != nil {
        log.Println("Failed to commit transaction:", err)
        return nil, err
    }

    return &bankAccountResponse, nil
}
