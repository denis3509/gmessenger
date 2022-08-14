package test

import (
	"log"
	"messenger/internal/config"
	"testing"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
)


func DB() *dbx.DB {
	cfg := config.GetConfig()
	db, err := dbx.MustOpen("postgres", cfg.DSN)
	if err !=nil {
		log.Fatal(err)
	}
	return db 
}

func ResetTables(t *testing.T, db *dbx.DB, tables ...string) {
	for _, table := range tables {
		_, err := db.TruncateTable(table).Execute()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}
