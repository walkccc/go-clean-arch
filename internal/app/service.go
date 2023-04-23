package app

import (
	"context"
	"log"
	"time"

	"github.com/walkccc/go-clean-arch/internal/controller"
	"github.com/walkccc/go-clean-arch/internal/repository"
	"github.com/walkccc/go-clean-arch/internal/usecase"
	"github.com/walkccc/go-clean-arch/internal/util"
	"github.com/walkccc/go-clean-arch/internal/util/token"
	pb "github.com/walkccc/go-clean-arch/pkg"
)

type MicroserviceServer struct {
	pb.UnimplementedMicroserviceServer
	createUserController controller.CreateUserController
	loginUserController  controller.LoginUserController
}

func NewMicroserviceServer(config util.Config, store repository.Store) *MicroserviceServer {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("Cannot create token maker:", err)
	}
	timeout := time.Second * 2
	createUserUsecase := usecase.NewCreateUserUsecase(store, timeout)
	createUserController := controller.NewCreateUserController(createUserUsecase)
	loginUserUsecase := usecase.NewLoginUserUsecase(store, tokenMaker, config.AccessTokenDuration, config.RefreshTokenDuration, timeout)
	loginUserController := controller.NewLoginUserController(loginUserUsecase)
	return &MicroserviceServer{
		createUserController: createUserController,
		loginUserController:  loginUserController,
	}
}

func (server *MicroserviceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return server.createUserController.CreateUser(ctx, req)
}

func (server *MicroserviceServer) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	return server.loginUserController.LoginUser(ctx, req)
}
