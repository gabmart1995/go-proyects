package main

import (
  "fmt"
  "example.com/greetings"
  "log"  // este paquete muestra los paquetes por consola
)

/*
  Para importar el archivo debes modificar el
  go.mod para que tenga una referencia del codigo
  fuente para realizar la ejecucion con el comando

  go mod edit -replace=example.com/greetings=../geetings

  esta opcion ejecuta la escritura del archivo del go.mod del
  paquete hello.

*/

func main() {

  // se preparan los logs de la aplicacion para responder los errores
  log.SetPrefix("greetings: ")
  log.SetFlags(0)

  // obtiene los valores de la funcion
  message, error := greetings.Hello("Gabriel")

  if error != nil {
    log.Fatal(error)
  }

  fmt.Println( message )
}
