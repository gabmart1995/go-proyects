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
	DEBUG := os.Getenv("DEBUG") == "true"

	// modo produccion
	if !DEBUG {
		gin.SetMode(gin.ReleaseMode)
	}

	// render app
	server := getApp()

	// render a angular app
	// server := getAppAngular()

	// render a react app
	// server := getAppReact()

	// configuracion de proxy para produccion
	if !DEBUG {
		server.SetTrustedProxies(nil)
	}

	log.Fatal(server.Run(":" + PORT))
}
