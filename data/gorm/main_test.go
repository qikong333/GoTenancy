package gorm

import (
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	ds := "test.db"
	conn, err := gorm.Open("sqlite3", ds)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//我们确认在开始测试前清除了所有数据

	db = conn

	retval := m.Run()
	os.Exit(retval)
}
