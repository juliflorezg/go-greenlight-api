package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

// This should return the JSON-encoded value for the movie
// runtime (in our case, it will return a string in the format "<runtime> mins").
func (rt Runtime) MarshalJSON() ([]byte, error) {
	quotedJSONValue := strconv.Quote(fmt.Sprintf("%d mins", rt))

	return []byte(quotedJSONValue), nil
}

// UnmarshalJSON converts a string in the format "<runtime> mins" to a Runtime type
func (rt *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt((parts[0]), 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*rt = Runtime(i)

	return nil
}
