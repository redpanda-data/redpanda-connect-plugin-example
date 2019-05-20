package main

import (
	"github.com/Jeffail/benthos/lib/service"

	// Add your plugin packages here
	_ "github.com/benthosdev/benthos-plugin-example/cache"
	_ "github.com/benthosdev/benthos-plugin-example/condition"
	_ "github.com/benthosdev/benthos-plugin-example/input"
	_ "github.com/benthosdev/benthos-plugin-example/output"
	_ "github.com/benthosdev/benthos-plugin-example/processor"
)

//------------------------------------------------------------------------------

func main() {
	service.Run()
}

//------------------------------------------------------------------------------
