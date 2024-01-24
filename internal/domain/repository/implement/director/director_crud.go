package repository_director

import (
	"grpc-movie/internal/constant"
	"grpc-movie/internal/domain/entity"
	objectvalue "grpc-movie/internal/domain/object-value"
	irepository "grpc-movie/internal/domain/repository/interface"
	pb "grpc-movie/internal/infra/proto/director"
	"grpc-movie/internal/utils"
	"net/http"

	"gorm.io/gorm"
)

type crud struct {
	DB *gorm.DB
}

func DirectorRepository(DB *gorm.DB) irepository.IDirectorCrud {
	return &crud{
		DB,
	}
}

func (u *crud) Insert(director *entity.Director) *objectvalue.ResponseValue {
	genreQuery := entity.Director{}

	err := u.DB.Model(&entity.Director{}).
		Where("name = ? AND state = ?", director.Name, constant.ActiveState).
		First(&genreQuery).Error
	if err != nil {
		utils.LogWarning("Registro no encontrado", "No se ha encontrado el registro", "Insert", director)
	}

	if genreQuery.ID == constant.NotExists {
		err := u.DB.Create(&director).Error
		if err != nil {
			message := utils.CheckErrorFromDB(err)
			utils.LogError("Error al guardar el registro", message, "Insert", http.StatusBadRequest, director)
			return objectvalue.BadResponseSingle(message)
		}

		utils.LogSuccess("Registro guardado", "Insert", director)
		return &objectvalue.ResponseValue{
			Title:   "¡Proceso exitoso!",
			IsOk:    true,
			Message: "La película se ha creado",
			Status:  http.StatusCreated,
			Value:   u.MarshalResponse(director),
		}
	}

	utils.LogError("Error al guardar el registro", "El género ya está creado o no hay datos existentes", "Insert", http.StatusBadRequest, director)
	return objectvalue.BadResponseSingle("El género ya está creado o no hay datos existentes")
}

func (u *crud) Update(director *entity.Director) *objectvalue.ResponseValue {
	directorMap := map[string]any{
		"name":      director.Name,
		"birthdate": director.Birthdate,
	}

	if utils.IsURL(director.Avatar) {
		directorMap["avatar"] = director.Avatar
	}

	err := u.DB.Model(&entity.Director{}).Where("id", director.ID).Updates(directorMap).Error
	if err != nil {
		message := utils.CheckErrorFromDB(err)
		utils.LogError("Error al actualizar el registro", message, "Update", http.StatusBadRequest, director)
		return objectvalue.BadResponseSingle(message)
	}

	utils.LogSuccess("Registro actualizado", "Update", director)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "Registro actualizado",
		Status:  http.StatusOK,
		Value:   u.MarshalResponse(director),
	}
}

func (u *crud) List(page, pageSize int) *objectvalue.ResponseValue {
	// Contar número de registros
	countResult := make(chan int64)
	defer close(countResult)

	go func() {
		var count int64
		u.DB.Model(&entity.Director{}).Where("state = ?", constant.ActiveState).Count(&count)
		countResult <- count
	}()

	// consulta para traer los registros
	var directors []*entity.Director
	var directorsPT []*pb.Director

	err := u.DB.Limit(pageSize).
		Offset(pageSize*page).
		Find(&directors, "state = ?", constant.ActiveState).Error
	if err != nil {
		utils.LogError("Error al listar los registros", "No es posible listar los registros, revisar base de datos", "List", http.StatusBadRequest)
		return &objectvalue.ResponseValue{
			Title:   "Proceso no exitoso",
			IsOk:    false,
			Message: "No se han podido listar los directores",
			Status:  http.StatusInternalServerError,
			Value:   directorsPT,
		}
	}

	totalCount := <-countResult / int64(pageSize)

	for _, movie := range directors {
		tempMovie := u.MarshalResponse(movie)
		directorsPT = append(directorsPT, tempMovie)
	}

	utils.LogSuccess("Listado generado exitosamente", "List", page, pageSize)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "Listado de directores",
		Status:  http.StatusOK,
		Value:   directorsPT,
		Count:   totalCount,
	}
}

func (u *crud) Delete(ID uint64) *objectvalue.ResponseValue {
	err := u.DB.Model(&entity.Director{}).Where("id", ID).Update("state", constant.InactiveState).Error
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

func (u *crud) MarshalResponse(director *entity.Director) *pb.Director {
	return &pb.Director{
		Id:        director.ID,
		Name:      director.Name,
		Birthdate: utils.FormatDate(director.Birthdate),
		Avatar:    director.Avatar,
		State:     director.State,
		CreatedAt: director.CreatedAt.String(),
		UpdatedAt: director.UpdatedAt.String(),
	}
}
