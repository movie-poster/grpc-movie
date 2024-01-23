package handler_actor

import (
	"context"
	"grpc-movie/internal/constant"
	"grpc-movie/internal/domain/entity"
	irepository "grpc-movie/internal/domain/repository/interface"
	"net/http"
	"time"

	pb "grpc-movie/internal/infra/proto/actor"
)

func NewServerActor(crud irepository.IActorCrud) *server {
	return &server{
		crud: crud,
	}
}

type server struct {
	crud irepository.IActorCrud
	pb.UnimplementedActorCrudServer
}

func (s *server) Insert(context context.Context, actor *pb.Actor) (*pb.ResponseActor, error) {
	date, err := time.Parse("2006-01-02", actor.GetBirthdate())
	if err != nil {
		return &pb.ResponseActor{
			Title:   "Error de conversión de fecha",
			IsOk:    false,
			Message: "No fue posible convertir la fecha",
			Status:  http.StatusBadRequest,
		}, nil
	}

	actorObject := &entity.Actor{
		Name:      actor.Name,
		Birthdate: date,
		State:     constant.ActiveState,
	}

	response := s.crud.Insert(actorObject)

	responsePB := pb.ResponseActor{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	if response.Value != nil {
		responsePB.Actor = response.Value.(*pb.Actor)
	}

	return &responsePB, nil
}

func (s *server) Update(context context.Context, actor *pb.Actor) (*pb.ResponseActor, error) {
	date, err := time.Parse("2006-01-02", actor.GetBirthdate())
	if err != nil {
		return &pb.ResponseActor{
			Title:   "Error de conversión de fecha",
			IsOk:    false,
			Message: "No fue posible convertir la fecha",
			Status:  http.StatusBadRequest,
		}, nil
	}

	actorObject := &entity.Actor{
		Name:      actor.Name,
		Birthdate: date,
	}

	response := s.crud.Update(actorObject)

	responsePB := pb.ResponseActor{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	return &responsePB, nil
}

func (s *server) List(context context.Context, req *pb.ListRequestActor) (*pb.ResponseActor, error) {
	response := s.crud.List(int(req.GetPage()), int(req.GetPageSize()))

	responsePB := pb.ResponseActor{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	if response.Value != nil {
		responsePB.Actors = response.Value.([]*pb.Actor)
	}

	return &responsePB, nil
}

func (s *server) Delete(context context.Context, req *pb.RequestByIdActor) (*pb.ResponseActor, error) {
	response := s.crud.Delete(req.GetId())

	responsePB := pb.ResponseActor{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	return &responsePB, nil
}

func (s *server) GetById(context context.Context, req *pb.RequestByIdActor) (*pb.ResponseActor, error) {
	return nil, nil
}
