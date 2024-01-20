package main

import (
	"log"
	"net"
	"template-grpc/cmd/config"

	"google.golang.org/grpc"
)

func main() {
	conf := config.GetConfig()
	listener, err := net.Listen("tcp", ":"+conf.Server.Port)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	s = config.Run(s, "")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
