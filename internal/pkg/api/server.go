package api

import (
	"fmt"
	"log"

	"rip/internal/pkg/repo"

	"github.com/gin-gonic/gin"
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

func (s *Server) StartServer(rep repo.Repository) {
	log.Println("Server start up")

	r := gin.Default()

	r.LoadHTMLGlob("static/template/*")
	r.Static("/static", "./static")
	r.Static("/css", "./static")
	r.Static("/img", "./static")

	r.GET("/", filterDataService(rep))

	r.GET("/service/*id", showDataService(rep))

	r.Run(fmt.Sprintf("%s:%d", s.host, s.port))

	log.Println("Server down")
}
