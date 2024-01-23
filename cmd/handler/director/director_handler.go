package handler_director

import (
	"context"
	"grpc-movie/internal/constant"
	"grpc-movie/internal/domain/entity"
	irepository "grpc-movie/internal/domain/repository/interface"
	"net/http"
	"time"

	pb "grpc-movie/internal/infra/proto/director"
)

func NewServerDirector(crud irepository.IDirectorCrud) *server {
	return &server{
		crud: crud,
	}
}

type server struct {
	crud irepository.IDirectorCrud
	pb.UnimplementedDirectorCrudServer
}

func (s *server) Insert(context context.Context, director *pb.Director) (*pb.ResponseDirector, error) {
	date, err := time.Parse("2006-01-02", director.GetBirthdate())
	if err != nil {
		return &pb.ResponseDirector{
			Title:   "Error de conversión de fecha",
			IsOk:    false,
			Message: "No fue posible convertir la fecha",
			Status:  http.StatusBadRequest,
		}, nil
	}

	directorObject := &entity.Director{
		Name:      director.Name,
		Birthdate: date,
		State:     constant.ActiveState,
	}

	response := s.crud.Insert(directorObject)

	responsePB := pb.ResponseDirector{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	if response.Value != nil {
		responsePB.Director = response.Value.(*pb.Director)
	}

	return &responsePB, nil
}

func (s *server) Update(context context.Context, director *pb.Director) (*pb.ResponseDirector, error) {
	date, err := time.Parse("2006-01-02", director.GetBirthdate())
	if err != nil {
		return &pb.ResponseDirector{
			Title:   "Error de conversión de fecha",
			IsOk:    false,
			Message: "No fue posible convertir la fecha",
			Status:  http.StatusBadRequest,
		}, nil
	}

	directorObject := &entity.Director{
		Name:      director.Name,
		Birthdate: date,
	}

	response := s.crud.Update(directorObject)

	responsePB := pb.ResponseDirector{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	return &responsePB, nil
}

func (s *server) List(context context.Context, req *pb.ListRequestDirector) (*pb.ResponseDirector, error) {
	response := s.crud.List(int(req.GetPage()), int(req.GetPageSize()))

	responsePB := pb.ResponseDirector{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	if response.Value != nil {
		responsePB.Directors = response.Value.([]*pb.Director)
	}

	return &responsePB, nil
}

func (s *server) Delete(context context.Context, req *pb.RequestByIdDirector) (*pb.ResponseDirector, error) {
	response := s.crud.Delete(req.GetId())

	responsePB := pb.ResponseDirector{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	return &responsePB, nil
}

func (s *server) GetById(context context.Context, req *pb.RequestByIdDirector) (*pb.ResponseDirector, error) {
	return nil, nil
}
