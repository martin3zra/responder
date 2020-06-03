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
In this example, we would respond with `http.StatusOK`, content-type of `application/json` and will wrote the data to the buffer.

Please check the [respond_test.go](https://github.com/martin3zra/respond/blob/master/respond_test.go) for more samples.