package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/** get a instance a server */
func getApp() *gin.Engine {

	// gin server
	router := gin.Default()
	response := gin.H{
		"name":  "Gabriel Martinez",
		"title": "Curso de GO",
	}

	// load html
	router.LoadHTMLGlob("./public/templates/**/*.html")

	// load static simple
	router.Static("/assets", "public/static") // css and js
	router.Static("/images", "public/images") // images

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", response)
	})

	router.GET("/generic", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "generic.html", response)
	})

	router.GET("/elements", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "elements.html", response)
	})

	/* 404 handler */
	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "404-page.html", nil)
	})

	return router
}

func getAppAngular() *gin.Engine {
	router := gin.Default()

	router.LoadHTMLFiles("angular-app/index.html")
	router.Static("/static", "angular-app/static")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	return router
}

func getAppReact() *gin.Engine {
	router := gin.Default()

	router.LoadHTMLFiles("react-app/index.html")
	router.Static("/static", "react-app/static")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	return router
}
