package dao

import (
	"context"
	"embed"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
)

//go:embed migrations/*.sql
var sqlFiles embed.FS

func ListSQLFiles() ([]string, error) {
	files, err := sqlFiles.ReadDir("migrations")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			names = append(names, file.Name())
		}
	}
	return names, nil
}

func CatSQLFile(filename string) (string, error) {
	fd, err := sqlFiles.Open("migrations/" + filename)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ExecSQLFile(filename string) error {
	content, err := CatSQLFile(filename)
	if err != nil {
		return err
	}
	return DaoWithTx(func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, content)
		return err
	})
}
