package db

import (
	"fmt"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// DB holds connection to database
type DB struct {
	*dbx.DB
}

func NewDB(db *dbx.DB) DB {
	return DB{db}
}

func MustOpen(driver, dsn string) (DB, error) {
	_db, err := dbx.MustOpen(driver, dsn)
	db := NewDB(_db)
	return db, err
}

func (db *DB) CloneDB(dbNewName, dbOriginalName, dbUser string) error {
	err := db.CloseAllDBConnections(dbOriginalName)
	if err != nil {
		return err
	}
	q := db.NewQuery(
		fmt.Sprintf(
			"CREATE DATABASE %s WITH TEMPLATE %s OWNER %s;",
			dbNewName, dbOriginalName, dbUser),
	)

	_, err  = q.Execute()
	return err
}

func (db *DB) CloseAllDBConnections(dbName string) error {
	q := db.NewQuery(
		`SELECT pg_terminate_backend(pg_stat_activity.pid)
					 FROM pg_stat_activity 
					 WHERE pg_stat_activity.datname = {:dbName}
					 AND pid <> pg_backend_pid();`,
	)
	q.Bind(dbx.Params{
		"dbName": dbName,
	})
	_, err := q.Execute()
	return err
}
