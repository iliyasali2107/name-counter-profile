package main

import (
	"fmt"

	"name-counter-profile/cmd/client/routes"
	"name-counter-profile/pkg/config"
	"name-counter-profile/pkg/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.ProfileServiceClient
}

func InitServiceClient(c *config.Config) pb.ProfileServiceClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.Port, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewProfileServiceClient(cc)
}

func (svc *ServiceClient) GetURL(ctx *gin.Context) {
	routes.GetURL(ctx, svc.Client)
}
