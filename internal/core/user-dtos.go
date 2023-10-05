package core

type LoginUserDto struct {
	Email    string `json:"email"`
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
