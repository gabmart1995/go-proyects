package read_file

import (
	"io/ioutil" // paquete de lectura de archivos
)

type Page struct {
	Title string
	Body []byte
}

// se le integra la referencia del puntero actual a la funcion
func ( p *Page ) Save() error {

	fileName := "./content/" + p.Title + ".txt"

	// el tercer parametro es la permisologia unix solo el
	// que crea el archivo puede escribir y leer el archivo

	return ioutil.WriteFile( fileName, p.Body, 0600 )
}


func LoadPage( title string ) ( *Page, error ) {

	fileName := "./content/" + title + ".txt"
	body, error := ioutil.ReadFile( fileName )

	if error != nil {
		return nil, error
	}

	return &Page{ Title: title, Body: body }, nil
}

/* func main() {

	// crea el nuevo puntero y crea el archivo
	p1 := &Page{ Title: "TestPage", Body: []byte("This a single page.") }
	p1.save()

	p2, _ := loadPage("TestPage")

	// se realiza un casting para leer el contenido del archivo
	fmt.Println( string( p2.Body ) )
} */
