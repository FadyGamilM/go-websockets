package core

import (
	"context"
	"time"

	"github.com/FadyGamilM/go-websockets/internal/business/auth/token"
	"github.com/FadyGamilM/go-websockets/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, *DbError)
	GetByID(ctx context.Context, id int64) (*models.User, *DbError)
	GetByUsername(ctx context.Context, username string) (*models.User, *DbError)
}

type UserService interface {
	Signup(ctx context.Context, dto *SignupUserDto) (*UserResponseDto, *BusinessError)
	Login(ctx context.Context, dto *LoginUserDto) (*LoginUserResponseDto, *BusinessError)
}

type TokenMaker interface {
	/*
		@Params:
			username => to be included in the payload
			expiration => to set a short-life time for the token
	*/
	Create(username string, expiration time.Duration) (string, *token.Payload, error)
	Verify(token string) (*token.Payload, error)
}
