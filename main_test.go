package GoTenancy

import (
	"log"
	"os"
	"testing"

	"github.com/snowlyg/GoTenancy/data"
)

var db *data.DB

func TestMain(m *testing.M) {
	ds := "./data/test.db"

	db = &data.DB{}
	if err := db.Open("postgres", ds); err != nil {
		log.Fatal(err)
	}

	//我们确认在开始测试前清除了所有数据
	//if _, err := db.Connection.Exec("DELETE FROM GoTenancy_accounts;"); err != nil {
	//	log.Fatal(err)
	//}

	retval := m.Run()
	os.Exit(retval)
}
