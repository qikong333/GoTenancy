package transformer

/*
加载系统配置
*/

type Conf struct {
	App      App
	Mysql    Mysql
	Mongodb  Mongodb
	Sqlite   Sqlite
	Redis    Redis
	Admin    Admin
	TestData TestData
}

type App struct {
	Name                   string
	Url                    string
	LoggerLevel            string
	DriverType             string
	CookieNameForSessionID string
	CreateSysData          bool
}

type Mysql struct {
	Connect string
	Name    string
	TName   string
}

type Mongodb struct {
	Connect string
}

type Sqlite struct {
	Connect  string
	TConnect string
}

type Redis struct {
	Addr     string
	Password string
	DB       string
}

type Admin struct {
	UserName string
	Name     string
	Pwd      string
}

type TestData struct {
	UserName string
	Name     string
	Pwd      string
}
