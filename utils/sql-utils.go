package utils

import (
	"database/sql"
	"fmt"
)

const CREATE_STMT = `CREATE TABLE IF NOT EXISTS %s (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	file_name TEXT NOT NULL UNIQUE,
	status TEXT NOT NULL,
	CHECK (status IN ('PENDING', 'COMPLETED', 'FAILED'))
)`

func InsertNewMigrationRecord(db *sql.DB, fileName string, status string, mTableName string) error {
	SQL := fmt.Sprintf("INSERT INTO %s (file_name, status) VALUES (?, ?)", mTableName)
	_, err := db.Exec(SQL, fileName, status)
	return err
}

func UpdateMigrationRecord(db *sql.DB, file_name string, status string, mTableName string) error {
	SQL := fmt.Sprintf("UPDATE %s SET status = ? WHERE file_name = ?", mTableName)
	_, err := db.Exec(SQL, status, file_name)
	return err
}

// RunMigration runs a migration on the target database
func RunMigration(tdb *sql.DB, path string) error {
	content, err := ReadContentFromFile(path)
	if err != nil {
		return err
	}
	_, err = tdb.Exec(content)
	return err
}

func CheckIfFileIsAlreadyMigrated(db *sql.DB, fileName string) bool {
	rows, err := db.Query("SELECT 1 FROM migrations WHERE file_name = ? AND STATUS = 'COMPLETED'", fileName)
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func CheckIfMigrationsCanBeRun(db *sql.DB, mTableName string, files []SQLFile) error {
	// Check if there are any pending or failed records

	rows, err := db.Query("SELECT 1 FROM migrations WHERE status IN ('PENDING', 'FAILED') LIMIT 1")
	if err != nil {
		return &MigratorError{
			SysErr: err.Error(),
			Code:   MIGRATION_TABLE_CANNOT_BE_READ,
			Hint:   "Migration table cannot be read",
		}
	}

	defer rows.Close()

	if !rows.Next() {
		return &MigratorError{
			SysErr: "None",
			Code:   NO_MIGRATIONS_TO_RUN,
			Hint:   "No migrations to run",
		}
	}

	// check if files are in pending or failed state
	for _, f := range files {
		rows, err := db.Query("SELECT 1 FROM migrations WHERE file_name = ? AND status IN ('PENDING', 'FAILED') LIMIT 1", f.FileName)
		if err != nil {
			return &MigratorError{
				SysErr: err.Error(),
				Code:   MIGRATION_TABLE_CANNOT_BE_READ,
				Hint:   "Migration table cannot be read",
			}
		}
		defer rows.Close()
		for !rows.Next() {
			return &MigratorError{
				SysErr: "None",
				Code:   NO_MIGRATIONS_TO_RUN,
				Hint:   "Please make sure all migrations are in PENDING or FAILED state - " + f.FileName,
			}
		}
	}

	return nil
}
