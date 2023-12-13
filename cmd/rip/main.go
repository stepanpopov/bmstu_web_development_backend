package main

import (
	"context"
	"log"
	"rip/internal/config"
	"rip/internal/dsn"
	"rip/internal/pkg/api"
	"rip/internal/pkg/redis"
	repo "rip/internal/pkg/repo/gorm"
	"time"

	"rip/internal/pkg/repo/s3"
	"rip/internal/pkg/s3/minio"
)

func main() {
	ctx := context.Background()

	conf, err := config.NewConfig(ctx)
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
		api.WithJWTConfig(api.JWTConfig{Secret: conf.JWT.Secret, ExpiresIn: conf.JWT.ExpiresIn}),
	)

	// TODO: env and config
	minioCl, err := minio.MakeS3MinioClient("localhost:9000", "minio", "minio124")
	if err != nil {
		log.Fatal("S3 connect error:", err)
	}
	minioCl.IsOnline()

	cancelHC, err := minioCl.HealthCheck(time.Second * 3)
	if err != nil {
		log.Fatal("S3 health check error:", err)
	}

	t := time.NewTimer(time.Second * 20)

loop:
	for {
		select {
		case <-t.C:
			log.Fatal("S3 health check failed")
		default:
			if minioCl.IsOnline() {
				cancelHC()
				break loop
			}
		}
	}

	avatar := s3.NewS3MinioAvatarSaver("avatars", minioCl)

	redis, err := redis.New(ctx, conf.Redis)
	if err != nil {
		log.Fatal(err)
	}

	serv.StartServer(repo, avatar, redis)
}
