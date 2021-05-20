package main

import (
  "fmt"
  "example.com/greetings"
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
  // obtiene el mensaje y lo imprime
  message := greetings.Hello("Gabriel")
  fmt.Println( message )
}
