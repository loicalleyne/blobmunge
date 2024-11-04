# blobmunge
blobmunge provides helper functions for using Redpanda Connect's
Bloblang mapping language to munge structured data. blobmunge uses 
[MapStructure/v2](github.com/go-viper/mapstructure/v2) to decode data prior to processing.

## Features

- blobmunge.New() creates a new BlobMunger from a Bloblang mapping.
- BlobMunger.UpdateMapping() parses bloblang and returns an executor to be applied to input data.
- BlobMunger.ApplyBloblangMapping() executes a bloblang mapping on structured input data.
- InputMap() converts a structured input (json string or []byte, Go data type) into a map[string]any.

## 🚀 Install

Using blobmunge is easy. First, use `go get` to install the latest version
of the library.

```sh
go get -u github.com/loicalleyne/blobmunge@latest
```

## 💡 Usage

You can import `blobmunge` using:

```go
import "github.com/loicalleyne/blobmunge"
```

Create a new BlobMunger with a Bloblang rule, then apply it to some data.
```go
var jsonS1 string = `{"id":1234,"dev":"12345ert"}`
var bloblString string = `root.id = this.id\nroot.device = this.dev`
u, _ := blobmunge.New(bloblString)
b, _ := u.ApplyBloblangMapping(jsonS1)
fmt.Printf(string(b))
// {"device":"12345ert","id":1234}
```

## 💫 Show your support

Give a ⭐️ if this project helped you!

## License

blobmunge is released under the Apache 2.0 license. See [LICENCE.txt](LICENCE.txt)