---
layout: default
title: GoTenancy routing
---

[back to main content](index.md)

# Routing

In the book we build our own server and routing engine in less than 100 lines of code.

You define your top level routes in your package like this:

```go
package main

import "github.com/dstpierre/GoTenancy"

func main() {
	routes := make(map[string]*GoTenancy.Route)
	routes[""] = pages.New()
	routes["tasks"] = tasks.New()
}
```

Lets address the WTH are `pages.New()` and `tasks.New()` first.

This is the `GoTenancy.Route` struct:

```go
type Route struct {
	// middleware
	WithDB           bool
	Logger           bool
	EnforceRateLimit bool
	AllowCrossOrigin bool

	// authorization
	MinimumRole model.Roles

	Handler http.Handler
}
```

We need to build a `map` of our top level URL / ressources. For instance, if 
your app have the following routes:

**/tasks, /tasks/123, /about, /**

You would need to define two `Route`. In the example above we have the "catch-all" 
route `routes[""] = pages.New()` in the `pages` package and the route 
`routes["tasks"] = tasks.New()` in the `tasks` package.

*We discuss where to define your routes in the [Defining your handlers](handlers.md) section.*

As a minimum working example, you could do this to test the library:

```go
package main

import (
	"net/http"
	"github.com/dstpierre/GoTenancy"
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

The important aspect to grasp here is that the routes `map` entries correspond 
to your top level URL, the rest of the routing will be done in the `http.Handler`'s 
`ServeHTTP` function.

### An example of the tasks package

Lets continue our example and examine the what could be our `tasks` package.

The `routes.go` file will implement the `http.Handler`'s `ServeHTTP` as well as 
defining the `Route` for this top level `/tasks` part of our application.

```go
package tasks

import (
	"net/http"

	"github.com/dstpierre/GoTenancy"
	"github.com/snowlyg/GoTenancy/model"
)

type Tasks struct{}

func New() *GoTenancy.Route {
	return &GoTenancy.Route{
		Handler:     Tasks{},
		Logger:      true,
		MinimumRole: model.RoleUser,
		WithDB:      true,
	}
}

func (t Tasks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.index(w, r)
}

func (t Tasks) index(w http.ResponseWriter, r *http.Request) {
	GoTenancy.ServePage(w, r, "index.html", nil)
}
```

We define a `Tasks` struct and attach the implementation of `ServeHTTP` function.

The `index` handler for now simply render a template named `index.html`. Lets 
update this file to handle the routing we needed above:

****/tasks, /tasks/123**

```go
func (t Tasks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = GoTenancy.ShiftPath(r.URL.Path)
	if head == "" {
		t.list(w, r)
	} else {
		i, err := strconv.ParseInt(head, 10, 64)
		if err != nil {
			GoTenancy.NewError(err, http.StatusNotFound).Handler.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), contextID, i)
		if r.Method == http.MethodPut {
			t.update(w, r.WithContext(ctx))
		} else if r.Method == http.MethodDelete {
			t.delete(w, r.WithContext(ctx))
		}
	}
}
```

You can see how we continue to handle the routing inside our `tasks` package. You 
can use the `ShiftPath` function to handle pieces of your URLs, like we are doing 
for the `/tasks/123` route to extract the ID and we insert it into the request Context.

Next topic is the [Request/Response](req-resp.md).