package utils

import (
	"database/sql"
	"fmt"
	"os"
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

	lastRecord := db.QueryRow("SELECT id FROM migrations ORDER BY id DESC LIMIT 1")
	var maxId int
	err := lastRecord.Scan(&maxId)
	if err != nil {
		maxId = 1
	}

	op = strings.ToUpper(op)
	desc = strings.ToUpper(desc)
	desc = strings.ReplaceAll(desc, " ", "_")

	if !isValidOperation(op) {
		return "", fmt.Errorf("invalid operation: %s, can only be one of CREATE | UPDATE | DELETE", op)
	}

	return fmt.Sprintf("%d_%s_%s", maxId, op, desc), nil
}

func GenerateANewMigrationDatabase(db *sql.DB) error {
	return nil
}

func GenerateANewMigrationSqlFile(db *sql.DB, op string, desc_short string, dir string) error {
	newName, err := GenerateUniqueName(db, op, desc_short)
	if err != nil {
		return err
	}
	os.Create(fmt.Sprintf("%s/%s.sql", dir, newName))
	return nil
}
