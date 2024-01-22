package utils

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func CheckErrorFromDB(err error) string {
	var mysqlErr *mysql.MySQLError

	switch {
	case errors.As(err, &mysqlErr) && mysqlErr.Number == 1062:
		return "Hay datos duplicados, ya registrados por otro usuario"
	case errors.As(err, &mysqlErr) && mysqlErr.Number == 1452:
		return "No se puede crear o actualizar el registro debido a datos inexistentes"
	case errors.As(err, &mysqlErr) && mysqlErr.Number == 1451:
		return "No se puede eliminar o actualizar el registro"
	case errors.As(err, &mysqlErr) && mysqlErr.Number == 1054:
		return "Intentas registrar una columna o campo desconocido"
	case errors.Is(err, gorm.ErrRecordNotFound):
		return "No se pudo encontrar el registro"
	default:
		return err.Error()
	}
}
