package main

import (
	"fmt"
	"log"
	"net"

	"github.com/tomasoares/fc2-grpc/pb"
	"github.com/tomasoares/fc2-grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("Initiating server...")

	lis, err := net.Listen("tcp", "localhost:50021")
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	fmt.Println("Listening!")

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, services.NewUserService())
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("could not serve: %v", err)
	}

}
