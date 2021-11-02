package main 

import (
	"fmt"
	"example.com/web-scrapper/config"
)

func main() {
	fmt.Println( "url to scrapping: " + config.GetUrl() )
	fmt.Println( "element: " + config.GetElement() )
	fmt.Println("scrapping ...")
}