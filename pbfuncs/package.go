package pbfuncs

import (
	b64 "encoding/base64"
	"net"

	"github.com/benthosdev/benthos/v4/public/bloblang"
)

func init() {
	ipv4StrAsBytesSpec := bloblang.NewPluginSpec().
		Param(bloblang.NewStringParam("key")).
		Description(`Returns an IPv4 address string as a 4 byte value similar to inet_aton.`).
		Example("", `root.srcip = ipv4str_as_bytes("65.66.67.68")`)
	err := bloblang.RegisterFunctionV2("ipv4str_as_bytes", ipv4StrAsBytesSpec, func(args *bloblang.ParsedParams) (bloblang.Function, error) {
		key, err := args.GetString("key")
		if err != nil {
			return nil, err
		}

		return func() (interface{}, error) {
			// Convert IP Address string to bytes
			possibleIP := net.IP.To4(net.ParseIP(key))
			ipBase64Enc := b64.StdEncoding.EncodeToString([]byte(possibleIP))
			return ipBase64Enc, nil
		}, nil
	})
	if err != nil {
		panic(err)
	}
}
