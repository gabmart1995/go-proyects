package main

import (
  "fmt"
  "bufio"
  "os"
)

func main()  {
  // code ...
  confirm := readFile()
  fmt.Println( confirm )
}

func readFile() bool {

  // se ejecuta cuando aplica la finalizacion de la funcion utilizando de defer

  file, error := os.Open("./hola.txt")
  var numberLine int

  // defer utiliza la ejecucion de una funcion anonima para su ejecucion
  defer func () {
    file.Close()
    fmt.Println("Defer")
  }()

  if error != nil {
    fmt.Println("Hubo un error")
    return false
  }

  // generamos un elemento de tipo scanner
  scanner := bufio.NewScanner( file )

  // generamos un bucle por cada linea del archivo
  for scanner.Scan() {

    numberLine++

    // obtenemos la linea
    line := scanner.Text()

    fmt.Println( line )
    fmt.Println( numberLine )
  }

  return true
}
