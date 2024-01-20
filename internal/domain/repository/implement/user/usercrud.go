package repository

import (
	"template-grpc/internal/domain/entity"
	objectvalue "template-grpc/internal/domain/object-value"
	ireposity "template-grpc/internal/domain/repository/interface"
)

type userCrud struct {
}

func NewRepository() ireposity.IUserCrud {
	return &userCrud{}
}

func (u *userCrud) Insert(entity.User) *objectvalue.Response {
	return nil
}
func (u *userCrud) Delete(id int32) *objectvalue.Response {
	return nil
}
func (u *userCrud) Update(entity.User) *objectvalue.Response {
	return nil
}
