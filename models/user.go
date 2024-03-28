package user

import (
	"social-media-app/helpers"
)

type UserRegisterRequest struct {
	Name     string `json:"name" binding:"required,min=5,max=30" validate:"required,min=5,max=30"`
	Password string `json:"password" binding:"required,min=5,max=30" validate:"required,min=5,max=30"`
	Email    string `json:"email" binding:"required,email"`
}

type UserLoginRequest struct {
	Password string `json:"password" binding:"required,min=5,max=30" validate:"required,min=5,max=30"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	ID          int    `json:"-"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	AccessToken string `json:"accessToken"`
}

func (user *UserRegisterRequest) HashPassword() error {
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return nil
}

func BeforeCreateUser(user *UserRegisterRequest) {
	user.HashPassword()
}
