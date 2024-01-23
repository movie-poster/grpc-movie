package config

import (
	"flag"

	handler_actor "grpc-movie/cmd/handler/actor"
	handler_director "grpc-movie/cmd/handler/director"
	handler_genre "grpc-movie/cmd/handler/genre"
	handler_movie "grpc-movie/cmd/handler/movie"

	repository_actor "grpc-movie/internal/domain/repository/implement/actor"
	repository_director "grpc-movie/internal/domain/repository/implement/director"
	repository_genre "grpc-movie/internal/domain/repository/implement/genre"
	repository_movie "grpc-movie/internal/domain/repository/implement/movie"

	"grpc-movie/internal/infra/proto/actor"
	"grpc-movie/internal/infra/proto/director"
	"grpc-movie/internal/infra/proto/genre"
	"grpc-movie/internal/infra/proto/movie"

	"google.golang.org/grpc"
)

func init() {
	var configPath = ""
	configPath = *flag.String("config", "", "")

	if configPath == "" {
		configPath = "../data/config.yml"
	}

	setConfiguration(configPath)
}

func setConfiguration(configPath string) {
	Setup(configPath)
}

func Run(s *grpc.Server, configPath string) *grpc.Server {

	conf := GetConfig()
	setupDB(conf)

	movie.RegisterMovieCrudServer(s, handler_movie.NewServerMovie(repository_movie.MovieRepository(DB)))
	genre.RegisterGenreCrudServer(s, handler_genre.NewServerGenre(repository_genre.GenreRepository(DB)))
	actor.RegisterActorCrudServer(s, handler_actor.NewServerActor(repository_actor.ActorRepository(DB)))
	director.RegisterDirectorCrudServer(s, handler_director.NewServerDirector(repository_director.DirectorRepository(DB)))

	return s
}
