escvpnetgo
======

escvpnetgo is a golang implementation for sending ESC/VP.NET protocol commands.

Examples
--------

```go
package main

import (
    "fmt"
    "os"
    "github.com/germanosin/escvpnetgo"
)

func main() {
    c, err := escvpnetgo.NewESCVPNET("10.150.136.32:3629")

    if err != nil {
        fmt.Printf("Unable to connect server %s\n", err.Error())
        os.Exit(1)
    }

    defer c.Close()

    result, err := c.Execute("LAMP?")

    fmt.Printf("Execute result: %s\n", result)
}
```

