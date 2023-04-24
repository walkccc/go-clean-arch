package controller

import (
	"context"

	"github.com/walkccc/go-clean-arch/internal/usecase"
	pb "github.com/walkccc/go-clean-arch/pkg"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateUserController interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
}

type createUserController struct {
	createUserUsecase usecase.CreateUserUsecase
}

func NewCreateUserController(createUserUsercase usecase.CreateUserUsecase) CreateUserController {
	return &createUserController{createUserUsecase: createUserUsercase}
}

func (c *createUserController) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := c.createUserUsecase.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		User: &pb.User{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt)}}, nil
}
