package user

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/locales/es"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	es_translations "github.com/go-playground/validator/v10/translations/es"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	Id        int64
	Nombre    string `validate:"required,min=2,max=100"`
	Apellidos string `validate:"required,min=2,max=255"`
	Email     string `validate:"required,email,min=2,max=255"`
	Password  string `validate:"required,min=8,max=20"`
	Rol       string
	Imagen    string
	db        *sql.DB
}

// metodo constructor
func New(db *sql.DB) User {
	return User{db: db}
}

func (u *User) GetId() int64 {
	return u.Id
}

func (u *User) SetId(id int64) {
	u.Id = id
}

func (u *User) GetNombre() string {
	return u.Nombre
}

func (u *User) SetNombre(nombre string) {
	u.Nombre = nombre
}

func (u *User) GetApellido() string {
	return u.Apellidos
}

func (u *User) SetApellido(apellido string) {
	u.Apellidos = apellido
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) SetPassword(password string, hashing bool) {

	if hashing {
		// hashing the password
		bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)

		if err != nil {
			u.Password = ""
		}

		u.Password = string(bytes)
		return
	}

	u.Password = password
}

func (u *User) GetRol() string {
	return u.Rol
}

func (u *User) SetRol(rol string) {
	u.Rol = rol
}

func (u *User) GetImagen() string {
	return u.Imagen
}

func (u *User) SetImagen(imagen string) {
	u.Imagen = imagen
}

func (u *User) Save() error {
	sql := "INSERT INTO usuarios(nombre, apellidos, email, password, rol) VALUES(?, ?, ?, ?, 'user')"
	stmt, err := u.db.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	// cambiamos los valores
	_, err = stmt.Exec(
		u.GetNombre(),
		u.GetApellido(),
		u.GetEmail(),
		u.GetPassword(),
	)

	if err != nil {
		return err
	}

	return nil
}

// send a empty slice and validate all struct
func (u *User) Validate(fields []string, UI string) (fiber.Map, error) {
	var err error

	validate := validator.New()
	es := es.New()
	uni := ut.New(es, es)
	trans, _ := uni.GetTranslator("es")
	es_translations.RegisterDefaultTranslations(validate, trans)

	// indica si la validacion es parcial
	if len(fields) > 0 {
		err = validate.StructPartial(u, fields...)

	} else { // sino toda la estructura
		err = validate.Struct(u)

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

func (u *User) Login() (User, error) {
	// comprobar si existe el usuario
	query := "SELECT * FROM usuarios WHERE email = ?;"
	stmt, err := u.db.Prepare(query)

	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	var user User
	err = stmt.QueryRow(u.Email).Scan(
		&user.Id,
		&user.Nombre,
		&user.Apellidos,
		&user.Email,
		&user.Password,
		&user.Rol,
		&user.Imagen,
	)

	if err != nil && err == sql.ErrNoRows {
		return User{}, errors.New("No se hallaron resultados")
	}

	// revisamos las contrasenas
	err = bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(u.GetPassword()))

	if err != nil {
		return User{}, err
	}

	return user, nil
}
