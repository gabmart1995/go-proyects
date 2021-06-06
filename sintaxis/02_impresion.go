package main

import (
  "fmt"
  "bufio"
  "os"
)

func main()  {

  /*
    En esta seccion se aplican los verbos para la
    impresion en consola
  */

  edad := 22

  // variables %v
  fmt.Printf( "Mi edad es %v\n", edad )

  // boleana
  flag := true
  fmt.Printf( "Mi edad es %t\n", flag )


  // float
  precio := 145.20
  fmt.Printf( "Mi edad es %f\n", precio )

  var nombre string;

  // scanf 
  fmt.Print("Ingresa tu nombre: ")
  fmt.Scanf("%s", &nombre )
  fmt.Printf( "Hola %s\n", nombre )

  // leer cadenas de texto
  reader := bufio.NewReader( os.Stdin )

  fmt.Print("Ingresa una fase: ")

  nombre, error := reader.ReadString('\n')

  // valida si existe el error
  if error != nil {
    fmt.Println( error )
  
  } else {
    fmt.Printf( "mi frase es: %s\n", nombre )

  }

}
