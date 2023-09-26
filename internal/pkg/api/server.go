package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type PlainData struct {
	ID   int
	Name string
	Blob string
}

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

func(s *Server) StartServer() {
	log.Println("Server start up")

	PlainDatas := []PlainData{
		PlainData{
			ID:   1,
			Name: "Encode your secrets",
			Blob: "secret",
		},
		PlainData{
			ID:   2,
			Name: "Decode your life",
			Blob: "01001001000100",
		},
	}

	r := gin.Default()

	/* r.SetFuncMap(template.FuncMap{
		"contains": strings.Contains,
	}) */

	r.LoadHTMLGlob("templates/*")

	r.GET("/", showAllPlainData(PlainDatas))

	r.GET("/filter", filterPlainDatas(PlainDatas))

	r.GET("/service/*id", showPlainData(PlainDatas))

	r.Static("/image", "./resources")

	r.Run(fmt.Sprintf("%s:%d", s.host, s.port))

	log.Println("Server down")
}
