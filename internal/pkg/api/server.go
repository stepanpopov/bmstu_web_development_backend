package api

import (
	"fmt"
	"log"
	"time"

	"rip/internal/pkg/repo"

	"github.com/gin-gonic/gin"
)

type JWTConfig struct {
	ExpiresIn time.Duration
	Secret    string
}

type Server struct {
	host      string
	port      int
	jwtConfig JWTConfig
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

func WithJWTConfig(c JWTConfig) func(*Server) {
	return func(s *Server) {
		s.jwtConfig = c
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
func (s *Server) StartServer(rep repo.Repository, avatar repo.Avatar) {
	log.Println("Server start up")

	r := gin.Default()

	dataService := r.Group("/dataService")
	dataService.GET("/", filterDataService(rep))
	dataService.GET("/:id", getDataServiceByID(rep))
	dataService.POST("/", createDataService(rep))
	dataService.DELETE("/:id", deleteDataService(rep, avatar))
	dataService.PUT("/", updateDataService(rep))
	dataService.POST("/:id/image", putImage(rep, avatar))

	dataService.POST("/draft/:id", addToDraft(rep))
	dataService.DELETE("/draft/:id", deleteFromDraft(rep)) //

	encDecRequest := r.Group("/encryptDecryptRequest")
	encDecRequest.GET("/filter", getEncryptDecryptRequests(rep))
	encDecRequest.GET("/:id", getEncryptDecryptRequestsByID(rep))
	encDecRequest.POST("/", createDraft(rep))
	encDecRequest.PUT("/form/:id", formEncryptDecryptRequest(rep))
	encDecRequest.PUT("/update_moderator/:id", updateModeratorEncryptDecryptRequest(rep))
	encDecRequest.DELETE("/:req_id", deleteEncryptDecryptRequest(rep))
	encDecRequest.DELETE("/:req_id/delete/:data_id", deleteDataFromEncryptDecryptRequest(rep))

	auth := r.Group("/auth")
	auth.POST("/login", login(rep, s.jwtConfig.Secret, s.jwtConfig.ExpiresIn))
	auth.POST("register", register(rep))

	// удаление услуги из заявки + мб тогда delete draft не нужен
	// TODO: get draft???

	r.Run(fmt.Sprintf("%s:%d", s.host, s.port))
	log.Println("Server down")
}
