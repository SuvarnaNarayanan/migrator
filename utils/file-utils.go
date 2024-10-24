package utils

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
	return files, err
}

func ReadContentFromFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
