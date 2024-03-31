package user

type User struct {
	Name     string `json:"name" binding:"required,min=5,max=30" validate:"required,min=5,max=30"`
	Password string `json:"password" binding:"required,min=5,max=30" validate:"required,min=5,max=30"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
}

type UserResponse struct {
	ID          int    `json:"-"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	AccessToken string `json:"accessToken"`
}
