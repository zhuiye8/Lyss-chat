package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	Id       int64
	Username string
	Password string
	Email    string
}

type UserRepo interface {
	CreateUser(context.Context, *User) (*User, error)
	GetUserByUsername(context.Context, string) (*User, error)
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUsecase) Register(ctx context.Context, u *User) (*User, error) {
	// In a real application, you would hash the password here
	return uc.repo.CreateUser(ctx, u)
}

func (uc *UserUsecase) Login(ctx context.Context, username, password string) (string, error) {
	// In a real application, you would validate the password and generate a JWT
	return "fake-jwt-token", nil
}
