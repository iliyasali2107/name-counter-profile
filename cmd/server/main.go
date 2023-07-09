package main

import (
	"fmt"
	"log"
	"net"

	"name-counter-url/pkg/config"
	"name-counter-url/pkg/db"
	"name-counter-url/pkg/pb"
	"name-counter-url/pkg/service"

	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to config", err)
	}

	storage := db.Init(c.DBUrl)

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln("Failed to listening")
	}

	fmt.Println("Auth service is on: ", c.Port)

	srv := service.NewService(storage)

	grpcServer := grpc.NewServer()

	pb.RegisterURLServiceServer(grpcServer, srv)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve: ", err)
	}
}
