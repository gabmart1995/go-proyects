package main

import (
  "fmt"
  "io/ioutil"
)

func main()  {
  // code ...
  file_data, error := ioutil.ReadFile("./hola.txt")

  if error != nil {
    fmt.Println( "Hubo un error al crear archivo" )
  }

  fmt.Println( string( file_data ) )
}
