package controllers

import (
	"fmt"
	"log"
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models/user"
	"social-media-app/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func handleRegisterRequest(ctx *gin.Context) (user.UserRegisterRequest, error) {
	request, ok := ctx.MustGet("request").(user.UserRegisterRequest)
	if !ok {
			return user.UserRegisterRequest{}, fmt.Errorf("failed to cast request to UserRegisterRequest")
	}
	return request, nil
}

func handleLoginRequest(ctx *gin.Context) (user.UserLoginRequest, error) {
	request, ok := ctx.MustGet("request").(user.UserLoginRequest)
	if !ok {
			return user.UserLoginRequest{}, fmt.Errorf("failed to cast request to UserLoginRequest")
	}
	return request, nil
}

func UserRegister(ctx *gin.Context) {
	DB, err := helpers.HandleDBContext(ctx)
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userRequest, err := handleRegisterRequest(ctx)
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	hashedPassword, err := helpers.HashPassword(userRequest.Password)
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userRequest.Password = hashedPassword

	createdUser, createError := services.CreateUser(DB, userRequest)
	if createError != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(createError.Error(), "email") {
			statusCode = http.StatusConflict
		}
		helpers.HandleErrorResponse(ctx, statusCode, createError.Error())
		return
	}

	log.Println("kesini ", createdUser.ID)
	token, err := helpers.GenerateToken(createdUser.ID)
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to generate token: %s", err))
		return
	}

	createdUser.AccessToken = token

	responseData := gin.H{
		"message": "User registered successfully",
		"data":    createdUser,
	}

	ctx.JSON(http.StatusCreated, responseData)
}

func UserLogin(ctx *gin.Context) {
	DB, err := helpers.HandleDBContext(ctx)
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userRequest, err := handleLoginRequest(ctx)
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	getUser, getError := services.GetUserById(DB, userRequest.Email)
	if getError != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(getError.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		helpers.HandleErrorResponse(ctx, statusCode, getError.Error())
		return
	}

	comparePass := helpers.ComparePassword([]byte(getUser.Password), []byte(userRequest.Password))
	if !comparePass {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid password")
		return
	}

	token, err := helpers.GenerateToken(getUser.ID)
	if err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to generate token: %s", err))
		return
	}

	getUser.AccessToken = token

	responseData := gin.H{
		"message": "User successfully logged",
		"data":    getUser,
	}

	ctx.JSON(http.StatusOK, responseData)
}
