package database

import (
	"fmt"
	"log"
	"time"
	"todolist/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Postgres *gorm.DB

func PostgresConnect() {
	dbURL := "postgres://postgres:postgres@localhost:5432/todos"

	// loc, _ := time.LoadLocation("Asia/Jakarta")
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			fmt.Println("Asia/Jakarta")
			ti, _ := time.LoadLocation("Asia/Jakarta")
			return time.Now().In(ti)
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&model.Todo{})

	log.Printf("✅ Database connected successfully")

	Postgres = db
}
