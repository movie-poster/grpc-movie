package handler_actor

import (
	"context"
	"grpc-movie/internal/constant"
	"grpc-movie/internal/domain/entity"
	irepository "grpc-movie/internal/domain/repository/interface"
	"grpc-movie/internal/utils"
	"net/http"
	"time"

	pb "grpc-movie/internal/infra/proto/actor"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func NewServerActor(crud irepository.IActorCrud, clientCloudinary *cloudinary.Cloudinary) *server {
	return &server{
		crud:       crud,
		cloudinary: clientCloudinary,
	}
}

type server struct {
	crud       irepository.IActorCrud
	cloudinary *cloudinary.Cloudinary
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

	uploadResult, err := s.cloudinary.Upload.Upload(
		context,
		"data:image/png;base64,"+actor.GetAvatar(),
		uploader.UploadParams{
			UploadPreset: "preset-actor",
			Folder:       "director",
			Format:       "png",
		},
	)
	if err != nil {
		return &pb.ResponseActor{
			Title:   "No fue posibile subir el archivo",
			IsOk:    false,
			Message: "Fotografía no subida",
			Status:  http.StatusBadRequest,
		}, nil
	}

	actorObject := &entity.Actor{
		Name:      actor.Name,
		Birthdate: date,
		Avatar:    uploadResult.SecureURL,
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

	if !utils.IsURL(actor.Avatar) && actor.Avatar != "" {
		uploadResult, err := s.cloudinary.Upload.Upload(
			context,
			"data:image/png;base64,"+actor.GetAvatar(),
			uploader.UploadParams{
				UploadPreset: "preset-actor",
				Folder:       "director",
				Format:       "png",
			},
		)
		if err != nil {
			return &pb.ResponseActor{
				Title:   "No fue posibile subir el archivo",
				IsOk:    false,
				Message: "Fotografía no subida",
				Status:  http.StatusBadRequest,
			}, nil
		}
		actor.Avatar = uploadResult.SecureURL
	}

	actorObject := &entity.Actor{
		Model:     entity.Model{ID: actor.GetId()},
		Name:      actor.Name,
		Birthdate: date,
		Avatar:    actor.Avatar,
	}

	response := s.crud.Update(actorObject)

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
		responsePB.TotalPages = uint64(response.Count)
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
