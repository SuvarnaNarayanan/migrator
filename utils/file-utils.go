package utils

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type SQLFile struct {
	Path     string
	Id       int
	FileName string
}

func SortFilesById(files []SQLFile) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Id < files[j].Id
	})
}

func ReadAllSQLFiles(mDir string, mdb *sql.DB) ([]SQLFile, error) {
	// Returns only unmigrated files in the migrations directory as specified in the config file

	var files []SQLFile
	err := filepath.Walk(mDir, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if !info.IsDir() && ext == ".sql" {
			strId := strings.Split(info.Name(), "_")[0]
			id, err := strconv.Atoi(strId)
			fileName := info.Name()
			if err != nil {
				return err
			}
			f := SQLFile{
				Path:     path,
				Id:       id,
				FileName: fileName,
			}
			if !CheckIfFileIsAlreadyMigrated(mdb, fileName) {
				files = append(files, f)
			}
		}
		return nil
	})
	SortFilesById(files)
	var errorString string
	if err != nil {
		errorString = err.Error()
		return nil, &MigratorError{
			SysErr: errorString,
			Code:   MIGRATION_SQL_FILE_CANNOT_BE_READ,
			Hint:   "SQL files cannot be read",
		}
	}
	return files, nil
}

func ReadContentFromFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func Init() error {

	// Create migration directory and empty example migration config file

	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		err = os.Mkdir("migrations", 0755)
		if err != nil {
			return &MigratorError{
				SysErr: err.Error(),
				Code:   MIGRATION_DIR_CANNOT_BE_CREATED,
				Hint:   "Please make sure you have the necessary permissions to create a directory in the root of your project",
			}
		}
	}

	m := &MigratorConfig{
		Migration: Migration{
			DbName:    "migrations",
			TableName: "migrations",
			Dir:       "migrations",
		},
		TargetDbInfo: TargetDbInfo{
			Driver:     "sqlite3 | mysql | postgres",
			DataSource: "",
			UserName:   "",
			Password:   "",
		},
	}

	if _, err := os.Stat("migrations.config.yaml"); os.IsNotExist(err) {
		mBytes, err := yaml.Marshal(m)
		if err != nil {
			return err
		}

		err = os.WriteFile("migrations.config.yaml", mBytes, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
