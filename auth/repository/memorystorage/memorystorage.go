package memorystorage

import (
	"context"
	"github.com/google/uuid"
	"github.com/lazy-bees/borsch/auth"
	"github.com/lazy-bees/borsch/auth/models"
	"sync"
)

type UserMemoryStorage struct {
	users map[string]*models.User
	mutex *sync.Mutex
}

func NewUserMemoryStorage() *UserMemoryStorage {
	return &UserMemoryStorage{
		users: make(map[string]*models.User),
		mutex: &sync.Mutex{},
	}
}

func (r *UserMemoryStorage) CreateUser(ctx context.Context, userName, userPwd string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.users[userName]; ok {
		return auth.ErrUserAlreadyExists
	}

	r.users[userName] = &models.User{
		ID:   uuid.New().String(),
		Name: userName,
		PWD:  userPwd,
	}

	return nil
}

func (r *UserMemoryStorage) GetUser(ctx context.Context, userName, userPwd string) (*models.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if u, ok := r.users[userName]; ok && u.PWD == userPwd {
		return u, nil
	} else {
		return nil, auth.ErrUserNotFound
	}
}
