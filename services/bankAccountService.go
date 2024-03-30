package services

import (
	"context"
	"database/sql"
	"sort"

	// "errors"
	"log"
	"social-media-app/models/bankAccount"

	// "social-media-app/models/transaction"
	"social-media-app/models"
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
        WHERE user_id = $1 AND account_number = $2 AND currency = $3
    `, userID, balanceRequest.SenderBankAccountNumber, balanceRequest.Currency).
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

func GetBalanceService(DB *pgxpool.Pool, userID int) ([]bankAccount.BankAccountResponse, error) {
	ctx := context.Background()
	tx, err := DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // Rollback transaction if not committed

	var balances []bankAccount.BankAccountResponse
	rows, err := tx.Query(ctx, `
        SELECT id, user_id, account_number, bank_name, balance, currency
        FROM bank_accounts
        WHERE user_id = $1
    `, userID)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var balance bankAccount.BankAccountResponse
		err = rows.Scan(&balance.ID, &balance.UserID, &balance.BankAccountNumber, &balance.BankName, &balance.Balance, &balance.Currency)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
		balances = append(balances, balance)
	}

	err = rows.Err()
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	// Sort balances by descending order
	sort.Slice(balances, func(i, j int) bool {
		if balances[i].Currency == balances[j].Currency {
			return balances[i].Balance > balances[j].Balance
		}
		return balances[i].Currency < balances[j].Currency
	})

	return balances, nil
}

func GetBalanceByHistoryService(DB *pgxpool.Pool, userID int, limit, offset int) (models.BalanceHistoryData, error) {
    ctx := context.Background()
    tx, err := DB.Begin(ctx)
    if err != nil {
        return models.BalanceHistoryData{}, err
    }
    defer tx.Rollback(ctx)

    var balanceHistory []models.BalanceHistoryResponse
    rows, err := tx.Query(ctx, `SELECT bi.id AS transactionId, ba.balance, ba.currency, bi.transfer_proof_img AS transferProofImg, ba.account_number AS bankAccountNumber, ba.bank_name AS bankName, bi.created_at AS createdAt FROM balance_income bi JOIN bank_accounts ba ON bi.bank_id = ba.id WHERE ba.user_id = $1 ORDER BY bi.created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset)
    if err != nil {
        tx.Rollback(ctx)
        return models.BalanceHistoryData{}, err
    }
    defer rows.Close()

    for rows.Next() {
        var history models.BalanceHistoryResponse
        var source models.Source
        var createdAt time.Time
        err = rows.Scan(&history.TransactionID, &history.Balance, &history.Currency, &history.TransferProofImg, &source.BankAccountNumber, &source.BankName, &createdAt)
        if err != nil {
            tx.Rollback(ctx)
            return models.BalanceHistoryData{}, err
        }

        history.CreatedAt = createdAt.UnixNano() / int64(time.Millisecond)
        // history.CreatedAtMillis = createdAt.UnixNano() / int64(time.Millisecond)
        history.Source = source
        balanceHistory = append(balanceHistory, history)
    }

    err = rows.Err()
    if err != nil {
        tx.Rollback(ctx)
        return models.BalanceHistoryData{}, err
    }

    err = tx.Commit(ctx)
    if err != nil {
        return models.BalanceHistoryData{}, err
    }

    total, err := getTotalBalanceHistoryCount(DB, userID)
    if err != nil {
        return models.BalanceHistoryData{}, err
    }

    return models.BalanceHistoryData{
        Message: "success",
        Data:    balanceHistory,
        Meta: models.BalanceHistoryMeta{
            Limit:  limit,
            Offset: offset,
            Total:  total,
        },
    }, nil
}

func getTotalBalanceHistoryCount(DB *pgxpool.Pool, userID int) (int, error) {
	ctx := context.Background()
	var total int
	err := DB.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM balance_income bi
		JOIN bank_accounts ba ON bi.bank_id = ba.id
		WHERE ba.user_id = $1
	`, userID).Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return total, nil
}
