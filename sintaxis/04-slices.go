package main

import (
	"fmt"
)

func main() {

	/*
		Slice es un tipo de dato que puede modificar si estructura en tiempo de 
		ejecucion, se declara como un array pero sin la longitud espacificada.
		devuelve un array sin datos iniciales como nil

		Un slice posee una estructura predefinida como
		- Puntero al arreglo
		- Longitud del arrglo
		- Capacidad
	*/

	matrix := []int{ 1, 2, 3, 4 }

	// es un array tiene datos predefinidos
	array := [3]int{ 1 ,2, 3 }

	// toma el valor del inicio hasta el segundo elemento
	slice := array[:2]

	// toma desde la segunda posicion hasta el final
	slice = array[1:]

	fmt.Println( matrix )
	fmt.Println( len( matrix ) )
	fmt.Println( slice )

	// Make permite la inicializacion de un slice de forma rapida

	// el segundo paramtro es la longitud del slice: cuantos elementos
	// tiene el array interno lleno de datos

	// el tercer parametro de la funcion es la capacidad del slice (5)
	// es el limite de asignacion de elementos al array interno

	// que devuelve un puntero
	
	newArrayMake := make( []int, 3, 5 ) 

	fmt.Println( "Nuevo array" )

	// reconstuye un nuevo slice con los datos
	newArrayMake = append( newArrayMake, 2 )

	// cap nos devuelve la capacidad del slice
	fmt.Println( newArrayMake )
	// fmt.Println( cap( newArrayMake ) )

}