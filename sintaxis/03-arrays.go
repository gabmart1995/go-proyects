package main

import (
	"fmt"
)

func main() {

	// arreglo estatico con 10 elementos: no se modifica el tiempo de ejecucion
	var array [10]int
	var matrix [3][2]int

	// asignacion directa con inicializacion de variables
	array2 := [3]int{ 1, 2, 3 }

	fmt.Println( array )
	fmt.Println( array2 )


	// asignacion en posicion
	array2[2] = 20

	matrix[0][1] = 1
	
	fmt.Println( matrix )

	// longitud
	for i := 0; i < len( array2 ); i++ {
		fmt.Println( array2[i] )
	}
}