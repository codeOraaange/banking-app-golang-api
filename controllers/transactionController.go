package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models/transaction"
	"social-media-app/services"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

func handlePostTransactionRequest(ctx *gin.Context) (transaction.TransactionRequest, error) {
	request, ok := ctx.MustGet("request").(transaction.TransactionRequest)
	if !ok {
		return transaction.TransactionRequest{}, fmt.Errorf("failed to cast request to BankAccountRequest")
	}
	return request, nil
}

func PostTransaction(ctx *gin.Context) {
	DB, err := helpers.HandleDBContext(ctx)
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	transactionRequest, err := handlePostTransactionRequest(ctx)
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userData := ctx.MustGet("userData").(jwt5.MapClaims)
	userID := int(userData["id"].(float64))

	createTransaction, err := services.PostTransactionService(DB, userID, transactionRequest)
	var ErrExpiredToken = errors.New("token has expired")
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInsufficientBalance):
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, "Balance is not enough")
			return
		case errors.Is(err, ErrExpiredToken):
			helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, "Token is missing or expired")
			return
		default:
			helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}

	responseTransaction := gin.H{
		"message": "successfully send balance",
		"data":    createTransaction,
	}
	ctx.JSON(http.StatusOK, responseTransaction)
}