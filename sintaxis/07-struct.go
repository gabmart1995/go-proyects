package main

import (
	"fmt"
)

type User struct {
	age int
	name, surname string
}

func ( this User ) getAllName() string {
	
	/*
		user es el identificador de la estructura

		User es una copia del estructura pero al final del
		dia no lo modfica directamente
	*/

	return this.name + " " + this.surname
}

// seter
func ( this *User ) setName( name string ) {
	this.name = name
}


func main() {

	// inicializacion de una estructura vacia
	var user User
	
	fmt.Println( user )

	user = User{ age: 24, name: "Gabriel", surname: "Martinez" }

	fmt.Println( user ) 

	// es la misma forma simplificada
	user2 := User{ age: 20, name: "Darianna", surname: "Martinez" }

	fmt.Println( user2 ) 

	// keyword new permite crear una estructura asociado a un puntero
	pointerUser := new( User )

	// las estructuras son mutables
	pointerUser.name = "Alfonso"
	pointerUser.surname = "Martinez"

	pointerUser.setName("Mauricio")

	// permite acceder al valor en memoria
	// fmt.Println( ( *pointerUser ).name )
	
	// fmt.Println( pointerUser.name )

	fmt.Println( pointerUser.getAllName() )
}