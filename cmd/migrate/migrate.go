package main

import (
	"fmt"
	"rip/internal/dsn"
	"rip/internal/pkg/repo"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(repo.All()...)
	if err != nil {
		fmt.Println("cant migrate db")
		panic(err)
	}
}
