// package blobmunge provides helper functions for using Redpanda Connect's
// bloblang mapping language to munge structured data. blobmunge uses 
// [MapStructure/v2] to decode data prior to processing.
// [MapStructure/v2]: github.com/go-viper/mapstructure/v2
package blobmunge

import (
	"github.com/go-viper/mapstructure/v2"
	json "github.com/goccy/go-json"
	"github.com/redpanda-data/benthos/v4/public/bloblang"
)

var (
	ErrUndefinedInput = errors.New("nil input")
	ErrUndefinedMapping = errors.New("nil mapping")
)

// ParseMapping parses bloblang and returns an excutor to be applied to
// input data.
func ParseMapping(bloblangMapping string) (*bloblang.Executor, error) {
	return bloblang.Parse(bloblangMapping)
}


// ApplyBloblangMapping executes a bloblang mapping on structured input data.
func ApplyBloblangMapping(input any, exe *bloblang.Executor) ([]byte, error) {
	if exe == nil {
		return nil, ErrUndefinedMapping
	}
	i, err := inputMap(input)
	if err != nil {
		return nil, err
	}
	// Execute the Bloblang mapping
	res, err := exe.Query(i)
	if err != nil {
		return nil, errors.New("bloblang error %w", err)
	}

	// Convert the result back into a JSON string
	jsonResult, err := json.Marshal(res)
	if err != nil {
		return nil, errors.New("json marshal error %w", err)
	}
	return jsonResult, nil
}

// InputMap takes structured input data and attempts to decode it to 
// map[string]any. Input data can be json in string or []byte, or any other
// Go data type which can be handled by [MapStructure/v2].
// [MapStructure/v2]: github.com/go-viper/mapstructure/v2
func InputMap(input any) (map[string]any, error) {
	m := map[string]any{}
	switch input := a.(type) {
	case nil:
		return nil, ErrUndefinedInput
	case map[string]any:
		input, nil
	case []byte:
		err := json.Unmarshal(input, &m)
		if err != nil {
			return nil, fmt.Errorf("%v : %v", ErrInvalidInput, err)
		}
	case string:
		err := json.Unmarshal([]byte(input), &m)
		if err != nil {
			return nil, fmt.Errorf("%v : %v", ErrInvalidInput, err)
		}
	default:
		err := mapstructure.Decode(a, &m)
		if err != nil {
			return nil, fmt.Errorf("%v : %v", ErrInvalidInput, err)
		}
	}
	return m, nil
}
