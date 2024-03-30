package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"social-media-app/helpers"
	"strconv"

	// "time"

	// "social-media-app/models/user"
	"social-media-app/models/bankAccount"
	"social-media-app/services"
	// "social-media-app/helpers"

	// "social-media-app/models"

	// "strings"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

func handlePostBalanceRequest(ctx *gin.Context) (bankAccount.BankAccountRequest, error) {
	request, ok := ctx.MustGet("request").(bankAccount.BankAccountRequest)
	if !ok {
		return bankAccount.BankAccountRequest{}, fmt.Errorf("failed to cast request to BankAccountRequest")
	}
	return request, nil
}

func PostBalance(ctx *gin.Context) {
	DB, err := helpers.HandleDBContext(ctx)
    if err != nil {
        helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
        return
    }

    balanceRequest, err := handlePostBalanceRequest(ctx)
    if err != nil {
        helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err.Error())
        return
    }

    userData := ctx.MustGet("userData").(jwt5.MapClaims)
    userID := int(userData["id"].(float64))

    createBalance, err := services.PostBalanceService(DB, balanceRequest, userID)
    var ErrExpiredToken = errors.New("token has expired")
	if err != nil {
        if errors.Is(err, ErrExpiredToken) { // Assuming ErrExpiredToken is defined in your helpers package
            helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, "Token is missing or expired")
            return
        }
        helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
        return
    }

	responseBalance := gin.H{
        "message": "Balance added successfully",
        "data":    createBalance,
    }

    ctx.JSON(http.StatusOK, responseBalance)
}

func GetBalance(ctx *gin.Context) {
	DB, err := helpers.HandleDBContext(ctx)
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userData := ctx.MustGet("userData").(jwt5.MapClaims)
	userID := int(userData["id"].(float64))

	balances, err := services.GetBalanceService(DB, userID)
	var ErrExpiredToken = errors.New("token has expired")
	if err != nil {
        if errors.Is(err, ErrExpiredToken) { // Assuming ErrExpiredToken is defined in your helpers package
            helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, "Token is missing or expired")
            return
        }
        helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
        return
    }

	response := gin.H{
		"message": "success",
		"data":    balances,
	}

	ctx.JSON(http.StatusOK, response)
}

func GetBalanceByHistory(ctx *gin.Context) {
	DB, err := helpers.HandleDBContext(ctx)
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userData := ctx.MustGet("userData").(jwt5.MapClaims)
	userID := int(userData["id"].(float64))

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "5"))
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	balanceByHistory, err := services.GetBalanceByHistoryService(DB, userID, limit, offset)
	var ErrExpiredToken = errors.New("token has expired")
	if err != nil {
        if errors.Is(err, ErrExpiredToken) { // Assuming ErrExpiredToken is defined in your helpers package
            helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, "Token is missing or expired")
            return
        }
        helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
        return
    }

	ctx.JSON(http.StatusOK, balanceByHistory)
}
