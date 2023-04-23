package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/walkccc/go-clean-arch/internal/app"
	"github.com/walkccc/go-clean-arch/internal/repository"
	"github.com/walkccc/go-clean-arch/internal/util"
	pb "github.com/walkccc/go-clean-arch/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	go runGrpcServer(config, store)
	runGrpcGatewayServer(config, store)
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

func runGrpcGatewayServer(config util.Config, store repository.Store) {
	microserviceServer := app.NewMicroserviceServer(store)

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterMicroserviceHandlerServer(ctx, grpcMux, microserviceServer)
	if err != nil {
		log.Fatal("Cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannot create listener:", err)
	}

	log.Printf("Start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Cannot start HTTP gateway server:", err)
	}
}
