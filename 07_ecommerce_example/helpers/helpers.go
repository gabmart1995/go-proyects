package helpers

import (
	"encoding/json"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"shirts-shop-golang/config"
	"shirts-shop-golang/models/cart"
	"shirts-shop-golang/models/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func GetSessionAndFlashMessages(c *fiber.Ctx) fiber.Map {
	data := fiber.Map{}

	// obtenemos el valor de la session
	sess, err := config.Store.Get(c)
	if err != nil {
		log.Fatal(err)
	}

	// verificamos los mensajes flash
	hasMessages := sess.Get("messages") != nil
	if hasMessages {
		messages := make(map[string]string)
		jsonStringMessages := sess.Get("messages").(string)

		err := json.Unmarshal([]byte(jsonStringMessages), &messages)
		if err != nil {
			log.Fatal(err)
		}

		// establece los errores
		for key, value := range messages {
			data[key] = value
		}

		sess.Delete("messages")
	}

	// extraemos los datos del usuario autenticado de la sesion
	// establecemos los nuevos campos
	// validamos la identidad del usuario
	data["IsLogged"] = sess.Get("identity") != nil
	if data["IsLogged"].(bool) {
		var user user.User
		jsonStringUserLogged := sess.Get("identity").(string)

		err := json.Unmarshal([]byte(jsonStringUserLogged), &user)
		if err != nil {
			log.Fatal(err)
		}

		data["UserId"] = user.GetId()
		data["UserName"] = user.GetNombre()
		data["UserSurname"] = user.GetApellido()
		data["UserRol"] = user.GetRol()
		data["UserEmail"] = user.GetEmail()
	}

	// si existe el carrito lo añadimos a data para mostrarlo en vista
	hasCart := sess.Get("Cart") != nil
	if hasCart {
		var cartItems []cart.Cart

		jsonString := sess.Get("Cart").(string)

		if err := json.Unmarshal([]byte(jsonString), &cartItems); err != nil {
			log.Fatal(err)
		}

		data["Cart"] = cartItems
	}

	if err = sess.Save(); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	return data
}

// establece los mensajes de errores en la interfaz
func SetSessionMessages(sess *session.Session, errorsField fiber.Map, nameFieldSession string) error {
	bytes, err := json.Marshal(errorsField)

	if err != nil {
		return err
	}

	// establecemos los datos de la sesion
	sess.Set(nameFieldSession, string(bytes))

	if err = sess.Save(); err != nil {
		return err
	}

	return nil
}

func SaveImage(c *fiber.Ctx, file *multipart.FileHeader) error {
	var filePath string

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// comprueba si existe el directorio
	filePath = filepath.Join(dir, "static", "uploads", "images")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {

		// crea directorio en forma recursiva
		if err = os.MkdirAll(filePath, 0777); err != nil {
			return err
		}
	}

	filePath = filepath.Join(filePath, file.Filename)

	// salvamos el archivo
	if err := c.SaveFile(file, filePath); err != nil {
		return err
	}

	return nil
}
