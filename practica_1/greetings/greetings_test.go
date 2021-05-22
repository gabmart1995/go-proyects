/* archivos de prueba unitaria */
package greetings

import (
	"testing"
	"regexp"
)


// TestHelloName llama greetings Hello con el nombre, chequeando y
// devolviendo el valor valido

/*
	Nota para ejecutar la prueba simplemente coloca
	por consola go test -v, para mostrar el procedimiento de 
	ejecucion
*/

func TestHelloName( t *testing.T ) {
	
	name := "Gladys"

	// se crea la expresion regular si falla la expresion 
	// mostrara un error fatal

	want := regexp.MustCompile(`\b`+name+`\b`)

	message, error := Hello("Gladys")

	if !want.MatchString( message ) || error != nil {
		t.Fatalf(`Hello("Gladys") = %q, %v want match for %#q, nil`, message, error, want )
	}
}

func TestHelloEmpty( t *testing.T ) {
	
	message, error := Hello("")

	if message != "" || error == nil {
		t.Fatalf(`Hello("") = %q, %v, want ""`, message, error )
	}
}