# Blockchain App

Esta aplicacion, corresponde a la creacion de una blockchain usando GO como lenguaje 
de desarrollo, con las bibliotecas estandar, boltDB (como base de datos) y testify para el entorno de pruebas

## Lecciones aprendidas
- Creación de la estructura basica del blockchain, se utiliza el patron generador
para obtener los bloques sin afectar el rendimiento. (rama part-1-blockchain)
- Desarrollo de prueba de trabajo, para reforzar la creacion de bloques (rama part-2-blockchain)
- Se adapta la funcionalidad a una linea de comandos para mas comodidad (part-3-blockchain)
- Se llega a implementar las primeras transacciones junto con la creacion de la prueba de 
trabajo (part-4-blockchain)
- Es la seccion mas compleja, corresponde a la creacion de las direcciones junto a las claves publicas y privadas, usando algoritmos matematicos base58 curva elliptica. Verificación y firmas de las transacciones usando el algoritmo antes descrito (part-5-blockchain)
- Optimizacion de consultas y generación de los conjuntos de transacciones. (part-6-blockchain)
- Desentralizacion de los nodos blockchain, cada uno posee su sincronizacion junto a su base de
datos distribuida. Se construye un socket tcp para enviar los datos entre los diferentes nodos. (part-7-blockchain)

### Para desarrollo
Usa `go mod tidy` para construir los modulos y descargar las dependencias.

### Despluiegue
Simplemente se distribuye un binario usa `GOOS=linux|windows GOARCH=amd64|i686 go build .`

## Nota:
Para probar en modo local compilar la app en rama 6, si deseas ver como funciona el blockchain
desentralizado usar la rama 7. Debido a que se implementa un servicio que sincroniza los nodos
con los demas miembros de la red.

### Pruebas en la rama 7 
Para realizar las pruebas en la rama 7 seguir los siguientes pasos
- Generar 3 instancias de terminal
- Establecer en sus 3 instancias los puertos de conexion tcp <br /> 
(`export NODE_ID=3000`) nodo principal<br />
(`export NODE_ID=3001`) nodo cliente<br />
(`export NODE_ID=3002`) nodo minero<br />
- En la instancia de `NODE_ID=3000` crea un wallet y un blockchain
- Copia el blockchain creado a blockchain genesis se usara para incializar<br> 
los otros nodos `cp blockchain_3000.db blockchain_genesis.db`
- En la instancia de `NODE_ID=3001` crea 3 wallets
- En la instancia de `NODE_ID=3000` envia monedas desde el nodo central a cualquiera de los 
3 wallets en la instancia `NODE_ID=3001` asegurate de activar el flag `-mine`
- Una vez finalizado levanta el node con startnode y dejarlo levantado.
- En la instancia de `NODE_ID=3001` Copia el blockchain `cp blockchain_genesis.db blockchain_3001.db`
- Inicia la primera sincronizacion con startnode en la instancia `NODE_ID=3001` 
automaticamente bajara los bloques que ya se encuentran minados. En este punto ya la 
blockchain esta sincronizada
- En la instancia de `NODE_ID=3002` Copia el blockchain `cp blockchain_genesis.db blockchain_3002.db`
- Inicializa el nodo minero con startnode y pasando su direccion como parametro, esto automaticamente se sincroniza y los nuevos bloques pasaran a ser minados en esta instancia.
- Una vez finalizado, se sincroniza el nodo central y distribuye la informacion a los hijos.

### Agradecimientos
- jeiwan (Ivan Kuznetsov) - Creador del [articulo](https://jeiwan.net/posts/building-blockchain-in-go-part-1/) y sentar el articulo detallado de desarrollo blockchain.
- yohilime - Por dar solución a un problema de encoding del algoritmo p256 en la parte 5

