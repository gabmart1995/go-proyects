package config

import (
	"flag"
)

type Option struct {
	name string
	value string
	description string
}

/* Opciones del CLI */
// array
var options = []Option{
	{
		name: "url",
		value: "https://es.lipsum.com/",
		description: "Url a para realizar web scrapping",
	},
}

// slices
var Flags []*string

func init() {

	// la funcion del init se ejecuta cuanto se importa el archivo se ejecuta una sola vez
	for _, option := range options  {
		Flags = append( Flags, flag.String( option.name, option.value, option.description, )) 	
	}

	
	// se registan los flags dentro del CLI
	flag.Parse()
}

func GetUrl() *string {
	return Flags[0]
}
