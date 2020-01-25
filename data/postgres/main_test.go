package postgres

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var db *sql.DB

func TestMain(m *testing.M) {
	ds := "user=postgres password=postgres dbname=test sslmode=disable"
	conn, err := sql.Open("postgres", ds)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}

	//我们确认在开始测试前清除了所有数据
	_, err = conn.Exec("DELETE FROM GoTenancy_accounts;")
	if err != nil {
		log.Fatal(err)
	}

	db = conn

	retval := m.Run()
	os.Exit(retval)
}
