package core

type LoginUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupUserDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
