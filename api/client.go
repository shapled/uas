package api

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"sync"
	"uas/settings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var dbClient *sqlx.DB
var dbClientInitLock sync.Mutex

const (
	DriverMySQL = "mysql"
	DriverSQLite = "sqlite3"
)

func initDBClient() error {
	if dbClient != nil {
		return nil
	}
	dbClientInitLock.Lock()
	defer dbClientInitLock.Unlock()
	if dbClient != nil {
		return nil
	}
	var err error
	switch settings.UASSettings.DB.Driver {
	case DriverMySQL:
		dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s",
			settings.UASSettings.DB.User,
			settings.UASSettings.DB.Password,
			settings.UASSettings.DB.Host,
			settings.UASSettings.DB.Port,
			settings.UASSettings.DB.Name)
		dbClient, err = sqlx.Connect("mysql", dsn)
	case DriverSQLite:
		dbClient, err = sqlx.Connect("sqlite3", settings.UASSettings.DB.File)
	default:
		logrus.Fatalf("unknown driver settings: %s", settings.UASSettings.DB.Driver)
	}
	return err
}

func Dao(handler func(context.Context, *sqlx.DB) error) error {
	if err := initDBClient(); err != nil {
		return err
	}
	return handler(context.Background(), dbClient)
}

func DaoWithTx(handler func(context.Context, *sqlx.Tx) error) error {
	if err := initDBClient(); err != nil {
		return err
	}
	tx, err := dbClient.Beginx()
	if err != nil {
		return err
	}
	if err = handler(context.Background(), tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
