package controller

import (
	"context"

	"github.com/walkccc/go-clean-arch/internal/usecase"
	pb "github.com/walkccc/go-clean-arch/pkg"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LoginUserController interface {
	LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error)
}

type loginUserController struct {
	loginUserUsecase usecase.LoginUserUsecase
}

func NewLoginUserController(loginUserUsecase usecase.LoginUserUsecase) LoginUserController {
	return &loginUserController{loginUserUsecase: loginUserUsecase}
}

func (c *loginUserController) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, loginMetadata, err := c.loginUserUsecase.LoginUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.LoginUserResponse{
		User: &pb.User{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt)},
		SessionId:             loginMetadata.SessionId,
		AccessToken:           loginMetadata.AccessToken,
		RefreshToken:          loginMetadata.RefreshToken,
		AccessTokenExpiresAt:  timestamppb.New(loginMetadata.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(loginMetadata.RefreshTokenExpiresAt)}, nil
}
