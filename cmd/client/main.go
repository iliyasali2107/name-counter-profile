package main

import (
	"log"

	"name-counter-url/pkg/config"

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

	routes.GET("/url/:id", svc.GetURL)

	routes.Run(c.ClientPort)
}
