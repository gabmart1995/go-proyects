package config

import (
	"flag"
)

type Option struct {
	Name string
	Value string
	Description string
	Flag *string
}


/* Opciones del CLI */

var options = []Option{
	{
		Name: "url",
		Value: "https://es.lipsum.com/",
		Description: "Url a solicitar para realizar web scrapping",
	},
	{
		Name: "element",
		Value: "html",
		Description: "Clase o elemento donde comienza la busqueda",
	},
}



func init() {

	for index, option := range options {
		options[index].Flag = flag.String( option.Name, option.Value, option.Description )
	}
	
	// se registan los flags dentro del CLI
	flag.Parse()
}


func GetUrl() string {

	// fmt.Printf(" %+v", options[0] )

	if options[0].Flag == nil {
		return "url is nil"
	}

	return *options[0].Flag
}


func GetElement() string {

	if options[0].Flag == nil {
		return "element is nil"
	}

	return *options[1].Flag
}