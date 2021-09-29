package main

import (
	"context"

	"github.com/Jeffail/benthos/v3/public/service"

	// Import all standard Benthos components
	_ "github.com/Jeffail/benthos/v3/public/components/all"

	// Add your plugin packages here
	_ "github.com/benthosdev/benthos-plugin-example/bloblang"
	_ "github.com/benthosdev/benthos-plugin-example/cache"
	_ "github.com/benthosdev/benthos-plugin-example/input"
	_ "github.com/benthosdev/benthos-plugin-example/output"
	_ "github.com/benthosdev/benthos-plugin-example/processor"
)

func main() {
	service.RunCLI(context.Background())
}
