package core

type LoginUserDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignupUserDto struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponseDto struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUserResponseDto struct {
	AccessToken string `json:"access_token"`
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
}
