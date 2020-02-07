package authusecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/lazy-bees/borsch/auth"
	"github.com/lazy-bees/borsch/auth/models"
	"github.com/lazy-bees/borsch/auth/repository"
	"time"
)

type authClaims struct {
	jwt.StandardClaims
	User *models.User `json:"user"`
}

type AuthUseCase struct {
	userRepo       repository.UserRepository
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthUseCase(userRepo repository.UserRepository, hashSalt string, signingKey []byte, tokenTTLSeconds time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTLSeconds,
	}
}

func (uc *AuthUseCase) SignUp(ctx context.Context, userName, userPwd string) error {
	pwd := sha1.New()
	pwd.Write([]byte(userPwd))
	pwd.Write([]byte(uc.hashSalt))

	return uc.userRepo.CreateUser(ctx, userName, fmt.Sprintf("%x", pwd.Sum(nil)))
}

func (uc *AuthUseCase) SignIn(ctx context.Context, userName, userPwd string) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(userPwd))
	pwd.Write([]byte(uc.hashSalt))
	userPwd = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := uc.userRepo.GetUser(ctx, userName, userPwd)
	if err != nil {
		return "", auth.ErrUserNotFound
	}

	claims := authClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(uc.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(uc.signingKey)
}

func (uc *AuthUseCase) GetUser(ctx context.Context, accessToken string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return uc.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*authClaims); ok && token.Valid {
		if user, err := uc.userRepo.GetUser(ctx, claims.User.Name, claims.User.PWD); err != nil {
			return nil, err
		} else {
			return user, nil
		}
	}

	return nil, auth.ErrInvalidAccessToken
}
