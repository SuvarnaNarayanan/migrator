package utils

import (
	"database/sql"
	"fmt"
)

const CREATE_STMT = `CREATE TABLE IF NOT EXISTS %s (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	file_name TEXT NOT NULL,
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
	migrated := false
	rows, err := db.Query("SELECT 1 FROM migrations WHERE file_name = ? AND STATUS = 'COMPLETED'", fileName)
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var m int
		rows.Scan(&m)
		if m == 1 {
			migrated = true
		}
	}
	return migrated
}

func CheckIfMigrationsCanBeRun(db *sql.DB, mTableName string) bool {
	// Check if there are any pending or failed records

	rows, err := db.Query("SELECT 1 FROM migrations WHERE status IN ('PENDING', 'FAILED') LIMIT 1")
	if err != nil {
		return false
	}

	defer rows.Close()

	var m int
	for rows.Next() {
		rows.Scan(&m)
		if m == 1 {
			return true
		}
	}

	return false

}
