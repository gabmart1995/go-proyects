package product

import (
	"database/sql"
	"log"
	"mime/multipart"
	"shirts-shop-golang/models/category"

	"github.com/go-playground/locales/es"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	es_translations "github.com/go-playground/validator/v10/translations/es"
	"github.com/gofiber/fiber/v2"
)

type Product struct {
	Id          int64
	CategoryId  int64   `validate:"gt=0"`
	Nombre      string  `validate:"required,max=100,min=2"`
	Descripcion string  `validate:"max=255"`
	Precio      float64 `validate:"gt=0"`
	Stock       int64   `validate:"gt=0"`
	Oferta      string
	Fecha       string
	Imagen      string
	db          *sql.DB
	Category    category.Category `validate:"omitempty"` // ignora la validacion
}

func New(db *sql.DB) Product {
	return Product{db: db}
}

func (p *Product) GetId() int64 {
	return p.Id
}

func (p *Product) SetId(id int64) {
	p.Id = id
}

func (p *Product) GetNombre() string {
	return p.Nombre
}

func (p *Product) SetNombre(nombre string) {
	p.Nombre = nombre
}

func (p *Product) GetDescripcion() string {
	return p.Descripcion
}

func (p *Product) SetDescripcion(descripcion string) {
	p.Descripcion = descripcion
}

func (p *Product) GetPrecio() float64 {
	return p.Precio
}

func (p *Product) SetPrecio(precio float64) {
	p.Precio = precio
}

func (p *Product) GetStock() int64 {
	return p.Stock
}

func (p *Product) SetStock(stock int64) {
	p.Stock = stock
}

func (p *Product) GetOferta() string {
	return p.Oferta
}

func (p *Product) SetOferta(oferta string) {
	p.Oferta = oferta
}

func (p *Product) GetFecha() string {
	return p.Fecha
}

func (p *Product) SetFecha(fecha string) {
	p.Fecha = fecha
}

func (p *Product) GetImagen() string {
	return p.Imagen
}

func (p *Product) SetImagen(imagen string) {
	p.Imagen = imagen
}

func (p *Product) GetCategoryId() int64 {
	return p.CategoryId
}

func (p *Product) SetCategoryId(categoryId int64) {
	p.CategoryId = categoryId
}

func (p *Product) GetAll() ([]Product, error) {
	var products []Product

	sql := "SELECT * FROM productos ORDER BY id DESC;"
	rows, err := p.db.Query(sql)

	if err != nil {
		return products, err
	}

	defer rows.Close()

	// realizamos la lectura
	for rows.Next() {
		var product Product

		err := rows.Scan(
			&product.Id,
			&product.CategoryId,
			&product.Nombre,
			&product.Descripcion,
			&product.Precio,
			&product.Stock,
			&product.Oferta,
			&product.Fecha,
			&product.Imagen,
		)

		if err != nil {
			log.Fatal(err)
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *Product) GetAllCategory() ([]Product, error) {
	var products []Product

	query := "SELECT p.*, c.nombre FROM productos p INNER JOIN categorias c ON c.id = p.categoria_id WHERE p.categoria_id = ?  ORDER BY id DESC;"
	rows, err := p.db.Query(query, p.GetCategoryId())

	if err != nil {
		return products, err
	}

	defer rows.Close()

	// realizamos la lectura
	for rows.Next() {
		var product Product

		err := rows.Scan(
			&product.Id,
			&product.CategoryId,
			&product.Nombre,
			&product.Descripcion,
			&product.Precio,
			&product.Stock,
			&product.Oferta,
			&product.Fecha,
			&product.Imagen,
			&product.Category.Nombre,
		)

		if err != nil {
			log.Fatal(err)
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *Product) GetRandom(limit int64) ([]Product, error) {
	var products []Product

	query := "SELECT * FROM productos ORDER BY RAND() LIMIT ?;"
	rows, err := p.db.Query(query, limit)

	if err != nil {
		return products, err
	}

	defer rows.Close()

	// realizamos la lectura
	for rows.Next() {
		var product Product

		err := rows.Scan(
			&product.Id,
			&product.CategoryId,
			&product.Nombre,
			&product.Descripcion,
			&product.Precio,
			&product.Stock,
			&product.Oferta,
			&product.Fecha,
			&product.Imagen,
		)

		if err != nil {
			log.Fatal(err)
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *Product) GetOne() (Product, error) {
	var product Product

	query := "SELECT * FROM productos WHERE id = ?"

	err := p.db.QueryRow(query, p.GetId()).Scan(
		&product.Id,
		&product.CategoryId,
		&product.Nombre,
		&product.Descripcion,
		&product.Precio,
		&product.Stock,
		&product.Oferta,
		&product.Fecha,
		&product.Imagen,
	)

	if err != nil {
		return product, err
	}

	return product, nil
}

// valida los campos de la estructura
func (p *Product) Validate(fields []string, UI string, file *multipart.FileHeader) (fiber.Map, error) {
	// constrimos el objeto de errores vacio
	var err error
	errors := fiber.Map{}

	validate := validator.New()
	es := es.New()
	uni := ut.New(es, es)
	trans, _ := uni.GetTranslator("es")
	es_translations.RegisterDefaultTranslations(validate, trans)

	// indica si la validacion es parcial
	if len(fields) > 0 {
		err = validate.StructPartial(p, fields...)

	} else { // sino toda la estructura
		err = validate.Struct(p)

	}

	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}

		// obtenemos los errores de las estructuras
		errs := err.(validator.ValidationErrors)

		// fmt.Println(errs.Translate(trans))
		// devuelve un mapa en cada campo tiene un error
		errsTranslate := errs.Translate(trans)
		// fmt.Println(errsTranslate["User.Email"])

		for _, err := range err.(validator.ValidationErrors) {
			errors[UI+err.Field()] = errsTranslate[err.StructNamespace()]
		}

		if len(errors) > 0 {
			return errors, nil
		}
	}

	// validamos la imagen
	if file != nil {
		// valida el tipo del archivo
		const MBSize = 1000000
		contains := func(mimeType string) bool {
			mimeTypes := []string{"image/jpeg", "image/png", "image/webp", "image/gif", "image/jpg"}

			for _, value := range mimeTypes {
				if value == mimeType {
					return true
				}
			}

			return false
		}

		// validamos la tipoMime del archivo
		contentType := file.Header["Content-Type"][0]

		if !contains(contentType) {
			errors[UI+"File"] = "El archivo subido no es una imagen"
			return errors, nil
		}

		if file.Size > MBSize {
			errors[UI+"File"] = "El archivo debe ser menor a 1MB"
			return errors, nil
		}
	}

	return nil, nil
}

func (p *Product) Save() error {
	query := "INSERT INTO productos(categoria_id, nombre, descripcion, precio, stock, oferta, fecha, imagen) VALUES(?, ?, ?, ?, ?, ?, CURDATE(), ?)"
	stmt, err := p.db.Prepare(query)

	if err != nil {
		return err

	}
	defer stmt.Close()

	_, err = stmt.Exec(
		p.GetCategoryId(),
		p.GetNombre(),
		p.GetDescripcion(),
		p.GetPrecio(),
		p.GetStock(),
		p.GetOferta(),
		p.GetImagen(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *Product) Update() error {
	var err error
	query := "UPDATE productos SET categoria_id = ?, nombre = ?, descripcion = ?, precio = ?, stock = ?"

	// en caso de llegar imagen concatena la propiedad
	if len(p.GetImagen()) > 0 {
		query += ", imagen = ?"
	}

	query += " WHERE id = ?;"

	stmt, err := p.db.Prepare(query)

	if err != nil {
		return err

	}
	defer stmt.Close()

	if len(p.GetImagen()) > 0 {
		_, err = stmt.Exec(
			p.GetCategoryId(),
			p.GetNombre(),
			p.GetDescripcion(),
			p.GetPrecio(),
			p.GetStock(),
			p.GetImagen(),
			p.GetId(),
		)

	} else {
		_, err = stmt.Exec(
			p.GetCategoryId(),
			p.GetNombre(),
			p.GetDescripcion(),
			p.GetPrecio(),
			p.GetStock(),
			p.GetId(),
		)
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *Product) Delete() error {
	query := "DELETE FROM productos WHERE id = ?;"
	stmt, err := p.db.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.Exec(p.GetId()); err != nil {
		return err
	}

	return nil
}
