package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"todolist/database"
	"todolist/model"
	"todolist/module"

	_ "todolist/docs"

	"github.com/go-co-op/gocron/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
)

// @title			Todos API
// @version			1.0
// @description		API for managing todos.
// @termsOfService	http://swagger.io/terms/

// @contact.name	Jefri Herdi Triyanto
// @contact.url		https://jefripunza.com
// @contact.email	hi@jefripunza.com

// @license.name	Apache 2.0
// @license.url		http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath		/
func main() {
	database.PostgresConnect()

	// scheduler due date
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}
	j, err := s.NewJob(
		gocron.DurationJob(
			1*time.Second,
		),
		gocron.NewTask(
			func() {
				now := time.Now()
				// timezone to UTC
				loc, _ := time.LoadLocation("UTC")
				now = now.In(loc)
				fmt.Println(now.Format("2006-01-02 15:04:05"))

				// list all tasks
				var todos []model.Todo
				find := database.Postgres.Where("due_date <= ? AND status = 'pending'", now).Find(&todos)
				if find.Error != nil {
					log.Fatal(find.Error)
				}

				var wg sync.WaitGroup
				wg.Add(len(todos))
				for i := 1; i <= len(todos); i++ {
					go worker(&wg, todos[i-1])
				}
				wg.Wait() // tunggu semua goroutine selesai
			},
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Job ID: %s\n", j.ID())
	s.Start()

	app := fiber.New(fiber.Config{
		ServerHeader:  "Todos App",
		CaseSensitive: true,
		BodyLimit:     10 * 1024 * 1024, // 10 MB / max file size
	})

	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		ExposeHeaders: "Content-Type,Authorization,Accept,X-Browser-ID,X-Balance",
		AllowOrigins:  "*",
	}))
	app.Use(requestid.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
	}))

	app.Use(logger.New())

	module.Todo{}.Route(app)

	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "endpoint not found!",
		})
	})

	log.Fatal(app.Listen(":3000"))
}

func worker(wg *sync.WaitGroup, todo model.Todo) {
	defer wg.Done()

	// update status menjadi completed
	todo.Status = "completed"
	update := database.Postgres.Save(&todo)
	if update.Error != nil {
		log.Fatal(update.Error)
	}
	log.Printf("Task %s updated successfully\n", todo.Title)
}
