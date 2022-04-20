package bloblang

import (
	"fmt"

	"github.com/benthosdev/benthos/v4/public/bloblang"
)

func init() {
	crazyObjectSpec := bloblang.NewPluginSpec().
		Param(bloblang.NewInt64Param("keys"))

	err := bloblang.RegisterFunctionV2("crazy_object", crazyObjectSpec, func(args *bloblang.ParsedParams) (bloblang.Function, error) {
		keys, err := args.GetInt64("keys")
		if err != nil {
			return nil, err
		}

		return func() (interface{}, error) {
			obj := map[string]interface{}{}
			for i := 0; i < int(keys); i++ {
				obj[fmt.Sprintf("key%v", i)] = fmt.Sprintf("value%v", i)
			}
			return obj, nil
		}, nil
	})
	if err != nil {
		panic(err)
	}

	intoObjectSpec := bloblang.NewPluginSpec().
		Param(bloblang.NewStringParam("key"))

	err = bloblang.RegisterMethodV2("into_object", intoObjectSpec, func(args *bloblang.ParsedParams) (bloblang.Method, error) {
		key, err := args.GetString("key")
		if err != nil {
			return nil, err
		}

		return func(v interface{}) (interface{}, error) {
			return map[string]interface{}{key: v}, nil
		}, nil
	})
	if err != nil {
		panic(err)
	}
}
