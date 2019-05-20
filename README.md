Benthos Plugin Example
======================

This project demonstrates the recommended way to build your own Benthos
component plugins and run them in a custom distribution.

## Build

Start by writing your plugins where ever you like, there are examples in this
repo for [inputs][inputs], [processors][processors], [conditions][conditions]
and [outputs][outputs] to copy from.

Next, author a main file that calls `service.Run()` and imports your plugins
[as shown in this example][plugin-main]:

``` go
package main

import (
	"github.com/Jeffail/benthos/lib/service"

	// Add your plugin packages here
	_ "github.com/benthosdev/benthos-plugin-example/input"
	_ "github.com/benthosdev/benthos-plugin-example/output"
	_ "github.com/benthosdev/benthos-plugin-example/processor"
)

func main() {
	service.Run()
}
```

Finally, build your custom main func:

`go build ./cmd/benthos-plugin-example`

## Run

The new service you've built will come with all of the usual Benthos components
plus all of your custom plugins, which you can use like any other type. The only
difference between your plugins and original Benthos components is that the
config field for plugin specific fields is always `plugin`.

For example, to use the example plugin components `gibberish`, `is_all_caps`,
`reverse` and `blue_stdout` our config might look like this:

``` yaml
input:
  type: gibberish
  plugin:
    length: 80

pipeline:
  processors:
  - type: throttle
    throttle:
      period: 1s
  - type: reverse
  - type: filter_parts
    filter_parts:
      type: not
      not:
        type: is_all_caps

output:
  type: blue_stdout
```

And you can run it like this:

`./benthos-plugin-example -c ./yourconfig.yaml`

[plugin-main]: ./cmd/benthos-plugin-example/main.go#L22
[inputs]: ./input
[processors]: ./processor
[conditions]: ./condition
[outputs]: ./output