This folder contains the tools to interact with the polygon reduction system.

There are two existing ways to use the system:

1. Use the provided command-line interface; or
2. Build an application on top of the provided [SDK](./polyreduce-sdk-go/).

To use the command line, one has to build the binary using [Go](https://go.dev):

```
$ go build -o infloop
$ ./infloop
The command line tool of polygon reduction service.

Version:     v0.0.1

Usage:
  polyred [command]

Available Commands:
  config      Config the simplification target
  download    Download simplified model from polyred service
  help        Help about any command
  ping        ping polyred service
  run         Trigger polygon reduction to specific model
  upload      Upload .fbx model to polyred service

Flags:
  -h, --help   help for polyred

Use "polyred [command] --help" for more information about a command.
```

To use the SDK, one can import this package:

```go
import "changkun.de/x/infloop/tools/polyreduce-sdk-go"
```

This folder may be updated subsequently to release more features both on the command-line tool and SDK.