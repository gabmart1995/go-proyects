package main

import (
	"fmt"
)

type Human struct {
	name string
}

type Tutor struct {
	Human
	Dummy
}

type Dummy struct {
	name string
}


func ( this Human ) talk() string {
	return "bla bla bla"
}

func ( this Tutor ) talk() string {

	// para acceder al metodo del padre se llama a partir de la referencia
	return this.Human.talk() + " Bienvendo"
}

func main() {
	

	/*
		Los campos anonimos nos permite realizar herencia entre
		estructuras, permite acceder a los campos del padre sin 
		tener que especificar el atributo a traves de la estructura 
		que lo implementa

	*/

	tutor := Tutor{ Human{ "Gabriel" }, Dummy{ "Dummy" } }

	fmt.Println( tutor.Human.name )
	fmt.Println( tutor.Dummy.name )

	fmt.Println( tutor.talk() )
}