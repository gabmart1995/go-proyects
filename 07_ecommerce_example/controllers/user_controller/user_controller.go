package user_controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"shirts-shop-golang/config"
	"shirts-shop-golang/helpers"
	"shirts-shop-golang/models/user"
	"shirts-shop-golang/services"

	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)

	return c.Render("templates/index", data, "templates/layouts/main")
}

func Create(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)

	return c.Render("templates/user/register", data, "templates/layouts/main")
}

func Save(c *fiber.Ctx, db *sql.DB) error {
	var data fiber.Map
	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	userModel := user.New(db)
	userModel.SetNombre(c.FormValue("name"))
	userModel.SetApellido(c.FormValue("surname"))
	userModel.SetEmail(c.FormValue("email"))
	userModel.SetPassword(c.FormValue("password"), false)

	// evalua todos los campos mandamos un slice
	// vacio para que evalue todos los campos
	errorsForm, err := userModel.Validate(
		make([]string, 0),
		"Register",
	)

	// evalua si hay un error critico
	if err != nil {
		log.Fatal(err)
	}

	if errorsForm != nil {
		helpers.SetSessionMessages(sess, errorsForm, "messages")
		return c.Redirect("/user_controller/register")
	}

	// ciframos la contrasena antes de mandar a BD
	userModel.SetPassword(userModel.GetPassword(), true)

	// en este punto guardamos en la base de datos
	if err = userModel.Save(); err != nil {
		helpers.SetSessionMessages(sess, fiber.Map{
			"DbError": err.Error(),
		}, "messages")
		return c.Redirect("/user_controller/register")
	}

	// usuario insertado en BD
	data = fiber.Map{
		"DbSuccess": "Usuario agregado con exito",
	}

	helpers.SetSessionMessages(sess, data, "messages")

	// mandamos el correo de registro del cliente
	services.SendCreateUser(userModel.GetEmail())

	// redireccionamos a la vista
	return c.Redirect("/user_controller/register")
}

func Login(c *fiber.Ctx, db *sql.DB) error {
	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	user := user.New(db)

	// obtener datos del formulario
	user.SetEmail(c.FormValue("email"))
	user.SetPassword(c.FormValue("password"), false)

	// validar campos
	errorsForm, err := user.Validate(
		[]string{"Email", "Password"},
		"Login",
	)

	if err != nil {
		log.Fatal(err)
	}

	if errorsForm != nil {
		helpers.SetSessionMessages(sess, errorsForm, "messages")
		return c.Redirect("/")
	}

	// consultar la base de datos para comprobar las credenciales
	userLogged, err := user.Login()

	if err != nil {
		helpers.SetSessionMessages(sess, fiber.Map{
			"LoginError": "Identificacion fallida",
		}, "messages")

		return c.Redirect("/")
	}

	// establecemos los valores de la struct en json string
	bytes, err := json.Marshal(userLogged)
	if err != nil {
		log.Fatal(err)
	}

	sess.Set("identity", string(bytes))

	// guardamos los datos de la session
	if err = sess.Save(); err != nil {
		log.Fatal(err)
	}

	return c.Redirect("/")
}

func Logout(c *fiber.Ctx, db *sql.DB) error {
	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	// comprueba la identidad del usuario
	// si existe lo elimina
	if sess.Get("identity") != nil {
		if err = sess.Destroy(); err != nil {
			log.Fatal(err)
		}
	}

	return c.Redirect("/")
}
