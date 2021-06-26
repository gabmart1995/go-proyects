package main

import (
  "net/http"
  "log"
  "fmt"
)

func main()  {

  // creacion del servidor de archivos estaticos
  fileServer := http.FileServer( http.Dir("public") )
  http.Handle( "/", http.StripPrefix( "/", fileServer ))

  fmt.Println("Servidor estatico operando en el puerto 8000")

  // si existe algun error desde la plataforma
  log.Fatal( http.ListenAndServe( ":8000", nil ) )
}
