package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lazy-bees/borsch/auth"
	"github.com/lazy-bees/borsch/auth/usecase"
	"net/http"
)

type Handler struct {
	uc usecase.UseCase
}

func NewHandler(uc usecase.UseCase) *Handler {
	return &Handler{uc: uc}
}

type signInput struct {
	UserName string `json:"user_name"`
	UserPwd  string `json:"user_pwd"`
}

func (h *Handler) SignUp(c *gin.Context) {
	inp := new(signInput)

	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.uc.SignUp(c.Request.Context(), inp.UserName, inp.UserPwd); err != nil {
		if err == auth.ErrUserAlreadyExists {
			c.AbortWithStatus(http.StatusConflict)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type signInResponse struct {
	Token string `json:"token"`
}

func (h *Handler) SignIn(c *gin.Context) {
	inp := new(signInput)

	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.uc.SignIn(c.Request.Context(), inp.UserName, inp.UserPwd)
	if err != nil {
		if err == auth.ErrUserNotFound {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, signInResponse{Token: token})
}

func (h *Handler) GetUser(c *gin.Context) {
	user, err := h.uc.GetUser(c.Request.Context(), c.Query("jwt"))

	if err != nil && err != auth.ErrUserNotFound {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}
