package filters

import (
	"fmt"
	"html/template"
	"shirts-shop-golang/models/cart"
	"strconv"
	"time"
)

/* modulo de filtros  en templates HTML */

// devuelve la fecha en formato string
func GetYear() string {
	t := time.Now()
	return strconv.Itoa(t.Year())
}

// filtra el contenido HTML
func UnescapeHTML(s string) template.HTML {
	return template.HTML(s)
}

func IsAdmin(role string) bool {
	return role == "admin"
}

// ============================
//
//	filtros del carrito
//
// ============================
func CountCart(cartItems []cart.Cart) string {
	total := 0

	// si existe el array del carrito
	if cartItems != nil {
		total = len(cartItems)
	}

	return fmt.Sprintf("%d", total)
}

func GetTotalCart(cartItems []cart.Cart) string {
	total := 0.00

	if cartItems != nil {
		for _, cartItem := range cartItems {
			total += float64(cartItem.Product.Precio) * float64(cartItem.Unidades)
		}
	}

	return fmt.Sprintf("%.2f", total)
}

// ============================
//
//	filtros de la orden
//
// ============================
func GetStatus(status string) string {
	switch status {
	case "confirm":
		return "Pendiente"

	case "preparation":
		return "En preparación"

	case "ready":
		return "Preparado para enviar"

	default:
		return "Enviado"
	}
}
