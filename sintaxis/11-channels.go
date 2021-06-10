package main

import (
  "fmt"
)

/*
  Introduccion a programas concurrentes

  Canales: nos permiten comunicar informacion
  entre las diversas rutinas de codigo.
*/

func main()  {
  channel := make( chan string )

  go func( channel chan string )  {

    for {
      var name string
      fmt.Scanln( &name )

      // para enviar la informacion en el Canal
      channel <- name
    }
  }( channel )

  msg := <- channel

  fmt.Println("Este es la info del canal " + msg )

  msg = <- channel

  fmt.Println("Este es la info del canal " + msg )
}
