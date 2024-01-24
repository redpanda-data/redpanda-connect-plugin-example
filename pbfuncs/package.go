package pbfuncs

import (
	"fmt"
	"net"

	"github.com/benthosdev/benthos/v4/public/bloblang"
)

func init() {
	ipv4StrAsBytesSpec := bloblang.NewPluginSpec().
		Param(bloblang.NewStringParam("key"))

	err := bloblang.RegisterFunctionV2("ipv4str_as_bytes", ipv4StrAsBytesSpec, func(args *bloblang.ParsedParams) (bloblang.Function, error) {
		key, err := args.GetString("key")
		if err != nil {
			return nil, err
		}

		return func() (interface{}, error) {
			// Convert IP Address string to bytes
			possibleIP := net.IP.To4(net.ParseIP(key))
			return fmt.Sprintf("%c%c%c%c", possibleIP[0], possibleIP[1], possibleIP[2], possibleIP[3]), nil
		}, nil
	})
	if err != nil {
		panic(err)
	}
}
