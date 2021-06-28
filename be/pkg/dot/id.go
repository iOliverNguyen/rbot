// Package cm (common) provides most used elements.
package dot

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

type IntID int64

const jsonNull = "null"

var rng = rand.Reader
var maxInt = big.NewInt(int64(1<<63 - 1))

// NewIntID mocks generating new id. In production, we need to use a more robust function.
func NewIntID() IntID {
	n, err := rand.Int(rng, maxInt)
	if err != nil {
		panic(err)
	}
	return IntID(n.Int64())
}

func (id IntID) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, 32)
	b = append(b, '"')
	b = strconv.AppendInt(b, int64(id), 10)
	b = append(b, '"')
	return b, nil
}

func (id *IntID) UnmarshalJSON(data []byte) error {
	if string(data) == jsonNull {
		*id = 0
		return nil
	}
	if data[0] == '"' {
		data = data[1 : len(data)-1]
	}
	i, err := strconv.ParseInt(string(data), 10, 64)
	*id = IntID(i)
	return err
}
