package ireposity

import (
	"grpc-movie/internal/domain/entity"
	objectvalue "grpc-movie/internal/domain/object-value"
)

type IGenreCrud interface {
	Insert(*entity.Genre) *objectvalue.ResponseValue
	List(page, pageSize uint64) *objectvalue.ResponseValue
	Delete(ID uint64) *objectvalue.ResponseValue
}
