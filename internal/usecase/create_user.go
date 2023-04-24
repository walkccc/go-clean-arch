package usecase

import (
	"context"
	"time"

	"github.com/walkccc/go-clean-arch/internal/repository"
	"github.com/walkccc/go-clean-arch/internal/util"
	pb "github.com/walkccc/go-clean-arch/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateUserUsecase interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*repository.User, error)
}

type createUserUsecase struct {
	store   repository.Store
	timeout time.Duration
}

func NewCreateUserUsecase(store repository.Store, timeout time.Duration) CreateUserUsecase {
	return &createUserUsecase{store: store, timeout: timeout}
}

func (u *createUserUsecase) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*repository.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password: %s", err)
	}
	user, err := u.store.CreateUser(ctx, repository.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	})
	return &user, err
}
