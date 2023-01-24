# Responder

Wrapper in top of [http package](https://golang.org/pkg/net/http/).

## Install

You might as well just copy the [responder.go](https://github.com/martin3zra/responder/blob/master/responder.go) file into your own project (and the [responder_test.go](https://github.com/martin3zra/respond/blob/master/responder_test.go) while you're at it for future generations) rather than adding a dependency.

But it is maintained as a Go module which you can get with:

```bash
go get github.com/martin3zra/responder
```

## Usage

Import it:

```go
import (
	"github.com/martin3zra/responder"
)
```

The package will show you a clean way and guide you through the necessary steps to respond from your handlers.

```go
func (w http.ResponseWriter, r *http.Request) {
	//do your thing and respond
	data := map[string]interface{}{
		"customer": "..."
	}

	respond := responder.New(w)
	respond.OK(data)

	// If you whish to flash any data
	respond.With("error", "resource can not be created!").OK(data)
}

// without flashing data
{
  "data": {"customer": "..."},
  "flash": {
	"success": nil,
	"error": nil
  },
}

// Flashing data
{
  "data": {"customer": "..."},
  "flash": {
	"success": nil,
	"error": "resource can not be created!",
  }
}
```

In this example, we would respond with `http.StatusOK`, content-type of `application/json` and will write the data to the buffer.

Please check the [respond_test.go](https://github.com/martin3zra/responder/blob/master/responder_test.go) for more samples.

We also add a **_error formatter_** interface so you can add more information and custom code for your projects. You only have to implement his methods `Code`, `Error` and this other one that are optional `Description` and `InfoURL`

we also an **_error formatter_** [error.go](https://github.com/martin3zra/responder/blob/master/error.go) interface so you can add more information and custom code for your project errors. You only have to implement his methods `Code`, `Error`, `Description`, `InfoURL`. The interface implements from the core [errors pkg](https://golang.org/pkg/errors/) interface.

```go
type NotFound struct {}

func (NotFound) Status() int {
    return http.StatusNotFound
}

func (NotFound) Code() int {
	return 40401
}

func (NotFound) Error() string {
	return "resource not found"
}

func (NotFound) Description() *string {
	s := "resource not found long description"
	return &s
}

func (NotFound) InfoURL() *string {
	s := "resource not found URL"
	return &s
}

//Then when you use any of the respond alias just have to pass the
NotFound struct as parameter, and the output will be next one.

{
  "code": 41,
  "message": "Invalid credentials",
  "description": "The requested service needs valid credentials",
  "info_url": "https://url.com/to-your-errors-table/41"
}
```

As the methods `Description` and `InfoURL` are optional the resolver will ignore them when they are `nil`. But you can be embedded to avoid create these methods directly in your custom error.

```go
type NotFound struct {
	respond.ErrorDescriptor
}

//If you want a specific status code just implement Status()
func (NotFound) Status() int {
    return http.StatusNotFound
}

func (NotFound) Code() int {
	return 40401
}

func (NotFound) Error() string {
	return "resource not found"
}

//Then when you use any of the respond alias just have to pass the
NotFound struct as parameter, and the output will be next one.

{
  "code": 41,
  "message": "Invalid credentials"
}
```
