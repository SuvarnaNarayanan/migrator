package utils

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isValidOperation(op string) bool {
	return op == "CREATE" || op == "UPDATE" || op == "DELETE"
}

func GenerateUniqueName(db *sql.DB, op string, desc string) (string, error) {

	// ID_DB-OPERATION_DESCRIPTOR
	// ID: incrementing number
	// DB-OPERATION: db operation - CREATE, UPDATE, DELETE
	// DESCRIPTOR: short description of the operation

	lastRecord := db.QueryRow("SELECT file_name FROM migrations ORDER BY id DESC LIMIT 1")
	var maxIdString string
	var maxId int
	err := lastRecord.Scan(&maxIdString)
	if err != nil {
		maxId = 1
	} else {
		maxIdString = strings.Split(maxIdString, "_")[0]
		maxId, err = strconv.Atoi(maxIdString)
		if err != nil {
			return "", &MigratorError{
				SysErr: err.Error(),
				Code:   IMPROPRER_MIGRATION_FILE_NAME,
				Hint:   "Please make sure the migration file name is in the correct format - ID_DB-OPERATION_DESCRIPTOR where ID is an integer, it is followed by an underscore, followed by the operation (CREATE, UPDATE, DELETE) and then the descriptor. \n Looks like you might have the offending file already as a record in the migrations table - consider updating the file name manually.",
			} // this should never happen
		}
	}

	op = strings.ToUpper(op)
	desc = strings.ToUpper(desc)
	desc = strings.Join(strings.Fields(desc), "_")

	if !isValidOperation(op) {
		return "", fmt.Errorf("invalid operation: %s, can only be one of CREATE | UPDATE | DELETE", op)
	}

	return fmt.Sprintf("%d_%s_%s", maxId+1, op, desc), nil
}

func GenerateANewMigrationDatabase(db *sql.DB) error {
	return nil
}

func GenerateANewMigrationSqlFile(db *sql.DB, op string, desc_short string, dir string) (*SQLFile, error) {
	newName, err := GenerateUniqueName(db, op, desc_short)
	if err != nil {
		return nil, err
	}
	os.Create(fmt.Sprintf("%s/%s.sql", dir, newName))
	return &SQLFile{
		FileName: newName,
		Id:       -1, // does not matter
		Path:     "", // does not matter
	}, nil
}
