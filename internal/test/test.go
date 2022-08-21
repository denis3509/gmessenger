package test

import (
	"fmt"
	"log"
	"messenger/internal/config"
	"testing"

	db "messenger/pkg/db"

	_ "github.com/lib/pq"
)

func NewTestDB(t *testing.T) db.DB {
	cfg := config.GetConfig().DB
	testCfg := cfg
	testCfg.Name = cfg.Name + "_test"
	db, err := db.MustOpen(testCfg.DriverName, testCfg.DSN())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return db
}

// DropAndCreateTestDB drop existing test database and creates clone
// of main database with suffix "_test"
func DropAndCreateTestDB() {
	dbCfg := config.GetConfig().DB
	db, err := db.MustOpen(dbCfg.DriverName, dbCfg.DSN())
	if err != nil {
		log.Fatal("connection error", err)
	}
	
	dbTestName := dbCfg.Name + "_test"
	
	err  = db.CloseAllDBConnections(dbTestName)
	if err != nil {
		log.Fatal("cannot close connection: ", err)
	}

	q := db.DB.NewQuery(
		fmt.Sprintf(`DROP DATABASE IF EXISTS %s`,dbTestName),
	) 
	_, err = q.Execute()
	if err != nil {
		log.Fatal("drop database error: ", err)
	}

	err = db.CloneDB(dbTestName, dbCfg.Name, dbCfg.User)
	if err != nil {
		log.Fatal("clone database error: ", err)
	}
	fmt.Printf("database '%s' created", dbTestName)

}

func ResetTables(t *testing.T, db *db.DB, tables ...string) {
	for _, table := range tables {
		_, err := db.TruncateTable(table).Execute()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}
