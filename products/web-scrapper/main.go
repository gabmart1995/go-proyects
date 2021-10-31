package main 

import (
	"fmt"
	"example.com/web-scrapper/config"
)

func main() {
	fmt.Println("values: " + *config.GetUrl())
}