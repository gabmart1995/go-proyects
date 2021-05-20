package greetings

import (
	"fmt"
	"errors"
	"math/rand"
  "time"
	// "log"
)

// para exportar la funcion la primera letra debe estar capitalizada
func Hello( name string ) ( string, error ) {

	if name == "" {
		return "", errors.New("Nombre vacio")
	}

	var message string
	message = fmt.Sprintf( randomFormat(), name )

	return message, nil
}

func randomFormat() string {
  formats := []string {
    "Hola, %v. Bienvenido",
    "Un placer verte, %v!",
    "Hasta la proxima, %v",
  }

  // log.Print( formats )

  return formats[ rand.Intn( len( formats )  ) ]
}

// setea los valores iniciales del paquete
func init() {

  // inicializa un numero aleatorio, usando el formato de tiempo de Unix
  // desde 1 enero del 1970
  rand.Seed( time.Now().UnixNano() )

  // log.Println( time.Now().UnixNano() )
}
