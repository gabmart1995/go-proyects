package category_controller

import (
	"database/sql"
	"log"
	"shirts-shop-golang/config"
	"shirts-shop-golang/helpers"
	"shirts-shop-golang/models/category"
	"shirts-shop-golang/models/product"

	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)

	// consultamos las categorias
	category := category.New(db)
	categories, err := category.GetAll("DESC")

	if err != nil {
		log.Fatal(err)
	}

	// pasamos las categorias por propiedad
	data["Categories"] = categories

	return c.Render("templates/category/index", data, "templates/layouts/main")
}

func Create(c *fiber.Ctx, db *sql.DB) error {
	data := helpers.GetSessionAndFlashMessages(c)
	return c.Render("templates/category/create", data, "templates/layouts/main")
}

func Save(c *fiber.Ctx, db *sql.DB) error {
	sess, _ := config.Store.Get(c)
	category := category.New(db)
	category.SetNombre(c.FormValue("name"))

	errorsForm, err := category.Validate(make([]string, 0), "CategoryCreate")

	if err != nil {
		panic(err)
	}

	if errorsForm != nil {
		helpers.SetSessionMessages(sess, errorsForm, "messages")
		return c.Redirect("/category_controller/create")
	}

	// almacenamos la categoria
	if err = category.Save(); err != nil {
		helpers.SetSessionMessages(sess, fiber.Map{
			"DbError": err.Error(),
		}, "errors")
		return c.Redirect("/category_controller/create")
	}

	return c.Redirect("/category_controller/")
}

func Show(c *fiber.Ctx, db *sql.DB) error {
	idCategory, err := c.ParamsInt("id")

	if err != nil {
		log.Fatal(err)
	}

	data := helpers.GetSessionAndFlashMessages(c)

	// obtenemos la categoria
	categoryModel := category.New(db)
	categoryModel.SetId(int64(idCategory))

	categoryDB, err := categoryModel.GetOne()

	if err != nil {
		log.Fatal(err)
	}

	data["Category"] = categoryDB

	// obtenemos los productos
	productModel := product.New(db)
	productModel.SetCategoryId(int64(idCategory))

	products, err := productModel.GetAllCategory()

	if err != nil {
		log.Fatal(err)
	}

	data["Products"] = products

	return c.Render("templates/category/show", data, "templates/layouts/main")
}
