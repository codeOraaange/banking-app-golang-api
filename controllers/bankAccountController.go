package controllers

import (
	"fmt"
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models"
	"social-media-app/services"
	// "strings"

	"github.com/gin-gonic/gin"
	// jwt5 "github.com/golang-jwt/jwt/v5"
)

func handlePostBalanceRequest(ctx *gin.Context) (models.BankAccountRequest, error) {
	request, ok := ctx.MustGet("request").(models.BankAccountRequest)
	if !ok {
		return models.BankAccountRequest{}, fmt.Errorf("failed to cast request to BankAccountRequest")
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
		helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// userData := ctx.MustGet("userData").(jwt5.MapClaims)
	// userID := int(userData["id"].(float64))
	userID := 1
	fmt.Println("line 57 controller")
	createBalance, err := services.PostBalanceService(DB, balanceRequest, userID)
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("line 63 controller")

	responseBalance := gin.H{
		"message": "Balance added successfully",
		"data":    createBalance,
	}
	
	ctx.JSON(http.StatusCreated, responseBalance)
}
