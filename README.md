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

* Web 服务可以服务 HTML 模版, 静态文件 。 以及一个 API JSON 。
* 解析和编码 type<->JSON 的简单辅助函数。
* 路由逻辑代码自定义。
* 中间件: 日志, 身份认证, 频率限制和节流。
* 用户 使用多重方式传递一个 token 和一个简单的角色基础认证实现身份认证和授权。
* 数据库未知数据层。使用 gorm 处理数据。
* 用户管理, 开票 (每个账号或者每个用户) 和 webhooks 管理. [开发中]
* 简单队列 (使用 Redis) 和队列任务的发布和订阅。
* Cron-like 定期任务的计划。

开发中的部分意味着，那些部分代码功能还未完成。 

## 快速开始

下面是一些帮助你快速启动项目的提示。

### 定义路由

使用 GoTenancy 处理路由， 你仅仅需要通过一个 `map[string]*GoTenancy.Route` 传入一个顶级路由。

例如, 如果在你的项目里面有如下路由:

`/task, /task/mine, /task/done, /ping`

你需要传递下面的 `map` 到 GoTenancy 的 `NewServer` 函数内:

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

`task` 和 `ping` 都是继承自 `http` 的 `ServeHTTP` 函数的类型, 例如:

```go
type Task struct{}

func (t *Task) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 在你自己的代码中处理路由逻辑的其余部分
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

你可以在它自己的包里面或者你的 `main` 包内定义 `Task`。

每个路由能选择加入特定的中间件，列表如下:

```go
// 路由表示具有可选中间件的 Web 处理程序。
type Route struct {
	// 中间件
	WithDB           bool // 增加 数据库连接到请求上下文
	Logger           bool // 写入请求信息输出
	EnforceRateLimit bool // 强制执行默认速率和限制限制

	// 授权
	MinimumRole model.Roles // 指明最小角色去访问这个路由

	Handler http.Handler // 这个 handler 将被执行
}
```

如何处理有参数的路由 `/task/detail/id-goes-here`:

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
	// 现在你可以调用数据和传入id(可以使用 AccountID 和 UserID)
	// 从请求上下文的 Auth 值
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

使用实例:

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

	// 连接数据库
	db := &data.DB{}

	if err := db.Open(*dn, *ds); err != nil {
		log.Fatal("unable to connect to the database:", err)
	}

	mux.DB = db

	isDev := false
	if *e == "dev" {
		isDev = true
	}

	// 如果 q 为 true，则设置为队列执行器的 pub/子订阅者
	executors := make(map[queue.TaskID]queue.TaskExecutor)
	// 如果你有自定义任务执行器，则可以使用你自己的实现填充此映射
    // 队列.任务执行器接口
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
* Redis 组件是 **必须的** 并且它和`队列` 包混在一起使用修改起来比较困难。
* 管理账号/用户控制器还没有完成
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