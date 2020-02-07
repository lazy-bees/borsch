package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lazy-bees/borsch/auth/usecase"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc usecase.UseCase) {
	h := NewHandler(uc)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up", h.SignUp)
		authEndpoints.POST("/sign-in", h.SignIn)
		authEndpoints.GET("/user", h.GetUser)
	}
}
