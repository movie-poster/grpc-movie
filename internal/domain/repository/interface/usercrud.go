package ireposity

import (
	"template-grpc/internal/domain/entity"
	objectvalue "template-grpc/internal/domain/object-value"
)

type IUserCrud interface {
	Insert(entity.User) *objectvalue.Response
	Delete(id int32) *objectvalue.Response
	Update(entity.User) *objectvalue.Response
}
