# Bittorrent client
Esta es la construcción de una aplicacion que permite descargar archivos
utilizando el protocolo bittorrent.

## Lecciones aprendidas
- Entendimiento sobre las redes desentralizadas P2P.
- Introduccion al protocolo bittorrent para descargar.
- Utilzacion del desplazamiento de bits para crear los protocolos de comunicacion
- Familiaridad con los estados de conexión TCP entre los peers
- Introduccion a las pruebas unitarias en GO.

## Limitaciones
- Solo se puede descargar un solo archivo
- Se utiliza un archivo .torrent para descargar.

### Compilar
`GOOS=linux|windows GOARCH=ia32|amd64 go build .`

### uso 
`bittorrent-client file.torrent file.iso`