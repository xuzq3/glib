package multicast

import (
	"github.com/google/uuid"
)

func NewSeqno() string {
	return uuid.New().String()
}
