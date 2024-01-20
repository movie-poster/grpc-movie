package config

import (
	"flag"
	"template-grpc/cmd/handler"
	repository "template-grpc/internal/domain/repository/implement/user"
	pb "template-grpc/internal/infra/proto"

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
	pb.RegisterUserCrudServer(s, handler.NewServerUser(repository.NewRepository()))
	return s

}
