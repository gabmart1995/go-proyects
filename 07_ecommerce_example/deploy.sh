#!/bin/bash
# Script que controla el build de produccion del aplicacion web
echo "¿Cuál del sistema operativo a desplegar? (linux, windows):"
read os

echo "¿Cuál es la architectura del servidor a desplegar? (amd64, 386):"
read arch

directory="./dist"

if [ -d $directory ]
then
    echo "El directorio ya existe, se actualiza el mismo"
    rm -r $directory

    mkdir -p $directory

    cp -a templates $directory
    cp -a static $directory
    cp database.sql $directory
else
    mkdir -p $directory

    cp -a templates $directory
    cp -a static $directory
    cp database.sql $directory
    cp .env $directory # copia el elemento de entorno
fi

# compila la aplicacion en la aquitectura seleccionada
GOOS=$os GOARCH=$arch go build -o $directory
