package database

import (
	"fmt"
	envutil "yelloment-api/env_util"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DB : db connection pool
var DB *sqlx.DB

// GetDB : db connection pool을 획득합니다.
func GetDB(dbType string) *sqlx.DB {
	var db *sqlx.DB
	switch dbType {
	case "test":
		db, _ = sqlx.Connect("sqlite", "test.db")
	case "operation_local_test":
		db, _ = sqlx.Connect("sqlite", "operation.db")
	case "operation_remote_test":
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			envutil.GetGoDotEnvVariable("OPERATION_REMOTE_TEST_DB_ID"),
			envutil.GetGoDotEnvVariable("OPERATION_REMOTE_TEST_DB_PW"),
			envutil.GetGoDotEnvVariable("OPERATION_REMOTE_TEST_DB_HOST"),
			envutil.GetGoDotEnvVariable("OPERATION_REMOTE_TEST_DB_NAME"),
		)
		db, _ = sqlx.Connect("mysql", dsn)
		fmt.Println("db", db)
	default:
		panic(fmt.Sprintf("dbType is not allowed. dbType: %s", dbType))
	}

	return db
}

// SetDB : db connection pool 변수를 설정합니다.
func SetDB(dbType string) {
	DB = GetDB(dbType)
}
