package requests

import (
	"encoding/json"
)

// Serializer is a function
// that can convert arbitrary data
// to bytes in some format.
type Serializer = func(v any) ([]byte, error)

// Deserializer is a function
// that can read data in some format
// and store the result in v.
type Deserializer = func(data []byte, v any) error

var (
	// JSONSerializer is used by BodyJSON and Builder.BodyJSON.
	// The default serializer may be changed in a future version of requests.
	JSONSerializer Serializer = json.Marshal
	// JSONDeserializer is used by ToJSON and Builder.ToJSON.
	// The default deserializer may be changed in a future version of requests.
	JSONDeserializer Deserializer = json.Unmarshal
)
