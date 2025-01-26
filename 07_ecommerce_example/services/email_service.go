package services

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"shirts-shop-golang/filters"
	"shirts-shop-golang/models/order"
	"time"

	"gopkg.in/gomail.v2"
)

// exportamos el canal para el servicio
var (
	ChannelEmail       chan *gomail.Message
	pathEmailTemplates string // variable privada
)

// inicializa el servicio del correo electronico
func InitDaemonMail() {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	// establecemos el path de templates del correo
	pathEmailTemplates = filepath.Join(wd, "templates", "email")
	ChannelEmail = make(chan *gomail.Message)

	// iniciamos la goroutine
	go sendEmail()
}

// go routine que envia los correos por el canal
func sendEmail() {
	var (
		sender gomail.SendCloser
		err    error
	)

	isDialOpen := false
	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_SENDER"), os.Getenv("SMTP_PASSWORD"))

	// creamos el loop para escuchar el servicio
	for {
		select {
		case mail, ok := <-ChannelEmail:

			// si el canal esta cerrado
			if !ok {
				// fmt.Println("canal cerrado.")
				return
			}

			// sino esta autenticado
			// realiza la autenticacion
			if !isDialOpen {
				// fmt.Println("autenticando ...")

				// el dial si falla reintenta la conexión
				// saldra del bucle si la misma es exitosa
				for {
					if sender, err = dialer.Dial(); err == nil {
						isDialOpen = true
						break
					}

					// fmt.Println("reintentando conexión ...")
				}
			}

			// fmt.Println("enviando correo ...")

			// enviamos el correo
			if err = gomail.Send(sender, mail); err != nil {
				log.Fatal(err)
			}

			// fmt.Println("correo enviado.")

		// despues de 1min de inactividad cierra el dial de conexion
		case <-time.After(30 * time.Second):

			if isDialOpen {
				// fmt.Println("cerrando conexión ...")

				if err = sender.Close(); err != nil {
					log.Fatal(err)
				}

				isDialOpen = false

				// fmt.Println("cerrado con éxito.")
			}
		}
	}
}

func SendCreateOrderSucess(email string, dataOrder order.Order) {
	var body bytes.Buffer
	temp := template.Must(
		template.ParseFiles(filepath.Join(pathEmailTemplates, "order_success.html")))

	// estructura de datos para la vista
	data := struct {
		Title   string
		Message string
		Order   order.Order
	}{
		Title:   "Pedido completado",
		Message: "Pedido creado con éxito",
		Order:   dataOrder,
	}

	if err := temp.Execute(&body, data); err != nil {
		log.Fatal(err)
	}

	// resultado
	// fmt.Println(body.String())

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Tienda de camiseta. Pedido completado")

	message.SetBody("text/html", body.String())

	// le pasamos el mensaje por el canal
	ChannelEmail <- message
}

func SendUpdateOrderSucess(email string, dataOrder order.Order) {
	var body bytes.Buffer
	temp := template.Must(
		template.ParseFiles(filepath.Join(pathEmailTemplates, "order_status.html")))

	// estructura de datos para la vista
	data := struct {
		Title     string
		Order     order.Order
		GetStatus func(status string) string
	}{
		Title:     "Pedido Actualizado con éxito",
		Order:     dataOrder,
		GetStatus: filters.GetStatus,
	}

	if err := temp.Execute(&body, data); err != nil {
		log.Fatal(err)
	}

	// resultado
	// fmt.Println(body.String())

	message := gomail.NewMessage()
	message.SetHeader("From", "gabmart1995@gmail.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Tienda de camisetas. Pedido actualizado")

	message.SetBody("text/html", body.String())

	// le pasamos el mensaje por el canal
	ChannelEmail <- message
}

func SendCreateUser(email string) {
	var body bytes.Buffer
	temp := template.Must(
		template.ParseFiles(filepath.Join(pathEmailTemplates, "user_create.html")))

	// estructura de datos para la vista
	data := struct {
		Title   string
		Message string
	}{
		Title:   "Exito",
		Message: "Usuario creado con éxito",
	}

	if err := temp.Execute(&body, data); err != nil {
		log.Fatal(err)
	}

	// resultado
	// fmt.Println(body.String())

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Tienda de camisetas. Bienvenido")

	message.SetBody("text/html", body.String())

	// le pasamos el mensaje por el canal
	ChannelEmail <- message
}
