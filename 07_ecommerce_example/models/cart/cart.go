package cart

import (
	"shirts-shop-golang/models/product"
)

type Cart struct {
	PedidoId   int64
	ProductoId int64
	Unidades   int64
	product.Product
}

func (c *Cart) SetUnidades(unidades int64) {
	c.Unidades = unidades
}

func (c *Cart) GetUnidades() int64 {
	return c.Unidades
}
