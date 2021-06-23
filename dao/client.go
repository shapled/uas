package dao

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"sync"
	"uas/settings"

	_ "github.com/go-sql-driver/mysql"
)

var dbClient *sqlx.DB
var dbClientInitLock sync.Mutex

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
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?multiStatements=true",
			settings.UASSettings.DB.User,
			settings.UASSettings.DB.Password,
			settings.UASSettings.DB.Host,
			settings.UASSettings.DB.Port,
			settings.UASSettings.DB.Name)
	dbClient, err = sqlx.Connect("mysql", dsn)
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
