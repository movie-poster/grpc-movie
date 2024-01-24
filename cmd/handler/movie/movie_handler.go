package handler_movie

import (
	"context"
	"grpc-movie/internal/constant"
	"grpc-movie/internal/domain/entity"
	irepository "grpc-movie/internal/domain/repository/interface"
	"net/http"

	pb "grpc-movie/internal/infra/proto/movie"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func NewServerMovie(crud irepository.IMovieCrud, clientCloudinary *cloudinary.Cloudinary) *server {
	return &server{
		crud:       crud,
		cloudinary: clientCloudinary,
	}
}

type server struct {
	crud       irepository.IMovieCrud
	cloudinary *cloudinary.Cloudinary
	pb.UnimplementedMovieCrudServer
}

func (s *server) Insert(context context.Context, movie *pb.Movie) (*pb.ResponseMovie, error) {
	uploadResult, err := s.cloudinary.Upload.Upload(
		context,
		"data:image/png;base64,"+movie.GetPoster(),
		uploader.UploadParams{
			UploadPreset: "preset-movie",
			Folder:       "movie",
			Format:       "png",
		},
	)
	if err != nil {
		return &pb.ResponseMovie{
			Title:   "No fue posibile subir el archivo",
			IsOk:    false,
			Message: "Fotografía no subida",
			Status:  http.StatusBadRequest,
		}, nil
	}

	movieObject := &entity.Movie{
		Title:      movie.GetTitle(),
		Synopsis:   movie.GetSynopsis(),
		Year:       movie.GetYear(),
		Rating:     float64(movie.GetRating()),
		Duration:   movie.GetDuration(),
		DirectorID: movie.GetDirectorId(),
		Poster:     uploadResult.SecureURL,
		State:      constant.ActiveState,
	}

	for _, actor := range movie.GetActors() {
		movieObject.Actors = append(movieObject.Actors, entity.Actor{
			Model: entity.Model{ID: actor.GetId()},
		})
	}

	for _, genre := range movie.GetGenres() {
		movieObject.Actors = append(movieObject.Actors, entity.Actor{
			Model: entity.Model{ID: genre.GetId()},
		})
	}

	response := s.crud.Insert(movieObject)

	responsePB := pb.ResponseMovie{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	if response.Value != nil {
		responsePB.Movie = response.Value.(*pb.Movie)
	}

	return &responsePB, nil
}

func (s *server) Update(context context.Context, movie *pb.Movie) (*pb.ResponseMovie, error) {
	movieObject := &entity.Movie{
		Title:      movie.GetTitle(),
		Synopsis:   movie.GetSynopsis(),
		Year:       movie.GetYear(),
		Rating:     float64(movie.GetRating()),
		Duration:   movie.GetDuration(),
		DirectorID: movie.GetDirectorId(),
		State:      constant.ActiveState,
	}

	response := s.crud.Update(movieObject)

	responsePB := pb.ResponseMovie{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	return &responsePB, nil
}

func (s *server) List(context context.Context, req *pb.ListRequestMovie) (*pb.ResponseMovie, error) {
	response := s.crud.List(req)

	responsePB := pb.ResponseMovie{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	if response.Value != nil {
		responsePB.Movies = response.Value.([]*pb.Movie)
	}

	return &responsePB, nil
}

func (s *server) Delete(context context.Context, req *pb.RequestByIdMovie) (*pb.ResponseMovie, error) {
	response := s.crud.Delete(req.GetId())

	responsePB := pb.ResponseMovie{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	return &responsePB, nil
}

func (s *server) GetById(context context.Context, req *pb.RequestByIdMovie) (*pb.ResponseMovie, error) {
	response := s.crud.GetById(req.GetId())

	responsePB := pb.ResponseMovie{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	if response.Value != nil {
		responsePB.Movie = response.Value.(*pb.Movie)
	}

	return &responsePB, nil
}
