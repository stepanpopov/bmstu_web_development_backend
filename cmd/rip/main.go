package main

import (
	"context"
	"log"
	"rip/internal/config"
	"rip/internal/dsn"
	"rip/internal/pkg/api"
	repo "rip/internal/pkg/repo/gorm"
)

func main() {
	conf, err := config.NewConfig(context.Background())
	if err != nil {
		log.Fatal("Config error:", err)
	}

	log.Print(conf.ServiceHost, conf.ServicePort)

	repo, err := repo.NewPostgres(dsn.FromEnv())
	if err != nil {
		log.Fatal("DB connect error:", err)
	}

	serv := api.NewServer(
		api.WithHost(conf.ServiceHost),
		api.WithPort(conf.ServicePort),
	)

	serv.StartServer(repo)
}
