package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config interface {
	GetDatabaseName() string
	GetMigrationsTableName() string
	ReadConfig() string
}

type Migration struct {
	DbName    string `yaml:"dbname"`
	TableName string `yaml:"tablename"`
	Dir       string `yaml:"dir"`
}

type TargetDbInfo struct {
	Driver     string `yaml:"driver"`
	DataSource string `yaml:"datasource"`
	UserName   string `yaml:"username"`
	Password   string `yaml:"password"`
}

type MigratorConfig struct {
	Migration    Migration    `yaml:"migration"`
	TargetDbInfo TargetDbInfo `yaml:"targetdb"`
}

func ReadConfig() (*MigratorConfig, error) {

	cFile, err := os.ReadFile("./migrations.config.yaml")
	if err != nil {
		return nil, NewConfigFileNotFoundError(err)
	}

	c := MigratorConfig{}

	err = yaml.Unmarshal(cFile, &c)
	if err != nil {
		return nil, &MigratorError{
			SysErr: err.Error(),
			Code:   YAML_UNMARSHAL_ERROR,
			Hint:   "Please make sure to follow the correct format in the migrations.config.yaml file",
		}
	}

	return &c, nil
}

func (c *MigratorConfig) GetMigrationsDatabaseName() (string, error) {
	if c.Migration.DbName == "" {
		return "", &MigratorError{
			SysErr: "None",
			Code:   DATABASE_NAME_NOT_FOUND,
			Hint:   "Please make sure to specify the database name in the migrations.config.yaml file",
		}
	}
	return c.Migration.DbName, nil
}

func (c *MigratorConfig) GetMigrationsTableName() (string, error) {
	if c.Migration.DbName == "" {
		return "", &MigratorError{
			SysErr: "None",
			Code:   MIGRATION_TABLE_NAME_NOT_FOUND,
			Hint:   "Please make sure to specify the migration table name in the migrations.config.yaml file",
		}
	}
	return c.Migration.TableName, nil
}

func (c *MigratorConfig) GetMigrationsDir() (string, error) {
	if c.Migration.DbName == "" {
		return "", &MigratorError{
			SysErr: "None",
			Code:   MIGRATION_TABLE_NAME_NOT_FOUND,
			Hint:   "Please make sure to specify the migration table name in the migrations.config.yaml file",
		}
	}
	return c.Migration.Dir, nil
}

func (c *MigratorConfig) GetTargetDbDataSource() (string, error) {
	if c.Migration.DbName == "" {
		return "", &MigratorError{
			SysErr: "None",
			Code:   TARGET_DB_INFO_NOT_FOUND,
			Hint:   "Please make sure to specify the migration table name in the migrations.config.yaml file",
		}
	}
	return c.TargetDbInfo.DataSource, nil
}

func (c *MigratorConfig) GetTargetDbDriver() (string, error) {
	if c.Migration.DbName == "" {
		return "", &MigratorError{
			SysErr: "None",
			Code:   TARGET_DB_INFO_NOT_FOUND,
			Hint:   "Please make sure to specify the migration table name in the migrations.config.yaml file",
		}
	}
	return c.TargetDbInfo.Driver, nil
}

func (c *MigratorConfig) GetTargetDbUserName() (string, error) {
	if c.Migration.DbName == "" {
		return "", &MigratorError{
			SysErr: "None",
			Code:   TARGET_DB_INFO_NOT_FOUND,
			Hint:   "Please make sure to specify the migration table name in the migrations.config.yaml file",
		}
	}
	return c.TargetDbInfo.UserName, nil
}
