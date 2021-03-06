/*
	此包用于获取配置，
	iris 框架本身的配置处理已经比较完善，
	增加这些方法主要是增加配置使用的灵活性
*/
package config

import (
	"sync"
	"time"

	"GoTenancy/backend/files"
	"GoTenancy/backend/transformer"
	"github.com/kataras/iris/v12"
	gf "github.com/snowlyg/gotransformer"
)

type config struct {
	Tc  *transformer.Conf
	Isc iris.Configuration
}

var cfg *config
var once sync.Once

func singleton() *config {
	once.Do(func() {
		path := files.GetAbsPath("./config/conf.tml")
		isc := iris.TOML(path) // 加载配置文件
		tc := getTfConf(isc)
		cfg = &config{Tc: tc, Isc: isc}
	})
	return cfg
}

func getTfConf(isc iris.Configuration) *transformer.Conf {

	app := transformer.App{}
	g := gf.NewTransform(&app, isc.Other["App"], time.RFC3339)
	_ = g.Transformer()

	db := transformer.Mysql{}
	g.OutputObj = &db
	g.InsertObj = isc.Other["Mysql"]
	_ = g.Transformer()

	mongodb := transformer.Mongodb{}
	g.OutputObj = &mongodb
	g.InsertObj = isc.Other["Mongodb"]
	_ = g.Transformer()

	redis := transformer.Redis{}
	g.OutputObj = &redis
	g.InsertObj = isc.Other["Redis"]
	_ = g.Transformer()

	sqlite := transformer.Sqlite{}
	g.OutputObj = &sqlite
	g.InsertObj = isc.Other["Sqlite"]
	_ = g.Transformer()

	admin := transformer.Admin{}
	g.OutputObj = &admin
	g.InsertObj = isc.Other["Admin"]
	_ = g.Transformer()

	testData := transformer.TestData{}
	g.OutputObj = &testData
	g.InsertObj = isc.Other["TestData"]
	_ = g.Transformer()

	return &transformer.Conf{
		App:      app,
		Mysql:    db,
		Mongodb:  mongodb,
		Redis:    redis,
		Sqlite:   sqlite,
		Admin:    admin,
		TestData: testData,
	}
}

func GetIrisConf() iris.Configuration {
	return singleton().Isc
}

func getTc() *transformer.Conf {
	return singleton().Tc
}

func GetAppName() string {
	return getTc().App.Name
}

func GetAppUrl() string {
	return getTc().App.Url
}

func GetAppLoggerLevel() string {
	return getTc().App.LoggerLevel
}

func GetAppDriverType() string {
	return getTc().App.DriverType
}

func GetAppCookieNameForSessionID() string {
	return getTc().App.CookieNameForSessionID
}

func GetAppCreateSysData() bool {
	return getTc().App.CreateSysData
}

func GetMysqlConnect() string {
	return getTc().Mysql.Connect
}

func GetMysqlName() string {
	return getTc().Mysql.Name
}

func GetMysqlTName() string {
	return getTc().Mysql.TName
}

func GetMongodbConnect() string {
	return getTc().Mongodb.Connect
}

func GetSqliteConnect() string {
	return files.GetAbsPath(getTc().Sqlite.Connect)
}

func GetSqliteTConnect() string {
	return files.GetAbsPath(getTc().Sqlite.TConnect)
}

func GetRedisAddr() string {
	return getTc().Redis.Addr
}

func GetRedisPwd() string {
	return getTc().Redis.Password
}

func GetRedisDb() string {
	return getTc().Redis.DB
}

func GetAdminUserName() string {
	return getTc().Admin.UserName
}

func GetAdminName() string {
	return getTc().Admin.Name
}

func GetAdminPwd() string {
	return getTc().Admin.Pwd
}

func GetTestDataUserName() string {
	return getTc().TestData.UserName
}

func GetTestDataName() string {
	return getTc().TestData.Name
}

func GetTestDataPwd() string {
	return getTc().TestData.Pwd
}
