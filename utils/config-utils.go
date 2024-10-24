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
		return nil, err
	}

	c := MigratorConfig{}

	err = yaml.Unmarshal(cFile, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *MigratorConfig) GetMigrationsDatabaseName() string {
	return c.Migration.DbName
}

func (c *MigratorConfig) GetMigrationsTableName() string {
	return c.Migration.TableName
}

func (c *MigratorConfig) GetMigrationsDir() string {
	return c.Migration.Dir
}

func (c *MigratorConfig) GetTargetDbDataSource() string {
	return c.TargetDbInfo.DataSource
}

func (c *MigratorConfig) GetTargetDbDriver() string {
	return c.TargetDbInfo.Driver
}

func (c *MigratorConfig) GetTargetDbUserName() string {
	return c.TargetDbInfo.UserName
}
