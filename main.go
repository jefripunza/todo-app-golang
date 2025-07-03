package main

import (
	"log"
	"todolist/database"
	"todolist/module"

	_ "todolist/docs"

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
