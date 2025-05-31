package vm

import (
	"encoding/json"
	"fmt"
	"os"
	"sqllite-go/internal/common/models"
)

const FILE_NAME = "DB.txt"

type VM interface {
	ExecuteByteCode() models.VMResult
}

// estructuramos una mini maquina virtual para estructura la
// ejecucion de los datos
type MiniVM struct {
	tableName string
	byteCode  *models.ByteCode
}

func NewVM(tableName string, byteCode *models.ByteCode) VM {
	return &MiniVM{tableName: tableName, byteCode: byteCode}
}

func (vm_ *MiniVM) ExecuteByteCode() models.VMResult {
	table, err := vm_.dbOpen()
	response := models.VMResult{}

	if err != nil {
		response.Err = err
		return response
	}

	// extraemos el tipo de operacion de acuerdo al tipo
	typeOperation := vm_.byteCode.Instructions[0].Type

	if typeOperation == models.ByteCodeOperationTypeINSERT {
		str, err := vm_.executeInsert(table)

		response.MSG = str
		response.Err = err
		response.Cursor = nil // el cursor es el apuntador de las paginas

		return response
	}

	if typeOperation == models.ByteCodeOperationTypeSELECT {
		cursor := vm_.executeSelect(table)
		response.Cursor = cursor

		return response
	}

	// en caso de no ser ninguna mostramos un error
	response.Err = fmt.Errorf("%s", "Operation not found")

	return response
}

// lee el archivo de la base de datos
func (vm_ *MiniVM) dbOpen() (*models.Table, error) {
	// sino existe el archivo de BD lo crea
	if !checkFileExists(FILE_NAME) {
		file, err := os.OpenFile(FILE_NAME, os.O_CREATE|os.O_APPEND, 0644)

		if err != nil {
			return nil, err
		}

		file.Close()
	}

	// procede a leer el archivo
	data, err := os.ReadFile(FILE_NAME)

	if err != nil {
		return nil, err
	}

	var pager models.Pager

	// asigna los datos a la interfaz pager
	if len(data) > 0 {
		if err := json.Unmarshal(data, &pager); err != nil {
			return nil, err
		}
	}

	table := models.Table{
		Name:    vm_.tableName,
		NumRows: len(pager.Pages),
		Pager:   &pager,
	}

	return &table, nil
}

func checkFileExists(path string) bool {
	_, err := os.Open(path)
	return err == nil
}

// escribe los datos en el disco
func (vm_ *MiniVM) write(pager *models.Pager) error {
	file, err := os.Create(FILE_NAME)

	if err != nil {
		return err
	}

	defer file.Close()

	// escribimos el archivo
	data, err := json.Marshal(pager)

	if err != nil {
		return err
	}

	file.Write(data)
	file.Sync() // una vez guardado hace flush de la data (gestion de memoria)

	return nil
}

// ejecuta la instruccion insert en la maquina virtual
func (vm_ *MiniVM) executeInsert(table *models.Table) (*string, error) {
	count := vm_.byteCode.Instructions[2].Count

	// en caso de que la tabla llegue vacia
	if table.NumRows == 0 {
		// armamos un array con los valores de las columnas
		cols := []string{}

		for i := 3; i < (count + 3); i++ {
			cols = append(cols, *vm_.byteCode.Instructions[i].Identifier)
		}

		table.Pager = &models.Pager{
			Columns: cols,
			Pages:   [][]models.Page{},
		}
	}

	index := count + 4
	HasPrimaryKey := true

	pages := []models.Page{
		{
			IsPrimaryKey: &HasPrimaryKey,
			IntValue:     vm_.byteCode.Instructions[index].IntValue,
		},
	}

	// procedemos a verificar las claves primarias en cada pagina
	if table.NumRows > 0 {
		for _, currentPage := range table.Pager.Pages {
			isDuplicatedKey := (currentPage[0].IntValue != nil &&
				*currentPage[0].IntValue == *vm_.byteCode.Instructions[index].IntValue)

			if isDuplicatedKey {
				return nil, fmt.Errorf(
					"duplicated key %d",
					*vm_.byteCode.Instructions[index].IntValue,
				)
			}
		}
	}

	// incrementamos al siguiente valor
	index++

	// extraemos las instrucciones del bytecode
	for index < len(vm_.byteCode.Instructions) {
		page := models.Page{
			IntValue:    vm_.byteCode.Instructions[index].IntValue,
			StringValue: vm_.byteCode.Instructions[index].StringValue,
		}

		pages = append(pages, page)
		index++
	}

	table.Pager.Pages = append(table.Pager.Pages, pages)
	table.NumRows = len(table.Pager.Pages)

	// por ultimo actualizamos el archivo
	if err := vm_.write(table.Pager); err != nil {
		return nil, err
	}

	str := "Record stored succesfully"

	return &str, nil
}

// ejecuta la instruccion select en la maquina virtual
func (vm_ *MiniVM) executeSelect(table *models.Table) *models.Pager {
	count := vm_.byteCode.Instructions[2].Count
	cols := []string{}

	for i := 3; i < (count + 3); i++ {
		cols = append(cols, *vm_.byteCode.Instructions[i].Identifier)
	}

	table.Pager.Columns = cols

	return table.Pager
}
