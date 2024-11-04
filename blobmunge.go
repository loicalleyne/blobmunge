// package blobmunge provides helper functions for using Redpanda Connect's
// bloblang mapping language to munge structured data. blobmunge uses
// [MapStructure/v2] to decode data prior to processing.
// [MapStructure/v2]: github.com/go-viper/mapstructure/v2
package blobmunge

import (
	"errors"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	json "github.com/goccy/go-json"
	"github.com/redpanda-data/benthos/v4/public/bloblang"
)

var (
	ErrUndefinedInput   = errors.New("nil input")
	ErrInvalidInput     = errors.New("invalid input")
	ErrUndefinedMapping = errors.New("nil mapping")
)

// BlobMunger
type BlobMunger struct {
	exe *bloblang.Executor
}

// New returns new BlobMunger or an error if mapping cannot be parsed.
func New(bloblangMapping string) (*BlobMunger, error) {
	var err error
	nb := &BlobMunger{}
	nb.exe, err = parseMapping(bloblangMapping)
	if err != nil {
		return nil, fmt.Errorf("%v, %w", ErrUndefinedMapping, err)
	}
	return nb, nil
}

// UpdateMapping method attempts to update the mapping rule, if an error is
// encountered it leaves the existing mapping in place.
func (b *BlobMunger) UpdateMapping(bloblangMapping string) error {
	var err error
	exe, err := parseMapping(bloblangMapping)
	if err != nil {
		return fmt.Errorf("error updating mapping, %w", ErrUndefinedMapping, err)
	}
	b.exe = exe
	return nil
}

// parseMapping parses bloblang and returns an executor to be applied to
// input data.
func parseMapping(bloblangMapping string) (*bloblang.Executor, error) {
	return bloblang.Parse(bloblangMapping)
}

// ApplyBloblangMapping executes a bloblang mapping on structured input data.
func (b *BlobMunger) ApplyBloblangMapping(input any) ([]byte, error) {
	i, err := InputMap(input)
	if err != nil {
		return nil, err
	}
	// Execute the Bloblang mapping
	res, err := b.exe.Query(i)
	if err != nil {
		return nil, errors.New("bloblang error " + err.Error())
	}

	// Convert the result back into a JSON string
	jsonResult, err := json.Marshal(res)
	if err != nil {
		return nil, errors.New("json marshal error " + err.Error())
	}
	return jsonResult, nil
}

// InputMap takes structured input data and attempts to decode it to
// map[string]any. Input data can be json in string or []byte, or any other
// Go data type which can be handled by [MapStructure/v2].
// [MapStructure/v2]: github.com/go-viper/mapstructure/v2
func InputMap(a any) (map[string]any, error) {
	m := map[string]any{}
	switch input := a.(type) {
	case nil:
		return nil, ErrUndefinedInput
	case map[string]any:
		return input, nil
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
