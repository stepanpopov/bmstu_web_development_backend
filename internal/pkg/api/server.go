package api

import (
	"fmt"
	"log"

	"rip/internal/pkg/repo"

	"github.com/gin-gonic/gin"
)

const (
	creatorID   = 0
	moderatorID = 1
)

type Server struct {
	host string
	port int
}

func WithHost(host string) func(*Server) {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port int) func(*Server) {
	return func(s *Server) {
		s.port = port
	}
}

func NewServer(options ...func(*Server)) *Server {
	srv := &Server{}
	for _, o := range options {
		o(srv)
	}
	return srv
}

// TODO: add abort() to errors
func (s *Server) StartServer(rep repo.Repository) {
	log.Println("Server start up")

	r := gin.Default()

	dataService := r.Group("/dataService")
	dataService.GET("/", getDataService(rep))
	dataService.GET("/:id", getDataServiceByID(rep))
	dataService.PUT("/", createDataService(rep))
	dataService.DELETE("/:id", deleteDataService(rep))
	dataService.POST("/", updateDataService(rep))
	dataService.PUT("/draft/:id", addToDraft(rep))

	encDecRequest := r.Group("/encryptDecryptReques")
	encDecRequest.POST("/filter", getEncryptDecryptRequests(rep))

	r.Run(fmt.Sprintf("%s:%d", s.host, s.port))
	log.Println("Server down")
}
