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
