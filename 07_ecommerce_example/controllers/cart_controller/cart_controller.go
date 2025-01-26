package cart_controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"shirts-shop-golang/config"
	"shirts-shop-golang/helpers"
	"shirts-shop-golang/models/cart"
	"shirts-shop-golang/models/product"

	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)

	if data["Cart"] != nil {
		cartSession := data["Cart"].([]cart.Cart)
		data["CartNotEmpty"] = len(cartSession) > 0
	}

	return c.Render("templates/cart/view", data, "templates/layouts/main")
}

func Add(c *fiber.Ctx, db *sql.DB) error {
	var cartItems []cart.Cart

	data := helpers.GetSessionAndFlashMessages(c)
	idProduct, err := c.ParamsInt("id")

	if err != nil {
		return c.Redirect("/")
	}

	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	if data["Cart"] != nil {
		cartItems = data["Cart"].([]cart.Cart)

		// funcion que determina si existe un articulo en carrito
		containItem := func() bool {
			// recorremos el carrito para ver si existe el producto
			// añadimos uno al carrito si lo encuentra

			for index, cart := range cartItems {
				if cart.Product.GetId() == int64(idProduct) {
					cartItems[index].Unidades++
					return true
				}
			}

			return false
		}

		if !containItem() {
			productModel := product.New(db)
			productModel.SetId(int64(idProduct))

			productDB, err := productModel.GetOne()

			if err != nil {
				log.Fatal(err)
			}

			cartItems = append(cartItems, cart.Cart{
				Unidades: 1,
				Product:  productDB,
			})
		}

	} else {
		// en caso de no poseer carrito
		// creamos el carrito, buscamos el producto
		productModel := product.New(db)
		productModel.SetId(int64(idProduct))

		productDB, err := productModel.GetOne()

		if err != nil {
			log.Fatal(err)
		}

		cartItems = make([]cart.Cart, 0)
		cartItems = append(cartItems, cart.Cart{
			Unidades: 1,
			Product:  productDB,
		})
	}

	// sustituir el valor y volver a agregar a la session
	// encodeamos a JSON para almacenar en session
	bytes, err := json.Marshal(cartItems)

	if err != nil {
		log.Fatal(err)
	}

	sess.Set("Cart", string(bytes))

	// guardamos la session
	if err = sess.Save(); err != nil {
		log.Fatal(err)
	}

	return c.Redirect("/cart_controller/")
}

func Remove(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)
	indexCart, err := c.ParamsInt("id")

	if err != nil {
		log.Fatal(err)
	}

	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	if data["Cart"] != nil {
		// obtenemos el slice borramos el elemento del indice y guardamos
		cartItems := data["Cart"].([]cart.Cart)

		// filtramos por el indice, crea una copia desde el punto
		// donde se encuentra el elemento y arranca en el siguiente elemento
		// por encima del seleccionado
		cartItems = append(cartItems[:indexCart], cartItems[indexCart+1:]...)

		// sustituir el valor y volver a agregar a la session
		// encodeamos a JSON para almacenar en session
		bytes, err := json.Marshal(cartItems)

		if err != nil {
			log.Fatal(err)
		}

		sess.Set("Cart", string(bytes))

		// guardamos la session
		if err = sess.Save(); err != nil {
			log.Fatal(err)
		}
	}

	return c.Redirect("/cart_controller/")
}

func DeleteAll(c *fiber.Ctx, db *sql.DB) error {
	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	cartSession := sess.Get("Cart")

	if cartSession != nil {
		sess.Delete("Cart")

		if err := sess.Save(); err != nil {
			log.Fatal(err)
		}
	}

	return c.Redirect("/cart_controller/")
}

// aumenta 1 el valor del carrito
func Up(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)
	indexCart, err := c.ParamsInt("id")

	if err != nil {
		log.Fatal(err)
	}

	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	if data["Cart"] != nil {
		// obtenemos el slice y aumentamos las unidades
		cartItems := data["Cart"].([]cart.Cart)
		cartItems[indexCart].Unidades++

		// sustituir el valor y volver a agregar a la session
		// encodeamos a JSON para almacenar en session
		bytes, err := json.Marshal(cartItems)

		if err != nil {
			log.Fatal(err)
		}

		sess.Set("Cart", string(bytes))

		// guardamos la session
		if err = sess.Save(); err != nil {
			log.Fatal(err)
		}
	}

	return c.Redirect("/cart_controller/")
}

// disminuye 1  el valor del carrito
func Down(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)
	indexCart, err := c.ParamsInt("id")

	if err != nil {
		log.Fatal(err)
	}

	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	if data["Cart"] != nil {
		// obtenemos el slice y aumentamos las unidades
		cartItems := data["Cart"].([]cart.Cart)
		cartItems[indexCart].Unidades--

		// si las unidad es 0 retiramos del carrito
		if cartItems[indexCart].Unidades == 0 {
			cartItems = append(cartItems[:indexCart], cartItems[indexCart+1:]...)
		}

		// sustituir el valor y volver a agregar a la session
		// encodeamos a JSON para almacenar en session
		bytes, err := json.Marshal(cartItems)

		if err != nil {
			log.Fatal(err)
		}

		sess.Set("Cart", string(bytes))

		// guardamos la session
		if err = sess.Save(); err != nil {
			log.Fatal(err)
		}
	}

	return c.Redirect("/cart_controller/")
}
