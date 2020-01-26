package data

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/snowlyg/GoTenancy/model"
)

var conn *gorm.DB
var err error

// Open 创建数据库连接并且初始化。
func (db *DB) Open(driverName, dataSource string) error {
	conn, err = gorm.Open(driverName, dataSource)
	if err != nil {
		return errors.New("gorm.Open error")
	}

	if db == nil {
		return errors.New("data.DB is nil")
	}
	//defer conn.Close()

	// Migrate the schema
	conn.AutoMigrate(
		&model.AccessToken{},
		&model.Account{},
		&model.User{},
		&model.APIRequest{},
		&model.Webhook{},
	)
	db.Users = &Users{DB: conn}
	db.Webhooks = &Webhooks{DB: conn}
	db.Connection = conn

	return nil
}

func (db *DB) Close() {
	db.Connection.Close()
}
