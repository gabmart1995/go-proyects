package order_controller

import (
	"database/sql"
	"fmt"
	"log"
	"shirts-shop-golang/config"
	"shirts-shop-golang/filters"
	"shirts-shop-golang/helpers"
	"shirts-shop-golang/models/cart"
	"shirts-shop-golang/models/order"
	"shirts-shop-golang/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Create(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)
	return c.Render("templates/order/create", data, "templates/layouts/main")
}

func Save(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)
	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	orderModel := order.New(db)
	orderModel.SetProvincia(c.FormValue("provincia", ""))
	orderModel.SetLocalidad(c.FormValue("localidad", ""))
	orderModel.SetDireccion(c.FormValue("direccion", ""))
	orderModel.SetUsuarioId(data["UserId"].(int64))
	orderModel.SetCart(data["Cart"].([]cart.Cart))

	cost, err := strconv.ParseFloat(
		filters.GetTotalCart(data["Cart"].([]cart.Cart)),
		64,
	)

	if err != nil {
		log.Fatal(err)
	}

	// establece el coste
	orderModel.SetCoste(cost)
	errorsForm, err := orderModel.Validate(make([]string, 0), "Order")

	if err != nil {
		log.Fatal(err)
	}

	if errorsForm != nil {
		helpers.SetSessionMessages(sess, errorsForm, "messages")

		return c.Redirect("/order_controller/create")
	}

	// guardamos en BD
	if err = orderModel.Save(); err != nil {
		helpers.SetSessionMessages(sess, fiber.Map{
			"DbError": err.Error(),
		}, "messages")

		return c.Redirect("/order_controller/create")
	}

	// salvamos la linea del pedido
	if err = orderModel.SaveOrderLine(); err != nil {
		helpers.SetSessionMessages(sess, fiber.Map{
			"DbError": err.Error(),
		}, "messages")

		return c.Redirect("/order_controller/create")
	}

	// mandamos el correo electronico
	services.SendCreateOrderSucess(
		data["UserEmail"].(string),
		orderModel,
	)

	// continuamos el renderizado
	helpers.SetSessionMessages(sess, fiber.Map{
		"OrderConfirm": "confirm",
	}, "messages")

	return c.Redirect("/order_controller/confirm")
}

func Confirm(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)

	// buscamos el ultimo pedido del usuario identificado
	orderModel := order.New(db)
	orderModel.SetUsuarioId(data["UserId"].(int64))

	orderDB, err := orderModel.GetOneByUser()

	if err != nil {
		log.Fatal(err)
	}

	orderModel.SetId(orderDB.GetId())
	cartItems, err := orderModel.GetProductsByOrder()

	if err != nil {
		log.Fatal(err)

	}

	orderDB.Cart = cartItems

	data["Order"] = orderDB

	return c.Render("templates/order/confirm", data, "templates/layouts/main")
}

func MyOrders(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)

	orderModel := order.New(db)
	orderModel.SetUsuarioId(data["UserId"].(int64))

	// sacamos los pedidos del usuario
	ordersDB, err := orderModel.GetAllByUser()

	if err != nil {
		log.Fatal(err)
	}

	data["Orders"] = ordersDB

	return c.Render("templates/order/my-orders", data, "templates/layouts/main")
}

func View(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)

	idOrder, err := c.ParamsInt("id")

	if err != nil {
		return c.Redirect("/order_controller/my-orders")
	}

	orderModel := order.New(db)
	orderModel.SetId(int64(idOrder))

	// sacar el pedido
	orderDB, err := orderModel.GetOne()

	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}

		return c.Render("templates/order/view", data, "templates/layouts/main")
	}

	// sacar los productos por orden
	cartDB, err := orderModel.GetProductsByOrder()

	if err != nil {
		log.Fatal(err)
	}

	// asignamos los productos dentro de la orden
	orderDB.Cart = cartDB

	data["Order"] = orderDB

	return c.Render("templates/order/view", data, "templates/layouts/main")
}

func Manage(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)

	orderModel := order.New(db)
	orderModel.SetUsuarioId(data["UserId"].(int64))

	ordersDB, err := orderModel.GetAll()

	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalln(err)
		}

	} else {
		data["Orders"] = ordersDB
		data["Manage"] = true
	}

	return c.Render("templates/order/my-orders", data, "templates/layouts/main")
}

func State(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)
	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	idOrder, err := strconv.Atoi(c.FormValue("id", "1"))

	if err != nil {
		log.Fatal(err)
	}

	orderModel := order.New(db)
	orderModel.SetId(int64(idOrder))
	orderModel.SetEstado(c.FormValue("estado", "confirm"))

	if err = orderModel.UpdateState(); err != nil {
		helpers.SetSessionMessages(sess, fiber.Map{
			"DbError": err.Error(),
		}, "messages")

		return c.Redirect(fmt.Sprintf("/order_controller/view/%d", idOrder))
	}

	// enviar el correo de actuaizacion
	services.SendUpdateOrderSucess(data["UserEmail"].(string), orderModel)

	return c.Redirect(fmt.Sprintf("/order_controller/view/%d", idOrder))
}
