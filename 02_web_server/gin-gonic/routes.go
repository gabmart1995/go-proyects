package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getApp(router *gin.Engine) *gin.Engine {

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"name":  "Gabriel Martinez",
			"title": "Curso de GO",
		})
	})

	router.GET("/generic", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "generic.html", gin.H{
			"name":  "Gabriel Martinez",
			"title": "Curso de GO",
		})
	})

	router.GET("/elements", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "elements.html", gin.H{
			"name":  "Gabriel Martinez",
			"title": "Curso de GO",
		})
	})

	/* 404 handler */
	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "404-page.html", nil)
	})

	return router
}
