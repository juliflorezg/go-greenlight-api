package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

// This should return the JSON-encoded value for the movie
// runtime (in our case, it will return a string in the format "<runtime> mins").
func (rt Runtime) MarshalJSON() ([]byte, error) {
	quotedJSONValue := strconv.Quote(fmt.Sprintf("%d mins", rt))

	return []byte(quotedJSONValue), nil
}
