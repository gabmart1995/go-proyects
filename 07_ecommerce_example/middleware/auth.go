package middleware

import (
	"encoding/json"
	"log"
	"shirts-shop-golang/config"
	"shirts-shop-golang/filters"
	"shirts-shop-golang/models/user"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	data := fiber.Map{}
	data["IsLogged"] = sess.Get("identity") != nil

	// redirecciona si no es admin o no esta autenticado
	if !data["IsLogged"].(bool) {
		return c.Redirect("/")
	}

	return c.Next()
}

func AdminMiddleware(c *fiber.Ctx) error {
	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	data := fiber.Map{}
	data["IsLogged"] = sess.Get("identity")
	var user user.User
	jsonStringUserLogged := sess.Get("identity").(string)

	err = json.Unmarshal([]byte(jsonStringUserLogged), &user)
	if err != nil {
		log.Fatal(err)
	}

	data["UserRol"] = user.GetRol()

	// redirecciona si no es admin o no esta autenticado
	if !filters.IsAdmin(data["UserRol"].(string)) {
		return c.Redirect("/")
	}

	return c.Next()
}
