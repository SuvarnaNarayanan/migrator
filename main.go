package main

import (
	"database/sql"
	"fmt"
	"migrator/utils"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type NewMigrationUserInfo struct {
	Operation string
	Desc      string
}

func main() {

	mConfig, err := utils.ReadConfig()
	if err != nil {
		panic(err)
	}

	mdbName := mConfig.GetMigrationsDatabaseName()
	mTableName := mConfig.GetMigrationsTableName()

	mdb, err := sql.Open("sqlite3", fmt.Sprintf("./%s.db", mdbName))
	if err != nil {
		panic(err)
	}
	tdb, err := sql.Open(mConfig.TargetDbInfo.Driver, mConfig.TargetDbInfo.DataSource)
	if err != nil {
		panic(err)
	}
	defer mdb.Close()
	defer tdb.Close()

	_, err = mdb.Exec(fmt.Sprintf(utils.CREATE_STMT, mTableName))

	if err != nil {
		panic(err)
	}

	var files []utils.SQLFile
	files, err = utils.ReadAllSQLFiles(mConfig.GetMigrationsDir(), mdb)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", files)

	MakeMigrations(mConfig, mdb, files)
	Migrate(mConfig, mdb, files, tdb)

}

func MakeMigrations(mConfig *utils.MigratorConfig, mdb *sql.DB, fileNames []utils.SQLFile) {
	for _, f := range fileNames {
		err := utils.InsertNewMigrationRecord(mdb, f.FileName, "PENDING", mConfig.GetMigrationsTableName())
		if err != nil {
			panic(err)
		}
	}
}

func Migrate(mConfig *utils.MigratorConfig, mdb *sql.DB, fileNames []utils.SQLFile, tdb *sql.DB) {

	canMigrate := utils.CheckIfMigrationsCanBeRun(mdb, mConfig.GetMigrationsTableName())
	if !canMigrate {
		fmt.Println("No migrations to run")
		return
	}

	for _, f := range fileNames {
		err := utils.RunMigration(tdb, f.Path)
		if err != nil {
			err := utils.UpdateMigrationRecord(mdb, f.FileName, "FAILED", mConfig.GetMigrationsTableName())
			if err != nil {
				panic(err)
			}
			continue
		}
		err = utils.UpdateMigrationRecord(mdb, f.FileName, "COMPLETED", mConfig.GetMigrationsTableName())
		if err != nil {
			panic(err)
		}
	}
}

func CreateMigrationSqlFile() {
	// GenerateANewMigrationSqlFile()
	// insertNewMigrationRecord(db, desc_long, newName, "PENDING")
}
