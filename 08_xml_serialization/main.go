package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

type Users struct {
	XMLName xml.Name `xml:"users"`
	Users   []User   `xml:"user"`
}

type User struct {
	XMLName xml.Name `xml:"user"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name"`
	Social  Social   `xml:"social"`
}

type Social struct {
	XMLName  xml.Name `xml:"social"`
	Facebook string   `xml:"facebook"`
	Twitter  string   `xml:"twitter"`
	Youtube  string   `xml:"youtube"`
}

// extrae XML desde fuentes externas
func extractXML() {
	xmlFile, err := os.Open("example.xml")

	if err != nil {
		log.Panic(err)
	}

	defer xmlFile.Close()

	var users Users

	// obtenemos los bytes del archivo
	b, err := io.ReadAll(xmlFile)

	if err := xml.Unmarshal(b, &users); err != nil {
		log.Panic(err)
	}

	// imprime los datos
	for _, user := range users.Users {
		fmt.Println("User type: " + user.Type)
		fmt.Println("User Name: " + user.Name)
		fmt.Println("Facebook Url: " + user.Social.Facebook)

		fmt.Printf("\n")
	}
}

// permite crear estructuras XML para alojar en archivos o respuestas HTTP
func createXML() {
	data := Users{
		Users: []User{
			{
				Name: "Test 1",
				Type: "admin",
				Social: Social{
					Facebook: "https://facebook.com",
					Twitter:  "https://twitter.com",
					Youtube:  "https://youtube.com",
				},
			},
			{
				Name: "Test 2",
				Type: "user",
				Social: Social{
					Facebook: "https://facebook.com",
					Twitter:  "https://twitter.com",
					Youtube:  "https://youtube.com",
				},
			},
			{
				Name: "Test 3",
				Type: "client",
				Social: Social{
					Facebook: "https://facebook.com",
					Twitter:  "https://twitter.com",
					Youtube:  "https://youtube.com",
				},
			},
		},
	}

	out, err := xml.Marshal(&data)

	if err != nil {
		log.Panic(err)
	}

	// imprime los resultados
	fmt.Println(string(out))
}

func main() {
	extractXML()
	createXML()
}
