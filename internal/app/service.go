package app

import (
	"context"
	"time"

	"github.com/walkccc/go-clean-arch/internal/controller"
	"github.com/walkccc/go-clean-arch/internal/repository"
	"github.com/walkccc/go-clean-arch/internal/usecase"
	pb "github.com/walkccc/go-clean-arch/pkg"
)

type MicroserviceServer struct {
	pb.UnimplementedMicroserviceServer
	createUserController controller.CreateUserController
}

func NewMicroserviceServer(store repository.Store) *MicroserviceServer {
	timeout := time.Second * 2
	createUserUsecase := usecase.NewCreateUserUsecase(store, timeout)
	createUserController := controller.NewCreateUserController(createUserUsecase)
	return &MicroserviceServer{
		createUserController: createUserController,
	}
}

func (server *MicroserviceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return server.createUserController.CreateUser(ctx, req)
}
