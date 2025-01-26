package product_controller

import (
	"database/sql"
	"fmt"
	"log"
	"shirts-shop-golang/config"
	"shirts-shop-golang/helpers"
	"shirts-shop-golang/models/category"
	"shirts-shop-golang/models/product"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)
	productModel := product.New(db)
	products, err := productModel.GetRandom(6)

	if err != nil {
		log.Fatal(err)
	}

	data["Products"] = products

	return c.Render("templates/index", data, "templates/layouts/main")
}

func Manage(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)

	productModel := product.New(db)
	products, err := productModel.GetAll()

	if err != nil {
		log.Fatal(err)
	}

	// establecemos el valor en la vista
	data["Products"] = products

	return c.Render("templates/product/manage", data, "templates/layouts/main")
}

func Create(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)

	// obtenemos las categorias
	categoryModel := category.New(db)
	categories, err := categoryModel.GetAll("ASC")

	if err != nil {
		log.Fatal(err)
	}

	data["Categories"] = categories

	return c.Render("templates/product/create", data, "templates/layouts/main")
}

func Save(c *fiber.Ctx, db *sql.DB) error {
	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	productModel := product.New(db)
	productModel.SetNombre(c.FormValue("nombre", ""))
	productModel.SetDescripcion(c.FormValue("descripcion", ""))

	price, _ := strconv.ParseFloat(c.FormValue("precio", "0"), 64)
	productModel.SetPrecio(price)

	stock, _ := strconv.ParseInt(c.FormValue("stock", "0"), 10, 64)
	productModel.SetStock(stock)

	categoryId, _ := strconv.ParseInt(c.FormValue("categoria_id", "0"), 10, 64)
	productModel.SetCategoryId(categoryId)

	// para obtener la imagen  parseamos el formualrio a formData
	file, _ := c.FormFile("imagen")

	errorsForm, err := productModel.Validate(
		make([]string, 0),
		"CreateProduct",
		file,
	)

	if err != nil {
		log.Fatal(err)
	}

	if errorsForm != nil {
		helpers.SetSessionMessages(sess, errorsForm, "messages")
		return c.Redirect("/product_controller/create")
	}

	// subir el archivo al servidor
	// comprobamos si llega el archivo
	if file != nil {
		if err := helpers.SaveImage(c, file); err != nil {
			log.Fatal(err)
		}

		productModel.SetImagen(file.Filename)
	}

	// guardamos en la base de datos
	if err = productModel.Save(); err != nil {
		helpers.SetSessionMessages(sess, fiber.Map{
			"DbError": err.Error(),
		}, "messages")

		return c.Redirect("/product_controller/create")
	}

	helpers.SetSessionMessages(sess, fiber.Map{
		"DbSuccess": "Producto agregado con exito",
	}, "messages")

	return c.Redirect("/product_controller/manage")
}

func Edit(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)
	idProduct, err := c.ParamsInt("id")

	if err != nil {
		return c.Redirect("/product_controller/manage")
	}
	// buscamos el producto
	productModel := product.New(db)
	productModel.SetId(int64(idProduct))

	productDB, err := productModel.GetOne()

	if err != nil {
		data["DbError"] = err.Error()

	} else {
		data["Product"] = productDB

	}

	// obtenemos las categorias
	categoryModel := category.New(db)
	categories, err := categoryModel.GetAll("ASC")

	if err != nil {
		log.Fatal(err)
	}

	data["Categories"] = categories
	data["Edit"] = true

	return c.Render("templates/product/create", data, "templates/layouts/main")
}

func Update(c *fiber.Ctx, db *sql.DB) error {
	idProduct, err := c.ParamsInt("id")

	if err != nil {
		return c.Redirect("/product_controller/manage")
	}

	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	productModel := product.New(db)
	productModel.SetId(int64(idProduct))

	productModel.SetNombre(c.FormValue("nombre", ""))
	productModel.SetDescripcion(c.FormValue("descripcion", ""))

	price, _ := strconv.ParseFloat(c.FormValue("precio", "0"), 64)
	productModel.SetPrecio(price)

	stock, _ := strconv.ParseInt(c.FormValue("stock", "0"), 10, 64)
	productModel.SetStock(stock)

	categoryId, _ := strconv.ParseInt(c.FormValue("categoria_id", "0"), 10, 64)
	productModel.SetCategoryId(categoryId)

	// para obtener la imagen  parseamos el formualrio a formData
	file, _ := c.FormFile("imagen")

	errorsForm, err := productModel.Validate(
		make([]string, 0),
		"CreateProduct",
		file,
	)

	if err != nil {
		log.Fatal(err)
	}

	if errorsForm != nil {
		helpers.SetSessionMessages(sess, errorsForm, "messages")
		return c.Redirect(fmt.Sprintf("/product_controller/edit/%d", idProduct))
	}

	// subir el archivo al servidor
	// comprobamos si llega el archivo
	if file != nil {
		if err := helpers.SaveImage(c, file); err != nil {
			log.Fatal(err)
		}

		productModel.SetImagen(file.Filename)
	}

	// guardamos en la base de datos
	if err = productModel.Update(); err != nil {
		helpers.SetSessionMessages(sess, fiber.Map{
			"DbError": err.Error(),
		}, "messages")

		return c.Redirect(fmt.Sprintf("/product_controller/edit/%d", idProduct))
	}

	helpers.SetSessionMessages(sess, fiber.Map{
		"DbSuccess": "Producto modificado con exito",
	}, "messages")

	return c.Redirect("/product_controller/manage")
}

func Delete(c *fiber.Ctx, db *sql.DB) error {
	idProduct, err := c.ParamsInt("id")

	if err != nil {
		log.Fatal(err)
	}

	sess, err := config.Store.Get(c)

	if err != nil {
		log.Fatal(err)
	}

	productModel := product.New(db)
	productModel.SetId(int64(idProduct))

	if err = productModel.Delete(); err != nil {
		helpers.SetSessionMessages(sess, fiber.Map{
			"DbError": err.Error(),
		}, "messages")

	} else {
		helpers.SetSessionMessages(sess, fiber.Map{
			"DbSuccess": "Producto borrado con exito",
		}, "messages")

	}

	return c.Redirect("/product_controller/manage")
}

func Show(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)

	idProduct, err := c.ParamsInt("id")

	if err != nil {
		log.Fatal(err)
	}

	productModel := product.New(db)
	productModel.SetId(int64(idProduct))

	productDB, err := productModel.GetOne()

	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}

		// en este punto no hay registros del producto

	} else {
		data["Product"] = productDB

	}

	return c.Render("templates/product/show", data, "templates/layouts/main")
}
