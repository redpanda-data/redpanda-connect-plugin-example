Redpanda Connect Plugin Example
===============================

This project demonstrates the recommended way to build your own Redpanda Connect component plugins and run them in a custom distribution.

## Build

Start by writing your plugins where ever you like, there are examples in this repo for [bloblang functions and methods][bloblang], [inputs][inputs], [processors][processors] and [outputs][outputs] to copy from.

Next, author a main file that calls `service.Run()` and imports your plugins [as shown in this example][plugin-main]:

```go
package main

import (
	"context"

	"github.com/redpanda-data/benthos/v4/public/service"

	// Import full suite of FOSS connect plugins
	_ "github.com/redpanda-data/connect/public/bundle/free/v4"

	// Add your plugin packages here
	_ "github.com/redpanda-data/redpanda-connect-plugin-example/bloblang"
	_ "github.com/redpanda-data/redpanda-connect-plugin-example/cache"
	_ "github.com/redpanda-data/redpanda-connect-plugin-example/input"
	_ "github.com/redpanda-data/redpanda-connect-plugin-example/output"
	_ "github.com/redpanda-data/redpanda-connect-plugin-example/processor"
)

func main() {
	service.RunCLI(context.Background())
}
```

Finally, build your custom main func:

```sh
go build
```

Alternatively build it as a Docker image with:

```sh
docker build . -t redpanda-connect-plugin-example
```

## Upgrading

In order to track the latest upstream changes to Redpanda Connect, upgrade the `github.com/redpanda-data/connect/v4` package and all other transient dependencies will be pulled up to the latest versions:

```sh
go get github.com/redpanda-data/connect/v4@latest
go mod tidy
```

## Testing

There are few examples of unit tests for plugin components in this repo. The notable examples are the [gibberish input tests][gibberish.input.tests] which demonstrates how to test config validation within your component constructors, and the [reverse processor tests][reverse.processor.tests] which tests the processor behaviour and also demonstrates testing a component that uses `*service.Logger` and `*service.Metrics`.

## Run

The new service you've built will come with all of the usual Redpanda Connect components plus all of your custom plugins, which you can use like any other type. For example, to use the example plugin components `gibberish`, `reverse` and `blue_stdout`, and our new Bloblang function `crazy_object` and method `into_object`, our config might look like this:

```yaml
input:
  gibberish:
    length: 80

pipeline:
  threads: 1
  processors:
  - sleep:
      duration: 1s
  - reverse: {}
  - bloblang: |
      root.gibberish = content()
      root.more_stuff = crazy_object(10).into_object("foo")

output:
  blue_stdout: {}
```

And you can run it like this:

```sh
./redpanda-connect-plugin-example -c ./yourconfig.yaml
```

For more examples on how to configure your plugins check out [`./config`](./config).

[plugin-main]: ./main.go#L15
[inputs]: ./input
[gibberish.input.tests]: ./input/gibberish_test.go
[processors]: ./processor
[reverse.processor.tests]: ./processor/reverse_test.go
[bloblang]: ./bloblang
[outputs]: ./output
