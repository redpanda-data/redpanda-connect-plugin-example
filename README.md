Benthos Plugin Example
======================

This project demonstrates how to build your own Benthos component plugins and
run them. Once your plugins are written there are two ways of adding them to
a Benthos service.

## Build with Custom Main Func

You can create your own service with your plugins by copying the Benthos
[main func](https://github.com/Jeffail/benthos/blob/master/cmd/benthos/main.go) 
and adding your plugin packages as dependencies
[at the top](./cmd/benthos-plugin-example/main.go#L22).

`go build ./cmd/benthos-plugin-example`

## EXPERIMENTAL: Build Plugins as a Library

You can also compile your plugins into a library:

`go build -buildmode=plugin -o ./lib/benthos-plugin-example.so`

And then have a Benthos service load it from a directory:

`benthos -plugins-dir ./yourplugindir -c ./yourconfig.yaml`

NOTE: This requires exactly matching the build process and dependencies of the
original Benthos binary, this can be quite difficult without isolated build
environments.

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