# Otel GO Starter
_A simple means of getting the Open Telemetry global instrumentation configure and started_


## Get Started

Using the _Otel GO Starter_ allows for easy configuration of the global open telemetry instrumentations.

```shell
> go get github.com/MovieStoreGuy/otel-go-starter@latest
```

Then all that is required to do within main is the following:
```golang
package main

import (
    "context"

    otelstarter "github.com/MovieStoreGuy/otel-go-starter"
)

func main() {
    ctx, done := context.WithCancel(context.Background())
    defer done()

    defer otelstarter.Start(ctx).Shutdown()

    // Start the remainder of the application
}
```