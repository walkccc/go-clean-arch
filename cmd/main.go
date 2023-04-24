package main

import (
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/walkccc/go-clean-arch/internal/app"
	"github.com/walkccc/go-clean-arch/internal/repository"
	"github.com/walkccc/go-clean-arch/internal/util"
	pb "github.com/walkccc/go-clean-arch/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	db := repository.NewDB(config.DBDriver, config.DBSource)
	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot ping DB: %v", err)
	}

	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	store := repository.NewStore(db)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store repository.Store) {
	microserviceServer := app.NewMicroserviceServer(store)
	grpcServer := grpc.NewServer()
	pb.RegisterMicroserviceServer(grpcServer, microserviceServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Cannot create listener for gRPC:", err)
	}

	log.Printf("Start gRPC server at %s,", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannot start gRPC server:", err)
	}
}
