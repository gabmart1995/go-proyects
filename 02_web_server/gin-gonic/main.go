package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("local.env")

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// get env variable PORT
	PORT := os.Getenv("PORT")

	// render app
	server := getApp()

	// render a angular app
	// server := getAppAngular()

	// render a react app
	// server := getAppReact()

	log.Fatal(server.Run(":" + PORT))
}
