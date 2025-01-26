package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"
	"shirts-shop-golang/config"
	"shirts-shop-golang/controllers/cart_controller"
	"shirts-shop-golang/controllers/category_controller"
	"shirts-shop-golang/controllers/order_controller"
	"shirts-shop-golang/controllers/product_controller"
	"shirts-shop-golang/controllers/user_controller"
	"shirts-shop-golang/middleware"
	"shirts-shop-golang/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

//go:embed templates/*
var templatesFS embed.FS

// construye el metodo para las consultas
func makeHandler(callback func(*fiber.Ctx, *sql.DB) error, db *sql.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		callback(c, db)
		return nil
	}
}

func main() {
	db, err := config.CreateConectionDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	engine := html.NewFileSystem(http.FS(templatesFS), ".html")
	config.ConfigEngine(engine, db)

	// iniciolizamos las variables de entorno
	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// iniciamos el demonio del correo electronico
	services.InitDaemonMail()
	defer close(services.ChannelEmail) // cuando termine el proceso cierra el canal

	app := fiber.New(fiber.Config{
		Views:        engine,
		ErrorHandler: config.ErrorHandler, // funcion majedora de errores
	})

	app.Static("/", config.Getwd())
	app.Get("/", makeHandler(product_controller.Index, db))

	// modulo usuarios
	userRouter := app.Group("/user_controller")
	userRouter.
		Get("/", makeHandler(user_controller.Index, db)).
		Get("/register", makeHandler(user_controller.Create, db)).
		Post("/save", makeHandler(user_controller.Save, db)).
		Post("/login", makeHandler(user_controller.Login, db)).
		Get("/logout", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			makeHandler(user_controller.Logout, db)}...,
		)

	// modulo categorias
	categoriesRouter := app.Group("/category_controller")
	categoriesRouter.
		Get("/", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			middleware.AdminMiddleware,
			makeHandler(category_controller.Index, db)}...,
		).
		Get("/create", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			middleware.AdminMiddleware,
			makeHandler(category_controller.Create, db)}...,
		).
		Post("/save", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			middleware.AdminMiddleware,
			makeHandler(category_controller.Save, db)}...,
		).
		Get("/show/:id<int>", makeHandler(category_controller.Show, db))

	// modulo productos
	productsRouter := app.Group("/product_controller")
	productsRouter.
		Get("/", makeHandler(product_controller.Index, db)).
		Get("/manage", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			middleware.AdminMiddleware,
			makeHandler(product_controller.Manage, db)}...,
		).
		Get("/create", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			middleware.AdminMiddleware,
			makeHandler(product_controller.Create, db)}...).
		Post("/save", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			middleware.AdminMiddleware,
			makeHandler(product_controller.Save, db)}...,
		).
		Get("/edit/:id<int>", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			middleware.AdminMiddleware,
			makeHandler(product_controller.Edit, db)}...,
		).
		Post("/update/:id<int>", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			middleware.AdminMiddleware,
			makeHandler(product_controller.Update, db)}...,
		).
		Get("/delete/:id<int>", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			middleware.AdminMiddleware,
			makeHandler(product_controller.Delete, db)}...,
		).
		Get("/show/:id<int>", makeHandler(product_controller.Show, db))

	// modulo ordenes
	orderRouter := app.Group("/order_controller")
	orderRouter.
		Get("/create", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			makeHandler(order_controller.Create, db)}...,
		).
		Post("/save", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			makeHandler(order_controller.Save, db)}...,
		).
		Get("/confirm", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			makeHandler(order_controller.Confirm, db)}...,
		).
		Get("/my-orders", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			makeHandler(order_controller.MyOrders, db)}...,
		).
		Get("/view/:id<int>", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			makeHandler(order_controller.View, db)}...,
		).
		Get("/manage", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			middleware.AdminMiddleware,
			makeHandler(order_controller.Manage, db)}...,
		).
		Post("/state", []func(*fiber.Ctx) error{
			middleware.AuthMiddleware,
			middleware.AdminMiddleware,
			makeHandler(order_controller.State, db)}...,
		)

	// carrito modulo
	cartRouter := app.Group("/cart_controller")
	cartRouter.
		Get("/", makeHandler(cart_controller.Index, db)).
		Get("/add/:id<int>", makeHandler(cart_controller.Add, db)).
		Get("/delete_all", makeHandler(cart_controller.DeleteAll, db)).
		Get("/remove/:id<int>", makeHandler(cart_controller.Remove, db)).
		Get("/up/:id<int>", makeHandler(cart_controller.Up, db)).
		Get("/down/:id<int>", makeHandler(cart_controller.Down, db))

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
