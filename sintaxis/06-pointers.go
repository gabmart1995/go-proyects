package main

import (
	"fmt"
)


func main() {

	/*
		1. Un puntero es una direccion en memoria, 
		2. En lugar del valor tienes la direccion en la que se enceuntra el valor.
	*/

	var x, y *int

	integer := 5

	// accede al valor de la direccion en memoria 
	x = &integer
	y = &integer

	// el asterisco el valor imprime el valor almacenado en la memorias
	*x = 6

	fmt.Println( *x )
	fmt.Println( *y )
}