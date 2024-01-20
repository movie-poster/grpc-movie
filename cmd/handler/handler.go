package handler

import (
	"context"
	"fmt"
	repository "template-grpc/internal/domain/repository/implement/user"
	irepository "template-grpc/internal/domain/repository/interface"

	pb "template-grpc/internal/infra/proto"
)

func NewServerUser(usercrud irepository.IUserCrud) *server {
	return &server{
		userCrud: repository.NewRepository(),
	}
}

type server struct {
	userCrud irepository.IUserCrud
	pb.UnimplementedUserCrudServer
}

func (s *server) Insert(context context.Context, user *pb.User) (*pb.Response, error) {

	fmt.Println(user)
	return &pb.Response{
		Id:      1,
		IsOk:    true,
		Message: "dasdasdsa",
	}, nil
}

func (s *server) Update(context.Context, *pb.User) (*pb.Response, error) {
	return nil, nil
}

func (s *server) List(context.Context, *pb.ListRequest) (*pb.Users, error) {
	return nil, nil
}

func (s *server) Delete(context.Context, *pb.User) (*pb.Response, error) {
	return nil, nil
}
