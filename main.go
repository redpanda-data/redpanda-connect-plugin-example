package main

import (
	"context"

	"github.com/redpanda-data/benthos/v4/public/service"

	// Import full suite of FOSS connect plugins
	_ "github.com/redpanda-data/connect/public/bundle/free/v4"

	// Or, in order to import both FOSS and enterprise plugins, replace the
	// above with:
	// _ "github.com/redpanda-data/connect/public/bundle/enterprise/v4"

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
