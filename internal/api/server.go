package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Service struct {
	ID   int
	Data string
}

func StartServer() {
	log.Println("Server start up")

	services := []Service{
		Service{
			ID:   1,
			Data: "AAAA",
		},
		Service{
			ID:   2,
			Data: "VVVV",
		},
	}

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", showAllServices(services))

	r.POST("/filter", filterServices(services))

	r.GET("/service/*id", showService(services))

	r.Static("/image", "./resources")

	r.Run()

	log.Println("Server down")
}
