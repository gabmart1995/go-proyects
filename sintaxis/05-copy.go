package main

import (
	"fmt"
)

func main() {

	/*
		Copy permite pasar datos entre slices

		Nota: copia el minimo de elemetos si la copia del make
		vale 0 copia esa cantidad, ojo con eso, para solicionarlo
		puedes pasar la longitud de la slice fuente
	*/

	slice := []int{1, 2, 3, 4, 5, 6}

	// se duplica la capacidad del array
	sliceCopy := make( []int, len( slice ), cap( slice ) * 2 )

	copy( sliceCopy, slice )

	sliceCopy = append( sliceCopy, 30 )

	fmt.Println( slice )
	fmt.Println( sliceCopy )
}