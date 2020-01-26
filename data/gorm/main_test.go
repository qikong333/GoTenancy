package gorm

import (
	"log"
	"os"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/snowlyg/GoTenancy/data"
	"github.com/snowlyg/GoTenancy/model"
)

var db *data.DB

func TestMain(m *testing.M) {
	ds := "../test.db"
	db = &data.DB{}
	err := db.Open("sqlite3", ds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//我们确认在开始测试前清除了所有数据
	db.Connection.DropTableIfExists(
		&model.AccessToken{},
		&model.Account{},
		&model.User{},
		&model.APIRequest{},
		&model.Webhook{},
	)

	retval := m.Run()
	os.Exit(retval)
}
