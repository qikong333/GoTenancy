package data

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/snowlyg/GoTenancy/model"
)

// Open 创建数据库连接并且初始化。
func (db *DB) Open(driverName, dataSource string) error {
	var err error
	db.Connection, err = gorm.Open(driverName, dataSource)
	if err != nil {
		return errors.New("failed to connect database")
	}
	//defer db.Connection.Close()

	// Migrate the schema
	db.Connection.AutoMigrate(
		&model.AccessToken{},
		&model.Account{},
		&model.User{},
		&model.APIRequest{},
		&model.Webhook{},
	)
	return nil
}

func (db *DB) Close() {
	db.Connection.Close()
}
