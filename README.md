# GoTenancy
#### 一个用 Go 构建多租户平台和 web 应用的库

---

### 快速使用实例

你可以创建 main 包和复制 `docker-compose.yml` 文件。该项目需要 Redis 和数据库支持才能工作.

```go
package main

import (
    "net/http"


    "github.com/snowlyg/GoTenancy"
    "github.com/snowlyg/GoTenancy/model"
)

func main() {
	routes := make(map[string]*GoTenancy.Route)
	routes["test"] = &GoTenancy.Route{
		Logger:      true,
		MinimumRole: model.RolePublic,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			GoTenancy.Respond(w, r, http.StatusOK, "hello world!")
		}),
	}

	mux := GoTenancy.NewServer(routes)
	http.ListenAndServe(":8080", mux)
}
```

启动 docker 容器和项目:

```shell
$> docker-compose up
$> go run main.go
```

然后访问 localhost:8080:

```shell
$> curl http://localhost:8080/test
"hello? world!"
```

## 使用文档

* [安装](#安装)
* [包括什么？](#包括什么)
* [快速启动](#快速启动)
	- [定义路由](#定义路由)
	- [如何处理数据库](#如何处理数据库)
	- [使用JSON或者HTML响应请求](#使用JSON或者HTML响应请求)
	- [解析JSON](#解析JSON)
	- [从请求上下文获取当前数据库和用户](#从请求上下文获取当前数据库和用户)
* [状态和贡献](#状态和贡献)
* [运行测试](#运行测试)
* [参考项目](#参考项目)
* [许可证](#许可证)

## 安装

`go get github.com/snowlyg/GoTenancy@latest`

## 包括什么

该项目包括下面这些内容:

* Web server capable of serving HTML templates, static files. Also JSON for an API.
* Easy helper functions for parsing and encoding type<->JSON.
* Routing logic in your own code.
* Middlewares: logging, authentication, rate limiting and throttling.
* User authentication and authorization using multiple ways to pass a token and a simple role based authorization.
* Database agnostic data layer. Currently handling PostgreSQL.
* User management, billing (per account or per user) and webhooks management. [in dev]
* Simple queue (using Redis) and Pub/Sub for queuing tasks.
* Cron-like scheduling for recurring tasks.

The in dev part means that those parts needs some refactoring compare to what was built 
in the book. The vast majority of the code is there and working, but it's not "library" friendly 
at the moment.

## 快速开始

下面是一些帮助你快速启动项目的提示。

### 定义路由

You only need to pass the top-level routes that GoTenancy needs to handle via a `map[string]*GoTenancy.Route`.

For example, if you have the following routes in your web application:

`/task, /task/mine, /task/done, /ping`

You would pass the following `map` to GoTenancy's `NewServer` function:

```go
routes := make(map[string]*GoTenancy.Route)
routes["task"] = &GoTenancy.Route{
	Logger: true,
	WithDB: true,
	handler: task,
	...
}
routes["ping"] = &GoTenancy.Route(
	Logger: true,
	Handler: ping,
)
```

Where `task` and `ping` are types that implement `http`'s `ServeHTTP` function, for instance:

```go
type Task struct{}

func (t *Task) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// you handle the rest of the routing logic in your own code
	var head string
	head, r.URL.Path = GoTenancy.ShiftPath(r.URL.Path)
	if head =="/" {
		t.list(w, r)
	} else if head == "mine" {
		t.mine(w, r)
	}
	...
}
```

You may define `Task` in its own package or inside your `main` package.

Each route can opt-in to include specific middleware, here's the list:

```go
// Route represents a web handler with optional middlewares.
type Route struct {
	// middleware
	WithDB           bool // Adds the database connection to the request Context
	Logger           bool // Writes to the stdout request information
	EnforceRateLimit bool // Enforce the default rate and throttling limits

	// authorization
	MinimumRole model.Roles // Indicates the minimum role to access this route

	Handler http.Handler // The handler that will be executed
}
```

This is how you would handle parameterized route `/task/detail/id-goes-here`:

```go
func (t *Task) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = GoTenancy.ShiftPath(r.URL.Path)
	if head == "detail" {
		t.detail(w, r)
	}
}

func (t *Task) detail(w http.ResponseWriter, r *http.Request) {
	id, _ := GoTenancy.ShiftPath(r.URL.Path)
	// id = "id-goes-here
	// and now you may call the database and passing this id (probably with the AccountID and UserID)
	// from the Auth value of the request Context
}
```

### 如何处理数据库
`data` 包集成了 gorm 包处理数据库链接，
`data` 包有一个包含 `Connection` 字段数据库端点的 `DB` 类型。
调用 `http.ListenAndServe` 之前你需要初始化 `Server` 的 `DB` 字段:

```go
db := &data.DB{}

if err := db.Open(*dn, *ds); err != nil {
	log.Fatal("unable to connect to the database:", err)
}

mux.DB = db
```

Where `*dn` and `*ds` are flags containing "postgres" and 
"user=postgres password=postgres dbname=postgres sslmode=disable" for example,
respectively which are the driver name and the datasource connection string.

This is an example of what your `main` function could be:

```go
func main() {
	dn := flag.String("driver", "postgres", "name of the database driver to use, only postgres is supported at the moment")
	ds := flag.String("datasource", "", "database connection string")
	q := flag.Bool("queue", false, "set as queue pub/sub subscriber and task executor")
	e := flag.String("env", "dev", "set the current environment [dev|staging|prod]")
	flag.Parse()

	if len(*dn) == 0 || len(*ds) == 0 {
		flag.Usage()
		return
	}

	routes := make(map[string]*GoTenancy.Route)
	routes["test"] = &GoTenancy.Route{
		Logger:      true,
		MinimumRole: model.RolePublic,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			GoTenancy.Respond(w, r, http.StatusOK, "hello? Worker!")
		}),
	}

	mux := GoTenancy.NewServer(routes)

	// open the database connection
	db := &data.DB{}

	if err := db.Open(*dn, *ds); err != nil {
		log.Fatal("unable to connect to the database:", err)
	}

	mux.DB = db

	isDev := false
	if *e == "dev" {
		isDev = true
	}

	// Set as pub/sub subscriber for the queue executor if q is true
	executors := make(map[queue.TaskID]queue.TaskExecutor)
	// if you have custom task executor you may fill this map with your own implementation 
	// of queue.taskExecutor interface
	cache.New(*q, isDev, executors)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Println(err)
	}

}
```

### 使用JSON或者HTML响应请求

 `GoTenancy` 包有两个非常有用的函数:

**Respond**: 返回 JSON 格式数据:

```go
GoTenancy.Respond(w, r, http.StatusOK, oneTask)
```

**ServePage**: 返回 HTML 模版内容:

```go
GoTenancy.ServePage(w, r, "template.html", data)
```

### 解析JSON

调用 `GoTenancy.ParseBody` 辅助方法可以处理 JSON 数据解析，这是一个典型的 http 处理程序:

```go
func (t Type) do(w http.ResponseWriter, r *http.Request) {
	var oneTask MyTask
	if err := GoTenancy.ParseBody(r.Body, &oneTask); err != nil {
		GoTenancy.Respond(w, r, http.StatusBadRequest, err)
		return
	}
	...
}
```

### 从请求上下文获取当前数据库和用户

你肯定需要获取当前数据库的引用和登录用户。 通过请求的 `Context`实现它：

```go
func (t Type) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value(GoTenancy.ContextDatabase).(data.DB)
	auth := ctx.Value(ContextAuth).(Auth)

	// you may use the db.Connection in your own data implementation
	tasks := Tasks{DB: db.Connection}
	list, err := tasks.List(auth.AccountID, auth.UserID)
	if err != nil {
		GoTenancy.Respond(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, list)
}
```

## 状态和贡献

以下方面的内容，仍然有点粗糙:

* 测试覆盖不足。
* Redis 组件是 **必须的** 并且修改起来比较困难，而且它和`队列` 包一起使用。
* 管理 account/user 控制器还没有完成
* 开票控制器需要优化。
* 控制器应该位于 `internal` 包内。
* 仍然不确定数据包的编写方式是否自用/易于理解。
* 授权无法颗粒化。例如，如果 /task 需要 `model.RoleUser` ，/task/delete 不能使用 `model.RoleAdmin` 作为 `MinimumRole`。

## 运行测试

```shell
$> go test -tags mem ./...
```

## 参考项目

[dstpierre/gosaas](https://github.com/dstpierre/gosaas) 

## 许可证

[MIT](https://github.com/snowlyg/GoTenancy/blob/master/LICENSE)