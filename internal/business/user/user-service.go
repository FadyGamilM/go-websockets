package user

import (
	"context"

	"github.com/FadyGamilM/go-websockets/config"
	authService "github.com/FadyGamilM/go-websockets/internal/business/auth"
	"github.com/FadyGamilM/go-websockets/internal/core"
	"github.com/FadyGamilM/go-websockets/internal/models"
)

type userService struct {
	userRepo  core.UserRepository
	tokenAuth core.TokenMaker
}

type UserServiceConfig struct {
	UserRepo  core.UserRepository
	TokenAuth core.TokenMaker
}

func NewUserService(usc *UserServiceConfig) core.UserService {
	return &userService{
		userRepo:  usc.UserRepo,
		tokenAuth: usc.TokenAuth,
	}
}

func (us *userService) Signup(ctx context.Context, dto *core.SignupUserDto) (*core.UserResponseDto, *core.BusinessError) {
	// hash the password
	hashedPass, err := authService.HashPassword(dto.Password)
	if err != nil {
		businessErr := core.New_Business_InternalLogicError("")
		return nil, businessErr
	}

	user, dbErr := us.userRepo.Create(ctx, &models.User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: hashedPass,
	})
	if err != nil {
		businessErr := core.New_Business_InternalLogicError(dbErr.Msg)
		return nil, businessErr
	}

	return &core.UserResponseDto{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (us *userService) Login(ctx context.Context, dto *core.LoginUserDto) (*core.LoginUserResponseDto, *core.BusinessError) {
	// check if there is a registered user in our system with this username
	registeredUser, dbErr := us.userRepo.GetByUsername(ctx, dto.Username)
	if dbErr != nil {
		if dbErr.Type == core.ERROR_FETCHING_USER {
			businessErr := core.New_Business_InternalLogicError(dbErr.Msg)
			return nil, businessErr
		} else if dbErr.Type == core.ERROR_NOT_FOUND_USER {
			businessErr := core.New_Business_NonExistingResourceError(dbErr.Msg)
			return nil, businessErr
		}
	}

	// check if the given password is the same as the stored hashed password
	IsCorrect := authService.CheckPassword(dto.Password, registeredUser.Password)
	if !IsCorrect {
		businessErr := core.New_Business_WrongAuthCredentialsError(dbErr.Msg)
		return nil, businessErr
	}

	// load the secret key from the config.yaml
	pasetoConfigs, err := config.LoadPasetoTokenConfig("./config")
	if err != nil {
		businessErr := core.New_Business_InternalLogicError("")
		return nil, businessErr
	}

	Token, _, err := us.tokenAuth.Create(dto.Username, pasetoConfigs.Paseto.Access_token_expiration)
	if err != nil {
		businessErr := core.New_Business_InternalLogicError("")
		return nil, businessErr
	}

	return &core.LoginUserResponseDto{
		AccessToken: Token,
		ID:          registeredUser.ID,
		Username:    registeredUser.Username,
		Email:       registeredUser.Email,
	}, nil

}
