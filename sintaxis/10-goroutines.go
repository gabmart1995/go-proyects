package main

import (
  "fmt"
  "time"
  "strings"
)

/*
  Introduccion a programas concurrentes
*/

func SlowName( name string )  {
  letras := strings.Split( name, "" )

  for _, letra := range letras {

    // se genera cada letra por cada segundo
    time.Sleep( 1000 * time.Millisecond )
    fmt.Println( letra )
  }
}

func main()  {
  // metodo sincrono espera a que el valor culmine para comenzar
  // el siguiente, para transformar en asincrono se utiliza las
  // go rutinas, se antepone a la funcion la palabra go

  go SlowName("Gabriel")
  fmt.Println("Que aburrido")

  var wait string
  fmt.Scanln( &wait )
}
