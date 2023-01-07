// Write the main test file for datebase connection.
package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/amit16110/simplebank/util"
	_ "github.com/lib/pq"
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://postgres:secret@localhost:5432/simplebank?sslmode=disable"
// )

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..") //it's mean go to parent folder
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	testDb, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	testQueries = New(testDb)

	// Run. The main testing start
	os.Exit(m.Run())

}
