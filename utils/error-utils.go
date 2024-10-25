package utils

const (
	CONFIG_FILE_NOT_FOUND = iota
	YAML_UNMARSHAL_ERROR
	DATABASE_NAME_NOT_FOUND
	MIGRATION_TABLE_NAME_NOT_FOUND
	MIGRATION_DIRECTORY_NOT_FOUND
	TARGET_DB_INFO_NOT_FOUND
	MIGRATION_DATABASE_CANNOT_BE_OPENED
	TARGET_DATABASE_CANNOT_BE_OPENED
	MIGRATION_TABLE_CANNOT_BE_CREATED
	MIGRATION_DIR_CANNOT_BE_CREATED
	MIGRATION_RECORD_CANNOT_BE_INSERTED
	MIGRATION_SQL_FILE_CANNOT_BE_READ
	MIGRATION_TABLE_CANNOT_BE_READ
	IMPROPRER_MIGRATION_FILE_NAME
	SQL_EXECUTION_ERROR
	NO_MIGRATIONS_TO_RUN
)

type MigratorError struct {
	SysErr string
	Code   int
	Hint   string
}

func (c *MigratorError) Error() string {
	return c.SysErr
}

func NewConfigFileNotFoundError(err error) error {
	return &MigratorError{
		SysErr: err.Error(),
		Code:   CONFIG_FILE_NOT_FOUND,
		Hint:   "Please make sure the migrations.config.yaml file is in the root of your project",
	}
}

func PrintError(err error) {
	switch err.(type) {
	case *MigratorError:
		mErr := err.(*MigratorError)
		println("Error: ", mErr.SysErr)
		println("Hint: ", mErr.Hint)
		println("Code: ", mErr.Code)
	default:
		println("Error: ", err.Error())
	}
}
