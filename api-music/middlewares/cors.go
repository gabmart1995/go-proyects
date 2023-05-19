package middlewares

import "github.com/gofiber/fiber/v2"

func CORSMiddleware(c *fiber.Ctx) error {

	c.Request().Header.Set("Access-Control-Allow-Origin", "*")
	c.Request().Header.Set("Access-Control-Allow-Credentials", "true")
	c.Request().Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Request().Header.Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	if c.Method() == "OPTIONS" {
		return c.SendStatus(204)
	}

	return c.Next()
}
