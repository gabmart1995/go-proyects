package config

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"shirts-shop-golang/filters"
	"shirts-shop-golang/models/category"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

var Store *session.Store = session.New(session.Config{
	Expiration: time.Hour,
})

// funcion en Prod obtiene la ruta real del ejecutable
func Getwd() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "static")
}

// anade funciones personalizadas al motor de plantillas
func ConfigEngine(engine *html.Engine, db *sql.DB) {
	engine.AddFunc("getYear", filters.GetYear)
	engine.AddFunc("unescape", filters.UnescapeHTML)
	engine.AddFunc("isAdmin", filters.IsAdmin)
	engine.AddFunc("countCart", filters.CountCart)
	engine.AddFunc("getTotalCart", filters.GetTotalCart)
	engine.AddFunc("getStatus", filters.GetStatus)
	engine.AddFunc("getCategories", func() []category.Category {
		categories := make([]category.Category, 0)
		query := "SELECT * FROM categorias ORDER BY id DESC LIMIT 10;"
		rows, err := db.Query(query)

		if err != nil {
			return categories
		}

		defer rows.Close()

		// realizamos la lectura
		for rows.Next() {
			var category category.Category

			err := rows.Scan(&category.Id, &category.Nombre)

			if err != nil {
				log.Fatal(err)
			}

			categories = append(categories, category)
		}

		return categories
	})
}

func CreateConectionDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "test:123456@tcp(localhost:3306)/tienda_master")

	if err != nil {
		return nil, err
	}

	// establece los caracteres latinos
	db.Query("SET NAMES 'utf-8'")

	return db, nil
}

// manejador de errores
func ErrorHandler(c *fiber.Ctx, err error) error {

	// obtenemos el codigo de error
	var errorFiber *fiber.Error
	code := fiber.StatusInternalServerError

	if errors.As(err, &errorFiber) {
		code = errorFiber.Code
	}

	// manejo de errores 404
	if code == fiber.StatusNotFound {
		return c.
			Status(code).
			Render("templates/404", nil, "templates/layouts/main")
	}

	return nil
}
