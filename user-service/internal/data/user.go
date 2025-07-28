package data

import (
	"context"

	"Lyss/user-service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) CreateUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	// In a real application, you would insert the user into the database here.
	// This is a mock implementation.
	r.log.WithContext(ctx).Infof("CreateUser: %v", u.Username)
	return u, nil
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*biz.User, error) {
	// In a real application, you would query the user from the database here.
	// This is a mock implementation.
	r.log.WithContext(ctx).Infof("GetUserByUsername: %v", username)
	if username == "test" {
		return &biz.User{Id: 1, Username: "test", Password: "password", Email: "test@example.com"}, nil
	}
	return nil, nil // user not found
}
