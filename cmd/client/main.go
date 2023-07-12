package main

import (
	"log"
	"url-redirecter-url/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to config", err)
	}

	svc := &ServiceClient{
		Client: InitServiceClient(&c),
	}
	routes := gin.Default()

	routes.GET("/urls/:id", svc.GetURL)
	routes.GET("/urls", svc.GetUserURLs)
	routes.POST("/urls", svc.AddURL)
	routes.POST("/urls/activate", svc.SetActiveURL)

	routes.Run(c.ClientPort)
}
