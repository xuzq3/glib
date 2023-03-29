package ws

import (
	"github.com/google/uuid"
	"math"
)

const abortIndex int8 = math.MaxInt8 / 2

func NewSeqno() string {
	return uuid.New().String()
}
