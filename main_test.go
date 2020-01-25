package GoTenancy

import (
	"log"
	"os"
	"testing"

	"github.com/snowlyg/GoTenancy/data"
)

var db *data.DB

func TestMain(m *testing.M) {
	ds := "user=postgres password=postgres dbname=test sslmode=disable"

	db = &data.DB{}
	if err := db.Open("postgres", ds); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Connection.Ping(); err != nil {
		log.Fatal(err)
	}

	// we make sure to clean everything before starting the tests
	if _, err := db.Connection.Exec("DELETE FROM GoTenancy_accounts;"); err != nil {
		log.Fatal(err)
	}

	retval := m.Run()
	os.Exit(retval)
}
