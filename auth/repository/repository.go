package repository

import (
	"context"
	"github.com/lazy-bees/borsch/auth/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userName, userPwd string) error
	GetUser(ctx context.Context, userName, userPwd string) (*models.User, error)
}
