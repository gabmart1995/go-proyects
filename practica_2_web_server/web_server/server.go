package main

import (
  "fmt"
  "log"
  "net/http"
  "example.com/read_file"
  "html/template"
)

func handlerExample( writer http.ResponseWriter, request *http.Request )  {
  /*
    [1:] es un segmento de ruta desde el primer caracter hasta el final
    fmt.Println( request.URL.Path[1:] )
  */

  fmt.Fprintf( writer, "Hola mundo desde el server, me gustan los %s!", request.URL.Path[1:] )
}

func viewHandler( writer http.ResponseWriter, request *http.Request )  {

  // extrae la subcadena desde el view en adelante
  title := request.URL.Path[ len("/view/"): ]
  page, error := read_file.LoadPage( title )

  if error != nil {
    http.Redirect( writer, request, "/edit/" + title, http.StatusFound )
    return
  }

  renderTemplate( writer, "view", page )

  /*
    templateHTML := (`<h1>%s</h1><div>%s</div>`)
    fmt.Fprintf( writer, templateHTML, page.Title, page.Body )
  */
}

func editHandler( writer http.ResponseWriter, request *http.Request )  {

  title := request.URL.Path[ len("/edit/"): ]
  page, error := read_file.LoadPage( title )

  if error != nil {
    page = &read_file.Page{ Title: title }
  }

  renderTemplate( writer, "edit", page )
}

func saveHandler( writer http.ResponseWriter, request *http.Request ) {

  title := request.URL.Path[ len("/save/"): ]

  // obtiene el body del formulario
  body := request.FormValue("body")

  // crea el puntero
  page := &read_file.Page{ Title: title, Body: []byte( body ) }
  page.Save()

  http.Redirect( writer, request, "/view/" + title, http.StatusFound )
}


func renderTemplate( writer http.ResponseWriter, templ string, page *read_file.Page )  {

  // retorna un *template.Template
  templateHTML, _ := template.ParseFiles( "./templates/" + templ + ".html" )
  templateHTML.Execute( writer, page )
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
