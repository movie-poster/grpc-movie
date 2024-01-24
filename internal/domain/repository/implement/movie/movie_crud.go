package repository_movie

import (
	"grpc-movie/internal/constant"
	"grpc-movie/internal/domain/entity"
	objectvalue "grpc-movie/internal/domain/object-value"
	irepository "grpc-movie/internal/domain/repository/interface"
	pa "grpc-movie/internal/infra/proto/actor"
	pd "grpc-movie/internal/infra/proto/director"
	pg "grpc-movie/internal/infra/proto/genre"
	pb "grpc-movie/internal/infra/proto/movie"
	"grpc-movie/internal/utils"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type crud struct {
	DB *gorm.DB
}

func MovieRepository(DB *gorm.DB) irepository.IMovieCrud {
	return &crud{
		DB,
	}
}

func (u *crud) Insert(movie *entity.Movie) *objectvalue.ResponseValue {
	movieQuery := entity.Movie{}

	err := u.DB.Model(&entity.Movie{}).
		Where("title = ? AND state = ?", movie.Title, constant.ActiveState).
		First(&movieQuery).Error
	if err != nil {
		utils.LogWarning("Registro no encontrado", "No se ha encontrado el registro", "Insert", movie)
	}

	if movieQuery.ID == constant.NotExists {
		err = u.DB.Create(&movie).Error
		if err != nil {
			message := utils.CheckErrorFromDB(err)
			utils.LogError("Error al guardar el registro", message, "Insert", http.StatusBadRequest, movie)
			return objectvalue.BadResponseSingle(message)
		}

		utils.LogSuccess("Registro guardado", "Insert", movie)
		return &objectvalue.ResponseValue{
			Title:   "¡Proceso exitoso!",
			IsOk:    true,
			Message: "La película se ha creado",
			Status:  http.StatusCreated,
			Value:   u.MarshalResponse(movie),
		}
	}

	utils.LogError("Error al guardar el registro", "El género ya está creado o no hay datos existentes", "Insert", http.StatusBadRequest, movie)
	return objectvalue.BadResponseSingle("El género ya está creado o no hay datos existentes")
}

func (u *crud) Update(movie *entity.Movie) *objectvalue.ResponseValue {
	movieMap := map[string]any{
		"title":       movie.Title,
		"synopsis":    movie.Synopsis,
		"year":        movie.Year,
		"rating":      movie.Rating,
		"duration":    movie.Duration,
		"director_id": movie.DirectorID,
	}

	err := u.DB.Model(entity.Movie{}).Where("id = ?", movie.ID).Updates(movieMap).Error
	if err != nil {
		message := utils.CheckErrorFromDB(err)
		utils.LogError("Error al actualizar el registro", message, "Update", http.StatusBadRequest, movie)
		return objectvalue.BadResponseSingle(message)
	}

	utils.LogSuccess("Registro actualizado", "Update", movie)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "la película fue editada correctamente",
		Status:  http.StatusOK,
	}
}

func (u *crud) List(req *pb.ListRequestMovie) *objectvalue.ResponseValue {
	var movies []*entity.Movie
	var moviesPT []*pb.Movie

	query := u.DB.Model(&entity.Movie{}).
		Preload("Director", "state = ?", constant.ActiveState).
		Preload("Actors", "state = ?", constant.ActiveState).
		Preload("Genres", "state = ?", constant.ActiveState)

	movieTitle := req.GetFilterCriteria().GetMovieTitle()
	if movieTitle != "" {
		query.Where("LOWER(title) LIKE ? AND state = ?", "%"+strings.ToLower(movieTitle)+"%", constant.ActiveState)
	}

	genreName := req.GetFilterCriteria().GetGenreName()
	if genreName != "" {
		query.Where("EXISTS (SELECT 1 FROM movie_genres mg WHERE mg.movie_id = movies.id AND mg.genre_id IN (SELECT id FROM genres WHERE LOWER(name) LIKE ? AND state = ?))", "%"+strings.ToLower(genreName)+"%", constant.ActiveState)
	}

	err := query.Limit(int(req.GetPageSize())).
		Offset(int(req.GetPage())*int(req.GetPageSize())).
		Find(&movies, "state = ?", constant.ActiveState).Error
	if err != nil {
		utils.LogError("Error al listar los registros", "No es posible listar los registros, revisar base de datos", "List", http.StatusBadRequest)
		return &objectvalue.ResponseValue{
			Title:   "Proceso no exitoso",
			IsOk:    false,
			Message: "No se han podido listar las películas",
			Status:  http.StatusInternalServerError,
			Value:   moviesPT,
		}
	}

	for _, movie := range movies {
		tempMovie := u.MarshalResponse(movie)
		moviesPT = append(moviesPT, tempMovie)
	}

	utils.LogSuccess("Listado generado exitosamente", "List", req)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "Listado de películas",
		Status:  http.StatusOK,
		Value:   moviesPT,
	}
}

func (u *crud) Delete(ID uint64) *objectvalue.ResponseValue {
	err := u.DB.Model(&entity.Movie{}).Where("id", ID).Update("state", constant.InactiveState).Error
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

func (u *crud) GetById(ID uint64) *objectvalue.ResponseValue {
	var movie *entity.Movie

	err := u.DB.Preload("Director", "state = ?", constant.ActiveState).
		Preload("Actors", "state = ?", constant.ActiveState).
		Preload("Genres", "state = ?", constant.ActiveState).
		Where("id = ? state = ?", ID, constant.ActiveState).
		First(&movie).Error
	if err != nil {
		message := utils.CheckErrorFromDB(err)
		utils.LogError("Error al buscar por ID", message, "GetById", http.StatusBadRequest, ID)
		return &objectvalue.ResponseValue{
			Title:   "Proceso no existoso",
			IsOk:    false,
			Message: message,
			Status:  http.StatusBadRequest,
			Value:   &pb.Movie{},
		}
	}

	utils.LogSuccess("Registro encontrado", "GetById", ID)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "Se ha encontrado la película con el ID",
		Status:  http.StatusCreated,
		Value:   u.MarshalResponse(movie),
	}
}

func (u *crud) MarshalResponse(movie *entity.Movie) *pb.Movie {
	moviePT := &pb.Movie{
		Id:         movie.ID,
		Title:      movie.Title,
		Synopsis:   movie.Synopsis,
		Year:       movie.Year,
		Rating:     float32(movie.Rating),
		Duration:   movie.Duration,
		DirectorId: movie.DirectorID,
		State:      movie.State,
		CreatedAt:  movie.CreatedAt.String(),
		UpdatedAt:  movie.UpdatedAt.String(),
		Director: &pd.Director{
			Id:        movie.Director.ID,
			Name:      movie.Director.Name,
			Birthdate: movie.Director.Birthdate.String(),
			Avatar:    movie.Director.Avatar,
			State:     movie.Director.State,
			CreatedAt: movie.Director.CreatedAt.String(),
			UpdatedAt: movie.Director.UpdatedAt.String(),
		},
	}

	for _, actor := range movie.Actors {
		moviePT.Actors = append(moviePT.Actors, &pa.Actor{
			Id:        actor.ID,
			Name:      actor.Name,
			Birthdate: actor.Birthdate.String(),
			Avatar:    actor.Avatar,
			State:     actor.State,
			CreatedAt: actor.CreatedAt.String(),
			UpdatedAt: actor.UpdatedAt.String(),
		})
	}

	for _, genre := range movie.Genres {
		moviePT.Genres = append(moviePT.Genres, &pg.Genre{
			Id:        genre.ID,
			Name:      genre.Name,
			State:     genre.State,
			CreatedAt: genre.CreatedAt.String(),
			UpdatedAt: genre.UpdatedAt.String(),
		})
	}

	return moviePT
}
