package handler_director

import (
	"context"
	"grpc-movie/internal/constant"
	"grpc-movie/internal/domain/entity"
	irepository "grpc-movie/internal/domain/repository/interface"
	pb "grpc-movie/internal/infra/proto/director"
	"grpc-movie/internal/utils"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func NewServerDirector(crud irepository.IDirectorCrud, clientCloudinary *cloudinary.Cloudinary) *server {
	return &server{
		crud:       crud,
		cloudinary: clientCloudinary,
	}
}

type server struct {
	crud       irepository.IDirectorCrud
	cloudinary *cloudinary.Cloudinary
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

	uploadResult, err := s.cloudinary.Upload.Upload(
		context,
		"data:image/png;base64,"+director.GetAvatar(),
		uploader.UploadParams{
			UploadPreset: "preset-director",
			Folder:       "director",
			Format:       "png",
		},
	)
	if err != nil {
		return &pb.ResponseDirector{
			Title:   "No fue posibile subir el archivo",
			IsOk:    false,
			Message: "Fotografía no subida",
			Status:  http.StatusBadRequest,
		}, nil
	}

	directorObject := &entity.Director{
		Name:      director.Name,
		Birthdate: date,
		Avatar:    uploadResult.SecureURL,
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

	if !utils.IsURL(director.Avatar) && director.Avatar != "" {
		uploadResult, err := s.cloudinary.Upload.Upload(
			context,
			"data:image/png;base64,"+director.GetAvatar(),
			uploader.UploadParams{
				UploadPreset: "preset-director",
				Folder:       "director",
				Format:       "png",
			},
		)
		if err != nil {
			return &pb.ResponseDirector{
				Title:   "No fue posibile subir el archivo",
				IsOk:    false,
				Message: "Fotografía no subida",
				Status:  http.StatusBadRequest,
			}, nil
		}
		director.Avatar = uploadResult.SecureURL
	}

	directorObject := &entity.Director{
		Model:     entity.Model{ID: director.GetId()},
		Name:      director.Name,
		Birthdate: date,
		Avatar:    director.Avatar,
	}

	response := s.crud.Update(directorObject)

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
		responsePB.TotalPages = uint64(response.Count)
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
