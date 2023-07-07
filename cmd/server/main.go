package main

import (
	"fmt"
	"log"
	"net"

	"name-counter-profile/pkg/config"
	"name-counter-profile/pkg/db"
	"name-counter-profile/pkg/pb"
	"name-counter-profile/pkg/service"

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

	pb.RegisterProfileServiceServer(grpcServer, srv)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve: ", err)
	}
}
