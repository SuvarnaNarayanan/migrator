package main

import (
	"database/sql"
	"flag"
	"fmt"
	"migrator/utils"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type NewMigrationUserInfo struct {
	Operation string
	Desc      string
}

type RunTimeOptions struct {
	Init            bool
	GenerateSQLFile bool
	MakeMigrations  bool
	Migrate         bool
	FakeMigrate     bool
}

func main() {

	rOptions := RunTimeOptions{}

	sysArg := os.Args[1]

	// make this into command line args
	flag.BoolVar(&rOptions.Init, "init", false, "Initialize the migrations database, table and folder if it does not exist")
	flag.BoolVar(&rOptions.GenerateSQLFile, "generate", false, "Generate a new migration sql file")
	flag.BoolVar(&rOptions.MakeMigrations, "makemigrations", false, "Make migrations")
	flag.BoolVar(&rOptions.Migrate, "migrate", false, "Migrate the pending migrations")
	flag.BoolVar(&rOptions.FakeMigrate, "fake", false, "Fake migrate the pending migrations")
	flag.Parse()

	if sysArg == "init" {
		rOptions.Init = true
	} else if sysArg == "generate" {
		rOptions.GenerateSQLFile = true
	} else if sysArg == "g" {
		rOptions.GenerateSQLFile = true
	} else if sysArg == "makemigrations" {
		rOptions.MakeMigrations = true
	} else if sysArg == "migrate" {
		rOptions.Migrate = true
	} else if sysArg == "m" {
		rOptions.Migrate = true
	} else if sysArg == "mm" {
		rOptions.MakeMigrations = true
	} else if sysArg == "i" {
		rOptions.Init = true
	} else if sysArg == "help" {
		fmt.Println(`
	Commands:
	init: Initialize the migrations database, table and folder if it does not exist
	generate: Generate a new migration sql file
	makemigrations: Make migrations
	migrate: Migrate the pending migrations
		`)
		return
	} else if sysArg == "fake" {
		rOptions.FakeMigrate = true
	}

	if rOptions.Init {
		utils.Init()
		return
	}

	mConfig, err := utils.ReadConfig()
	if err != nil {
		utils.PrintError(err)
		return
	}

	mdbName, err := mConfig.GetMigrationsDatabaseName()
	if err != nil {
		utils.PrintError(err)
		return
	}
	mTableName, err := mConfig.GetMigrationsTableName()
	if err != nil {
		utils.PrintError(err)
		return
	}

	mdb, err := sql.Open("sqlite3", fmt.Sprintf("./%s.db", mdbName))
	if err != nil {
		utils.PrintError(&utils.MigratorError{
			SysErr: err.Error(),
			Code:   utils.MIGRATION_DATABASE_CANNOT_BE_OPENED,
			Hint:   "Migration database cannot be opened",
		})
		return
	}
	tdb, err := sql.Open(mConfig.TargetDbInfo.Driver, mConfig.TargetDbInfo.DataSource)
	if err != nil {
		utils.PrintError(&utils.MigratorError{
			SysErr: err.Error(),
			Code:   utils.TARGET_DATABASE_CANNOT_BE_OPENED,
			Hint:   "Target database cannot be opened",
		})
	}
	defer mdb.Close()
	defer tdb.Close()

	_, err = mdb.Exec(fmt.Sprintf(utils.CREATE_STMT, mTableName))
	if err != nil {
		utils.PrintError(&utils.MigratorError{
			SysErr: err.Error(),
			Code:   utils.MIGRATION_TABLE_CANNOT_BE_CREATED,
			Hint:   "Migration table cannot be created",
		})
		panic(err)
	}

	if rOptions.GenerateSQLFile {

		// need to get input here from user

		var op string
		var desc string

		fmt.Println("Enter operation: ")
		fmt.Scanln(&op)
		fmt.Println("Enter short file name: ")
		fmt.Scanln(&desc)

		dir, err := mConfig.GetMigrationsDir()
		if err != nil {
			utils.PrintError(err)
			return
		}

		utils.GenerateANewMigrationSqlFile(mdb, op, desc, dir)
		return
	}

	if rOptions.MakeMigrations {
		var files []utils.SQLFile
		dir, err := mConfig.GetMigrationsDir()
		if err != nil {
			utils.PrintError(err)
			return
		}
		files, err = utils.ReadAllSQLFiles(dir, mdb)
		if err != nil {
			utils.PrintError(err)
			return
		}
		files = utils.FilterFiles(files, mdb)
		if len(files) == 0 {
			fmt.Println("No migrations to initialize")
			return
		}
		MakeMigrations(mConfig, mdb, files)
		return
	}

	if rOptions.Migrate {
		var files []utils.SQLFile
		dir, err := mConfig.GetMigrationsDir()
		if err != nil {
			utils.PrintError(err)
			return
		}
		files, err = utils.ReadAllSQLFiles(dir, mdb)
		if err != nil {
			utils.PrintError(err)
			return
		}
		Migrate(mConfig, mdb, files, tdb, false)
		return
	}

	if rOptions.FakeMigrate {
		var files []utils.SQLFile
		dir, err := mConfig.GetMigrationsDir()
		if err != nil {
			utils.PrintError(err)
			return
		}
		files, err = utils.ReadAllSQLFiles(dir, mdb)
		if err != nil {
			utils.PrintError(err)
			return
		}
		Migrate(mConfig, mdb, files, tdb, true)
		return
	}

	fmt.Println("No command specified")
}

func MakeMigrations(mConfig *utils.MigratorConfig, mdb *sql.DB, fileNames []utils.SQLFile) {
	for _, f := range fileNames {
		tableName, err := mConfig.GetMigrationsTableName()
		if err != nil {
			panic(err) // This should never happen
		}
		err = utils.InsertNewMigrationRecord(mdb, f.FileName, "PENDING", tableName)
		if err != nil {
			utils.PrintError(&utils.MigratorError{
				SysErr: err.Error(),
				Code:   utils.MIGRATION_RECORD_CANNOT_BE_INSERTED,
				Hint:   "Migration record cannot be inserted",
			})
			return
		}
		fmt.Println("file: ", f.FileName, " initialized")
	}
}

func Migrate(mConfig *utils.MigratorConfig, mdb *sql.DB, fileNames []utils.SQLFile, tdb *sql.DB, fake bool) {

	tableName, err := mConfig.GetMigrationsTableName()
	if err != nil {
		panic(err) // This should never happen
	}
	canMigrate := utils.CheckIfMigrationsCanBeRun(mdb, tableName, fileNames)
	if canMigrate != nil {
		utils.PrintError(canMigrate)
		return
	}

	for _, f := range fileNames {
		if !fake {
			err := utils.RunMigration(tdb, f.Path)
			if err != nil {
				utils.PrintError(&utils.MigratorError{
					SysErr: err.Error(),
					Code:   utils.SQL_EXECUTION_ERROR,
					Hint:   "SQL execution error - please check the SQL file for errors " + f.FileName,
				})
				err := utils.UpdateMigrationRecord(mdb, f.FileName, "FAILED", tableName)
				if err != nil {
					panic(err)
				}
				return
			}
		}
		err = utils.UpdateMigrationRecord(mdb, f.FileName, "COMPLETED", tableName)
		if err != nil {
			panic(err)
		}
		fmt.Println("file: ", f.FileName, " migrated")
	}
}
