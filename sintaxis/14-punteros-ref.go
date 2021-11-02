package main

/*
	Tutorial sobre punteros y referencias GO

	el "&" se coloca antes de la desclaracion de la variable
	para obtener la posicion en memoria almacenandose en una variable
	*string o *int ...

	para mostrar el contenido de la posicion en memoria incluye un * en la impresion 
	de la variable eliminandose asi la referencia de memoria

	Es más común usarlos al definir los argumentos de una función y los valores de retorno o 
	al definir métodos de tipos personalizados.

	Funciones

	al escribir una funcion puede definir los argumentos que se transmitiran pueden ser de 2 tipos:
	
	por valor: envia una copia de ese valor a la funcion y cualquier cambio dentro de la misma solo tiene efecto en 
	dicha funcion
	
	por referencia: cuando se pasa un puntero puede cambiar el valor de la variable en posicion en memoria 
	que se le paso por parametro.
*/

import (
	"fmt"
)

type Foods struct {
	name string
}

/*
func main() {


	var food string = "Pizza"

	// si estableces el signo "&" delante de la variable indicas que necesitas la posicion de la memoria
	// donde se encuenta la variable
	var pointer *string = &food

	fmt.Println( "food = ", food ) // pizza
	fmt.Println( "pointer = ", pointer ) // 0xc000040240  // direccion en memoria
	
	fmt.Println( "*pointer =", *pointer ) // pizza 


	// podemos sustituir el valor del puntero con el nuevo valor
	*pointer = "Ice Cream"
	fmt.Println( "*pointer =", *pointer ) // ice cream

	
	// si cambias el valor del puntero cambia directamente
	// el valor de la variable food al que se hace la referencia
	fmt.Println( "food = ", food ) // ice cream
}*/

// receptor de variables a traves de una funcion
func ( f Foods ) Reset() {
	f.name = ""
 
	// recibe una copia sin modificar el valor 
	// original alamcenado pasado por receptor
}

func ( f *Foods ) ResetPointer() {
	f.name = ""

	// modifica el valor de la posicion en memoria gracias al puntero
}

func main() {

	var food Foods = Foods{ name: "Pizza" }
	
	// se declara un puntero pero sin referencia a una variable el valor del puntero sera nil

	// var foodPointer *Foods

	// se asigna el valor de la variable y al mismo tiempo se obtiene la posicion en memoria
	
	// foodPointer = &Foods{ name: "Pizza" }

	// se elimina la referencia a la variable
	// fmt.Println( *foodPointer )

	// permite imprimir la clave y valode de la variable food
	fmt.Printf( "1) %+v\n", food )

	// para anadir metodos se le asigna al valor no al punterio 
	food.ResetPointer()
	
	// changeFoodPointer( foodPointer )

	// al volver a la funcion main el valor volvera a ser pizza
	fmt.Printf( "2) %+v\n", food )
}

// variable de valor
func changeFood( food Foods ) {

	// crea una copia de la variable que se pasa por parametro sin modificar el original 
	// almacenado en memoria en la funcion main

	food.name = "Arepa"
	fmt.Printf( "2) %+v\n", food )
}

// funcion que recibe por variable un tipo puntero no variable de valor
// error muy comun en go
func changeFoodPointer( food *Foods ) {

	// Atencion: esto modifica el valor almacenado en memoria en la variable funcion main

	// se valida que el puntero haga referencia a una variable sino devuelve nil
	if food == nil {
		fmt.Println("Food is nil")
		return
	}

	food.name = "Arepa"
	fmt.Printf( "2) %+v\n", food )
}

