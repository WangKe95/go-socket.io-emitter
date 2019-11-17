go-socket.io-emitter
===

A Golang implementation of socket.io-emitter.

Installation and development
---
```sh
$ go get github.com/WangKe95/go-socket.io-emitter
```

Usage
---

Example:

```go
package main

import (
    socketIO "github.com/WangKe95/go-socket.io-emitter"
)

func main() {
    emitter, _ := socketIO.NewEmitter(&socketIO.NewEmitterOpts{
        Host: "localhost",
        Port: 6379,
    })
    
    emitter.Emit("message", "Hello World!")
}
```


