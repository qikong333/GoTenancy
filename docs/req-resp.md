---
layout: default
title: GoTenancy requests / responses
---

[back to main content](index.md)

# Requests & Responses

There's not much to say since the library uses the `net/http` package as-is. 
There are some useful functions that you might want to use.

The library uses the request context to store two important pieces: the 
database and the authentication types.

Here's an example of a typical request handler:

```go
func (t Tasks) index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(GoTenancy.ContextAuth).(GoTenancy.Auth)
	db := ctx.Value(GoTenancy.ContextDatabase).(*data.DB)

	...
}
```

We will see later how we can use those variables. 

For now, let's focus on how to parse incoming JSON and how to respond to 
requests.

### GoTenancy.ParseBody

If you need to turn JSON into a struct, you may use the ParseBody function 
like this.

```go
type task struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Done bool `json:"done"`
}
func (t Tasks) index(w http.ResponseWriter, r *http.Request) { {
	var task Task
	if err := GoTenancy.ParseBody(r.Body, &task); err != nil {
		GoTenancy.Respond(w, r, http.StatusBadRequest, err)
		return
	}
	...
}
```

It leads us to the Respond function, which takes a struct and returns JSON.

### GoTenancy.Respond

This function takes the `http.ResponseWriter`, the *`http.Request`, the status 
and the struct to turn into JSON.

```go
GoTenancy.Respond(w, r, http.StatusOK, task)
```

### GoTenancy.ServePage

If you need to render an HTML template, you may use the `GoTenancy.ServePage` 
function. You need to have your templates saved in the `templates` directory at 
the root of your project.

```go
func (t Task) index(w http.ResponseWriter, r *http.Request) {
	GoTenancy.ServePage(w, r, "tasks_index.html", nil)
}
```

The last parameter is the struct that will be sent to your template.

Next topic is the [Defining your handlers](handlers.md).