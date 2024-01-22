package ireposity

import (
	"grpc-movie/internal/domain/entity"
	objectvalue "grpc-movie/internal/domain/object-value"
	"grpc-movie/internal/infra/proto/movie"
)

type IMovieCrud interface {
	Insert(*entity.Movie) *objectvalue.ResponseValue
	Update(*entity.Movie) *objectvalue.ResponseValue
	List(*movie.ListRequestMovie) *objectvalue.ResponseValue
	Delete(ID uint64) *objectvalue.ResponseValue
	GetById(ID uint64) *objectvalue.ResponseValue
}
