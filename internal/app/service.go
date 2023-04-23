package app

import (
	"github.com/walkccc/go-clean-arch/internal/repository"
	pb "github.com/walkccc/go-clean-arch/pkg"
)

type MicroserviceServer struct {
	pb.UnimplementedMicroserviceServer
}

func NewMicroserviceServer(store repository.Store) *MicroserviceServer {
	return &MicroserviceServer{}
}
