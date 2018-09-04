Benthos Plugin Example
======================

## WARNING: EXPERIMENTAL

This project demonstrates how to build a Benthos plugin. Once built, these
plugins can be dynamically loaded into Benthos by placing it into a directory
and pointing Benthos to that directory with the flag `plugins-dir`
(`/usr/lib/benthos/plugins`, by default.)

## Build

`go build -buildmode=plugin`

## Run

Once Benthos is run with your plugin loaded you can use it like any other type,
with the only exception being that the config field for plugin specific fields
is always `plugin`. For example, our example plugin config would look like this:

``` yaml
input:
  type: example
  plugin:
    length: 100
```

You can also print documentation from Benthos for any loaded plugins with
`--list-<component>-plugins`. For example, to print documentation any input
plugins that are loaded:

`benthos --list-input-plugins`
