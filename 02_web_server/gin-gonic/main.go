package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("local.env")

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// get env variable PORT
	PORT := os.Getenv("PORT")

	// gin server
	router := gin.Default()

	// load html
	router.LoadHTMLGlob("./public/templates/**/*.html")

	// load static simple
	router.Static("/assets", "public/static") // css and js
	router.Static("/images", "public/images") // images

	// router.StaticFS("/static", http.Dir("./public/static")) FS

	server := getApp(router)

	log.Fatal(server.Run(":" + PORT))
}
