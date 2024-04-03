package common

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-solana/core"
)

// UnmarshalDataByEncoding Unmarshal data to string by encoding
func UnmarshalDataByEncoding(input []byte) ([]byte, string, error) {
	var (
		err      error
		data     interface{}
		encoding string
	)
	// Unmarshal
	if err = json.Unmarshal(input, &data); err != nil {
		return input, "", err
	}
	// get active type
	switch v := data.(type) {
	case string:
		input = core.DecodeBase58Str(v)
	// slice
	case []interface{}:
		// none data
		if len(v) == 0 {
			return nil, "", err
		}
		// decode to string
		switch v[1] {
		case "base58":
			encoding = "base58"
			input = core.DecodeBase58Str(v[0].(string))
		case "base64":
			encoding = "base64"
			input, _ = base64.StdEncoding.DecodeString(v[0].(string))
		default:
			return nil, "", fmt.Errorf("unsupported type: %s", v[1])
		}
	}
	return input, encoding, err
}
