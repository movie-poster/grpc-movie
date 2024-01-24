package handler_genre

import (
	"context"
	"grpc-movie/internal/constant"
	"grpc-movie/internal/domain/entity"
	irepository "grpc-movie/internal/domain/repository/interface"

	pb "grpc-movie/internal/infra/proto/genre"
)

func NewServerGenre(crud irepository.IGenreCrud) *server {
	return &server{
		crud: crud,
	}
}

type server struct {
	crud irepository.IGenreCrud
	pb.UnimplementedGenreCrudServer
}

func (s *server) Insert(context context.Context, genre *pb.Genre) (*pb.ResponseGenre, error) {
	genreObject := &entity.Genre{
		Name:  genre.Name,
		State: constant.ActiveState,
	}

	response := s.crud.Insert(genreObject)

	responsePB := pb.ResponseGenre{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	if response.Value != nil {
		responsePB.Genre = response.Value.(*pb.Genre)
	}

	return &responsePB, nil
}

func (s *server) List(context context.Context, req *pb.ListRequestGenre) (*pb.ResponseGenre, error) {
	response := s.crud.List(int(req.GetPage()), int(req.GetPageSize()))

	responsePB := pb.ResponseGenre{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	if response.Value != nil {
		responsePB.Genres = response.Value.([]*pb.Genre)
		responsePB.TotalPages = uint64(response.Count)
	}

	return &responsePB, nil
}

func (s *server) Delete(context context.Context, req *pb.RequestByIdGenre) (*pb.ResponseGenre, error) {
	response := s.crud.Delete(req.GetId())

	responsePB := pb.ResponseGenre{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	return &responsePB, nil
}
