package usecase

import (
	"context"
	"github.com/lazy-bees/borsch/auth/models"
)

type UseCase interface {
	SignUp(ctx context.Context, userName, userPwd string) error
	SignIn(ctx context.Context, userName, userPwd string) (string, error)
	GetUser(ctx context.Context, accessToken string) (*models.User, error)
}
