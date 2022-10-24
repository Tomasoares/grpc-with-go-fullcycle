package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/tomasoares/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "João",
		Email: "j@j.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println("Added user: " + res.String())
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "João",
		Email: "j@j.com",
	}

	stream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive message: %v", err)
		}

		fmt.Println("Status:", res.Status)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "",
			Name:  "Tomas",
			Email: "tomas@gmail.com",
		},
		&pb.User{
			Id:    "",
			Name:  "Tomas2",
			Email: "tomas2@gmail.com",
		},
		&pb.User{
			Id:    "",
			Name:  "Tomas3",
			Email: "tomas3@gmail.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)

		fmt.Println("Sent", req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "",
			Name:  "Tomas",
			Email: "tomas@gmail.com",
		},
		&pb.User{
			Id:    "",
			Name:  "Tomas2",
			Email: "tomas2@gmail.com",
		},
		&pb.User{
			Id:    "",
			Name:  "Tomas3",
			Email: "tomas3@gmail.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error")
				break
			}
			fmt.Println("Received", res)
		}

		close(wait)
	}()

	<-wait
}
