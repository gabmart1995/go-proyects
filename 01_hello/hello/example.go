package main

import (
  "fmt"
  "bufio"
  "os"
  "strings"
)


func example()  {
  menu := (`
    Bienvenido:
    [ 1 ] Pizza
    [ 2 ] Tacos

    Cual prefieres??
  `)

  reader := bufio.NewReader( os.Stdin )

  fmt.Println( menu )

  input, _ := reader.ReadString('\n') // lee el caracter hasta el salto de linea

  seleccion := strings.TrimRight( input, "\r\n" )

  switch seleccion {
    case "1":
      fmt.Println("Prefieres Pizza")
      break
    case "2":
      fmt.Println("Prefieres Tacos")
      break
    default:
      fmt.Println("No se selecciona una opcion valida")
      break
  }
}

func main()  {

  var name string

  fmt.Printf("Ingresa tu nombre: ")

  // Scanf escanea hasta el espacio
  fmt.Scanf("%s", &name )

  fmt.Printf("La palabra %s contiene %d numero de palabras", name, len( name ))
}
