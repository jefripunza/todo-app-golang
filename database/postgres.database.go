package database

import (
	"log"
	"todolist/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Postgres *gorm.DB

func PostgresConnect() {
	dbURL := "postgres://postgres:postgres@localhost:5432/todos"

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&model.Todo{})

	log.Printf("✅ Database connected successfully")

	Postgres = db
}
