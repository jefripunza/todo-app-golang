package database

import (
	"log"
	"os"
	"todolist/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Postgres *gorm.DB

func PostgresConnect() {
	dbURL := "postgres://postgres:postgres@localhost:5432/todos"
	envDbUrl := os.Getenv("DB_URL")
	if envDbUrl != "" {
		dbURL = envDbUrl
	}

	// loc, _ := time.LoadLocation("Asia/Jakarta")
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&model.Todo{})

	log.Printf("✅ Database connected successfully")

	Postgres = db
}
