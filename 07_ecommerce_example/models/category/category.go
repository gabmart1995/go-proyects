package category

import (
	"database/sql"
	"log"

	"github.com/go-playground/locales/es"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	es_translations "github.com/go-playground/validator/v10/translations/es"
	"github.com/gofiber/fiber/v2"
)

type Category struct {
	db     *sql.DB
	Id     int64
	Nombre string `validate:"required,min=2,max=100"`
}

func New(db *sql.DB) Category {
	return Category{db: db}
}

func (c *Category) SetNombre(nombre string) {
	c.Nombre = nombre
}

func (c *Category) GetNombre() string {
	return c.Nombre
}

func (c *Category) SetId(id int64) {
	c.Id = id
}

func (c *Category) GetId() int64 {
	return c.Id
}

func (c *Category) GetAll(order string) ([]Category, error) {
	var categories []Category
	var query string

	query = "SELECT * FROM categorias"

	if order == "DESC" {
		query += " ORDER BY id DESC"
	}

	query += ";"

	rows, err := c.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// realizamos la lectura
	for rows.Next() {
		var category Category

		err := rows.Scan(&category.Id, &category.Nombre)

		if err != nil {
			log.Fatal(err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (c *Category) GetOne() (Category, error) {
	var category Category

	query := "SELECT * FROM categorias WHERE id = ?"
	err := c.db.QueryRow(query, c.GetId()).Scan(
		&category.Id,
		&category.Nombre,
	)

	if err != nil {
		return category, err
	}

	return category, nil
}

// send a empty slice and validate all struct
func (c *Category) Validate(fields []string, UI string) (fiber.Map, error) {
	var err error

	validate := validator.New()
	es := es.New()
	uni := ut.New(es, es)
	trans, _ := uni.GetTranslator("es")
	es_translations.RegisterDefaultTranslations(validate, trans)

	// indica si la validacion es parcial
	if len(fields) > 0 {
		err = validate.StructPartial(c, fields...)

	} else { // sino toda la estructura
		err = validate.Struct(c)

	}

	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}

		// constrimos el objeto de errores vacio
		errors := fiber.Map{}

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

	return nil, nil
}

func (c *Category) Save() error {
	sql := "INSERT INTO categorias(nombre) VALUES(?);"
	stmt, err := c.db.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	// cambiamos los valores
	_, err = stmt.Exec(c.GetNombre())

	if err != nil {
		return err
	}

	return nil
}
