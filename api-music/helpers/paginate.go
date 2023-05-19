package helpers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/* -- decorador de paginacion -- */
func Paginate(c *fiber.Ctx, param string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		var query string

		if query = c.Query(param); len(query) == 0 {
			query = "0"
		}

		start, err := strconv.Atoi(query)

		if err != nil {
			log.Fatalln(err)
		}

		return db.Offset(start).Limit(10)
	}
}
