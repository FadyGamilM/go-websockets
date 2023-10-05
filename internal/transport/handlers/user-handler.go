package handlers

import (
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

func NewUserHandler(uhc UserHandlerConfig) *userHandler {
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
	
}

func (uh *userHandler) HandleSignup(c *gin.Context) {

}

func (uh *userHandler) HandleLogout(c *gin.Context) {

}
