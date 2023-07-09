package main

import (
	"fmt"

	"name-counter-url/cmd/client/routes"
	"name-counter-url/pkg/config"
	"name-counter-url/pkg/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.URLServiceClient
}

func InitServiceClient(c *config.Config) pb.URLServiceClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.Port, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewURLServiceClient(cc)
}

func (svc *ServiceClient) GetURL(ctx *gin.Context) {
	routes.GetURL(ctx, svc.Client)
}
