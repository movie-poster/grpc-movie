package repository_actor

import (
	"grpc-movie/internal/constant"
	"grpc-movie/internal/domain/entity"
	objectvalue "grpc-movie/internal/domain/object-value"
	irepository "grpc-movie/internal/domain/repository/interface"
	pb "grpc-movie/internal/infra/proto/actor"
	"grpc-movie/internal/utils"
	"net/http"

	"gorm.io/gorm"
)

type crud struct {
	DB *gorm.DB
}

func ActorRepository(DB *gorm.DB) irepository.IActorCrud {
	return &crud{
		DB,
	}
}

func (u *crud) Insert(actor *entity.Actor) *objectvalue.ResponseValue {
	genreQuery := entity.Actor{}

	err := u.DB.Model(&entity.Actor{}).
		Where("name = ? AND state = ?", actor.Name, constant.ActiveState).
		First(&genreQuery).Error
	if err != nil {
		utils.LogWarning("Registro no encontrado", "No se ha encontrado el registro", "Insert", actor)
	}

	if genreQuery.ID == constant.NotExists {
		err := u.DB.Create(&actor).Error
		if err != nil {
			message := utils.CheckErrorFromDB(err)
			utils.LogError("Error al guardar el registro", message, "Insert", http.StatusBadRequest, actor)
			return objectvalue.BadResponseSingle(message)
		}

		utils.LogSuccess("Registro guardado", "Insert", actor)
		return &objectvalue.ResponseValue{
			Title:   "¡Proceso exitoso!",
			IsOk:    true,
			Message: "El actor se ha creado",
			Status:  http.StatusCreated,
			Value:   u.MarshalResponse(actor),
		}
	}

	utils.LogError("Error al guardar el registro", "El género ya está creado o no hay datos existentes", "Insert", http.StatusBadRequest, actor)
	return objectvalue.BadResponseSingle("El género ya está creado o no hay datos existentes")
}

func (u *crud) Update(actor *entity.Actor) *objectvalue.ResponseValue {
	actorMap := map[string]any{
		"name":      actor.Name,
		"birthdate": actor.Birthdate,
	}

	if utils.IsURL(actor.Avatar) {
		actorMap["avatar"] = actor.Avatar
	}

	err := u.DB.Model(&entity.Actor{}).Where("id", actor.ID).Updates(actorMap).Error
	if err != nil {
		message := utils.CheckErrorFromDB(err)
		utils.LogError("Error al actualizar el registro", message, "Update", http.StatusBadRequest, actor)
		return objectvalue.BadResponseSingle(message)
	}

	utils.LogSuccess("Registro actualizado", "Update", actor)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "Registro actualizado",
		Status:  http.StatusOK,
		Value:   u.MarshalResponse(actor),
	}
}

func (u *crud) List(page, pageSize int) *objectvalue.ResponseValue {
	// Contar número de registros
	countResult := make(chan int64)
	defer close(countResult)

	go func() {
		var count int64
		u.DB.Model(&entity.Actor{}).Where("state = ?", constant.ActiveState).Count(&count)
		countResult <- count
	}()

	// consulta para traer los registros
	var actors []*entity.Actor
	var actorsPT []*pb.Actor

	err := u.DB.Limit(pageSize).
		Offset(pageSize*page).
		Find(&actors, "state = ?", constant.ActiveState).Error
	if err != nil {
		utils.LogError("Error al listar los registros", "No es posible listar los registros, revisar base de datos", "ListActors", http.StatusBadRequest)
		return &objectvalue.ResponseValue{
			Title:   "Proceso no exitoso",
			IsOk:    false,
			Message: "No se han podido listar los actores",
			Status:  http.StatusInternalServerError,
			Value:   actorsPT,
		}
	}

	totalCount := <-countResult / int64(pageSize)

	for _, actor := range actors {
		tempActor := u.MarshalResponse(actor)
		actorsPT = append(actorsPT, tempActor)
	}

	utils.LogSuccess("Listado generado exitosamente", "ListActors", page, pageSize)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "Listado de actores",
		Status:  http.StatusOK,
		Value:   actorsPT,
		Count:   totalCount,
	}
}

func (u *crud) Delete(ID uint64) *objectvalue.ResponseValue {
	err := u.DB.Model(&entity.Actor{}).Where("id", ID).Update("state", constant.InactiveState).Error
	if err != nil {
		message := utils.CheckErrorFromDB(err)
		utils.LogError("Error al eliminar el registro", message, "Delete", http.StatusBadRequest, ID)
		return objectvalue.BadResponseSingle(message)
	}

	utils.LogSuccess("Registro eliminado", "Delete", ID)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "Se eliminó correctamente",
		Status:  http.StatusOK,
	}
}

func (u *crud) MarshalResponse(actor *entity.Actor) *pb.Actor {
	return &pb.Actor{
		Id:        actor.ID,
		Name:      actor.Name,
		State:     actor.State,
		Birthdate: utils.FormatDate(actor.Birthdate),
		Avatar:    actor.Avatar,
		CreatedAt: actor.CreatedAt.String(),
		UpdatedAt: actor.UpdatedAt.String(),
	}
}
