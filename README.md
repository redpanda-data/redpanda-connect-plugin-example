Benthos Plugin Example
======================

## WARNING: EXPERIMENTAL

This project demonstrates how to build your own Benthos component plugins into a
new service. There are [input](./input), [processor](./processor) and
[condition](./condition) examples.

## Build

`go build`

## Run

This new service comes with all the usual Benthos components plus all of your
custom plugins which you can use like any other type, with the only exception
being that the config field for plugin specific fields is always `plugin`. For
example, to use the example plugin components `gibberish`, `is_all_caps` and
`reverse` our config might look like this:

``` yaml
input:
  type: gibberish
  plugin:
    length: 100

pipeline:
  processors:
  - type: reverse
  - type: filter_parts
    filter_parts:
      type: is_all_caps 

output:
  type: stdout
```
