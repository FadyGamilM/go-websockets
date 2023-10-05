package handlers

import (
	"net/http"

	"github.com/FadyGamilM/go-websockets/internal/core"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService core.UserService
}

type UserHandlerConfig struct {
	R           *gin.Engine
	UserService core.UserService
}

func NewUserHandler(uhc *UserHandlerConfig) *userHandler {
	handler := &userHandler{
		userService: uhc.UserService,
	}
	userRoutes := uhc.R.Group("/api/users")
	userRoutes.POST("/login", handler.HandleLogin)
	userRoutes.POST("/signup", handler.HandleSignup)
	userRoutes.POST("/logout", handler.HandleLogout)
	return handler
}

func (uh *userHandler) HandleLogin(c *gin.Context) {
	loginDto := new(core.LoginUserDto)
	if err := c.ShouldBindJSON(loginDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
	}

	response, err := uh.userService.Login(c, loginDto)
	if err != nil {

		switch err.Type {
		case core.Error_Internal_Logic:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		case core.Error_WRONG_AUTH_CREDENTIAL:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		case core.Error_Non_Existing_Resource:
			c.JSON(http.StatusNotFound, gin.H{
				"error": err,
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	}

	c.SetCookie("token", response.AccessToken, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"data": &core.UserResponseDto{
			ID:       response.ID,
			Email:    response.Email,
			Username: response.Username,
		},
	})
}

func (uh *userHandler) HandleSignup(c *gin.Context) {
	dto := new(core.SignupUserDto)
	if err := c.ShouldBindJSON(dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
	}

	response, err := uh.userService.Signup(c, dto)
	if err != nil {

		switch err.Type {
		case core.Error_Internal_Logic:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		case core.Error_WRONG_AUTH_CREDENTIAL:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		case core.Error_Non_Existing_Resource:
			c.JSON(http.StatusNotFound, gin.H{
				"error": err,
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": response,
	})

}

func (uh *userHandler) HandleLogout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"data": "you are not authenticated",
	})
}
