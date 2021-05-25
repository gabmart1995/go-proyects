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


func Hellos( names []string ) ( map[ string ]string, error )  {

  // make crea un array con los resultados pasados por parametro
  // se crea un array asociativo con map que contiene tanto claves string y valores del mismo tipo

  messages := make( map[ string ] string )


  // el valor de _ correponde a variables que se ignoran por el compilador

  for _, name := range names {

    // fmt.Println( _, name ) mostrara un error

    message, error := Hello( name )

    if error != nil {
      return nil, error

    }

    // asigna los valores al array asociativo
    messages[name] = message
  }

  return messages, nil
}
