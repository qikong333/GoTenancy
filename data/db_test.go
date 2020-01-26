package data

import (
	"testing"
)

func Test_DB_Open(t *testing.T) {
	db := DB{}
	ds := "./test.db"

	if err := db.Open("sqlite3", ds); err != nil {
		t.Fatal("unable to connect to sqlite3", err)
	}

	defer db.Close()
}
