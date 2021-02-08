Benthos Plugin Example
======================

This project demonstrates the recommended way to build your own Benthos component plugins and run them in a custom distribution.

## Build

Start by writing your plugins where ever you like, there are examples in this repo for [bloblang functions and methods][bloblang], [inputs][inputs], [processors][processors] and [outputs][outputs] to copy from.

Next, author a main file that calls `service.Run()` and imports your plugins [as shown in this example][plugin-main]:

```go
package main

import (
	"github.com/Jeffail/benthos/v3/lib/service"

	// Add your plugin packages here
	_ "github.com/benthosdev/benthos-plugin-example/bloblang"
	_ "github.com/benthosdev/benthos-plugin-example/input"
	_ "github.com/benthosdev/benthos-plugin-example/output"
	_ "github.com/benthosdev/benthos-plugin-example/processor"
)

func main() {
	service.Run()
}
```

Finally, build your custom main func:

```sh
go build ./cmd/benthos-plugin-example
```

Alternatively build it as a Docker image with:

```sh
go mod vendor
docker build . -t benthos-plugin-example
```

### Serverless

In order to build a serverless distribution of Benthos you follow the same steps, but instead of importing the `service` package you import a serverless package corresponding with your target platform [as shown in this AWS Lambda example][plugin-lambda-main].

Then build and zip:

```sh
GOOS=linux go build ./cmd/benthos-lambda-plugin-example
zip -m -j benthos-lambda.zip ./benthos-lambda-plugin-example
```

## Run

The new service you've built will come with all of the usual Benthos components plus all of your custom plugins, which you can use like any other type. The only difference between your plugins and original Benthos components is that the config field for plugin specific fields is always `plugin`.

For example, to use the example plugin components `gibberish`, `reverse` and `blue_stdout`, and our new Bloblang function `crazy_object` and method `into_object`, our config might look like this:

```yaml
input:
  type: gibberish
  plugin:
    length: 80

pipeline:
  processors:
  - throttle:
      period: 1s
  - type: reverse
  - bloblang: |
      root.gibberish = content()
      root.more_stuff = crazy_object(10).into_object("foo")

output:
  type: blue_stdout
```

And you can run it like this:

```sh
./benthos-plugin-example -c ./yourconfig.yaml
```

For more examples on how to configure your plugins check out [`./config`](./config).

[plugin-main]: ./cmd/benthos-plugin-example/main.go#L22
[plugin-lambda-main]: ./cmd/benthos-lambda-plugin-example/main.go#L22
[inputs]: ./input
[processors]: ./processor
[bloblang]: ./bloblang
[outputs]: ./output
