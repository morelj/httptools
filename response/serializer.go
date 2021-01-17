package response

import "fmt"

func DefaultSerializer(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case []byte:
		return v, nil

	case string:
		return []byte(v), nil

	default:
		return nil, fmt.Errorf("Unsupported type %T", v)
	}
}
