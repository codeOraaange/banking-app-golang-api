package services

import (
	"context"
	"errors"
	"log"
	"social-media-app/models/user"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(DB *pgxpool.Pool, userRegister user.User) (*user.UserResponse, error) {
	ctx := context.Background()

	var user user.UserResponse
	err := DB.QueryRow(ctx, `
		WITH email_check AS (
			SELECT COUNT(1) AS count FROM users WHERE email = $1 LIMIT 1
		),
		insert_user AS (
			INSERT INTO users (name, password, email, created_at, updated_at)
			SELECT $2, $3, $4, $5, $6
			WHERE (SELECT count FROM email_check) = 0
			RETURNING id, name, email
		)
		SELECT id, name, email FROM insert_user
	`, userRegister.Email, userRegister.Name, userRegister.Password, userRegister.Email, time.Now(), time.Now()).
		Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		log.Println("Failed to insert user:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("email already exists")
		}
		return nil, err
	}

	return &user, nil
}

func GetUserById(DB *pgxpool.Pool, email string) (*user.UserResponse, error) {
	ctx := context.Background()

	var user user.UserResponse
	err := DB.QueryRow(ctx, `
		SELECT id, name, email, password
		FROM users
		WHERE email = $1
	`, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		log.Println("Failed to get user:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found by ID")
		}
		return nil, err
	}

	return &user, nil
}