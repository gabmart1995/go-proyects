package main

import (
  "fmt"
  "log"
  "net/http"
  "example.com/read_file"
  "html/template"
  "regexp"
  "errors"
)

func handlerExample( writer http.ResponseWriter, request *http.Request )  {
  
  /*
    [1:] es un segmento de ruta desde el primer caracter hasta el final
    fmt.Println( request.URL.Path[1:] )
  */

  fmt.Fprintf( writer, "Hola mundo desde el server, me gustan los %s!", request.URL.Path[1:] )
}

func viewHandler( writer http.ResponseWriter, request *http.Request )  {

  title, errorTitle := getTitle( writer, request )
  
  if errorTitle != nil {
    return 
  }

  page, error := read_file.LoadPage( title )

  if error != nil {
    
    http.Redirect( writer, request, "/edit/" + title, http.StatusFound )

    return
  }

  renderTemplate( writer, "view", page )
}

func editHandler( writer http.ResponseWriter, request *http.Request )  {

  // title := request.URL.Path[ len("/edit/"): ] before
  
  title, errorTitle := getTitle( writer, request )

  if errorTitle != nil {
    return 
  }

  page, error := read_file.LoadPage( title )

  if error != nil {
    page = &read_file.Page{ Title: title }
  }

  renderTemplate( writer, "edit", page )
}

func saveHandler( writer http.ResponseWriter, request *http.Request ) {

  title, errorTitle := getTitle( writer, request )

  if errorTitle != nil {
    return 
  }

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

  fmt.Println( error )

  if error != nil {

    http.Error( writer, error.Error(), http.StatusInternalServerError )
  }

}

func getTitle( writer http.ResponseWriter, request *http.Request ) ( string, error ) {

  match := validPath.FindStringSubmatch( request.URL.Path )

  if match == nil {
    
    http.NotFound( writer, request )

    return "", errors.New("Titulo de página inválida")
  }

  // fmt.Println( match )

  return match[2], nil
}

// ====================================================================
// main
// ====================================================================
func main()  {

  http.HandleFunc( "/view/", viewHandler )
  http.HandleFunc( "/edit/", editHandler )
  http.HandleFunc( "/save/", saveHandler )

  fmt.Println("Servidor escuchando en el puerto 8080")

  log.Fatal( http.ListenAndServe( ":8080", nil ) )
}


// ====================================================================
// variables globales
// ====================================================================

// variables globales crea una insancia de templates cacheadas en el servidor retorna un *Template

var templates = template.Must( template.ParseFiles( "./templates/edit.html", "./templates/view.html" ) )

// validacion de url

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
