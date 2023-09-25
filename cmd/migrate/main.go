package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"RIP/internal/app/ds"
	"RIP/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	env, err2 := dsn.FromEnv()
	if err2 != nil {
		panic("Error from reading env")
	}
	db, err := gorm.Open(postgres.Open(env), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.AutoMigrate(
		&ds.Application{},
		&ds.Moderator{},
		&ds.Tenders{},
		&ds.TenderApplication{},
	); err != nil {
		panic("cant migrate db:" + err.Error())
	}
}
