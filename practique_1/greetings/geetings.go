package greetings

import "fmt"

// para exportar la funcion la primera letra debe estar capitalizada
func Hello( name string ) string {

	// var message string
	// message = fmt.Sprintf( "Hola, %v. Bienvenido!", name )

	// es lo mismo, en una sola linea
	message := fmt.Sprintf( "Hola, %v. Welcome!", name )

	return message
}
