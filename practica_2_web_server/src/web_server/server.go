package main

import (
  "fmt"
  "log"
  "net/http"
  "example.com/read_file"
  "html/template"
  "regexp"
  // "errors"
)

func handlerExample( writer http.ResponseWriter, request *http.Request )  {
  
  /*
    [1:] es un segmento de ruta desde el primer caracter hasta el final
    fmt.Println( request.URL.Path[1:] )
  */

  fmt.Fprintf( writer, "Hola mundo desde el server, me gustan los %s!", request.URL.Path[1:] )
}

func viewHandler( writer http.ResponseWriter, request *http.Request, title string )  {

  page, error := read_file.LoadPage( title )

  if error != nil {
    
    http.Redirect( writer, request, "/edit/" + title, http.StatusFound )

    return
  }

  renderTemplate( writer, "view", page )
}

func editHandler( writer http.ResponseWriter, request *http.Request, title string )  {

  // title := request.URL.Path[ len("/edit/"): ] before
  
  page, error := read_file.LoadPage( title )

  if error != nil {
    page = &read_file.Page{ Title: title }
  }

  renderTemplate( writer, "edit", page )
}

func saveHandler( writer http.ResponseWriter, request *http.Request, title string ) {

  body := request.FormValue("body")

  page := &read_file.Page{ Title: title, Body: []byte( body ) }
  error := page.Save()

  if error != nil {

    http.Error( writer, error.Error(), http.StatusInternalServerError )
    
    return 
  }

  http.Redirect( writer, request, "/view/" + title, http.StatusFound )
}


func renderTemplate( writer http.ResponseWriter, templ string, page *read_file.Page )  {

  // se le pasa el nombre del archivo

  error := templates.ExecuteTemplate( writer, templ + ".html", page )

  if error != nil {

    http.Error( writer, error.Error(), http.StatusInternalServerError )
  }

}


// closures metodo poderoso para capturar los errores antes de ejecutar el codigo
// HandlerFunc es el tipo de no confundir con HandleFunc que es la funcion

func makeHandler( callback func( http.ResponseWriter, *http.Request, string )) ( http.HandlerFunc ) {
  
  return func( writer http.ResponseWriter, request *http.Request ) {
    
    match := validPath.FindStringSubmatch( request.URL.Path )

    if match == nil {
    
      http.NotFound( writer, request )

      return
    }

    callback( writer, request, match[2] )
  }
}

// ====================================================================
// main
// ====================================================================
func main()  {

  http.HandleFunc( "/view/", makeHandler( viewHandler ) )
  http.HandleFunc( "/edit/", makeHandler( editHandler ) )
  http.HandleFunc( "/save/", makeHandler( saveHandler ) )

  fmt.Println("Servidor escuchando en el puerto 8080")

  log.Fatal( http.ListenAndServe( ":8080", nil ) )
}


// ====================================================================
// variables globales
// ====================================================================

// variables globales crea una insancia de templates cacheadas en el servidor retorna un *Template

var templates = template.Must( template.ParseFiles( "./temp/edit.html", "./temp/view.html" ) )

// validacion de url

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
