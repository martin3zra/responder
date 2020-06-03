# Respond

Wrapper in top of [http package](https://golang.org/pkg/net/http/).

## Install

You might as well just copy the [respond.go](https://github.com/martin3zra/respond/blob/master/respond.go) file into your own project (and the [respond_test.go](https://github.com/martin3zra/respond/blob/master/respond_test.go) while you're at it for future generations) rather than adding a dependency.

But it is maintained as a Go module which you can get with:

```bash
go get github.com/martin3zra/respond
```


## Usage
Import it:

```go
import (
	"github.com/martin3zra/respond"
)
```
The package will show you a clean way and guide you through the necessary steps to respond from your handlers.

```go
func (w http.ResponseWriter, r *http.Request) {
	//do your thing and respond
	respond.OK(w, data)
}
```
In this example, we would respond with `http.StatusOK`, content-type of `application/json` and will write the data to the buffer.

Please check the [respond_test.go](https://github.com/martin3zra/respond/blob/master/respond_test.go) for more samples.

We also add a ***error formatter*** interface so you can add more information and custom code for your projects. You only have to implement his methods `Code`, `Error` and this other one that are optional `Description` and `InfoURL`

we also an ***error formatter*** [error.go](https://github.com/martin3zra/respond/blob/master/error.go) interface so you can add more information and custom code for your project errors. You only have to implement his methods `Code`, `Error`, `Description`, `InfoURL`. The interface implements from the core [errors pkg](https://golang.org/pkg/errors/) interface.

```go
type NotFound struct {}

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