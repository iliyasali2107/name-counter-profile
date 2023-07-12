package main

import (
	"fmt"
	"url-redirecter-url/cmd/client/routes"
	"url-redirecter-url/pkg/config"
	"url-redirecter-url/pkg/pb"

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

func (svc *ServiceClient) AddURL(ctx *gin.Context) {
	routes.AddURL(ctx, svc.Client)
}

func (svc *ServiceClient) SetActiveURL(ctx *gin.Context) {
	routes.SetActiveURL(ctx, svc.Client)
}

func (svc *ServiceClient) GetUserURLs(ctx *gin.Context) {
	routes.GetUserURLs(ctx, svc.Client)
}
