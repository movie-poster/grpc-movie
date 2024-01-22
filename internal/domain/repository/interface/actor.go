package ireposity

import (
	"grpc-movie/internal/domain/entity"
	objectvalue "grpc-movie/internal/domain/object-value"
)

type IActorCrud interface {
	Insert(*entity.Actor) *objectvalue.ResponseValue
	Update(*entity.Actor) *objectvalue.ResponseValue
	List(page, pageSize uint64) *objectvalue.ResponseValue
	Delete(ID uint64) *objectvalue.ResponseValue
}
