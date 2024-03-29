package services

import (
	"context"
	// "errors"
	"log"
	"social-media-app/models/bankAccount"
	"time"

	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PostBalanceService(DB *pgxpool.Pool, balanceRequest bankAccount.BankAccountRequest, userID int) (*bankAccount.BankAccountResponse, error) {
	ctx := context.Background()
	tx, err := DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // Rollback transaction if not committed

	var bankAccountResponse bankAccount.BankAccountResponse

	// Check if the account number exists
	err = tx.QueryRow(ctx, `
        SELECT id, user_id, account_number, bank_name, balance, currency
        FROM bank_accounts
        WHERE user_id = $1 AND account_number = $2
    `, userID, balanceRequest.SenderBankAccountNumber).
		Scan(&bankAccountResponse.ID, &bankAccountResponse.UserID, &bankAccountResponse.BankAccountNumber, &bankAccountResponse.BankName, &bankAccountResponse.Balance, &bankAccountResponse.Currency)

	if err != nil {
		if err != pgx.ErrNoRows {
			tx.Rollback(ctx)
			log.Println("Failed to check account number:", err)
			return nil, err
		}

		// Account number doesn't exist, create a new one
		err = tx.QueryRow(ctx, `
            INSERT INTO bank_accounts (user_id, account_number, bank_name, balance, currency, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING id, user_id, account_number, bank_name, balance, currency
        `, userID, balanceRequest.SenderBankAccountNumber, balanceRequest.SenderBankName, balanceRequest.AddedBalance, balanceRequest.Currency, time.Now(), time.Now()).
			Scan(&bankAccountResponse.ID, &bankAccountResponse.UserID, &bankAccountResponse.BankAccountNumber, &bankAccountResponse.BankName, &bankAccountResponse.Balance, &bankAccountResponse.Currency)
		if err != nil {
			tx.Rollback(ctx)
			log.Println("Failed to insert new balance:", err)
			return nil, err
		}
	} else {
		// Account number exists, update the balance
		_, err = tx.Exec(ctx, `
            UPDATE bank_accounts
            SET balance = balance + $1, updated_at = $2
            WHERE id = $3
        `, balanceRequest.AddedBalance, time.Now(), bankAccountResponse.ID)
		if err != nil {
			tx.Rollback(ctx)
			log.Println("Failed to update balance:", err)
			return nil, err
		}

		// Update the balance in the response
		bankAccountResponse.Balance += int(balanceRequest.AddedBalance)
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
