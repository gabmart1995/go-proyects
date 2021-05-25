package main

import (
  "fmt"
  "log"
  "net/http"
  "html/template"
  "regexp"
  "io/ioutil" 
)

type Page struct {
  Title string
  Body []byte
  Link template.HTML
}

// ====================================================================
// variables globales
// ====================================================================

var templates = template.Must( template.ParseFiles( "./temp/front.html", "./temp/edit.html", "./temp/view.html" ) )
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// ====================================================================
// main
// ====================================================================
func main()  {

  http.HandleFunc( "/view/", makeHandler( viewHandler ) )
  http.HandleFunc( "/edit/", makeHandler( editHandler ) )
  http.HandleFunc( "/save/", makeHandler( saveHandler ) )

  http.HandleFunc( "/", frontHandler )

  // asignacion de directorios estaticos para los estilos de la wiki (.css, .js)
  http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer( http.Dir("assets") )))

  fmt.Println("Servidor escuchando en el puerto 8080")

  log.Fatal( http.ListenAndServe( ":8080", nil ) )
}


// se le integra la referencia del puntero actual a la funcion
func ( p *Page ) Save() error {

  fileName := "./data/" + p.Title + ".txt"

  // el tercer parametro es la permisologia unix solo el
  // que crea el archivo puede escribir y leer el archivo

  return ioutil.WriteFile( fileName, p.Body, 0600 )
}


func LoadPage( title string ) ( *Page, error ) {

  fileName := "./data/" + title + ".txt"
  body, error := ioutil.ReadFile( fileName )

  if error != nil {
    return nil, error
  }

  return &Page{ Title: title, Body: body }, nil
}

func handlerExample( writer http.ResponseWriter, request *http.Request )  {
  
  /*
    [1:] es un segmento de ruta desde el primer caracter hasta el final
    fmt.Println( request.URL.Path[1:] )
  */

  fmt.Fprintf( writer, "Hola mundo desde el server, me gustan los %s!", request.URL.Path[1:] )
}

func viewHandler( writer http.ResponseWriter, request *http.Request, title string )  {

  page, error := LoadPage( title )

  if error != nil {
    
    http.Redirect( writer, request, "/edit/" + title, http.StatusFound )

    return
  }

  // le decimos a go que es seguro este html y lo asigna a la estuctura
  page.Link = template.HTML( createLink( page ) )

  renderTemplate( writer, "view", page )
}

func createLink( page *Page ) []byte {

  link := []byte(`<a href="/view/[PageName]">[PageName]</a>`)
  
  linkRexp := regexp.MustCompile("\\[([a-zA-Z0-9])+\\]")

  link = linkRexp.ReplaceAllFunc( link, func( s []byte ) []byte {

    title := linkRexp.ReplaceAllString( string(s), page.Title )
     
    return []byte( title )
  }) 

  return link
}

func editHandler( writer http.ResponseWriter, request *http.Request, title string )  {

  // title := request.URL.Path[ len("/edit/"): ] before
  
  page, error := LoadPage( title )

  if error != nil {
    page = &Page{ Title: title }
  }

  renderTemplate( writer, "edit", page )
}

func saveHandler( writer http.ResponseWriter, request *http.Request, title string ) {

  body := request.FormValue("body")

  page := &Page{ Title: title, Body: []byte( body ) }
  error := page.Save()

  if error != nil {

    http.Error( writer, error.Error(), http.StatusInternalServerError )
    
    return 
  }

  http.Redirect( writer, request, "/view/" + title, http.StatusFound )
}


func renderTemplate( writer http.ResponseWriter, templ string, page *Page )  {

  // se le pasa el nombre del archivo

  error := templates.ExecuteTemplate( writer, templ + ".html", page )

  if error != nil {

    http.Error( writer, error.Error(), http.StatusInternalServerError )
  }

}

func frontHandler( writer http.ResponseWriter, request *http.Request ) {
  renderTemplate( writer, "front", nil )
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