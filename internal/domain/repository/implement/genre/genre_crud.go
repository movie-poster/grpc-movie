package repository_genre

import (
	"grpc-movie/internal/constant"
	"grpc-movie/internal/domain/entity"
	objectvalue "grpc-movie/internal/domain/object-value"
	irepository "grpc-movie/internal/domain/repository/interface"
	pb "grpc-movie/internal/infra/proto/genre"
	"grpc-movie/internal/utils"
	"net/http"

	"gorm.io/gorm"
)

type crud struct {
	DB *gorm.DB
}

func GenreRepository(DB *gorm.DB) irepository.IGenreCrud {
	return &crud{
		DB,
	}
}

func (u *crud) Insert(genre *entity.Genre) *objectvalue.ResponseValue {
	genreQuery := entity.Genre{}

	err := u.DB.Model(&entity.Genre{}).
		Where("name = ? AND state = ?", genre.Name, constant.ActiveState).
		First(&genreQuery).Error
	if err != nil {
		utils.LogWarning("Registro no encontrado", "No se ha encontrado el registro", "Insert", genre)
	}

	if genreQuery.ID == constant.NotExists {
		err := u.DB.Create(&genre).Error
		if err != nil {
			message := utils.CheckErrorFromDB(err)
			utils.LogError("Error al guardar el registro", message, "Insert", http.StatusBadRequest, genre)
			return objectvalue.BadResponseSingle(message)
		}

		utils.LogSuccess("Registro guardado", "Insert", genre)
		return &objectvalue.ResponseValue{
			Title:   "¡Proceso exitoso!",
			IsOk:    true,
			Message: "La película se ha creado",
			Status:  http.StatusCreated,
			Value:   u.MarshalResponse(genre),
		}
	}

	utils.LogError("Error al guardar el registro", "El género ya está creado o no hay datos existentes", "Insert", http.StatusBadRequest, genre)
	return objectvalue.BadResponseSingle("El género ya está creado o no hay datos existentes")
}

func (u *crud) List(page, pageSize int) *objectvalue.ResponseValue {
	var movies []*entity.Genre
	var moviesPT []*pb.Genre

	err := u.DB.Limit(pageSize).
		Offset(pageSize*page).
		Find(&movies, "state = ?", constant.ActiveState).Error
	if err != nil {
		utils.LogError("Error al listar los registros", "No es posible listar los registros, revisar base de datos", "List", http.StatusBadRequest)
		return &objectvalue.ResponseValue{
			Title:   "Proceso no exitoso",
			IsOk:    false,
			Message: "No se han podido listar los géneros",
			Status:  http.StatusBadRequest,
			Value:   moviesPT,
		}
	}

	for _, movie := range movies {
		tempMovie := u.MarshalResponse(movie)
		moviesPT = append(moviesPT, tempMovie)
	}

	utils.LogSuccess("Listado generado exitosamente", "List", page, pageSize)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "Listado de géneros",
		Status:  http.StatusOK,
		Value:   moviesPT,
	}
}

func (u *crud) Delete(ID uint64) *objectvalue.ResponseValue {
	err := u.DB.Model(&entity.Genre{}).Where("id", ID).Update("state", constant.InactiveState).Error
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

func (u *crud) MarshalResponse(genre *entity.Genre) *pb.Genre {
	return &pb.Genre{
		Id:        genre.ID,
		Name:      genre.Name,
		State:     genre.State,
		CreatedAt: genre.CreatedAt.String(),
		UpdatedAt: genre.UpdatedAt.String(),
	}
}
