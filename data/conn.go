package data

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/snowlyg/GoTenancy/data/postgres"
)

// Open 创建数据库连接并且初始化 postgres 服务。
func (db *DB) Open(driverName, dataSource string) error {
	conn, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}

	if err := conn.Ping(); err != nil {
		return err
	}

	db.Users = &postgres.Users{DB: conn}
	db.Webhooks = &postgres.Webhooks{DB: conn}

	db.Connection = conn

	db.DatabaseName = "GoTenancy"
	return nil
}

func (db *DB) Close() {
	db.Connection.Close()
}
