package services

import (
	"context"
	"log"
	"banking-app-golang-api/models"
	"banking-app-golang-api/models/post"
	"time"
	"fmt"
    "net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gin-gonic/gin"
)

func GetUserByIdWithFriendCount(DB *pgxpool.Pool, userID int) (*post.PostCreatorResponse, error) {
	ctx := context.Background()

	tx, err := DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Println("Failed to rollback transaction:", rollbackErr)
			}
			panic(p)
		} else if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Println("Failed to rollback transaction:", rollbackErr)
			}
			log.Println("Transaction aborted:", err)
		}
	}()

	var user post.PostCreatorResponse
	err = tx.QueryRow(ctx, `
			SELECT u.id, u.name, u.image_url, COUNT(f.friend_id) AS friend_count, u.created_at 
			FROM users u
			LEFT JOIN friendship f ON u.id = f.user_id
			WHERE u.id = $1
			GROUP BY u.id, u.name, u.image_url, u.created_at
	`, userID).
			Scan(&user.UserId, &user.Name, &user.ImageUrl, &user.FriendCount, &user.CreatedAt)
	if err != nil {
			log.Println("Failed to query user:", err)
			return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return nil, err
	}

	return &user, nil
}

func GetUserByUserId(DB *pgxpool.Pool, userID int) (*models.Users, error) {
	ctx := context.Background()

	tx, err := DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Println("Failed to rollback transaction:", rollbackErr)
			}
			panic(p)
		} else if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Println("Failed to rollback transaction:", rollbackErr)
			}
			log.Println("Transaction aborted:", err)
		}
	}()

	var user models.Users
	var nullableImageUrl *string
	var nullableEmail *string
	var nullablePhone *string
	err = tx.QueryRow(ctx, `
			SELECT id, name, password, email, phone, image_url, credential_type, created_at, updated_at
			FROM users
			WHERE id = $1
	`, userID).
			Scan(&user.ID, &user.Name, &user.Password, &nullableEmail, &nullablePhone, &nullableImageUrl, &user.CredentialType, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
			log.Println("Failed to query user by ID:", err)
			return nil, err
	}

	if nullableImageUrl != nil {
		user.ImageURL = *nullableImageUrl
	}

	if nullableEmail != nil {
		user.Email = *nullableEmail
	}

	if nullablePhone != nil {
		user.Phone = *nullablePhone
	}

	if err = tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return nil, err
	}

	return &user, nil
}
