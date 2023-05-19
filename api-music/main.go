package main

import (
	"api-music/controllers/user"
	"api-music/middlewares"
	"api-music/migrations"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDatabase() *gorm.DB {
	dsn := "host=localhost user=test password=123456 dbname=music-api port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Conexión BD activa")

	migrations.Execute(db)

	return db
}

func main() {

	db := connectDatabase()
	sqlInstance, err := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	// cierra la conexion cuando finaliza la app
	defer sqlInstance.Close()

	app := fiber.New()

	// middlewares
	app.Use(middlewares.CORSMiddleware)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	api := app.Group("/api")

	api.Get("/user", makeHandler(user.GetAll, db))
	api.Get("/user/:id", makeHandler(user.Get, db))
	api.Post("/user", makeHandler(user.Create, db))
	api.Put("/user/:id", makeHandler(user.Update, db))
	api.Delete("/user/:id", makeHandler(user.Delete, db))

	// api.Post("/album")
	/*api.Get("/song/:id", song.Get)
	api.Get("/song", song.GetAll)
	api.Post("/song", song.Create)
	api.Put("/song/:id", song.Update)*/

	log.Fatal(app.Listen(":3000"))
}

// crea un manejador con el apuntador de la base de datos
func makeHandler(callback func(*fiber.Ctx, *gorm.DB) error, db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return callback(c, db)
	}
}
