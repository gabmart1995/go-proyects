package main

import (
  "fmt"
  "strconv"
)

func main()  {
  /*
    var x, y, z int
    var cadena string
    var flag bool
    var cadenas []string
  */

  const (
    xConst = "Prueba"
    dev = "Desarollo"
  )

  // numero a string
  x := "23"

  edad_int, _ := strconv.Atoi( x )

  fmt.Println( edad_int + 10 )
  fmt.Println( xConst )
  fmt.Println( dev )
}
