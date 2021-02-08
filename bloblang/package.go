package bloblang

import (
	"fmt"

	"github.com/Jeffail/benthos/v3/public/bloblang"
)

func init() {
	bloblang.RegisterFunction("crazy_object", func(args ...interface{}) (bloblang.Function, error) {
		// Expects a single integer argument
		var keys int

		if err := bloblang.NewArgSpec().IntVar(&keys).Extract(args); err != nil {
			return nil, err
		}

		return func() (interface{}, error) {
			obj := map[string]interface{}{}
			for i := 0; i < keys; i++ {
				obj[fmt.Sprintf("key%v", i)] = fmt.Sprintf("value%v", i)
			}
			return obj, nil
		}, nil
	})

	bloblang.RegisterMethod("into_object", func(args ...interface{}) (bloblang.Method, error) {
		// Expects a single string argument
		var key string

		if err := bloblang.NewArgSpec().StringVar(&key).Extract(args); err != nil {
			return nil, err
		}

		return func(v interface{}) (interface{}, error) {
			return map[string]interface{}{key: v}, nil
		}, nil
	})
}
