package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// load html
	router.LoadHTMLGlob("./public/templates/**/*.html")

	// load static simple
	router.Static("/assets", "public/static") // css and js
	router.Static("/images", "public/images") // images

	// router.StaticFS("/static", http.Dir("./public/static")) FS
	server := getApp(router)

	log.Fatal(server.Run(":8080"))
}
