package services

import (
	"context"
	"fmt"

	"github.com/tomasoares/fc2-grpc/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {

	//Insert - Database
	fmt.Println(req.Name)
	id := "123"

	return &pb.User{
		Id:    id,
		Name:  req.Name,
		Email: req.Email}, nil
}
