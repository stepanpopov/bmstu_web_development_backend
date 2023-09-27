package main

import (
	"context"
	"log"
	"rip/internal/config"
	"rip/internal/pkg/api"
	"rip/internal/pkg/repo/mock"
)

func main() {
	conf, err := config.NewConfig(context.Background())
	if err != nil {
		log.Fatal("Config error:", err)
	}

	log.Print(conf.ServiceHost, conf.ServicePort)

	repo := mock.New()

	serv := api.NewServer(
		api.WithHost(conf.ServiceHost),
		api.WithPort(conf.ServicePort),
	)

	serv.StartServer(repo)
}
