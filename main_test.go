package GoTenancy

import (
	"log"
	"os"
	"testing"

	"github.com/snowlyg/GoTenancy/data"
	"github.com/snowlyg/GoTenancy/model"
)

var db *data.DB

func TestMain(m *testing.M) {
	ds := "./data/test.db"
	db = &data.DB{}
	if err := db.Open("sqlite3", ds); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	retval := m.Run()

	//我们确认在开始测试前清除了所有数据
	db.Connection.DropTableIfExists(
		&model.AccessToken{},
		&model.Account{},
		&model.User{},
		&model.APIRequest{},
		&model.Webhook{},
	)

	os.Exit(retval)
}
