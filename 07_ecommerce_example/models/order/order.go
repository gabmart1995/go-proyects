package order

import (
	"database/sql"
	"errors"
	"shirts-shop-golang/models/cart"
	"strconv"

	"github.com/go-playground/locales/es"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	es_translations "github.com/go-playground/validator/v10/translations/es"
	"github.com/gofiber/fiber/v2"
)

type Order struct {
	Id        int64
	UsuarioId int64  `validate:"gt=0"`
	Provincia string `validate:"required,max=100,min=2"`
	Localidad string `validate:"required,max=100,min=2"`
	Direccion string `validate:"required,max=100,min=2"`
	Coste     float64
	Estado    string
	Fecha     string
	Hora      string
	Cart      []cart.Cart // herencia desde: orden <- carrito <- producto
	db        *sql.DB
}

func New(db *sql.DB) Order {
	return Order{db: db}
}

func (o *Order) GetId() int64 {
	return o.Id
}

func (o *Order) SetId(id int64) {
	o.Id = id
}

func (o *Order) GetUsuarioId() int64 {
	return o.UsuarioId
}

func (o *Order) SetUsuarioId(usuarioId int64) {
	o.UsuarioId = usuarioId
}

func (o *Order) GetProvincia() string {
	return o.Provincia
}

func (o *Order) SetProvincia(provincia string) {
	o.Provincia = provincia
}

func (o *Order) GetLocalidad() string {
	return o.Localidad
}

func (o *Order) SetLocalidad(localidad string) {
	o.Localidad = localidad
}

func (o *Order) GetDireccion() string {
	return o.Direccion
}

func (o *Order) SetDireccion(direccion string) {
	o.Direccion = direccion
}

func (o *Order) GetCoste() float64 {
	return o.Coste
}

func (o *Order) SetCoste(coste float64) {
	o.Coste = coste
}

func (o *Order) GetEstado() string {
	return o.Estado
}

func (o *Order) SetEstado(estado string) {
	o.Estado = estado
}

func (o *Order) GetFecha() string {
	return o.Fecha
}

func (o *Order) SetFecha(fecha string) {
	o.Fecha = fecha
}

func (o *Order) GetHora() string {
	return o.Hora
}

func (o *Order) SetHora(hora string) {
	o.Hora = hora
}

func (o *Order) GetCart() []cart.Cart {
	return o.Cart
}

func (o *Order) SetCart(cart []cart.Cart) {
	o.Cart = cart
}

// busca un pedido por usuario
func (o *Order) GetOne() (Order, error) {
	var order Order

	query := "SELECT * FROM pedidos WHERE id = ?;"
	err := o.db.QueryRow(query, o.GetId()).Scan(
		&order.Id,
		&order.UsuarioId,
		&order.Provincia,
		&order.Localidad,
		&order.Direccion,
		&order.Coste,
		&order.Estado,
		&order.Fecha,
		&order.Hora,
	)

	if err != nil {
		return order, err
	}

	return order, nil
}

func (o *Order) GetAll() ([]Order, error) {
	var orders []Order

	query := "SELECT * FROM pedidos ORDER BY id DESC;"
	rows, err := o.db.Query(query)

	if err != nil {
		return orders, err
	}

	defer rows.Close()

	for rows.Next() {
		var order Order

		err = rows.Scan(
			&order.Id,
			&order.UsuarioId,
			&order.Provincia,
			&order.Localidad,
			&order.Direccion,
			&order.Coste,
			&order.Estado,
			&order.Fecha,
			&order.Hora,
		)

		if err != nil {
			return orders, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (o *Order) GetAllByUser() ([]Order, error) {
	var orders []Order

	query := "SELECT * FROM pedidos WHERE usuario_id = ? ORDER BY id DESC;"
	rows, err := o.db.Query(query, o.GetUsuarioId())

	if err != nil {
		return orders, err
	}

	defer rows.Close()

	for rows.Next() {
		var order Order

		err = rows.Scan(
			&order.Id,
			&order.UsuarioId,
			&order.Provincia,
			&order.Localidad,
			&order.Direccion,
			&order.Coste,
			&order.Estado,
			&order.Fecha,
			&order.Hora,
		)

		if err != nil {
			return orders, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// busca un pedido por usuario
func (o *Order) GetOneByUser() (Order, error) {
	var order Order

	query := "SELECT * FROM pedidos WHERE usuario_id = ? ORDER BY id DESC LIMIT 1;"
	err := o.db.QueryRow(query, o.GetUsuarioId()).Scan(
		&order.Id,
		&order.UsuarioId,
		&order.Provincia,
		&order.Localidad,
		&order.Direccion,
		&order.Coste,
		&order.Estado,
		&order.Fecha,
		&order.Hora,
	)

	if err != nil {
		return order, err
	}

	return order, nil
}

func (o *Order) GetProductsByOrder() ([]cart.Cart, error) {
	var cartItems []cart.Cart

	// subconsulta
	query := "SELECT lp.unidades, pr.* FROM lineas_pedidos lp INNER JOIN productos pr ON pr.id = lp.producto_id WHERE lp.pedido_id = ?;"
	rows, err := o.db.Query(query, o.GetId())

	if err != nil {
		return cartItems, err
	}

	for rows.Next() {
		var cartItem cart.Cart

		err := rows.Scan(
			&cartItem.Unidades,
			&cartItem.Product.Id,
			&cartItem.Product.CategoryId,
			&cartItem.Product.Nombre,
			&cartItem.Product.Descripcion,
			&cartItem.Product.Precio,
			&cartItem.Product.Stock,
			&cartItem.Product.Oferta,
			&cartItem.Product.Fecha,
			&cartItem.Product.Imagen,
		)

		if err != nil {
			return cartItems, err
		}

		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil
}

func (o *Order) Save() error {
	query := "INSERT INTO pedidos(usuario_id, provincia, localidad, direccion, coste, estado, fecha, hora) VALUES(?, ?, ?, ?, ?, 'confirm', CURDATE(), CURTIME());"
	stmt, err := o.db.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		o.GetUsuarioId(),
		o.GetProvincia(),
		o.GetLocalidad(),
		o.GetDireccion(),
		o.GetCoste(),
	)

	if err != nil {
		return err
	}

	return nil
}

// salva la linea de pedido
func (o *Order) SaveOrderLine() error {
	var order Order

	// obtenemos el ultimo elemento insertado
	query := "SELECT LAST_INSERT_ID();"
	err := o.db.QueryRow(query).Scan(&order.Id)

	if err != nil {
		return err
	}

	cart := o.GetCart()

	if len(cart) == 0 {
		return errors.New("el carrito de compra esta vacio")
	}

	// generamos la consulta
	query = "INSERT INTO lineas_pedidos(pedido_id, producto_id, unidades) VALUES "

	for index, cartItem := range cart {
		query += "("
		query += strconv.Itoa(int(order.GetId()))
		query += ", "
		query += strconv.Itoa(int(cartItem.Product.GetId()))
		query += ", "
		query += strconv.Itoa(int(cartItem.GetUnidades()))
		query += ")"

		if index == (len(cart) - 1) {
			query += ";"

		} else {
			query += ","

		}
	}

	// ejecutamos la consulta
	_, err = o.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func (o *Order) Validate(fields []string, UI string) (fiber.Map, error) {
	var err error

	validate := validator.New()
	es := es.New()
	uni := ut.New(es, es)
	trans, _ := uni.GetTranslator("es")
	es_translations.RegisterDefaultTranslations(validate, trans)

	// indica si la validacion es parcial
	if len(fields) > 0 {
		err = validate.StructPartial(o, fields...)

	} else { // sino toda la estructura
		err = validate.Struct(o)

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

func (o *Order) UpdateState() error {
	var err error
	query := "UPDATE pedidos SET estado = ? WHERE id = ?;"

	stmt, err := o.db.Prepare(query)

	if err != nil {
		return err

	}
	defer stmt.Close()

	_, err = stmt.Exec(
		o.GetEstado(),
		o.GetId(),
	)

	if err != nil {
		return err
	}

	return nil
}
