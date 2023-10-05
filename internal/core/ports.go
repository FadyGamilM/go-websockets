package core

import (
	"context"

	"github.com/FadyGamilM/go-websockets/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type UserService interface {
	Signup(ctx context.Context, dto *SignupUserDto) (*models.User, error)
	Login(ctx context.Context, dto *LoginUserDto) (*models.User, error)
}
